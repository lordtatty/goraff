package goraff

import (
	"fmt"
	"sync"
)

// Scaff represents blueprint of blocks
// When it runs, it will create a graph of data
type Scaff struct {
	entrypoint *Block
	state      *Graph
	joins      *Joins
	blocks     *Blocks
}

func NewScaff() *Scaff {
	return &Scaff{}
}

func (g *Scaff) Blocks() *Blocks {
	if g.blocks == nil {
		g.blocks = &Blocks{}
	}
	return g.blocks
}

func (g *Scaff) Joins() *Joins {
	if g.joins == nil {
		g.joins = &Joins{
			Blocks: g.Blocks(),
		}
	}
	return g.joins
}

func (g *Scaff) SetEntrypoint(name string) {
	n := g.blocks.Get(name)
	g.entrypoint = n
}

func (g *Scaff) Go(graph *Graph) error {
	if graph == nil {
		return fmt.Errorf("graph not provided")
	}
	err := g.validate()
	if err != nil {
		return fmt.Errorf("error validating graph: %w", err)
	}
	g.state = graph
	return g.flowMgr()
}

func (g *Scaff) validate() error {
	if g.entrypoint == nil {
		return fmt.Errorf("entrypoint not set")
	}
	// check blocks
	err := g.Blocks().Validate()
	if err != nil {
		return fmt.Errorf("error validating blocks: %w", err)
	}
	// check joins
	err = g.Joins().Validate()
	if err != nil {
		return fmt.Errorf("error validating joins: %w", err)
	}
	return nil
}

type nextBlock struct {
	Block        *Block
	triggeringNS *Node
}

func (g *Scaff) flowMgr() error {
	if g.entrypoint == nil {
		return fmt.Errorf("entrypoint not set")
	}

	completedCh := make(chan nextBlock, 10)
	var wg sync.WaitGroup

	completedCh <- nextBlock{
		Block:        g.entrypoint,
		triggeringNS: nil,
	}
	wg.Add(1) // Increment for the initial node

	fmt.Println("starting block", g.entrypoint.Name)
	var foundErr error
	mut := sync.Mutex{}
	go func() {
		for n := range completedCh {
			go func(n nextBlock) {
				fmt.Println("completed block", n.Block.Name)
				defer wg.Done() // Ensure we mark this goroutine as done on finish
				if foundErr != nil {
					return
				}
				var tr *ReadableNode = nil
				if n.triggeringNS != nil {
					tr = n.triggeringNS.Reader()
				}
				nextBlocks, compeltedState, err := g.runBlock(n.Block, tr)
				if err != nil {
					fmt.Printf("error running block %s, letting all active blocks drain: %s \n", n.Block.Name, err.Error())
					mut.Lock()
					foundErr = fmt.Errorf("error running block: %w", err)
					mut.Unlock()
					return
				}
				for _, next := range nextBlocks {
					fmt.Println("adding block", next.Name)
					wg.Add(1) // Increment for each new block
					completedCh <- nextBlock{
						Block:        next,
						triggeringNS: compeltedState,
					}
				}
			}(n)
		}
	}()

	wg.Wait()          // Wait for all goroutines to finish
	close(completedCh) // Safe to close here as no more writes will happen
	return foundErr
}

func (g *Scaff) runBlock(n *Block, triggeringNS *ReadableNode) ([]*Block, *Node, error) {
	s := g.state.NewNode(n.Name, nil)
	r := NewReadableGraph(g.state)
	err := n.Action.Do(s, r, triggeringNS)
	if err != nil {
		return nil, nil, err
	}
	nextBlocks := []*Block{}
	s.MarkDone()
	if joins, ok := g.Joins().Get(n.Name); ok {
		for _, e := range joins {
			t, err := e.TriggersMet(r)
			if err != nil {
				return nil, nil, fmt.Errorf("error checking join condition: %w", err)
			}
			if t {
				nextBlocks = append(nextBlocks, e.To)
			}
		}
	}
	return nextBlocks, s, nil
}
