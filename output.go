package goraff

type NodeOutput struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Vals       []NodeOutputVal `json:"vals"`
	SubStateID string          `json:"substate_id"`
}

type NodeOutputVal struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Output struct {
	PrimaryStateID string        `json:"primary_state_id"`
	States         []StateOutput `json:"states"`
	Nodes          []NodeOutput  `json:"nodes"`
}

type StateOutput struct {
	ID      string   `json:"id"`
	NodeIDs []string `json:"node_ids"`
}

type Outputter struct {
}

func (o *Outputter) Output(s *StateReadOnly) *Output {
	out := &Output{
		PrimaryStateID: s.ID(),
		States:         o.allStates(s),
		Nodes:          o.allNodes(s),
	}
	return out
}

func (o *Outputter) allStates(s *StateReadOnly) []StateOutput {
	states := []StateOutput{
		*o.state(s),
	}
	for _, ns := range s.NodeIDs() {
		n := s.NodeState(ns)
		if n.ns.subState == nil {
			continue
		}
		a := o.allStates(n.ns.subState.Reader())
		states = append(states, a...)
	}
	return states
}

func (o *Outputter) state(s *StateReadOnly) *StateOutput {
	nodeIDs := []string{}
	for _, ns := range s.NodeIDs() {
		n := s.NodeState(ns)
		nodeIDs = append(nodeIDs, n.Name())
	}
	return &StateOutput{
		ID:      s.ID(),
		NodeIDs: nodeIDs,
	}
}

func (o *Outputter) allNodes(s *StateReadOnly) []NodeOutput {
	nodes := []NodeOutput{}
	for _, ns := range s.NodeIDs() {
		n := s.NodeState(ns)
		nodes = append(nodes, *o.node(n))
		if n.ns.subState == nil {
			continue
		}
		a := o.allNodes(n.ns.subState.Reader())
		nodes = append(nodes, a...)
	}
	return nodes
}

func (o *Outputter) node(ns *StateNodeReader) *NodeOutput {
	vals := []NodeOutputVal{}
	for k, v := range ns.ns.state {
		vals = append(vals, NodeOutputVal{Name: k, Value: string(v)})
	}
	subID := ""
	if ns.ns.subState != nil {
		subID = ns.ns.subState.Reader().ID()
	}
	return &NodeOutput{
		ID:         ns.ID(),
		Name:       ns.Name(),
		Vals:       vals,
		SubStateID: subID,
	}
}
