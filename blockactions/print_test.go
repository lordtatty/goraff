package blockactions_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/lordtatty/goraff/blockactions"
	"github.com/stretchr/testify/assert"
)

func TestPrint_Do(t *testing.T) {
	assert := assert.New(t)
	sut := blockactions.Print{}
	s := &goraff.Node{}
	r := &goraff.ReadableGraph{}
	err := sut.Do(s, r, nil)
	assert.NoError(err)
}

func TestPrint_DoWithTriggeringNode(t *testing.T) {
	assert := assert.New(t)
	sut := blockactions.Print{}
	s := &goraff.Node{}
	r := &goraff.ReadableGraph{}
	tn := &goraff.Node{}
	triggeringNode := tn.Get()
	err := sut.Do(s, r, triggeringNode)
	assert.NoError(err)
}
