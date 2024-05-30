package goraff

import (
	"sync"

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
	mut      sync.Mutex
}

func (n *Node) SetSubGraph(s *Graph) {
	n.mut.Lock()
	defer n.mut.Unlock()
	s.notifier = n.notifier
	n.subGraph = s
}

func (n *Node) MarkDone() {
	n.done = true
}

func (n *Node) Set(key string, value []byte) {
	n.mut.Lock()
	if n.state == nil {
		n.state = make(map[string][]byte)
	}
	n.state[key] = value
	n.mut.Unlock()
	if n.notifier != nil {
		n.notifier.Notify(StateChangeNotification{NodeID: n.id})
	}
}

func (n *Node) SetStr(key, value string) {
	n.Set(key, []byte(value))
}

func (n *Node) Reader() *ReadableNode {
	return &ReadableNode{node: n}
}

// ReadableNode is a read only view of a node state
type ReadableNode struct {
	node *Node
}

func (n *ReadableNode) State() map[string][]byte {
	n.node.mut.Lock()
	defer n.node.mut.Unlock()
	if n.node == nil {
		return map[string][]byte{}
	}
	// Copy the state
	state := make(map[string][]byte)
	for k, v := range n.node.state {
		state[k] = v
	}
	return state
}

func (n *ReadableNode) SubGraph() *ReadableGraph {
	if n.node.subGraph == nil {
		return nil
	}
	return &ReadableGraph{n.node.subGraph}
}

func (s *ReadableNode) Get(key string) []byte {
	if s.node.state == nil {
		return []byte{}
	}
	s.node.mut.Lock()
	defer s.node.mut.Unlock()
	return s.node.state[key]
}

func (s *ReadableNode) GetStr(key string) string {
	return string(s.Get(key))
}

func (s *ReadableNode) ID() string {
	if s.node.id == "" {
		s.node.id = uuid.NewString()
	}
	return s.node.id
}

func (s *ReadableNode) Name() string {
	return s.node.name
}

func (s *ReadableNode) Done() bool {
	return s.node.done
}
