package goraff_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestNodeState_SubState(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.StateNode{}
	s := &goraff.StateGraph{}
	n.SetSubGraph(s)

	sn := s.NewNodeState("subnode")
	sn.SetStr("key1", "value1")

	subGraph := n.Reader().SubGraph()
	node := subGraph.NodeStateByID(sn.Reader().ID())
	assert.Equal("value1", node.GetStr("key1"))
}

func TestStateNode_Reader(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.StateNode{}
	r := n.Reader()
	assert.Equal(n.Reader().ID(), r.ID())
}
