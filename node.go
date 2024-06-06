package goraff

import (
	"sync"

	"github.com/google/uuid"
)

// Node state represents a key value store for an individual node
type Node struct {
	id          string
	name        string
	state       map[string][][]byte
	done        bool
	notifier    ChangeNotifier
	subGraphs   []*ReadableGraph
	mut         sync.Mutex
	triggeredBy []*ReadableNode
}

func (n *Node) AddSubGraph(s *Graph) {
	n.mut.Lock()
	defer n.mut.Unlock()
	s.Notifier = n.notifier
	r := NewReadableGraph(s)
	n.subGraphs = append(n.subGraphs, r)
}

func (n *Node) MarkDone() {
	n.done = true
}

func (n *Node) Add(key string, value []byte) {
	n.mut.Lock()
	if n.state == nil {
		n.state = make(map[string][][]byte)
	}
	n.state[key] = append(n.state[key], value)
	n.mut.Unlock()
	if n.notifier != nil {
		n.notifier.Notify(GraphChangeNotification{NodeID: n.name})
	}
}

func (n *Node) AddStr(key, value string) {
	n.Add(key, []byte(value))
}

func (n *Node) Set(key string, value []byte) {
	n.mut.Lock()
	if n.state == nil {
		n.state = make(map[string][][]byte)
	}
	n.state[key] = [][]byte{value}
	n.mut.Unlock()
	if n.notifier != nil {
		n.notifier.Notify(GraphChangeNotification{NodeID: n.name})
	}
}

func (n *Node) SetStr(key, value string) {
	n.Set(key, []byte(value))
}

func (n *Node) Get() *ReadableNode {
	return &ReadableNode{node: n}
}

// ReadableNode is a read only view of a node state
type ReadableNode struct {
	node *Node
}

func (n *ReadableNode) State() map[string][][]byte {
	n.node.mut.Lock()
	defer n.node.mut.Unlock()
	if n.node == nil {
		return map[string][][]byte{}
	}
	// Copy the state
	state := make(map[string][][]byte)
	for k, v := range n.node.state {
		state[k] = v
	}
	return state
}

func (n *ReadableNode) SubGraph() []*ReadableGraph {
	if n.node.subGraphs == nil {
		return nil
	}
	return n.node.subGraphs
}

func (s *ReadableNode) First(key string) []byte {
	if s.node.state == nil {
		return []byte{}
	}
	s.node.mut.Lock()
	defer s.node.mut.Unlock()
	if s.node.state[key] == nil {
		return []byte{}
	}
	return s.node.state[key][0]
}

func (s *ReadableNode) FirstStr(key string) string {
	return string(s.First(key))
}

func (s *ReadableNode) All(key string) [][]byte {
	if s.node.state == nil {
		return [][]byte{}
	}
	s.node.mut.Lock()
	defer s.node.mut.Unlock()
	if s.node.state[key] == nil {
		return [][]byte{}
	}
	return s.node.state[key]
}

func (s *ReadableNode) AllStr(key string) []string {
	vals := []string{}
	for _, v := range s.All(key) {
		vals = append(vals, string(v))
	}
	return vals
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

func (n *ReadableNode) TriggeredBy() []*ReadableNode {
	return n.node.triggeredBy
}
