package goraff

import (
	"fmt"

	"github.com/google/uuid"
)

// StateGraph manages the state of all nodes in the graph
type StateGraph struct {
	id         string
	nodeStates []*StateNode
	notifier   *StateNotifier
}

func (s *StateGraph) Notifier() *StateNotifier {
	if s.notifier == nil {
		s.notifier = &StateNotifier{}
	}
	return s.notifier
}

func (s *StateGraph) NewNodeState(name string) *StateNode {
	// Else create a new node state
	ns := &StateNode{name: name, notifier: s.notifier}
	s.nodeStates = append(s.nodeStates, ns)
	return ns
}

func (s *StateGraph) NodeStateByName(name string) []*StateNode {
	// First see if we have this node state
	result := []*StateNode{}
	for _, ns := range s.nodeStates {
		if ns.name == name {
			result = append(result, ns)
		}
	}
	return result
}

func (s *StateGraph) FirstNodeStateByName(name string) *StateNode {
	// First see if we have this node state
	for _, ns := range s.nodeStates {
		if ns.name == name {
			return ns
		}
	}
	return nil
}

func (s *StateGraph) NodeStateByID(id string) *StateNode {
	// First see if we have this node state
	for _, ns := range s.nodeStates {
		if ns.Reader().ID() == id {
			return ns
		}
	}
	return nil
}

func (s *StateGraph) Reader() *GraphStateReader {
	return &GraphStateReader{s}
}

// GraphStateReader is a read only view of the state
type GraphStateReader struct {
	state *StateGraph
}

func (s *GraphStateReader) NodeStateByID(id string) (*StateNodeReader, error) {
	r := s.state.NodeStateByID(id)
	if r == nil {
		return nil, fmt.Errorf("Node state with id %s not found", id)
	}
	return &StateNodeReader{r}, nil
}

func (s *GraphStateReader) FirstNodeStateByName(name string) (*StateNodeReader, error) {
	st := s.state.FirstNodeStateByName(name)
	if st == nil {
		return nil, fmt.Errorf("Node state with name %s not found", name)
	}
	return &StateNodeReader{st}, nil
}

func (s *GraphStateReader) NodeState(id string) (*StateNodeReader, error) {
	r := s.state.NodeStateByID(id)
	if r == nil {
		return nil, fmt.Errorf("Node state with id %s not found", id)
	}
	return &StateNodeReader{r}, nil
}

func (s *GraphStateReader) NodeIDs() []string {
	ids := []string{}
	for _, ns := range s.state.nodeStates {
		ids = append(ids, ns.Reader().ID())
	}
	return ids
}

func (s *GraphStateReader) ID() string {
	if s.state.id == "" {
		id := uuid.New().String()
		s.state.id = id
	}
	return s.state.id
}
