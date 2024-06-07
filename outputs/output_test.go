package outputs_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/outputs"
	"github.com/stretchr/testify/assert"
)

func loadFixtureStr(filename string, replacements map[string]string) string {
	baseDir := "./fixtures"
	filePath := fmt.Sprintf("%s/%s", baseDir, filename)
	b, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	resultStr := string(b)
	for k, v := range replacements {
		key := fmt.Sprintf("<<%s>>", k)
		resultStr = strings.ReplaceAll(resultStr, key, v)
	}
	return resultStr
}

func TestOutputter(t *testing.T) {
	assert := assert.New(t)

	// Subgraph1
	subgraph := &goraff.Graph{}
	subnode := subgraph.NewNode("subnode", nil)
	subnode.SetStr("key1", "value1")
	// subgraphReadable := goraff.NewReadableGraph(subgraph)

	//  Subgraph2
	subgraph2 := &goraff.Graph{}
	subnode2 := subgraph2.NewNode("subnode2", nil)
	subnode2.SetStr("key3", "value3")
	// subgraphReadable2 := goraff.NewReadableGraph(subgraph2)

	// Main Graph
	g := &goraff.Graph{}
	// Node1 has two subgraaphs
	n1 := g.NewNode("node1", nil)
	n1.SetStr("key2", "value2")
	n1.AddSubGraph(subgraph)
	n1.AddSubGraph(subgraph2)

	n2 := g.NewNode("node2", nil)
	n2.AddStr("key", "value0")
	n2.AddStr("key", "value1")
	n2.AddStr("key", "value2")

	r := goraff.NewReadableGraph(g)

	sut := &outputs.Outputter{}
	result := sut.Output(r)

	// Load expected output
	// We need to replace the placehodler IDs with the actual IDs
	r1 := goraff.NewReadableGraph(g)
	r2 := goraff.NewReadableGraph(subgraph)
	r3 := goraff.NewReadableGraph(subgraph2)
	want := loadFixtureStr("testoutputter.json", map[string]string{
		"PRIMARY_GRAPH_ID": r1.ID(),
		"GRAPH2_ID":        r2.ID(),
		"GRAPH3_ID":        r3.ID(),
		"NODE1_ID":         n1.Get().ID(),
		"NODE2_ID":         n2.Get().ID(),
		"SUBNODE_ID":       subnode.Get().ID(),
		"SUBNODE2_ID":      subnode2.Get().ID(),
	})

	// result to json string
	b, err := json.Marshal(result)
	assert.Nil(err)
	resultStr := string(b)

	assert.JSONEq(want, resultStr)
}
