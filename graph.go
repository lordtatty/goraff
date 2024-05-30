package goraff

import (
	"fmt"

	"github.com/google/uuid"
)

// Graph manages the state of all nodes in the graph
type Graph struct {
	id       string
	nodes    []*Node
	notifier *GraphNotifier
}

func (s *Graph) Notifier() *GraphNotifier {
	if s.notifier == nil {
		s.notifier = &GraphNotifier{}
	}
	return s.notifier
}

func (s *Graph) NewNodeState(name string) *Node {
	// Else create a new node state
	ns := &Node{name: name, notifier: s.notifier}
	s.nodes = append(s.nodes, ns)
	return ns
}

func (s *Graph) NodeStateByName(name string) []*Node {
	// First see if we have this node state
	result := []*Node{}
	for _, ns := range s.nodes {
		if ns.name == name {
			result = append(result, ns)
		}
	}
	return result
}

func (s *Graph) FirstNodeStateByName(name string) *Node {
	// First see if we have this node state
	for _, ns := range s.nodes {
		if ns.name == name {
			return ns
		}
	}
	return nil
}

func (s *Graph) NodeByID(id string) *Node {
	// First see if we have this node state
	for _, ns := range s.nodes {
		if ns.Reader().ID() == id {
			return ns
		}
	}
	return nil
}

func (s *Graph) Reader() *ReadableGraph {
	return &ReadableGraph{s}
}

// ReadableGraph is a read only view of the state
type ReadableGraph struct {
	state *Graph
}

func (s *ReadableGraph) NodeByID(id string) (*ReadableNode, error) {
	r := s.state.NodeByID(id)
	if r == nil {
		return nil, fmt.Errorf("Node state with id %s not found", id)
	}
	return &ReadableNode{node: r}, nil
}

func (s *ReadableGraph) FirstNodeStateByName(name string) (*ReadableNode, error) {
	st := s.state.FirstNodeStateByName(name)
	if st == nil {
		return nil, fmt.Errorf("Node state with name %s not found", name)
	}
	return &ReadableNode{node: st}, nil
}

func (s *ReadableGraph) Node(id string) (*ReadableNode, error) {
	r := s.state.NodeByID(id)
	if r == nil {
		return nil, fmt.Errorf("Node state with id %s not found", id)
	}
	return &ReadableNode{node: r}, nil
}

func (s *ReadableGraph) NodeIDs() []string {
	ids := []string{}
	for _, ns := range s.state.nodes {
		ids = append(ids, ns.Reader().ID())
	}
	return ids
}

func (s *ReadableGraph) ID() string {
	if s.state.id == "" {
		id := uuid.New().String()
		s.state.id = id
	}
	return s.state.id
}
