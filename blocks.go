package goraff

import "fmt"

type Blocks struct {
	blocks []*Block
}

func (b *Blocks) Add(name string, a BlockAction) string {
	n := &Block{Action: a, Name: name}
	b.blocks = append(b.blocks, n)
	return n.Name
}

func (b *Blocks) Get(name string) *Block {
	for _, n := range b.blocks {
		if n.Name == name {
			return n
		}
	}
	return nil
}

func (b *Blocks) All() []*Block {
	return b.blocks
}

func (b *Blocks) Validate() error {
	names := map[string]struct{}{}
	for _, n := range b.blocks {
		if _, ok := names[n.Name]; ok {
			return fmt.Errorf("block name not unique: %s", n.Name)
		}
		names[n.Name] = struct{}{}
	}
	return nil
}

type BlockAction interface {
	Do(s *Node, r *ReadableGraph, triggeringNS *ReadableNode) error
}

// Block represents a node in the graph
type Block struct {
	Action BlockAction
	Name   string
}
