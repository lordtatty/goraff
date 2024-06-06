package blockactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type FanOut struct {
	Scaff  *goraff.Scaff
	InKey  string
	OutKey string
}

func (f *FanOut) Do(s *goraff.Node, r *goraff.ReadableGraph, previousNode *goraff.ReadableNode) error {
	fmt.Println("Running Scaff Node")
	outputs := previousNode.AllStr("result")
	if f.OutKey == "" {
		f.OutKey = "result"
	}
	for _, output := range outputs {
		graph := &goraff.Graph{}
		graph.NewNode(f.InKey, nil).AddStr("result", output)
		s.AddSubGraph(graph)
		err := f.Scaff.Go(graph)
		if err != nil {
			return fmt.Errorf("error running subgraph: %s", err.Error())
		}
		outNode := graph.FirstNodeByName(f.OutKey)
		if outNode == nil {
			return fmt.Errorf("could not find out node for node name: %s", f.OutKey)
		}
		result := outNode.Get().AllStr("result")
		for _, r := range result {
			s.AddStr("result", r)
		}
	}
	return nil
}
