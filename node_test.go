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

	sn := s.NewNode("subnode", nil)
	sn.SetStr("key1", "value1")

	subGraph := n.Reader().SubGraph()
	node, err := subGraph.NodeByID(sn.Reader().ID())
	assert.Nil(err)
	assert.Equal("value1", node.FirstStr("key1"))
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
	r := goraff.NewReadableGraph(s)
	n.SetSubGraph(s)

	subGraph := n.Reader().SubGraph()
	assert.NotNil(subGraph)
	assert.Equal(r.ID(), subGraph.ID())
}

func TestNode_State(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	n.Set("key1", []byte("value1"))
	n.Set("key2", []byte("value2"))

	state := n.Reader().State()
	assert.Equal(2, len(state))

	assert.Equal([][]byte{[]byte("value1")}, state["key1"])
	assert.Equal([][]byte{[]byte("value2")}, state["key2"])
}

func TestNode_SetGet(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	key := "testKey"
	expected := []byte("testValue")

	n.Set(key, expected)

	result := n.Reader().First(key)
	assert.Equal(expected, result)
}

func TestNode_SetStrGetStr(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	key := "testKeyStr"
	expected := "testValueStr"

	n.SetStr(key, expected)

	result := n.Reader().FirstStr(key)
	assert.Equal(expected, result)
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
	assert.Equal(value, r.First(key))
}

func TestReadableNode_GetStr(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	key := "testKeyGetStr"
	value := "testValueGetStr"

	n.SetStr(key, value)

	r := n.Reader()
	assert.Equal(value, r.FirstStr(key))
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
		result := n.Reader().First(key + strconv.Itoa(i))
		assert.Equal(value, result)
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
			n.Reader().First(key + strconv.Itoa(i))
		}(i)
	}

	wg.Wait()

	state := n.Reader().State()
	assert.Equal(100, len(state))
	for i := 0; i < 100; i++ {
		result := n.Reader().First(key + strconv.Itoa(i))
		assert.Equal(value, result)
	}
}

func TestNode_TriggeredBy(t *testing.T) {
	assert := assert.New(t)

	graph := &goraff.Graph{}
	graph.NewNode("parent1", nil)
	graph.NewNode("parent2", nil)
	graph.NewNode("parent3", nil)

	rGraph := goraff.NewReadableGraph(graph)

	rParent1, err := rGraph.FirstNodeByName("parent1")
	assert.Nil(err)
	rParent2, err := rGraph.FirstNodeByName("parent2")
	assert.Nil(err)
	rParent3, err := rGraph.FirstNodeByName("parent3")
	assert.Nil(err)

	graph.NewNode("sut_node", []*goraff.ReadableNode{rParent1, rParent2, rParent3})

	sut, err := rGraph.FirstNodeByName("sut_node")
	assert.Nil(err)

	triggeredBy := sut.TriggeredBy()
	assert.Equal(3, len(triggeredBy))
	assert.Equal(rParent1.ID(), triggeredBy[0].ID())
	assert.Equal(rParent2.ID(), triggeredBy[1].ID())
	assert.Equal(rParent3.ID(), triggeredBy[2].ID())
}

func TestNode_Add(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	key := "testKey"
	value1 := []byte("testValue")
	value2 := []byte("testValue2")
	value3 := []byte("testValue3")

	n.Add(key, value1)
	n.Add(key, value2)
	n.Add(key, value3)

	state := n.Reader().State()
	assert.Equal(3, len(state[key]))

	assert.Equal([][]byte{value1, value2, value3}, state[key])

	all := n.Reader().All(key)
	assert.Equal(3, len(all))
	assert.Equal(value1, all[0])
	assert.Equal(value2, all[1])
	assert.Equal(value3, all[2])
}

func TestNode_AddStr(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	key := "testKey"
	value1 := "testValue"
	value2 := "testValue2"
	value3 := "testValue3"

	n.AddStr(key, value1)
	n.AddStr(key, value2)
	n.AddStr(key, value3)

	state := n.Reader().State()
	assert.Equal(3, len(state[key]))

	assert.Equal([][]byte{[]byte(value1), []byte(value2), []byte(value3)}, state[key])

	all := n.Reader().AllStr(key)
	assert.Equal(3, len(all))
	assert.Equal(value1, all[0])
	assert.Equal(value2, all[1])
	assert.Equal(value3, all[2])
}

// MockNotifier is a mock implementation of the Notifier interface for testing purposes
type MockNotifier struct {
	Notified bool
}

func (m *MockNotifier) Notify(notification goraff.GraphChangeNotification) {
	m.Notified = true
}
