package goraff

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
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
	edges      map[string][]*Edge
}

func NewScaff() *Scaff {
	return &Scaff{}
}

func (g *Scaff) StateReadOnly() *ReadableGraph {
	if g.state == nil {
		g.state = &Graph{}
	}
	return g.state.Reader()
}

func (g *Scaff) Len() int {
	return len(g.blocks)
}

func (g *Scaff) AddBlock(a BlockAction) string {
	id := uuid.New().String()
	return g.AddBlockWithName(id, a)
}

func (g *Scaff) AddBlockWithName(name string, a BlockAction) string {
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
	return "node not found: " + e.ID
}

func (g *Scaff) AddEdge(fromID, toID string, condition FollowIf) error {
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
	if g.edges == nil {
		g.edges = make(map[string][]*Edge)
	}
	e := &Edge{From: from, To: to, Condition: condition}
	if _, ok := g.edges[from.Name]; !ok {
		g.edges[from.Name] = []*Edge{}
	}
	g.edges[from.Name] = append(g.edges[from.Name], e)
	return nil
}

func (g *Scaff) Go(graph *Graph) error {
	if graph == nil {
		return fmt.Errorf("graph not provided")
	}
	g.state = graph
	return g.flowMgr()
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
	s := g.state.NewNodeState(n.Name)
	r := g.state.Reader()
	err := n.Action.Do(s, r, triggeringNS)
	if err != nil {
		return nil, nil, err
	}
	nextBlocks := []*Block{}
	s.MarkDone()
	if edges, ok := g.edges[n.Name]; ok {
		for _, e := range edges {
			t, err := e.TriggersMet(r)
			if err != nil {
				return nil, nil, fmt.Errorf("error checking edge condition: %w", err)
			}
			if t {
				nextBlocks = append(nextBlocks, e.To)
			}
		}
	}
	return nextBlocks, s, nil
}
