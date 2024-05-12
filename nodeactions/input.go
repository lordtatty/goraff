package nodeactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type Input struct {
	Value string
}

func (l *Input) Do(s *goraff.StateNode, r *goraff.GraphStateReader, t *goraff.StateNodeReader) error {
	fmt.Println("Running Input Node")
	s.SetStr("result", l.Value)
	return nil
}
