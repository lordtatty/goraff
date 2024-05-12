package outputs

import (
	"encoding/json"
	"fmt"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/websocket"
)

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

func (o *Outputter) Output(s *goraff.GraphStateReader) *Output {
	out := &Output{
		PrimaryStateID: s.ID(),
		States:         o.allStates(s),
		Nodes:          o.allNodes(s),
	}
	return out
}

func (o *Outputter) allStates(s *goraff.GraphStateReader) []StateOutput {
	states := []StateOutput{
		*o.state(s),
	}
	for _, ns := range s.NodeIDs() {
		n := s.NodeState(ns)
		if n.SubGraph() == nil {
			continue
		}
		a := o.allStates(n.SubGraph())
		states = append(states, a...)
	}
	return states
}

func (o *Outputter) state(s *goraff.GraphStateReader) *StateOutput {
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

func (o *Outputter) allNodes(s *goraff.GraphStateReader) []NodeOutput {
	nodes := []NodeOutput{}
	for _, ns := range s.NodeIDs() {
		n := s.NodeState(ns)
		nodes = append(nodes, *o.node(n))
		if n.SubGraph() == nil {
			continue
		}
		a := o.allNodes(n.SubGraph())
		nodes = append(nodes, a...)
	}
	return nodes
}

func (o *Outputter) node(ns *goraff.StateNodeReader) *NodeOutput {
	vals := []NodeOutputVal{}
	for k, v := range ns.State() {
		vals = append(vals, NodeOutputVal{Name: k, Value: string(v)})
	}
	subID := ""
	if ns.SubGraph() != nil {
		subID = ns.SubGraph().ID()
	}
	return &NodeOutput{
		ID:         ns.ID(),
		Name:       ns.Name(),
		Vals:       vals,
		SubStateID: subID,
	}
}

// TODO - test this
func BroadcastChanges(g *goraff.Graph, ws *websocket.WebSocketServer) {
	g.State().AddOnUpdate(func(s *goraff.GraphStateReader) {
		out := Outputter{}
		o := out.Output(s)
		snd, err := json.Marshal(o)
		if err != nil {
			fmt.Println("error marshalling state")
			return
		}
		ws.Send(string(snd))
	})
}

// TODOD - test this
func PrintUpdatesToConsole(g *goraff.Graph) {
	g.State().AddOnUpdate(func(s *goraff.GraphStateReader) {
		out := Outputter{}
		o := out.Output(s)
		fmt.Println("##########################################")
		fmt.Println("##########################################")
		fmt.Println("##########################################")
		b, _ := json.MarshalIndent(o, "", "  ")
		fmt.Println(string(b))
	})
}
