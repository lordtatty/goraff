package goraff_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestOutputter(t *testing.T) {
	assert := assert.New(t)
	substate := &goraff.GraphState{}
	subnode := substate.NewNodeState("subnode")
	subnode.SetStr("key1", "value1")

	s := goraff.GraphState{}
	n1 := s.NewNodeState("node1")
	n1.SetStr("key2", "value2")
	n1.SetSubState(substate)

	r := s.Reader()

	sut := &goraff.Outputter{}
	outputs := sut.Output(r)

	want := &goraff.Output{
		PrimaryStateID: s.Reader().ID(),
		States: []goraff.StateOutput{
			{
				ID:      s.Reader().ID(),
				NodeIDs: []string{n1.Reader().Name()},
			},
			{
				ID:      substate.Reader().ID(),
				NodeIDs: []string{subnode.Reader().Name()},
			},
		},
		Nodes: []goraff.NodeOutput{
			{
				ID:   n1.Reader().ID(),
				Name: n1.Reader().Name(),
				Vals: []goraff.NodeOutputVal{
					{Name: "key2", Value: "value2"},
				},
				SubStateID: substate.Reader().ID(),
			},
			{
				ID:   subnode.Reader().ID(),
				Name: subnode.Reader().Name(),
				Vals: []goraff.NodeOutputVal{
					{Name: "key1", Value: "value1"},
				},
				SubStateID: "",
			},
		},
	}
	assert.Equal(want, outputs)
}
