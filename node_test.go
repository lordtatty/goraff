package goraff_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestNodeState_SubState(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	s := &goraff.Graph{}
	n.SetSubGraph(s)

	sn := s.NewNodeState("subnode")
	sn.SetStr("key1", "value1")

	subGraph := n.Reader().SubGraph()
	node, err := subGraph.NodeByID(sn.Reader().ID())
	assert.Nil(err)
	assert.Equal("value1", node.GetStr("key1"))
}

func TestStateNode_Reader(t *testing.T) {
	assert := assert.New(t)
	n := &goraff.Node{}
	r := n.Reader()
	assert.Equal(n.Reader().ID(), r.ID())
}
