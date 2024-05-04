package goraff

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type NodeAction interface {
	Do(s *NodeState, r *StateReadOnly, triggeringNodeID string)
}

// Node represents a node in the graph
type Node struct {
	id     string
	Action NodeAction
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

	nodeCh := make(chan nextNode, 10)
	var wg sync.WaitGroup

	nodeCh <- nextNode{
		Node:         g.entrypoint,
		triggeringID: "",
	}
	wg.Add(1) // Increment for the initial node

	fmt.Println("starting node", g.entrypoint.ID())
	go func() {
		for n := range nodeCh {
			go func(n nextNode) {
				fmt.Println("completed node", n.Node.ID())
				defer wg.Done() // Ensure we mark this goroutine as done on finish
				nextNodes := g.runNode(n.Node, n.triggeringID)
				for _, next := range nextNodes {
					fmt.Println("adding node", next.ID())
					wg.Add(1) // Increment for each new node
					nodeCh <- nextNode{
						Node:         next,
						triggeringID: n.Node.ID(),
					}
				}
			}(n)
		}
	}()

	wg.Wait()     // Wait for all goroutines to finish
	close(nodeCh) // Safe to close here as no more writes will happen
	return nil
}

func (g *Graph) runNode(n *Node, triggeringNodeID string) []*Node {
	s := g.state.NodeStateUpsert(n.ID())
	r := g.state.ReadOnly()
	n.Action.Do(s, r, triggeringNodeID)
	nextNodes := []*Node{}
	s.MarkDone()
	if edges, ok := g.edges[n.ID()]; ok {
		for _, e := range edges {
			if e.TriggersMet(r) {
				nextNodes = append(nextNodes, e.To)
			}
		}
	}
	return nextNodes
}
