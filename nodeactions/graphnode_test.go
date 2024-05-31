package nodeactions_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/nodeactions"
	"github.com/stretchr/testify/assert"
)

func TestGraphNode_Do(t *testing.T) {
	assert := assert.New(t)

	subScaff := &goraff.Scaff{}
	input1 := subScaff.AddBlock("input1", &nodeactions.Input{Value: "value1"})
	subScaff.SetEntrypoint(input1)

	sut := &nodeactions.ScaffNode{
		Scaff: subScaff,
	}

	scaff := goraff.NewScaff()
	n1 := scaff.AddBlock("sut_block", sut)
	scaff.SetEntrypoint(n1)
	scaff.AddEdge(n1, input1, nil)

	graph := &goraff.Graph{}
	err := scaff.Go(graph)
	assert.Nil(err)

	node := graph.FirstNodeStateByName(n1)
	sub := node.Reader().SubGraph()
	n, err := sub.FirstNodeStateByName(input1)
	assert.Nil(err)
	assert.Equal("value1", n.GetStr("result"))
}
