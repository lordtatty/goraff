package blockactions

import (
	"fmt"
	"sync"

	"github.com/lordtatty/goraff"
)

type FanOut struct {
	Scaff  *goraff.Scaff
	InKey  string
	OutKey string
}

func (f *FanOut) Do(s *goraff.Node, r *goraff.ReadableGraph, previousNode *goraff.ReadableNode) error {
	fmt.Println("Running Scaff Node")
	outputs := previousNode.AllStr("result")
	if f.OutKey == "" {
		f.OutKey = "result"
	}
	errCh := make(chan error, len(outputs))
	wg := sync.WaitGroup{}
	wg.Add(len(outputs))
	for _, output := range outputs {
		graph := &goraff.Graph{}
		graph.NewNode(f.InKey, nil).AddStr("result", output)
		s.AddSubGraph(graph)
		go func(g *goraff.Graph) {
			defer wg.Done()
			err := f.runScaff(g)
			if err != nil {
				errCh <- fmt.Errorf("error running graph: %s", err.Error())
				return
			}
		}(graph)
	}
	wg.Wait()
	close(errCh)
	// check for errors
	errs := []error{}
	for err := range errCh {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors running graph: %v", errs)
	}
	// Build final node result
	for _, r := range s.Get().SubGraph() {
		resultNode, err := r.FirstNodeByName(f.OutKey)
		if err != nil {
			return fmt.Errorf("could not find out node for node name: %s", f.OutKey)
		}
		for _, result := range resultNode.AllStr("result") {
			s.AddStr("result", result)
		}
	}
	return nil
}

func (f *FanOut) runScaff(g *goraff.Graph) error {
	err := f.Scaff.Go(g)
	if err != nil {
		return fmt.Errorf("error running subgraph: %s", err.Error())
	}
	r := goraff.NewReadableGraph(g)
	nodeNames := r.NodeNames()
	fmt.Println(r.ID()+"  - Node Names: ", nodeNames)
	outNode := g.FirstNodeByName(f.OutKey)
	if outNode == nil {
		return fmt.Errorf("could not find out node for node name: %s", f.OutKey)
	}
	return nil
}
