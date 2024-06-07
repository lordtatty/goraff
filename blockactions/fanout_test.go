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
	InKey  string
	OutKey string
}

func (m *MockActionReverse) Do(n *goraff.Node, r *goraff.ReadableGraph, previousNode *goraff.ReadableNode) error {
	inNode, err := r.FirstNodeByName(m.InNode)
	if err != nil {
		return fmt.Errorf("Could not find node %s", m.InNode)
	}
	outputs := inNode.AllStr(m.InKey)
	for _, output := range outputs {
		reversed := reverseStr(output)
		n.AddStr(m.OutKey, reversed)
	}
	return nil
}

func TestFanOut_Do(t *testing.T) {
	assert := assert.New(t)

	inNodeName := "input_node"
	InNodeKey := "in_result"
	outNodeName := "out_result_node"
	outNodeKey := "out_result_key"

	// Scaff
	scaff := &goraff.Scaff{}
	// scaff.Blocks().Add("input1", &blockactions.Input{Value: "value1"})
	reverseBlock := scaff.Blocks().Add(outNodeName, &MockActionReverse{
		InNode: inNodeName,
		InKey:  InNodeKey,
		OutKey: outNodeKey,
	})
	scaff.SetEntrypoint(reverseBlock)

	// Parent Graph should have one node (which is the "previous node" triggering this action)
	graph := &goraff.Graph{}
	n1 := graph.NewNode(inNodeName, nil)
	n1.AddStr(InNodeKey, "value1")
	n1.AddStr(InNodeKey, "value2")
	n1.AddStr(InNodeKey, "value3")

	// Create a node for the SUT action
	rGraph := goraff.NewReadableGraph(graph)
	sutNode := graph.NewNode("sut_block", []*goraff.ReadableNode{n1.Get()})

	// Create and Run the SUT
	sut := blockactions.FanOut{
		Scaff:   scaff,
		InNode:  n1.Get().Name(),
		InKey:   InNodeKey,
		OutNode: outNodeName,
		OutKey:  outNodeKey,
	}
	err := sut.Do(sutNode, rGraph, n1.Get())
	assert.Nil(err)

	// Asserations
	assert.Len(sutNode.Get().SubGraph(), 3)
	subGraphs := sutNode.Get().SubGraph()
	subInNode1, _ := subGraphs[0].FirstNodeByName(sut.InNode)
	subInNode2, _ := subGraphs[1].FirstNodeByName(sut.InNode)
	subInNode3, _ := subGraphs[2].FirstNodeByName(sut.InNode)

	// Assert the subgraphs have been given the correct input values
	assert.Equal("value1", subInNode1.FirstStr(sut.InKey))
	assert.Equal("value2", subInNode2.FirstStr(sut.InKey))
	assert.Equal("value3", subInNode3.FirstStr(sut.InKey))

	// Assert the inputs have been correctly reversed within each subgraph
	subRevNode1, _ := subGraphs[0].FirstNodeByName(outNodeName)
	subRevNode2, _ := subGraphs[1].FirstNodeByName(outNodeName)
	subRevNode3, _ := subGraphs[2].FirstNodeByName(outNodeName)
	assert.Equal("1eulav", subRevNode1.FirstStr(outNodeKey))
	assert.Equal("2eulav", subRevNode2.FirstStr(outNodeKey))
	assert.Equal("3eulav", subRevNode3.FirstStr(outNodeKey))

	// Assert final output node result
	assert.Len(sutNode.Get().All("result"), 3)
	assert.Equal("1eulav", sutNode.Get().FirstStr("result"))
	assert.Equal("2eulav", sutNode.Get().AllStr("result")[1])
	assert.Equal("3eulav", sutNode.Get().AllStr("result")[2])

}
