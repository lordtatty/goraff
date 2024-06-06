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

type nextJoin struct {
	Join         *Join
	previousNode *Node
}

func (g *Scaff) flowMgr() error {
	if g.entrypoint == nil {
		return fmt.Errorf("entrypoint not set")
	}

	completedCh := make(chan nextJoin, 10)
	var wg sync.WaitGroup

	completedCh <- nextJoin{
		Join:         &Join{From: nil, To: g.entrypoint},
		previousNode: nil,
	}
	wg.Add(1) // Increment for the initial node

	fmt.Println("starting block", g.entrypoint.Name)
	var foundErr error
	mut := sync.Mutex{}
	go func() {
		for n := range completedCh {
			// check Trigger before launching goroutine to prevent join race conditions
			if n.previousNode != nil {
				n.previousNode.MarkDone()
			}
			if n.Join == nil {
				wg.Done()
				continue
			}
			fmt.Println("considering block", n.Join.To.Name)
			r := NewReadableGraph(g.state)
			t, err := n.Join.TriggersMet(r)
			if err != nil {
				fmt.Printf("error checking join condition: %s\n", err.Error())
				wg.Done()
				continue
			}
			if !t {
				fmt.Printf("join condition not met To: %s\n", n.Join.To.Name)
				wg.Done()
				continue
			}
			fmt.Printf("join condition met To: %s\n", n.Join.To.Name)
			// launch goroutine
			go func(n nextJoin) {
				defer wg.Done() // Ensure we mark this goroutine as done on finish
				// run block
				block := n.Join.To
				defer fmt.Printf("finished block %s\n", n.Join.To.Name)
				fmt.Println("starting block", block.Name)
				if foundErr != nil {
					return
				}
				var tr *ReadableNode = nil
				if n.previousNode != nil {
					tr = n.previousNode.Get()
				}
				completedNode, err := g.runBlock(block, tr)
				if err != nil {
					fmt.Printf("error running block %s, letting all active blocks drain: %s \n", block.Name, err.Error())
					mut.Lock()
					foundErr = fmt.Errorf("error running block: %w", err)
					mut.Unlock()
					return
				}
				joins := g.Joins().Get(block.Name)
				for _, j := range joins {
					fmt.Println("queueing block join", j.To.Name)
					wg.Add(1) // Increment for each new block
					completedCh <- nextJoin{
						previousNode: completedNode,
						Join:         j,
					}
				}
				if len(joins) == 0 {
					wg.Add(1) // Increment
					completedCh <- nextJoin{
						previousNode: completedNode,
						Join:         nil,
					}
				}
			}(n)
		}
	}()

	wg.Wait()          // Wait for all goroutines to finish
	close(completedCh) // Safe to close here as no more writes will happen
	return foundErr
}

func (g *Scaff) runBlock(n *Block, triggeringNS *ReadableNode) (*Node, error) {
	s := g.state.NewNode(n.Name, nil)
	r := NewReadableGraph(g.state)
	err := n.Action.Do(s, r, triggeringNS)
	if err != nil {
		return nil, err
	}
	// s.MarkDone()
	return s, nil
}
