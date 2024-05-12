package nodeactions

import (
	"fmt"

	"github.com/lordtatty/goraff"
)

type Input struct {
	Value string
}

func (l *Input) Do(s *goraff.NodeState, r *goraff.StateReadOnly, t *goraff.NodeState) error {
	fmt.Println("Running Input Node")
	s.SetStr("result", l.Value)
	return nil
}
