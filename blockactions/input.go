package blockactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type Input struct {
	Value string
}

func (l *Input) Do(s *goraff.Node, r *goraff.ReadableGraph, t *goraff.ReadableNode) error {
	fmt.Println("Running Input Node")
	s.SetStr("result", l.Value)
	return nil
}
