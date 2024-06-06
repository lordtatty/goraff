package blockactions_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/blockactions"
	"github.com/stretchr/testify/assert"
)

func TestGraphNode_Do(t *testing.T) {
	assert := assert.New(t)

	// SubScaff
	subScaff := &goraff.Scaff{}
	input1 := subScaff.Blocks().Add("input1", &blockactions.Input{Value: "value1"})
	subScaff.SetEntrypoint(input1)

	// The SUT
	sut := &blockactions.ScaffNode{
		Scaff: subScaff,
	}

	// Main Scaff
	scaff := goraff.NewScaff()
	n1 := scaff.Blocks().Add("sut_block", sut)
	scaff.SetEntrypoint(n1)

	// Run the Scaff
	graph := &goraff.Graph{}
	err := scaff.Go(graph)
	assert.Nil(err)

	// Check the output
	node := graph.FirstNodeByName(n1)
	subGraphs := node.Get().SubGraph()
	assert.Len(subGraphs, 1)
	sub := subGraphs[0]
	n, err := sub.FirstNodeByName(input1)
	assert.Nil(err)
	assert.Equal("value1", n.FirstStr("result"))
}
