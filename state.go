package goraff

// State manages the state of all nodes in the graph
type State struct {
	nodeStates []*NodeState
}

func (s *State) NodeState(id string) *NodeState {
	// First see if we have this node state
	for _, ns := range s.nodeStates {
		if ns.nodeID == id {
			return ns
		}
	}
	// Else create a new node state
	ns := &NodeState{nodeID: id}
	s.nodeStates = append(s.nodeStates, ns)
	return ns
}

func (s *State) ReadOnly() *StateReadOnly {
	return &StateReadOnly{s}
}

// StateReadOnly is a read only view of the state
type StateReadOnly struct {
	state *State
}

func (s *StateReadOnly) NodeState(id string) *NodeStateReadOnly {
	r := s.state.NodeState(id)
	return &NodeStateReadOnly{r}
}

// Node state represents a key value store for an individual node
type NodeState struct {
	nodeID string
	state  map[string]string
	done   bool
}

func (n *NodeState) MarkDone() {
	n.done = true
}

func (n *NodeState) Get(key string) string {
	if n.state == nil {
		return ""
	}
	return n.state[key]
}

func (n *NodeState) Set(key, value string) {
	if n.state == nil {
		n.state = make(map[string]string)
	}
	n.state[key] = value
}

func (n *NodeState) ID() string {
	return n.nodeID
}

// NodeStateReadOnly is a read only view of a node state
type NodeStateReadOnly struct {
	state *NodeState
}

func (s *NodeStateReadOnly) Get(key string) string {
	return s.state.Get(key)
}

func (s *NodeStateReadOnly) ID() string {
	return s.state.ID()
}

func (s *NodeStateReadOnly) Done() bool {
	return s.state.done
}
