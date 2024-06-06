package blockactions_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/blockactions"
	"github.com/stretchr/testify/assert"
)

func TestFanOut_Do(t *testing.T) {
	assert := assert.New(t)

	// Scaff
	scaff := &goraff.Scaff{}
	scaff.Blocks().Add("input1", &blockactions.Input{Value: "value1"})

	// SUT
	sut := blockactions.FanOut{
		Scaff: scaff,
	}

	// Graph should contain some input vals
	graph := &goraff.Graph{}
	n1 := graph.NewNode("input_node", nil)
	n1.SetStr("input1", "value1")

	// Create a node for the SUT action
	rGraph := goraff.NewReadableGraph(graph)
	sutNode := graph.NewNode("sut_block", []*goraff.ReadableNode{n1.Get()})

	// Run the SUT
	sut.Do(sutNode, rGraph, nil)

	// Check the output
	assert.Equal("result1", sutNode.Get().FirstStr("result"))
}
