package outputs_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/outputs"
	"github.com/stretchr/testify/assert"
)

func TestOutputter(t *testing.T) {
	assert := assert.New(t)
	substate := &goraff.StateGraph{}
	subnode := substate.NewNodeState("subnode")
	subnode.SetStr("key1", "value1")

	s := goraff.StateGraph{}
	n1 := s.NewNodeState("node1")
	n1.SetStr("key2", "value2")
	n1.SetSubGraph(substate)

	r := s.Reader()

	sut := &outputs.Outputter{}
	result := sut.Output(r)

	want := &outputs.Output{
		PrimaryStateID: s.Reader().ID(),
		States: []outputs.GraphOutput{
			{
				ID:      s.Reader().ID(),
				NodeIDs: []string{n1.Reader().ID()},
			},
			{
				ID:      substate.Reader().ID(),
				NodeIDs: []string{subnode.Reader().ID()},
			},
		},
		Nodes: []outputs.NodeOutput{
			{
				ID:   n1.Reader().ID(),
				Name: n1.Reader().Name(),
				Vals: []outputs.NodeOutputVal{
					{Name: "key2", Value: "value2"},
				},
				SubGraphID: substate.Reader().ID(),
			},
			{
				ID:   subnode.Reader().ID(),
				Name: subnode.Reader().Name(),
				Vals: []outputs.NodeOutputVal{
					{Name: "key1", Value: "value1"},
				},
				SubGraphID: "",
			},
		},
	}
	assert.Equal(want, result)
}
