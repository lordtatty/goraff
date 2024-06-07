package outputs

import (
	"encoding/json"
	"fmt"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/websocket"
)

type NodeOutput struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Vals        []NodeOutputVal `json:"vals"`
	SubGraphIDs []string        `json:"subgraph_ids"`
}

type NodeOutputVal struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type Output struct {
	PrimaryStateID string        `json:"primary_state_id"`
	States         []GraphOutput `json:"states"`
	Nodes          []NodeOutput  `json:"nodes"`
}

type GraphOutput struct {
	ID      string   `json:"id"`
	NodeIDs []string `json:"node_ids"`
}

type Outputter struct {
}

func (o *Outputter) Output(s *goraff.ReadableGraph) *Output {
	st, err := o.allStates(s)
	if err != nil {
		fmt.Println("error getting node state")
		return nil
	}
	n, err := o.allNodes(s)
	if err != nil {
		fmt.Println("error getting node state")
		return nil
	}
	out := &Output{
		PrimaryStateID: s.ID(),
		States:         st,
		Nodes:          n,
	}
	return out
}

func (o *Outputter) allStates(s *goraff.ReadableGraph) ([]GraphOutput, error) {
	st, err := o.state(s)
	if err != nil {
		return nil, fmt.Errorf("error getting node state: %w", err)
	}
	states := []GraphOutput{
		*st,
	}
	for _, ns := range s.NodeIDs() {
		n, err := s.Node(ns)
		if err != nil {
			return nil, fmt.Errorf("error getting node state: %w", err)
		}
		if n.SubGraph() == nil {
			continue
		}
		for _, g := range n.SubGraph() {
			a, err := o.allStates(g)
			if err != nil {
				return nil, fmt.Errorf("error getting node state: %w", err)
			}
			states = append(states, a...)
		}
	}
	return states, nil
}

func (o *Outputter) state(s *goraff.ReadableGraph) (*GraphOutput, error) {
	nodeIDs := []string{}
	for _, ns := range s.NodeIDs() {
		n, err := s.Node(ns)
		if err != nil {
			return nil, fmt.Errorf("error getting node state: %w", err)
		}
		nodeIDs = append(nodeIDs, n.ID())
	}
	return &GraphOutput{
		ID:      s.ID(),
		NodeIDs: nodeIDs,
	}, nil
}

func (o *Outputter) allNodes(s *goraff.ReadableGraph) ([]NodeOutput, error) {
	nodes := []NodeOutput{}
	for _, ns := range s.NodeIDs() {
		n, err := s.Node(ns)
		if err != nil {
			return nil, fmt.Errorf("error getting node state: %w", err)
		}
		nodes = append(nodes, *o.node(n))
		if n.SubGraph() == nil {
			continue
		}
		for _, g := range n.SubGraph() {
			a, err := o.allNodes(g)
			if err != nil {
				return nil, fmt.Errorf("error getting node state: %w", err)
			}
			nodes = append(nodes, a...)
		}
	}
	return nodes, nil
}

func (o *Outputter) node(ns *goraff.ReadableNode) *NodeOutput {
	vals := []NodeOutputVal{}
	for _, key := range ns.Keys() {
		vals = append(vals, NodeOutputVal{
			Name:   key,
			Values: ns.AllStr(key),
		})
	}
	subIDs := []string{}
	for _, ns := range ns.SubGraph() {
		subIDs = append(subIDs, ns.ID())
	}
	return &NodeOutput{
		ID:          ns.ID(),
		Name:        ns.Name(),
		Vals:        vals,
		SubGraphIDs: subIDs,
	}
}

type ChangeListener interface {
	Listen(func(goraff.GraphChangeNotification))
}

// TODO - test this
func BroadcastChanges(l ChangeListener, r *goraff.ReadableGraph, ws *websocket.WebSocketServer) {
	l.Listen(func(c goraff.GraphChangeNotification) {
		out := Outputter{}
		o := out.Output(r)
		snd, err := json.Marshal(o)
		if err != nil {
			fmt.Println("error marshalling state")
			return
		}
		ws.Send(string(snd))
	})
}

// TODO - test this
func PrintUpdatesToConsole(l ChangeListener, r *goraff.ReadableGraph) {
	l.Listen(func(c goraff.GraphChangeNotification) {
		out := Outputter{}
		o := out.Output(r)
		fmt.Println("##########################################")
		fmt.Println("##########################################")
		fmt.Println("##########################################")
		b, _ := json.MarshalIndent(o, "", "  ")
		fmt.Println(string(b))
	})
}
