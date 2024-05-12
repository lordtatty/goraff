package goraff

// State manages the state of all nodes in the graph
type State struct {
	id         string
	nodeStates []*StateNode
	OnUpdate   []func(s *StateReadOnly)
}

func (s *State) AddOnUpdate(f func(s *StateReadOnly)) {
	s.OnUpdate = append(s.OnUpdate, f)
}

func (s *State) onUpdate() {
	if s.OnUpdate == nil {
		return
	}
	for _, f := range s.OnUpdate {
		f(s.Reader())
	}
}

func (s *State) NewNodeState(name string) *StateNode {
	// Else create a new node state
	ns := &StateNode{name: name, onUpdate: s.onUpdate}
	s.nodeStates = append(s.nodeStates, ns)
	return ns
}

func (s *State) NodeStateByName(name string) []*StateNode {
	// First see if we have this node state
	result := []*StateNode{}
	for _, ns := range s.nodeStates {
		if ns.name == name {
			result = append(result, ns)
		}
	}
	return result
}

func (s *State) FirstNodeStateByName(name string) *StateNode {
	// First see if we have this node state
	for _, ns := range s.nodeStates {
		if ns.name == name {
			return ns
		}
	}
	return nil
}

func (s *State) NodeStateByID(id string) *StateNode {
	// First see if we have this node state
	for _, ns := range s.nodeStates {
		if ns.Reader().ID() == id {
			return ns
		}
	}
	return nil
}

func (s *State) Reader() *StateReadOnly {
	return &StateReadOnly{s}
}

// StateReadOnly is a read only view of the state
type StateReadOnly struct {
	state *State
}

func (s *StateReadOnly) FirstNodeStateByName(name string) *StateNodeReader {
	st := s.state.FirstNodeStateByName(name)
	if st == nil {
		return nil
	}
	return &StateNodeReader{st}
}

func (s *StateReadOnly) NodeState(id string) *StateNodeReader {
	r := s.state.NodeStateByID(id)
	if r == nil {
		return nil
	}
	return &StateNodeReader{r}
}

func (s *StateReadOnly) NodeIDs() []string {
	ids := []string{}
	for _, ns := range s.state.nodeStates {
		ids = append(ids, ns.Reader().ID())
	}
	return ids
}

func (s *StateReadOnly) ID() string {
	return s.state.id
}
