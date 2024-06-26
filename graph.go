package goraff

import (
	"fmt"

	"github.com/google/uuid"
)

type GraphChangeNotification struct {
	NodeID string
}

type ChangeNotifier interface {
	Notify(notification GraphChangeNotification)
}

// Graph manages the state of all nodes in the graph
type Graph struct {
	id       string
	nodes    []*Node
	Notifier ChangeNotifier
}

func (s *Graph) NewNode(name string, trigeredBy []*ReadableNode) *Node {
	ns := &Node{name: name, notifier: s.Notifier, triggeredBy: trigeredBy}
	s.nodes = append(s.nodes, ns)
	return ns
}

func (s *Graph) NodeByName(name string) []*Node {
	result := []*Node{}
	for _, ns := range s.nodes {
		if ns.name == name {
			result = append(result, ns)
		}
	}
	return result
}

// Gets the first node
func (s *Graph) FirstNodeByName(name string) *Node {
	for _, ns := range s.nodes {
		if ns.name == name {
			return ns
		}
	}
	return nil
}

func (s *Graph) NodeByID(id string) *Node {
	for _, ns := range s.nodes {
		if ns.Get().ID() == id {
			return ns
		}
	}
	return nil
}

func NewReadableGraph(g *Graph) *ReadableGraph {
	return &ReadableGraph{g}
}

// ReadableGraph is a read only view of the state
type ReadableGraph struct {
	graph *Graph
}

func (s *ReadableGraph) NodeByID(id string) (*ReadableNode, error) {
	r := s.graph.NodeByID(id)
	if r == nil {
		return nil, fmt.Errorf("Node with id %s not found", id)
	}
	return &ReadableNode{node: r}, nil
}

func (s *ReadableGraph) FirstNodeByName(name string) (*ReadableNode, error) {
	st := s.graph.FirstNodeByName(name)
	if st == nil {
		return nil, fmt.Errorf("Node with name %s not found", name)
	}
	return &ReadableNode{node: st}, nil
}

func (s *ReadableGraph) Node(id string) (*ReadableNode, error) {
	r := s.graph.NodeByID(id)
	if r == nil {
		return nil, fmt.Errorf("Node with id %s not found", id)
	}
	return &ReadableNode{node: r}, nil
}

func (s *ReadableGraph) NodeIDs() []string {
	ids := []string{}
	for _, ns := range s.graph.nodes {
		ids = append(ids, ns.Get().ID())
	}
	return ids
}

func (s *ReadableGraph) NodeNames() []string {
	names := []string{}
	for _, ns := range s.graph.nodes {
		names = append(names, ns.name)
	}
	return names
}

func (s *ReadableGraph) ID() string {
	if s.graph.id == "" {
		id := uuid.New().String()
		s.graph.id = id
	}
	return s.graph.id
}
