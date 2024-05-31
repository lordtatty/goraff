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
	graph := &goraff.Graph{}
	graph.NewNode("node1").SetStr("key1", "value1")
	readable := goraff.NewReadableGraph(graph)
	assert.True(edge.TriggersMet(readable))
}

func TestEdgeCondition_KeyMatches_Fails(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfKeyMatches("node1", "key1", "value1")
	edge := &goraff.Edge{}
	edge.Condition = sut
	graph := &goraff.Graph{}
	graph.NewNode("node1").SetStr("key1", "value2")
	readable := goraff.NewReadableGraph(graph)
	assert.False(edge.TriggersMet(readable))
}

func TestEdgeCondition_NodesCompleted(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfNodesCompleted("node1", "node2")
	edge := &goraff.Edge{}
	edge.Condition = sut
	graph := &goraff.Graph{}
	graph.NewNode("node1").MarkDone()
	graph.NewNode("node2").MarkDone()
	readable := goraff.NewReadableGraph(graph)
	assert.True(edge.TriggersMet(readable))
}

func TestEdgeCondition_NodesCompleted_Fails(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfNodesCompleted("node1", "node2")
	edge := &goraff.Edge{}
	edge.Condition = sut
	graph := &goraff.Graph{}
	graph.NewNode("node1").MarkDone()
	graph.NewNode("node2") // not marked done
	readable := goraff.NewReadableGraph(graph)
	assert.False(edge.TriggersMet(readable))
}

func TestEdgeCondition_NodesCompleted_NodeStateIsNil(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.FollowIfNodesCompleted("node1", "node2")
	edge := &goraff.Edge{}
	edge.Condition = sut
	graph := &goraff.Graph{}
	// No nodes upserted
	readable := goraff.NewReadableGraph(graph)
	assert.False(edge.TriggersMet(readable))
}
