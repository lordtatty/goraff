package goraff_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestNodeState(t *testing.T) {
	assert := assert.New(t)
	n := goraff.NodeState{}
	n.Set("key1", "value1")
	n.Set("key2", "value2")
	assert.Equal("value1", n.Get("key1"))
	assert.Equal("value2", n.Get("key2"))
	n.Set("key3", "value3")
	assert.Equal("value1", n.Get("key1"))
	assert.Equal("value2", n.Get("key2"))
	assert.Equal("value3", n.Get("key3"))
}

func TestNodeState_GetOnUninitialisedState(t *testing.T) {
	assert := assert.New(t)
	n := goraff.NodeState{}
	assert.Equal("", n.Get("key1"))
}

func TestState_NodeState(t *testing.T) {
	assert := assert.New(t)
	s := goraff.State{}
	// Test that this creates a new node state
	n := s.NodeState("node1")
	n.Set("key1", "value1")
	n.Set("key2", "value2")
	assert.Equal("value1", n.Get("key1"))
	assert.Equal("value2", n.Get("key2"))
	// Test that this returns the same already-created state
	n2 := s.NodeState("node1")
	assert.Equal("value1", n2.Get("key1"))
	assert.Equal("value2", n2.Get("key2"))
	// Test the id
	assert.Equal("node1", n.ID())
}

func TestState_StateReadOnly(t *testing.T) {
	assert := assert.New(t)
	s := goraff.State{}
	// Test that this creates a new node state
	n := s.NodeState("node1")
	n.Set("key1", "value1")
	n.Set("key2", "value2")
	r := s.ReadOnly()
	nr := r.NodeState("node1")
	assert.Equal("value1", nr.Get("key1"))
	assert.Equal("value2", nr.Get("key2"))
}
