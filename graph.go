package goraff

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type NodeAction interface {
	Do(s *NodeState, r *StateReadOnly, triggeringNodeID string) error
}

// Node represents a node in the graph
type Node struct {
	id     string
	Action NodeAction
	Name   string
}

func (n *Node) ID() string {
	if n.id == "" {
		n.id = uuid.New().String()
	}
	return n.id
}

// Graph represents a graph of nodes
type Graph struct {
	nodes      []*Node
	entrypoint *Node
	state      *State
	edges      map[string][]*Edge
}

func New() *Graph {
	return &Graph{}
}

func NewWithState(s *State) *Graph {
	return &Graph{
		state: s,
	}
}

func (g *Graph) State() *StateReadOnly {
	return g.state.ReadOnly()
}

func (g *Graph) Len() int {
	return len(g.nodes)
}

func (g *Graph) AddNode(a NodeAction) string {
	n := &Node{Action: a}
	g.nodes = append(g.nodes, n)
	return n.ID()
}

func (g *Graph) AddNodeWithName(a NodeAction, name string) string {
	n := &Node{Action: a, Name: name}
	g.nodes = append(g.nodes, n)
	return n.ID()
}

func (g *Graph) SetEntrypoint(id string) {
	n := g.nodeByID(id)
	g.entrypoint = n
}

func (g *Graph) nodeByID(id string) *Node {
	for _, n := range g.nodes {
		if n.ID() == id {
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
	if _, ok := g.edges[from.ID()]; !ok {
		g.edges[from.ID()] = []*Edge{}
	}
	g.edges[from.ID()] = append(g.edges[from.ID()], e)
	return nil
}

func (g *Graph) Go() error {
	if g.state == nil {
		g.state = &State{}
	}
	return g.flowMgr()
}

type nextNode struct {
	Node         *Node
	triggeringID string
}

func (g *Graph) flowMgr() error {
	if g.entrypoint == nil {
		return fmt.Errorf("entrypoint not set")
	}

	completedCh := make(chan nextNode, 10)
	var wg sync.WaitGroup

	completedCh <- nextNode{
		Node:         g.entrypoint,
		triggeringID: "",
	}
	wg.Add(1) // Increment for the initial node

	fmt.Println("starting node", g.entrypoint.ID())
	var foundErr error
	mut := sync.Mutex{}
	go func() {
		for n := range completedCh {
			go func(n nextNode) {
				fmt.Println("completed node", n.Node.ID())
				defer wg.Done() // Ensure we mark this goroutine as done on finish
				if foundErr != nil {
					return
				}
				nextNodes, err := g.runNode(n.Node, n.triggeringID)
				if err != nil {
					fmt.Println("error running node, letting all nodes drain: ", n.Node.ID())
					mut.Lock()
					foundErr = fmt.Errorf("error running node: %w", err)
					mut.Unlock()
					return
				}
				for _, next := range nextNodes {
					fmt.Println("adding node", next.ID())
					wg.Add(1) // Increment for each new node
					completedCh <- nextNode{
						Node:         next,
						triggeringID: n.Node.ID(),
					}
				}
			}(n)
		}
	}()

	wg.Wait()          // Wait for all goroutines to finish
	close(completedCh) // Safe to close here as no more writes will happen
	return foundErr
}

func (g *Graph) runNode(n *Node, triggeringNodeID string) ([]*Node, error) {
	s := g.state.NodeStateUpsert(n.ID())
	s.Set("name", n.Name)
	r := g.state.ReadOnly()
	err := n.Action.Do(s, r, triggeringNodeID)
	if err != nil {
		return nil, err
	}
	nextNodes := []*Node{}
	s.MarkDone()
	if edges, ok := g.edges[n.ID()]; ok {
		for _, e := range edges {
			if e.TriggersMet(r) {
				nextNodes = append(nextNodes, e.To)
			}
		}
	}
	return nextNodes, nil
}
