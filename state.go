package goraff

import "github.com/google/uuid"

// State manages the state of all nodes in the graph
type State struct {
	id         string
	nodeStates []*NodeState
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

func (s *State) NewNodeState(name string) *NodeState {
	// Else create a new node state
	ns := &NodeState{name: name, onUpdate: s.onUpdate}
	s.nodeStates = append(s.nodeStates, ns)
	return ns
}

func (s *State) NodeStateByName(name string) []*NodeState {
	// First see if we have this node state
	result := []*NodeState{}
	for _, ns := range s.nodeStates {
		if ns.name == name {
			result = append(result, ns)
		}
	}
	return result
}

func (s *State) FirstNodeStateByName(name string) *NodeState {
	// First see if we have this node state
	for _, ns := range s.nodeStates {
		if ns.name == name {
			return ns
		}
	}
	return nil
}

func (s *State) NodeStateByID(id string) *NodeState {
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

func (s *StateReadOnly) FirstNodeStateByName(name string) *NodeStateReader {
	st := s.state.FirstNodeStateByName(name)
	if st == nil {
		return nil
	}
	return &NodeStateReader{st}
}

func (s *StateReadOnly) NodeState(id string) *NodeStateReader {
	r := s.state.NodeStateByID(id)
	if r == nil {
		return nil
	}
	return &NodeStateReader{r}
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

// Node state represents a key value store for an individual node
type NodeState struct {
	id       string
	name     string
	state    map[string][]byte
	done     bool
	onUpdate func()
	subState *State
}

func (n *NodeState) SetSubState(s *State) {
	s.AddOnUpdate(func(s *StateReadOnly) {
		if n.onUpdate != nil {
			n.onUpdate()
		}
	})
	n.subState = s
}

func (n *NodeState) SubState() *State {
	return n.subState
}

func (n *NodeState) MarkDone() {
	n.done = true
}

func (n *NodeState) Set(key string, value []byte) {
	if n.state == nil {
		n.state = make(map[string][]byte)
	}
	n.state[key] = value
	if n.onUpdate != nil {
		n.onUpdate()
	}
}

func (n *NodeState) SetStr(key, value string) {
	n.Set(key, []byte(value))
}

func (n *NodeState) Reader() *NodeStateReader {
	return &NodeStateReader{n}
}

// NodeStateReader is a read only view of a node state
type NodeStateReader struct {
	ns *NodeState
}

func (s *NodeStateReader) Get(key string) []byte {
	if s.ns.state == nil {
		return []byte{}
	}
	return s.ns.state[key]
}

func (s *NodeStateReader) GetStr(key string) string {
	return string(s.Get(key))
}

func (s *NodeStateReader) ID() string {
	if s.ns.id == "" {
		s.ns.id = uuid.NewString()
	}
	return s.ns.id
}

func (s *NodeStateReader) Name() string {
	return s.ns.name
}

func (s *NodeStateReader) Done() bool {
	return s.ns.done
}
