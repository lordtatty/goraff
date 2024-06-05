package goraff_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNodeState(t *testing.T) {
	assert := assert.New(t)
	n := goraff.Node{}
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := n.Reader()
	assert.Equal("value1", r.FirstStr("key1"))
	assert.Equal("value2", r.FirstStr("key2"))
	n.SetStr("key3", "value3")
	assert.Equal("value1", r.FirstStr("key1"))
	assert.Equal("value2", r.FirstStr("key2"))
	assert.Equal("value3", r.FirstStr("key3"))
}

func TestNodeState_GetOnUninitialisedState(t *testing.T) {
	assert := assert.New(t)
	n := goraff.Node{}
	assert.Equal("", n.Reader().FirstStr("key1"))
}

func TestNodeState_ID(t *testing.T) {
	assert := assert.New(t)
	n := goraff.Node{}
	assert.Regexp("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", n.Reader().ID())
}

func TestState_NodeStateByName(t *testing.T) {
	assert := assert.New(t)
	s := goraff.Graph{}
	// Test that this creates a new node state
	n := s.NewNode("node1", nil)
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := n.Reader()
	assert.Equal("value1", r.FirstStr("key1"))
	assert.Equal("value2", r.FirstStr("key2"))
	// Test that this returns the same already-created state
	n2 := s.NewNode("node1", nil)
	assert.Equal("value1", r.FirstStr("key1"))
	assert.Equal("value2", r.FirstStr("key2"))
	// Test the id
	r2 := n2.Reader()
	assert.Equal("node1", r2.Name())
}

func TestState_NodeStateByID(t *testing.T) {
	assert := assert.New(t)
	s := goraff.Graph{}
	// Test that this creates a new node state
	n := s.NewNode("node1", nil)
	n.SetStr("key1", "value1")
	n.SetStr("key", "value2")
	r := n.Reader()
	assert.Equal("value1", r.FirstStr("key1"))
	assert.Equal("value2", r.FirstStr("key"))
	// Test that this returns the same already-created state
	n2 := s.NodeByID(n.Reader().ID())
	assert.Equal("value1", r.FirstStr("key1"))
	assert.Equal("value2", r.FirstStr("key"))
	// Test the id
	r2 := n2.Reader()
	assert.Equal("node1", r2.Name())
}

func TestState_StateReadOnly(t *testing.T) {
	assert := assert.New(t)
	s := &goraff.Graph{}
	// Test that this creates a new node state
	n := s.NewNode("node1", nil)
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := goraff.NewReadableGraph(s)
	nr, err := r.Node(n.Reader().ID())
	assert.Nil(err)
	assert.Equal("value1", nr.FirstStr("key1"))
	assert.Equal("value2", nr.FirstStr("key2"))
}

func TestState_StateReadOnly_ID(t *testing.T) {
	assert := assert.New(t)
	s := &goraff.Graph{}
	// Test that this creates a new node state
	n := s.NewNode("node1", nil)
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := goraff.NewReadableGraph(s)
	nr, err := r.Node(n.Reader().ID())
	assert.Nil(err)
	assert.Equal("node1", nr.Name())
}

func TestState_Notifier(t *testing.T) {
	// Make sure this fires when updating a node
	mNotifier := mocks.NewChangeNotifier(t)
	mNotifier.EXPECT().Notify(goraff.GraphChangeNotification{NodeID: "node1"}).Times(1)
	// And let's check a second node too to be sure
	mNotifier.EXPECT().Notify(goraff.GraphChangeNotification{NodeID: "node2"}).Times(1)

	// Create the SUT and trigger the first node to update
	s := &goraff.Graph{Notifier: mNotifier}
	n := s.NewNode("node1", nil)
	n.SetStr("key", "value")

	// Trigger the second node to update
	n2 := s.NewNode("node2", nil)
	n2.SetStr("key", "value")
}

func TestStateReader_ID(t *testing.T) {
	assert := assert.New(t)
	s := &goraff.Graph{}
	r := goraff.NewReadableGraph(s)
	assert.Regexp("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", r.ID())
}
