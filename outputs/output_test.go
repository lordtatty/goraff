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
	subnode := subgraph.NewNode("subnode")
	subnode.SetStr("key1", "value1")
	subgraphReadable := goraff.NewReadableGraph(subgraph)

	g := &goraff.Graph{}
	n1 := g.NewNode("node1")
	n1.SetStr("key2", "value2")
	n1.SetSubGraph(subgraph)

	r := goraff.NewReadableGraph(g)

	sut := &outputs.Outputter{}
	result := sut.Output(r)

	want := &outputs.Output{
		PrimaryStateID: r.ID(),
		States: []outputs.GraphOutput{
			{
				ID:      r.ID(),
				NodeIDs: []string{n1.Reader().ID()},
			},
			{
				ID:      subgraphReadable.ID(),
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
				SubGraphID: subgraphReadable.ID(),
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
