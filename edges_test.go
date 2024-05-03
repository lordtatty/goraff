package goraff_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestEdgeCondition_KeyMatches(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.EdgeConditionKeyMatches("node1", "key1", "value1")
	edge := &goraff.Edge{}
	edge.Condition = sut
	state := &goraff.State{}
	state.NodeState("node1").Set("key1", "value1")
	assert.True(edge.Match(state.ReadOnly()))
}

func TestEdgeCondition_KeyMatches_Fails(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.EdgeConditionKeyMatches("node1", "key1", "value1")
	edge := &goraff.Edge{}
	edge.Condition = sut
	state := &goraff.State{}
	state.NodeState("node1").Set("key1", "value2")
	assert.False(edge.Match(state.ReadOnly()))
}
