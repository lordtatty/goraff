package goraff_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestJoinCondition_KeyMatches(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfKeyMatches("node1", "key1", "value1")
	join := &goraff.Join{}
	join.Condition = sut
	graph := &goraff.Graph{}
	graph.NewNode("node1", nil).SetStr("key1", "value1")
	readable := goraff.NewReadableGraph(graph)
	assert.True(join.TriggersMet(readable))
}

func TestJoinCondition_KeyMatches_Fails(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfKeyMatches("node1", "key1", "value1")
	join := &goraff.Join{}
	join.Condition = sut
	graph := &goraff.Graph{}
	graph.NewNode("node1", nil).SetStr("key1", "value2")
	readable := goraff.NewReadableGraph(graph)
	assert.False(join.TriggersMet(readable))
}

func TestJoinCondition_NodesCompleted(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfNodesCompleted("node1", "node2")
	join := &goraff.Join{}
	join.Condition = sut
	graph := &goraff.Graph{}
	graph.NewNode("node1", nil).MarkDone()
	graph.NewNode("node2", nil).MarkDone()
	readable := goraff.NewReadableGraph(graph)
	assert.True(join.TriggersMet(readable))
}

func TestJoinCondition_NodesCompleted_Fails(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfNodesCompleted("node1", "node2")
	join := &goraff.Join{}
	join.Condition = sut
	graph := &goraff.Graph{}
	graph.NewNode("node1", nil).MarkDone()
	graph.NewNode("node2", nil) // not marked done
	readable := goraff.NewReadableGraph(graph)
	assert.False(join.TriggersMet(readable))
}

func TestJoinCondition_NodesCompleted_NodeStateIsNil(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfNodesCompleted("node1", "node2")
	join := &goraff.Join{}
	join.Condition = sut
	graph := &goraff.Graph{}
	// No nodes upserted
	readable := goraff.NewReadableGraph(graph)
	assert.False(join.TriggersMet(readable))
}
