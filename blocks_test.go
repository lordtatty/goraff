package goraff_test

import (
	"testing"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlocks_Add(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.Blocks{}

	assert.Len(sut.All(), 0)

	sut.Add("block1", &actionMock{name: "action1"})
	assert.Len(sut.All(), 1)

	sut.Add("block2", &actionMock{name: "action2"})
	assert.Len(sut.All(), 2)
}

func TestBlocks_Get(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	sut := goraff.Blocks{}

	assert.Nil(sut.Get("block1"))

	sut.Add("block1", &actionMock{name: "action1"})
	sut.Add("block2", &actionMock{name: "action2"})
	sut.Add("block3", &actionMock{name: "action3"})

	result := sut.Get("block2")
	require.NotNil(result)
	assert.Equal("block2", result.Name)
}

func TestAll(t *testing.T) {
	assert := assert.New(t)
	sut := goraff.Blocks{}

	assert.Len(sut.All(), 0)

	sut.Add("block1", &actionMock{name: "action1"})
	sut.Add("block2", &actionMock{name: "action2"})
	sut.Add("block3", &actionMock{name: "action3"})

	result := sut.All()
	assert.Len(result, 3)
	assert.Equal("block1", result[0].Name)
	assert.Equal("block2", result[1].Name)
	assert.Equal("block3", result[2].Name)
}
