package goraff

import "github.com/google/uuid"

// Node state represents a key value store for an individual node
type StateNode struct {
	id       string
	name     string
	state    map[string][]byte
	done     bool
	onUpdate func()
	subGraph *StateGraph
}

func (n *StateNode) SetSubGraph(s *StateGraph) {
	s.AddOnUpdate(func(s *GraphStateReader) {
		if n.onUpdate != nil {
			n.onUpdate()
		}
	})
	n.subGraph = s
}

func (n *StateNode) MarkDone() {
	n.done = true
}

func (n *StateNode) Set(key string, value []byte) {
	if n.state == nil {
		n.state = make(map[string][]byte)
	}
	n.state[key] = value
	if n.onUpdate != nil {
		n.onUpdate()
	}
}

func (n *StateNode) SetStr(key, value string) {
	n.Set(key, []byte(value))
}

func (n *StateNode) Reader() *StateNodeReader {
	if n == nil {
		return nil
	}
	return &StateNodeReader{ns: n}
}

// StateNodeReader is a read only view of a node state
type StateNodeReader struct {
	ns *StateNode
}

func (n *StateNodeReader) State() map[string][]byte {
	return n.ns.state
}

func (n *StateNodeReader) SubGraph() *GraphStateReader {
	if n.ns.subGraph == nil {
		return nil
	}
	return &GraphStateReader{n.ns.subGraph}
}

func (s *StateNodeReader) Get(key string) []byte {
	if s.ns.state == nil {
		return []byte{}
	}
	return s.ns.state[key]
}

func (s *StateNodeReader) GetStr(key string) string {
	return string(s.Get(key))
}

func (s *StateNodeReader) ID() string {
	if s.ns.id == "" {
		s.ns.id = uuid.NewString()
	}
	return s.ns.id
}

func (s *StateNodeReader) Name() string {
	return s.ns.name
}

func (s *StateNodeReader) Done() bool {
	return s.ns.done
}
