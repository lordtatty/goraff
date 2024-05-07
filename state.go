package goraff

// State manages the state of all nodes in the graph
type State struct {
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
		f(s.ReadOnly())
	}
}

func (s *State) NodeStateUpsert(id string) *NodeState {
	// First see if we have this node state
	ns := s.NodeState(id)
	if ns != nil {
		return ns
	}
	// Else create a new node state
	ns = &NodeState{nodeID: id, onUpdate: s.onUpdate}
	s.nodeStates = append(s.nodeStates, ns)
	return ns
}

func (s *State) NodeState(id string) *NodeState {
	// First see if we have this node state
	for _, ns := range s.nodeStates {
		if ns.nodeID == id {
			return ns
		}
	}
	return nil
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
	if r == nil {
		return nil
	}
	return &NodeStateReadOnly{r}
}

func (s *StateReadOnly) Outputs() []NodeOutput {
	var outputs []NodeOutput
	for _, ns := range s.state.nodeStates {
		n := NodeStateReadOnly{ns}
		outputs = append(outputs, n.Outputs())
	}
	return outputs
}

// Node state represents a key value store for an individual node
type NodeState struct {
	nodeID   string
	state    map[string][]byte
	done     bool
	onUpdate func()
	subState *State
}

func (n *NodeState) SetSubState(s *State) {
	n.subState = s
}

func (n *NodeState) SubState() *State {
	return n.subState
}

func (n *NodeState) MarkDone() {
	n.done = true
}

func (n *NodeState) Get(key string) []byte {
	if n.state == nil {
		return []byte{}
	}
	return n.state[key]
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

func (n *NodeState) GetStr(key string) string {
	return string(n.Get(key))
}

func (n *NodeState) ID() string {
	return n.nodeID
}

// NodeStateReadOnly is a read only view of a node state
type NodeStateReadOnly struct {
	ns *NodeState
}

func (s *NodeStateReadOnly) Get(key string) []byte {
	return s.ns.Get(key)
}

func (s *NodeStateReadOnly) GetStr(key string) string {
	return s.ns.GetStr(key)
}

func (s *NodeStateReadOnly) ID() string {
	return s.ns.ID()
}

func (s *NodeStateReadOnly) Done() bool {
	return s.ns.done
}

type NodeOutput struct {
	ID   string          `json:"id"`
	Vals []NodeOutputVal `json:"vals"`
}

type NodeOutputVal struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (s *NodeStateReadOnly) Outputs() NodeOutput {
	output := NodeOutput{ID: s.ID()}
	for key, value := range s.ns.state {
		output.Vals = append(output.Vals, NodeOutputVal{Name: key, Value: string(value)})
	}
	return output
}
