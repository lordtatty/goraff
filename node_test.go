package goraff_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestNodeState_SubState(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	s := &goraff.Graph{}
	n.SetSubGraph(s)

	sn := s.NewNodeState("subnode")
	sn.SetStr("key1", "value1")

	subGraph := n.Reader().SubGraph()
	node, err := subGraph.NodeByID(sn.Reader().ID())
	assert.Nil(err)
	assert.Equal("value1", node.GetStr("key1"))
}

func TestStateNode_Reader(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	r := n.Reader()
	assert.Equal(n.Reader().ID(), r.ID())
}

func TestNode_SetSubGraph(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	s := &goraff.Graph{}
	n.SetSubGraph(s)

	subGraph := n.Reader().SubGraph()
	assert.NotNil(subGraph)
	assert.Equal(s.Reader().ID(), subGraph.ID())
}

func TestNode_Set(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	key := "testKey"
	value := []byte("testValue")

	n.Set(key, value)

	state := n.Reader().State()
	assert.Equal(value, state[key])
}

func TestNode_SetStr(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	key := "testKeyStr"
	value := "testValueStr"

	n.SetStr(key, value)

	state := n.Reader().State()
	assert.Equal([]byte(value), state[key])
}

func TestNode_MarkDone(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}

	n.MarkDone()

	assert.True(n.Reader().Done())
}

func TestReadableNode_Get(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	key := "testKeyGet"
	value := []byte("testValueGet")

	n.Set(key, value)

	r := n.Reader()
	assert.Equal(value, r.Get(key))
}

func TestReadableNode_GetStr(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	key := "testKeyGetStr"
	value := "testValueGetStr"

	n.SetStr(key, value)

	r := n.Reader()
	assert.Equal(value, r.GetStr(key))
}

func TestReadableNode_ID(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	r := n.Reader()

	id := r.ID()
	assert.NotEmpty(id)
	assert.Equal(id, r.ID()) // ID should remain consistent
}

func TestNode_ConcurrentSet(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	key := "testKey"
	value := []byte("testValue")

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			n.Set(key+strconv.Itoa(i), value)
		}(i)
	}

	wg.Wait()

	state := n.Reader().State()
	assert.Equal(100, len(state))
	for i := 0; i < 100; i++ {
		assert.Equal(value, state[key+strconv.Itoa(i)])
	}
}

func TestNode_ConcurrentReadWrite(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	key := "testKey"
	value := []byte("testValue")

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			n.Set(key+strconv.Itoa(i), value)
		}(i)

		go func(i int) {
			defer wg.Done()
			n.Reader().Get(key + strconv.Itoa(i))
		}(i)
	}

	wg.Wait()

	state := n.Reader().State()
	assert.Equal(100, len(state))
	for i := 0; i < 100; i++ {
		assert.Equal(value, state[key+strconv.Itoa(i)])
	}
}

// MockNotifier is a mock implementation of the Notifier interface for testing purposes
type MockNotifier struct {
	Notified bool
}

func (m *MockNotifier) Notify(notification goraff.StateChangeNotification) {
	m.Notified = true
}
