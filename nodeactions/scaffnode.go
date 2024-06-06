package nodeactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type ScaffNode struct {
	Scaff *goraff.Scaff
}

func (g *ScaffNode) Do(s *goraff.Node, r *goraff.ReadableGraph, triggeringNS *goraff.ReadableNode) error {
	fmt.Println("Running Scaff Node")
	graph := &goraff.Graph{}
	s.AddSubGraph(graph)
	g.Scaff.Go(graph)
	return nil
}
