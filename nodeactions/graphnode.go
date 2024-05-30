package nodeactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type ScaffNode struct {
	Scaff *goraff.Scaff
}

func (g *ScaffNode) Do(s *goraff.Node, r *goraff.ReadableGraph, triggeringNS *goraff.ReadableNode) error {
	fmt.Println("Running Graph Node")
	graph := &goraff.Graph{}
	s.SetSubGraph(graph)
	g.Scaff.Go(graph)
	return nil
}
