package nodeactions_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/nodeactions"
	"github.com/stretchr/testify/assert"
)

func TestGraphNode_Do(t *testing.T) {
	assert := assert.New(t)

	subgraph := &goraff.Graph{}
	substate := subgraph.State()
	input1 := subgraph.AddNodeWithName("input1", &nodeactions.Input{Value: "value1"})
	subgraph.SetEntrypoint(input1)

	sut := &nodeactions.GraphNode{
		Graph: subgraph,
	}

	graph := goraff.New()
	n1 := graph.AddNode(sut)
	graph.SetEntrypoint(n1)

	err := graph.Go()
	assert.Nil(err)

	node := graph.State().FirstNodeStateByName(n1)
	sub := node.Reader().SubGraph()
	assert.Equal(substate.Reader(), sub)
	n, err := sub.FirstNodeStateByName(input1)
	assert.Nil(err)
	assert.Equal("value1", n.GetStr("result"))
}
