package nodeactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type FanOut struct {
	Scaff          *goraff.Scaff
	IncludeOutputs []string
}

func (g *FanOut) Do(s *goraff.Node, r *goraff.ReadableGraph, triggeringNS *goraff.ReadableNode) error {
	fmt.Println("Running Scaff Node")
	// graph := &goraff.Graph{}
	// s.AddSubGraph(graph)
	// g.Scaff.Go(graph)
	s.AddStr("result", "result1")
	return nil
}
