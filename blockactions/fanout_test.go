package blockactions_test

import (
	"fmt"
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/blockactions"
	"github.com/stretchr/testify/assert"
)

func reverseStr(s string) string {
	// Convert string to a slice of runes to handle Unicode correctly
	runes := []rune(s)
	// Get the length of the slice
	n := len(runes)

	// Reverse the slice
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}

	// Convert the slice of runes back to a string
	return string(runes)
}

type MockActionReverse struct {
	InNode string
}

func (m *MockActionReverse) Do(s *goraff.Node, r *goraff.ReadableGraph, previousNode *goraff.ReadableNode) error {
	inNode, err := r.FirstNodeByName(m.InNode)
	if err != nil {
		return fmt.Errorf("Could not find node %s", m.InNode)
	}
	outputs := inNode.AllStr("result")
	for _, output := range outputs {
		reversed := reverseStr(output)
		s.AddStr("result", reversed)
	}
	return nil
}

func TestFanOut_Do(t *testing.T) {
	assert := assert.New(t)

	// Scaff
	scaff := &goraff.Scaff{}
	// scaff.Blocks().Add("input1", &blockactions.Input{Value: "value1"})
	reverseBlock := scaff.Blocks().Add("result", &MockActionReverse{InNode: "input_node"})
	scaff.SetEntrypoint(reverseBlock)

	// Parent Graph should have one node (which is the "previous node" triggering this action)
	graph := &goraff.Graph{}
	n1 := graph.NewNode("input_node", nil)
	n1.AddStr("result", "value1")
	n1.AddStr("result", "value2")
	n1.AddStr("result", "value3")

	// Create a node for the SUT action
	rGraph := goraff.NewReadableGraph(graph)
	sutNode := graph.NewNode("sut_block", []*goraff.ReadableNode{n1.Get()})

	// Create and Run the SUT
	sut := blockactions.FanOut{
		Scaff: scaff,
		InKey: n1.Get().Name(),
	}
	err := sut.Do(sutNode, rGraph, n1.Get())
	assert.Nil(err)

	// Asserations
	assert.Len(sutNode.Get().SubGraph(), 3)
	subGraphs := sutNode.Get().SubGraph()
	subInNode1, _ := subGraphs[0].FirstNodeByName(sut.InKey)
	subInNode2, _ := subGraphs[1].FirstNodeByName(sut.InKey)
	subInNode3, _ := subGraphs[2].FirstNodeByName(sut.InKey)

	// Assert the subgraphs have been given the correct input values
	assert.Equal("value1", subInNode1.FirstStr("result"))
	assert.Equal("value2", subInNode2.FirstStr("result"))
	assert.Equal("value3", subInNode3.FirstStr("result"))

	// Assert the inputs have been correctly reversed
	subRevNode1, _ := subGraphs[0].FirstNodeByName("result")
	subRevNode2, _ := subGraphs[1].FirstNodeByName("result")
	subRevNode3, _ := subGraphs[2].FirstNodeByName("result")
	assert.Equal("1eulav", subRevNode1.FirstStr("result"))
	assert.Equal("2eulav", subRevNode2.FirstStr("result"))
	assert.Equal("3eulav", subRevNode3.FirstStr("result"))

	// Assert final node result
	assert.Len(sutNode.Get().All("result"), 3)
	assert.Equal("1eulav", sutNode.Get().FirstStr("result"))
	assert.Equal("2eulav", sutNode.Get().AllStr("result")[1])
	assert.Equal("3eulav", sutNode.Get().AllStr("result")[2])

}
