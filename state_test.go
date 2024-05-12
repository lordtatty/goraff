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
	s := goraff.State{}
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
	s := goraff.State{}
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
	s := goraff.State{}
	// Test that this creates a new node state
	n := s.NewNodeState("node1")
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := s.Reader()
	nr := r.NodeState(n.Reader().ID())
	assert.Equal("value1", nr.GetStr("key1"))
	assert.Equal("value2", nr.GetStr("key2"))
}

func TestState_StateReadOnly_ID(t *testing.T) {
	assert := assert.New(t)
	s := goraff.State{}
	// Test that this creates a new node state
	n := s.NewNodeState("node1")
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := s.Reader()
	nr := r.NodeState(n.Reader().ID())
	assert.Equal("node1", nr.Name())
}

func TestState_OnUpdate(t *testing.T) {
	assert := assert.New(t)
	updated := false
	nsID := ""
	s := goraff.State{
		OnUpdate: []func(s *goraff.StateReadOnly){
			func(s *goraff.StateReadOnly) {
				assert.Equal("value", s.NodeState(nsID).GetStr("key"))
				updated = true
			},
		},
	}
	n := s.NewNodeState("node1")
	nsID = n.Reader().ID()
	assert.False(updated)
	n.SetStr("key", "value")
	assert.True(updated)
}
