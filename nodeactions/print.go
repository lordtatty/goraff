package nodeactions

import (
	"log"

	"github.com/lordtatty/goraff"
)

type Print struct {
}

func (p *Print) Do(s *goraff.Node, r *goraff.ReadableGraph, triggeringNode *goraff.ReadableNode) error {
	if triggeringNode == nil {
		log.Println("Node triggered by nil")
		return nil
	}
	log.Printf("Node triggered by %s\n", triggeringNode.ID())
	return nil
}
