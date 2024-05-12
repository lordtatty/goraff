package nodeactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type GraphNode struct {
	Graph *goraff.Graph
}

func (g *GraphNode) Do(s *goraff.NodeState, r *goraff.StateReadOnly, triggeringNS *goraff.NodeState) error {
	fmt.Println("Running Graph Node")
	sub := g.Graph.State()
	s.SetSubState(sub)
	g.Graph.Go()
	return nil
}
