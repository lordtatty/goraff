package nodeactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type GraphNode struct {
	Graph *goraff.Graph
}

func (g *GraphNode) Do(s *goraff.NodeState, r *goraff.StateReadOnly, triggeringNodeID string) error {
	fmt.Println("Running Graph Node")
	g.Graph.Go()
	s.SetStr("result", "Graph Node Result")
	return nil
}
