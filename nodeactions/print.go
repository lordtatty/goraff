package nodeactions

import (
	"log"

	"github.com/lordtatty/goraff"
)

type Print struct {
}

func (p *Print) Do(s *goraff.StateNode, r *goraff.GraphStateReader, triggeringNode *goraff.StateNodeReader) error {
	log.Printf("Node triggered by %s\n", triggeringNode.ID())
	return nil
}
