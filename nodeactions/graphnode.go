package nodeactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type GraphNode struct {
	Graph *goraff.Scaff
}

func (g *GraphNode) Do(s *goraff.Node, r *goraff.ReadableGraph, triggeringNS *goraff.ReadableNode) error {
	fmt.Println("Running Graph Node")
	sub := g.Graph.Graph()
	s.SetSubGraph(sub)
	g.Graph.Go()
	return nil
}
