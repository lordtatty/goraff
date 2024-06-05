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
	s := &goraff.Node{}
	r := &goraff.ReadableGraph{}
	err := sut.Do(s, r, nil)
	assert.NoError(err)
}

func TestPrint_DoWithTriggeringNode(t *testing.T) {
	assert := assert.New(t)
	sut := nodeactions.Print{}
	s := &goraff.Node{}
	r := &goraff.ReadableGraph{}
	tn := &goraff.Node{}
	triggeringNode := tn.Get()
	err := sut.Do(s, r, triggeringNode)
	assert.NoError(err)
}
