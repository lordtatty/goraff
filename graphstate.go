package goraff

// GraphState manages the state of all nodes in the graph
type GraphState struct {
	id         string
	nodeStates []*StateNode
	OnUpdate   []func(s *GraphStateReader)
}

func (s *GraphState) AddOnUpdate(f func(s *GraphStateReader)) {
	s.OnUpdate = append(s.OnUpdate, f)
}

func (s *GraphState) onUpdate() {
	if s.OnUpdate == nil {
		return
	}
	for _, f := range s.OnUpdate {
		f(s.Reader())
	}
}

func (s *GraphState) NewNodeState(name string) *StateNode {
	// Else create a new node state
	ns := &StateNode{name: name, onUpdate: s.onUpdate}
	s.nodeStates = append(s.nodeStates, ns)
	return ns
}

func (s *GraphState) NodeStateByName(name string) []*StateNode {
	// First see if we have this node state
	result := []*StateNode{}
	for _, ns := range s.nodeStates {
		if ns.name == name {
			result = append(result, ns)
		}
	}
	return result
}

func (s *GraphState) FirstNodeStateByName(name string) *StateNode {
	// First see if we have this node state
	for _, ns := range s.nodeStates {
		if ns.name == name {
			return ns
		}
	}
	return nil
}

func (s *GraphState) NodeStateByID(id string) *StateNode {
	// First see if we have this node state
	for _, ns := range s.nodeStates {
		if ns.Reader().ID() == id {
			return ns
		}
	}
	return nil
}

func (s *GraphState) Reader() *GraphStateReader {
	return &GraphStateReader{s}
}

// GraphStateReader is a read only view of the state
type GraphStateReader struct {
	state *GraphState
}

func (s *GraphStateReader) FirstNodeStateByName(name string) *StateNodeReader {
	st := s.state.FirstNodeStateByName(name)
	if st == nil {
		return nil
	}
	return &StateNodeReader{st}
}

func (s *GraphStateReader) NodeState(id string) *StateNodeReader {
	r := s.state.NodeStateByID(id)
	if r == nil {
		return nil
	}
	return &StateNodeReader{r}
}

func (s *GraphStateReader) NodeIDs() []string {
	ids := []string{}
	for _, ns := range s.state.nodeStates {
		ids = append(ids, ns.Reader().ID())
	}
	return ids
}

func (s *GraphStateReader) ID() string {
	return s.state.id
}
