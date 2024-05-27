package goraff_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestNodeState(t *testing.T) {
	assert := assert.New(t)
	n := goraff.StateNode{}
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := n.Reader()
	assert.Equal("value1", r.GetStr("key1"))
	assert.Equal("value2", r.GetStr("key2"))
	n.SetStr("key3", "value3")
	assert.Equal("value1", r.GetStr("key1"))
	assert.Equal("value2", r.GetStr("key2"))
	assert.Equal("value3", r.GetStr("key3"))
}

func TestNodeState_GetOnUninitialisedState(t *testing.T) {
	assert := assert.New(t)
	n := goraff.StateNode{}
	assert.Equal("", n.Reader().GetStr("key1"))
}

func TestNodeState_ID(t *testing.T) {
	assert := assert.New(t)
	n := goraff.StateNode{}
	assert.Regexp("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", n.Reader().ID())
}

func TestState_NodeStateByName(t *testing.T) {
	assert := assert.New(t)
	s := goraff.StateGraph{}
	// Test that this creates a new node state
	n := s.NewNodeState("node1")
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := n.Reader()
	assert.Equal("value1", r.GetStr("key1"))
	assert.Equal("value2", r.GetStr("key2"))
	// Test that this returns the same already-created state
	n2 := s.NewNodeState("node1")
	assert.Equal("value1", r.GetStr("key1"))
	assert.Equal("value2", r.GetStr("key2"))
	// Test the id
	r2 := n2.Reader()
	assert.Equal("node1", r2.Name())
}

func TestState_NodeStateByID(t *testing.T) {
	assert := assert.New(t)
	s := goraff.StateGraph{}
	// Test that this creates a new node state
	n := s.NewNodeState("node1")
	n.SetStr("key1", "value1")
	n.SetStr("key", "value2")
	r := n.Reader()
	assert.Equal("value1", r.GetStr("key1"))
	assert.Equal("value2", r.GetStr("key"))
	// Test that this returns the same already-created state
	n2 := s.NodeStateByID(n.Reader().ID())
	assert.Equal("value1", r.GetStr("key1"))
	assert.Equal("value2", r.GetStr("key"))
	// Test the id
	r2 := n2.Reader()
	assert.Equal("node1", r2.Name())
}

func TestState_StateReadOnly(t *testing.T) {
	assert := assert.New(t)
	s := goraff.StateGraph{}
	// Test that this creates a new node state
	n := s.NewNodeState("node1")
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := s.Reader()
	nr, err := r.NodeState(n.Reader().ID())
	assert.Nil(err)
	assert.Equal("value1", nr.GetStr("key1"))
	assert.Equal("value2", nr.GetStr("key2"))
}

func TestState_StateReadOnly_ID(t *testing.T) {
	assert := assert.New(t)
	s := goraff.StateGraph{}
	// Test that this creates a new node state
	n := s.NewNodeState("node1")
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := s.Reader()
	nr, err := r.NodeState(n.Reader().ID())
	assert.Nil(err)
	assert.Equal("node1", nr.Name())
}

func TestState_Notifier(t *testing.T) {
	assert := assert.New(t)
	s := goraff.StateGraph{}
	assert.NotNil(s.Notifier())
	// assert that the returned notifier is of type StateNotifier
	n := s.Notifier()
	assert.IsType(&goraff.StateNotifier{}, n)

}

func TestStateReader_ID(t *testing.T) {
	assert := assert.New(t)
	s := goraff.StateGraph{}
	r := s.Reader()
	assert.Regexp("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", r.ID())
}
