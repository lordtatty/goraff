package goraff

import (
	"github.com/google/uuid"
)

// Node state represents a key value store for an individual node
type Node struct {
	id       string
	name     string
	state    map[string][]byte
	done     bool
	notifier *GraphNotifier
	subGraph *Graph
}

func (n *Node) SetSubGraph(s *Graph) {
	s.notifier = n.notifier
	n.subGraph = s
}

func (n *Node) MarkDone() {
	n.done = true
}

func (n *Node) Set(key string, value []byte) {
	if n.state == nil {
		n.state = make(map[string][]byte)
	}
	n.state[key] = value
	if n.notifier != nil {
		n.notifier.Notify(StateChangeNotification{NodeID: n.id})
	}
}

func (n *Node) SetStr(key, value string) {
	n.Set(key, []byte(value))
}

func (n *Node) Reader() *ReadableNode {
	return &ReadableNode{ns: n}
}

// ReadableNode is a read only view of a node state
type ReadableNode struct {
	ns *Node
}

func (n *ReadableNode) State() map[string][]byte {
	if n.ns == nil {
		return map[string][]byte{}
	}
	return n.ns.state
}

func (n *ReadableNode) SubGraph() *ReadableGraph {
	if n.ns.subGraph == nil {
		return nil
	}
	return &ReadableGraph{n.ns.subGraph}
}

func (s *ReadableNode) Get(key string) []byte {
	if s.ns.state == nil {
		return []byte{}
	}
	return s.ns.state[key]
}

func (s *ReadableNode) GetStr(key string) string {
	return string(s.Get(key))
}

func (s *ReadableNode) ID() string {
	if s.ns.id == "" {
		s.ns.id = uuid.NewString()
	}
	return s.ns.id
}

func (s *ReadableNode) Name() string {
	return s.ns.name
}

func (s *ReadableNode) Done() bool {
	return s.ns.done
}
