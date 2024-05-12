package nodeactions_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/nodeactions"
	"github.com/stretchr/testify/assert"
)

func TestPrint_Do(t *testing.T) {
	assert := assert.New(t)
	sut := nodeactions.Print{}
	s := &goraff.StateNode{}
	r := &goraff.GraphStateReader{}
	err := sut.Do(s, r, nil)
	assert.NoError(err)
}

func TestPrint_DoWithTriggeringNode(t *testing.T) {
	assert := assert.New(t)
	sut := nodeactions.Print{}
	s := &goraff.StateNode{}
	r := &goraff.GraphStateReader{}
	tn := &goraff.StateNode{}
	triggeringNode := tn.Reader()
	err := sut.Do(s, r, triggeringNode)
	assert.NoError(err)
}
