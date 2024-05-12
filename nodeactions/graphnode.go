package nodeactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type GraphNode struct {
	Graph *goraff.Graph
}

func (g *GraphNode) Do(s *goraff.StateNode, r *goraff.GraphStateReader, triggeringNS *goraff.StateNodeReader) error {
	fmt.Println("Running Graph Node")
	sub := g.Graph.State()
	s.SetSubGraph(sub)
	g.Graph.Go()
	return nil
}
