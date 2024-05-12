package nodeactions

import (
	"log"

	"github.com/lordtatty/goraff"
)

type Print struct {
}

func (p *Print) Do(s *goraff.NodeState, r *goraff.StateReadOnly, triggeringNode *goraff.NodeState) error {
	log.Printf("Node triggered by %s\n", triggeringNode.Reader().ID())
	return nil
}
