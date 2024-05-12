package nodeactions

import (
	"log"

	"github.com/lordtatty/goraff"
)

type Print struct {
}

func (p *Print) Do(s *goraff.StateNode, r *goraff.GraphStateReader, triggeringNode *goraff.StateNodeReader) error {
	if triggeringNode == nil {
		log.Println("Node triggered by nil")
		return nil
	}
	log.Printf("Node triggered by %s\n", triggeringNode.ID())
	return nil
}
