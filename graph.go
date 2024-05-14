package goraff

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type NodeAction interface {
	Do(s *StateNode, r *GraphStateReader, triggeringNS *StateNodeReader) error
}

// Node represents a node in the graph
type Node struct {
	Action NodeAction
	Name   string
}

// Graph represents a graph of nodes
type Graph struct {
	nodes      []*Node
	entrypoint *Node
	state      *StateGraph
	edges      map[string][]*Edge
}

func New() *Graph {
	return &Graph{}
}

func NewWithState(s *StateGraph) *Graph {
	return &Graph{
		state: s,
	}
}

func (g *Graph) StateReadOnly() *GraphStateReader {
	if g.state == nil {
		g.state = &StateGraph{}
	}
	return g.state.Reader()
}

func (g *Graph) State() *StateGraph {
	if g.state == nil {
		g.state = &StateGraph{}
	}
	return g.state
}

func (g *Graph) Len() int {
	return len(g.nodes)
}

func (g *Graph) AddNode(a NodeAction) string {
	id := uuid.New().String()
	return g.AddNodeWithName(id, a)
}

func (g *Graph) AddNodeWithName(name string, a NodeAction) string {
	n := &Node{Action: a, Name: name}
	g.nodes = append(g.nodes, n)
	return n.Name
}

func (g *Graph) SetEntrypoint(id string) {
	n := g.nodeByID(id)
	g.entrypoint = n
}

func (g *Graph) nodeByID(id string) *Node {
	for _, n := range g.nodes {
		if n.Name == id {
			return n
		}
	}
	return nil
}

type ErrNodeNotFound struct {
	ID string
}

func (e ErrNodeNotFound) Error() string {
	return "node not found: " + e.ID
}

func (g *Graph) AddEdge(fromID, toID string, condition FollowIf) error {
	from := g.nodeByID(fromID)
	if from == nil {
		return ErrNodeNotFound{
			ID: fromID,
		}
	}
	to := g.nodeByID(toID)
	if to == nil {
		return ErrNodeNotFound{
			ID: toID,
		}
	}
	if g.edges == nil {
		g.edges = make(map[string][]*Edge)
	}
	e := &Edge{From: from, To: to, Condition: condition}
	if _, ok := g.edges[from.Name]; !ok {
		g.edges[from.Name] = []*Edge{}
	}
	g.edges[from.Name] = append(g.edges[from.Name], e)
	return nil
}

func (g *Graph) Go() error {
	if g.state == nil {
		g.state = &StateGraph{}
	}
	return g.flowMgr()
}

type nextNode struct {
	Node         *Node
	triggeringNS *StateNode
}

func (g *Graph) flowMgr() error {
	if g.entrypoint == nil {
		return fmt.Errorf("entrypoint not set")
	}

	completedCh := make(chan nextNode, 10)
	var wg sync.WaitGroup

	completedCh <- nextNode{
		Node:         g.entrypoint,
		triggeringNS: nil,
	}
	wg.Add(1) // Increment for the initial node

	fmt.Println("starting node", g.entrypoint.Name)
	var foundErr error
	mut := sync.Mutex{}
	go func() {
		for n := range completedCh {
			go func(n nextNode) {
				fmt.Println("completed node", n.Node.Name)
				defer wg.Done() // Ensure we mark this goroutine as done on finish
				if foundErr != nil {
					return
				}
				nextNodes, compeltedState, err := g.runNode(n.Node, n.triggeringNS.Reader())
				if err != nil {
					fmt.Println("error running node, letting all nodes drain: ", n.Node.Name)
					mut.Lock()
					foundErr = fmt.Errorf("error running node: %w", err)
					mut.Unlock()
					return
				}
				for _, next := range nextNodes {
					fmt.Println("adding node", next.Name)
					wg.Add(1) // Increment for each new node
					completedCh <- nextNode{
						Node:         next,
						triggeringNS: compeltedState,
					}
				}
			}(n)
		}
	}()

	wg.Wait()          // Wait for all goroutines to finish
	close(completedCh) // Safe to close here as no more writes will happen
	return foundErr
}

func (g *Graph) runNode(n *Node, triggeringNS *StateNodeReader) ([]*Node, *StateNode, error) {
	s := g.state.NewNodeState(n.Name)
	r := g.state.Reader()
	err := n.Action.Do(s, r, triggeringNS)
	if err != nil {
		return nil, nil, err
	}
	nextNodes := []*Node{}
	s.MarkDone()
	if edges, ok := g.edges[n.Name]; ok {
		for _, e := range edges {
			if e.TriggersMet(r) {
				nextNodes = append(nextNodes, e.To)
			}
		}
	}
	return nextNodes, s, nil
}
