package goraff_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestNodeState(t *testing.T) {
	assert := assert.New(t)
	n := goraff.NodeState{}
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	assert.Equal("value1", n.GetStr("key1"))
	assert.Equal("value2", n.GetStr("key2"))
	n.SetStr("key3", "value3")
	assert.Equal("value1", n.GetStr("key1"))
	assert.Equal("value2", n.GetStr("key2"))
	assert.Equal("value3", n.GetStr("key3"))
}

func TestNodeState_GetOnUninitialisedState(t *testing.T) {
	assert := assert.New(t)
	n := goraff.NodeState{}
	assert.Equal("", n.GetStr("key1"))
}

func TestState_NodeState(t *testing.T) {
	assert := assert.New(t)
	s := goraff.State{}
	// Test that this creates a new node state
	n := s.NodeStateUpsert("node1")
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	assert.Equal("value1", n.GetStr("key1"))
	assert.Equal("value2", n.GetStr("key2"))
	// Test that this returns the same already-created state
	n2 := s.NodeStateUpsert("node1")
	assert.Equal("value1", n2.GetStr("key1"))
	assert.Equal("value2", n2.GetStr("key2"))
	// Test the id
	assert.Equal("node1", n.ID())
}

func TestState_StateReadOnly(t *testing.T) {
	assert := assert.New(t)
	s := goraff.State{}
	// Test that this creates a new node state
	n := s.NodeStateUpsert("node1")
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := s.ReadOnly()
	nr := r.NodeState("node1")
	assert.Equal("value1", nr.GetStr("key1"))
	assert.Equal("value2", nr.GetStr("key2"))
}

func TestState_StateReadOnly_ID(t *testing.T) {
	assert := assert.New(t)
	s := goraff.State{}
	// Test that this creates a new node state
	n := s.NodeStateUpsert("node1")
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	r := s.ReadOnly()
	nr := r.NodeState("node1")
	assert.Equal("node1", nr.ID())
}

func TestState_OnUpdate(t *testing.T) {
	assert := assert.New(t)
	updated := false
	s := goraff.State{
		OnUpdate: []func(s *goraff.StateReadOnly){
			func(s *goraff.StateReadOnly) {
				assert.Equal("value", s.NodeState("node1").GetStr("key"))
				updated = true
			},
		},
	}
	n := s.NodeStateUpsert("node1")
	assert.False(updated)
	n.SetStr("key", "value")
	assert.True(updated)
}

func TestStateReadOnly_Outputs(t *testing.T) {
	assert := assert.New(t)
	s := goraff.State{}
	n := s.NodeStateUpsert("node1")
	n.SetStr("key1", "value1")
	n.SetStr("key2", "value2")
	n2 := s.NodeStateUpsert("node2")
	n2.SetStr("key1", "value1")
	n2.SetStr("key2", "value2")
	r := s.ReadOnly()
	outputs := r.Outputs()
	want := []goraff.NodeOutput{
		{ID: "node1",
			Vals: []goraff.NodeOutputVal{
				{Name: "key1", Value: "value1"},
				{Name: "key2", Value: "value2"},
			},
		},
		{ID: "node2",
			Vals: []goraff.NodeOutputVal{
				{Name: "key1", Value: "value1"},
				{Name: "key2", Value: "value2"},
			},
		},
	}
	assert.Len(outputs, 2)
	assert.ElementsMatch(want, outputs)

	// change a value
	n.SetStr("key1", "valueNEW")

	want = []goraff.NodeOutput{
		{ID: "node1",
			Vals: []goraff.NodeOutputVal{
				{Name: "key1", Value: "valueNEW"},
				{Name: "key2", Value: "value2"},
			},
		},
		{ID: "node2",
			Vals: []goraff.NodeOutputVal{
				{Name: "key1", Value: "value1"},
				{Name: "key2", Value: "value2"},
			},
		},
	}

	outputs = r.Outputs()
	assert.ElementsMatch(want, outputs)
}

func TestNodeState_SubState(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.NodeState{}
	s := &goraff.State{}
	n.SetSubState(s)

	sn := s.NodeStateUpsert("subnode")
	sn.SetStr("key1", "value1")

	assert.Equal("value1", n.SubState().NodeState("subnode").GetStr("key1"))

}
