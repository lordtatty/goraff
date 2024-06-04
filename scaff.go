package goraff

import (
	"fmt"
	"sync"
)

type BlockAction interface {
	Do(s *Node, r *ReadableGraph, triggeringNS *ReadableNode) error
}

// Block represents a node in the graph
type Block struct {
	Action BlockAction
	Name   string
}

// Scaff represents blueprint of blocks
// When it runs, it will create a graph of data
type Scaff struct {
	blocks     []*Block
	entrypoint *Block
	state      *Graph
	joins      map[string][]*Join
}

func NewScaff() *Scaff {
	return &Scaff{}
}

func (g *Scaff) AddBlock(name string, a BlockAction) string {
	n := &Block{Action: a, Name: name}
	g.blocks = append(g.blocks, n)
	return n.Name
}

func (g *Scaff) SetEntrypoint(id string) {
	n := g.blockByID(id)
	g.entrypoint = n
}

func (g *Scaff) blockByID(id string) *Block {
	for _, n := range g.blocks {
		if n.Name == id {
			return n
		}
	}
	return nil
}

type ErrBlockNotFound struct {
	ID string
}

func (e ErrBlockNotFound) Error() string {
	return "block not found: " + e.ID
}

func (g *Scaff) AddJoin(fromID, toID string, condition FollowIf) error {
	from := g.blockByID(fromID)
	if from == nil {
		return ErrBlockNotFound{
			ID: fromID,
		}
	}
	to := g.blockByID(toID)
	if to == nil {
		return ErrBlockNotFound{
			ID: toID,
		}
	}
	if g.joins == nil {
		g.joins = make(map[string][]*Join)
	}
	e := &Join{From: from, To: to, Condition: condition}
	if _, ok := g.joins[from.Name]; !ok {
		g.joins[from.Name] = []*Join{}
	}
	g.joins[from.Name] = append(g.joins[from.Name], e)
	return nil
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
	// check block names are unique
	names := map[string]struct{}{}
	for _, b := range g.blocks {
		if _, ok := names[b.Name]; ok {
			return fmt.Errorf("block name not unique: %s", b.Name)
		}
		names[b.Name] = struct{}{}
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
					fmt.Println("error running block, letting all active blocks drain: ", n.Block.Name)
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
	s := g.state.NewNode(n.Name)
	r := NewReadableGraph(g.state)
	err := n.Action.Do(s, r, triggeringNS)
	if err != nil {
		return nil, nil, err
	}
	nextBlocks := []*Block{}
	s.MarkDone()
	if joins, ok := g.joins[n.Name]; ok {
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
