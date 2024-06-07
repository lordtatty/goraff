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
		graph := f.createSubGraph(output)
		s.AddSubGraph(graph)
		go f.runSubGraph(graph, &wg, errCh)
	}

	wg.Wait()
	close(errCh)

	if errs := f.collectErrors(errCh); len(errs) > 0 {
		return fmt.Errorf("errors running graph: %v", errs)
	}

	return f.buildFinalNodeResult(s)
}

func (f *FanOut) createSubGraph(output string) *goraff.Graph {
	graph := &goraff.Graph{}
	graph.NewNode(f.InKey, nil).AddStr("result", output)
	return graph
}

func (f *FanOut) runSubGraph(graph *goraff.Graph, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	if err := f.runScaff(graph); err != nil {
		errCh <- fmt.Errorf("error running graph: %s", err.Error())
	}
}

func (f *FanOut) collectErrors(errCh <-chan error) []error {
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return errs
}

func (f *FanOut) buildFinalNodeResult(s *goraff.Node) error {
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
	if err := f.Scaff.Go(g); err != nil {
		return fmt.Errorf("error running subgraph: %s", err.Error())
	}
	r := goraff.NewReadableGraph(g)
	nodeNames := r.NodeNames()
	fmt.Println(r.ID()+"  - Node Names: ", nodeNames)
	if g.FirstNodeByName(f.OutKey) == nil {
		return fmt.Errorf("could not find out node for node name: %s", f.OutKey)
	}
	return nil
}
