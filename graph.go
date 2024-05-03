package goraff

import "github.com/google/uuid"

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

func (g *Graph) AddEdge(fromID, toID string, condition EdgeCondition) error {
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

func (g *Graph) Go() {
	if g.state == nil {
		g.state = &State{}
	}
	if g.entrypoint != nil {
		g.runNode(g.entrypoint, "")
	}
}

func (g *Graph) runNode(n *Node, triggeringNodeID string) {
	s := g.state.NodeState(n.ID())
	r := g.state.ReadOnly()
	n.Action.Do(s, r, triggeringNodeID)
	nextNodes := []*Node{}
	if edges, ok := g.edges[n.ID()]; ok {
		for _, e := range edges {
			if e.Match(r) {
				nextNodes = append(nextNodes, e.To)
			}
		}
	}
	if len(nextNodes) == 0 {
		return
	}
	for _, nextNode := range nextNodes {
		g.runNode(nextNode, n.ID())
	}
}
