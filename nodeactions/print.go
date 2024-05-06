package nodeactions

import (
	"log"

	"github.com/lordtatty/goraff"
)

type Print struct {
}

func (p *Print) Do(s *goraff.NodeState, r *goraff.StateReadOnly, triggeringNodeID string) error {
	log.Printf("Node triggered by %s\n", triggeringNodeID)
	return nil
}
