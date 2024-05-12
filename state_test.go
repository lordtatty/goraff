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
	n := goraff.NodeState{}
	assert.Equal("", n.Reader().GetStr("key1"))
}

func TestNodeState_ID(t *testing.T) {
	assert := assert.New(t)
	n := goraff.NodeState{}
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

// func TestStateReadOnly_Outputs(t *testing.T) {
// 	assert := assert.New(t)
// 	s := goraff.State{}
// 	n := s.NodeStateUpsert("node1")
// 	n.SetStr("key1", "value1")
// 	n.SetStr("key2", "value2")
// 	n2 := s.NodeStateUpsert("node2")
// 	n2.SetStr("key1", "value1")
// 	n2.SetStr("key2", "value2")
// 	r := s.ReadOnly()
// 	outputs := r.Outputs()
// 	want := []goraff.NodeOutput{
// 		{ID: "node1",
// 			Vals: []goraff.NodeOutputVal{
// 				{Name: "key1", Value: "value1"},
// 				{Name: "key2", Value: "value2"},
// 			},
// 		},
// 		{ID: "node2",
// 			Vals: []goraff.NodeOutputVal{
// 				{Name: "key1", Value: "value1"},
// 				{Name: "key2", Value: "value2"},
// 			},
// 		},
// 	}
// 	assert.Len(outputs, 2)
// 	assert.ElementsMatch(want, outputs)

// 	// change a value
// 	n.SetStr("key1", "valueNEW")

// 	want = []goraff.NodeOutput{
// 		{ID: "node1",
// 			Vals: []goraff.NodeOutputVal{
// 				{Name: "key1", Value: "valueNEW"},
// 				{Name: "key2", Value: "value2"},
// 			},
// 		},
// 		{ID: "node2",
// 			Vals: []goraff.NodeOutputVal{
// 				{Name: "key1", Value: "value1"},
// 				{Name: "key2", Value: "value2"},
// 			},
// 		},
// 	}

// 	outputs = r.Outputs()
// 	assert.ElementsMatch(want, outputs)
// }

func TestNodeState_SubState(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.NodeState{}
	s := &goraff.State{}
	n.SetSubState(s)

	sn := s.NewNodeState("subnode")
	sn.SetStr("key1", "value1")

	assert.Equal("value1", n.SubState().NodeStateByID(sn.Reader().ID()).Reader().GetStr("key1"))
}
