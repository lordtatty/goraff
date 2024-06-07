package blockactions

import (
	"fmt"
	"sync"

	"github.com/lordtatty/goraff"
)

type FanOut struct {
	Scaff   *goraff.Scaff
	InNode  string
	InKey   string
	OutNode string
	OutKey  string
}

// Do runs the FanOut process on the given node, graph, and previous node.
func (f *FanOut) Do(n *goraff.Node, r *goraff.ReadableGraph, prevNode *goraff.ReadableNode) error {
	fmt.Println("Running Scaff Node")
	if f.InKey == "" {
		f.InKey = "result"
	}
	if f.OutNode == "" {
		f.OutNode = "result"
	}
	if f.OutKey == "" {
		f.OutKey = "result"
	}

	results, err := f.getInputs(r, prevNode)
	if err != nil {
		return fmt.Errorf("error getting inputs: %s", err.Error())
	}

	errCh := make(chan error, len(results))
	var wg sync.WaitGroup
	wg.Add(len(results))

	for _, result := range results {
		subGraph := f.newSubGraph(result)
		n.AddSubGraph(subGraph)
		go f.processSubGraph(subGraph, &wg, errCh)
	}

	wg.Wait()
	close(errCh)

	if errs := f.collectErrors(errCh); len(errs) > 0 {
		return fmt.Errorf("errors running graph: %v", errs)
	}

	return f.combineResults(n)
}

func (f *FanOut) getInputs(r *goraff.ReadableGraph, prevNode *goraff.ReadableNode) ([][]byte, error) {
	if f.InNode != "" {
		n, err := r.FirstNodeByName(f.InNode)
		if err != nil {
			return nil, fmt.Errorf("could not find in node with name: %s", f.InNode)
		}
		return n.All(f.InKey), nil
	}
	if prevNode == nil {
		return nil, fmt.Errorf("no previous node provided, and no InNode specified")
	}
	return prevNode.All(f.InKey), nil
}

// newSubGraph initializes a new graph for the given result.
func (f *FanOut) newSubGraph(result []byte) *goraff.Graph {
	graph := &goraff.Graph{}
	graph.NewNode(f.InNode, nil).Add(f.InKey, result)
	return graph
}

// processSubGraph runs the sub-graph and handles any errors.
func (f *FanOut) processSubGraph(graph *goraff.Graph, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	if err := f.runScaff(graph); err != nil {
		errCh <- fmt.Errorf("error running graph: %s", err.Error())
	}
}

// collectErrors gathers errors from the error channel into a slice.
func (f *FanOut) collectErrors(errCh <-chan error) []error {
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return errs
}

// combineResults combines the results from all sub-graphs into the main node.
func (f *FanOut) combineResults(n *goraff.Node) error {
	for _, subGraph := range n.Get().SubGraph() {
		outNode, err := subGraph.FirstNodeByName(f.OutNode)
		if err != nil {
			return fmt.Errorf("could not find out node with name: %s", f.OutNode)
		}
		for _, result := range outNode.AllStr(f.OutKey) {
			n.AddStr("result", result)
		}
	}
	return nil
}

// runScaff runs the scaffolding process on the provided graph.
func (f *FanOut) runScaff(g *goraff.Graph) error {
	if err := f.Scaff.Go(g); err != nil {
		return fmt.Errorf("error running subgraph: %s", err.Error())
	}
	r := goraff.NewReadableGraph(g)
	nodeNames := r.NodeNames()
	fmt.Println(r.ID()+"  - Node Names: ", nodeNames)
	if g.FirstNodeByName(f.OutNode) == nil {
		return fmt.Errorf("could not find out node with name: %s", f.OutNode)
	}
	return nil
}
