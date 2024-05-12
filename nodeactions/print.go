package nodeactions

import (
	"log"

	"github.com/lordtatty/goraff"
)

type Print struct {
}

func (p *Print) Do(s *goraff.StateNode, r *goraff.StateReadOnly, triggeringNode *goraff.StateNode) error {
	log.Printf("Node triggered by %s\n", triggeringNode.Reader().ID())
	return nil
}
