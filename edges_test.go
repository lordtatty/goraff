package goraff_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestEdgeCondition_KeyMatches(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfKeyMatches("node1", "key1", "value1")
	edge := &goraff.Edge{}
	edge.Condition = sut
	state := &goraff.State{}
	state.NewNodeState("node1").SetStr("key1", "value1")
	assert.True(edge.TriggersMet(state.Reader()))
}

func TestEdgeCondition_KeyMatches_Fails(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfKeyMatches("node1", "key1", "value1")
	edge := &goraff.Edge{}
	edge.Condition = sut
	state := &goraff.State{}
	state.NewNodeState("node1").SetStr("key1", "value2")
	assert.False(edge.TriggersMet(state.Reader()))
}

func TestEdgeCondition_NodesCompleted(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfNodesCompleted("node1", "node2")
	edge := &goraff.Edge{}
	edge.Condition = sut
	state := &goraff.State{}
	state.NewNodeState("node1").MarkDone()
	state.NewNodeState("node2").MarkDone()
	assert.True(edge.TriggersMet(state.Reader()))
}

func TestEdgeCondition_NodesCompleted_Fails(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfNodesCompleted("node1", "node2")
	edge := &goraff.Edge{}
	edge.Condition = sut
	state := &goraff.State{}
	state.NewNodeState("node1").MarkDone()
	state.NewNodeState("node2") // not marked done
	assert.False(edge.TriggersMet(state.Reader()))
}

func TestEdgeCondition_NodesCompleted_NodeStateIsNil(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfNodesCompleted("node1", "node2")
	edge := &goraff.Edge{}
	edge.Condition = sut
	state := &goraff.State{}
	// No nodes upserted
	assert.False(edge.TriggersMet(state.Reader()))
}
