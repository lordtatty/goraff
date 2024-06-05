package outputs_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/outputs"
	"github.com/stretchr/testify/assert"
)

func TestOutputter(t *testing.T) {
	assert := assert.New(t)
	subgraph := &goraff.Graph{}
	subnode := subgraph.NewNode("subnode", nil)
	subnode.SetStr("key1", "value1")
	subgraphReadable := goraff.NewReadableGraph(subgraph)

	g := &goraff.Graph{}
	n1 := g.NewNode("node1", nil)
	n1.SetStr("key2", "value2")
	n1.SetSubGraph(subgraph)

	n2 := g.NewNode("node2", nil)
	n2.AddStr("key", "value0")
	n2.AddStr("key", "value1")
	n2.AddStr("key", "value2")

	r := goraff.NewReadableGraph(g)

	sut := &outputs.Outputter{}
	result := sut.Output(r)

	want := &outputs.Output{
		PrimaryStateID: r.ID(),
		States: []outputs.GraphOutput{
			{
				ID:      r.ID(),
				NodeIDs: []string{n1.Get().ID(), n2.Get().ID()},
			},
			{
				ID:      subgraphReadable.ID(),
				NodeIDs: []string{subnode.Get().ID()},
			},
		},
		Nodes: []outputs.NodeOutput{
			{
				ID:   n1.Get().ID(),
				Name: n1.Get().Name(),
				Vals: []outputs.NodeOutputVal{
					{Name: "key2", Value: "value2"},
				},
				SubGraphID: subgraphReadable.ID(),
			},
			{
				ID:   subnode.Get().ID(),
				Name: subnode.Get().Name(),
				Vals: []outputs.NodeOutputVal{
					{Name: "key1", Value: "value1"},
				},
				SubGraphID: "",
			},
			{
				ID:   n2.Get().ID(),
				Name: n2.Get().Name(),
				Vals: []outputs.NodeOutputVal{
					{Name: "key_0", Value: "value0"},
					{Name: "key_1", Value: "value1"},
					{Name: "key_2", Value: "value2"},
				},
				SubGraphID: "",
			},
		},
	}
	assert.Equal(want, result)
}
