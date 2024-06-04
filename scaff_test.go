package goraff_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/lordtatty/goraff"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	g := goraff.NewScaff()
	// must be an empty goraff.Graph{}
	assert.Equal(goraff.Scaff{}, *g)
}

type actionMock struct {
	name        string
	lastName    string
	expectNoRun bool
	t           *testing.T
	delay       time.Duration
	err         error
}

func (a *actionMock) Do(s *goraff.Node, r *goraff.ReadableGraph, triggeringNS *goraff.ReadableNode) error {
	if a.expectNoRun {
		a.t.Error("Action should not have run")
	}
	// Wait if a delay is set
	if a.delay > 0 {
		time.Sleep(a.delay)
	}
	if a.err != nil {
		return a.err
	}
	// Set the key to the name of the action
	key := fmt.Sprintf("%s_key", a.name)
	if a.lastName == "" {
		s.SetStr(key, a.name)
		return nil
	}
	lastKey := fmt.Sprintf("%s_key", a.lastName)
	lastVal := triggeringNS.GetStr(lastKey)
	// split string on " :: " and take the last element
	parts := strings.Split(lastVal, " :: ")
	lastVal = parts[len(parts)-1]
	val := fmt.Sprintf("%s :: %s", lastVal, a.name)
	s.SetStr(key, val)
	return nil
}

func TestGraph_Go_NoEdges(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Scaff{}
	a1 := &actionMock{name: "action1"}
	a2 := &actionMock{name: "action2", lastName: "action1", expectNoRun: true, t: t}

	n1 := g.AddBlock("action1", a1)
	_ = g.AddBlock("action2", a2)
	g.SetEntrypoint(n1)
	graph := &goraff.Graph{}
	err := g.Go(graph)
	assert.NoError(err)

	// Should only be one state for action1, as it should only have run once, and the key should be set to the action name
	states := graph.NodeByName("action1")
	assert.Len(states, 1)
	assert.Equal("action1", states[0].Reader().GetStr("action1_key"))
	// action2 should not have run
	assert.Len(graph.NodeByName("action2"), 0)
}

func TestGraph_NodeHasError(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Scaff{}
	a1 := &actionMock{name: "action1"}
	a2 := &actionMock{name: "action2", lastName: "action1", err: fmt.Errorf("error"), t: t}

	n1 := g.AddBlock("action1", a1)
	n2 := g.AddBlock("action2", a2)

	g.AddJoin(n1, n2, nil)

	g.SetEntrypoint(n1)
	graph := &goraff.Graph{}
	err := g.Go(graph)
	assert.Error(err)
	assert.Equal("error running block: error", err.Error())
}

func TestGraph_Go_WithEdges(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Scaff{}

	a1 := &actionMock{name: "action1"}
	n1 := g.AddBlock("action1", a1)
	a2 := &actionMock{name: "action2", lastName: "action1"}
	n2 := g.AddBlock("action2", a2)
	a3 := &actionMock{name: "action3", lastName: "action2"}
	n3 := g.AddBlock("action3", a3)
	a4 := &actionMock{name: "action4", expectNoRun: true, t: t}
	g.AddBlock("action4", a4) // thi should not run

	g.SetEntrypoint(n1)
	// with no condition we always follow the join
	g.AddJoin(n1, n2, nil)
	g.AddJoin(n2, n3, nil)
	// No join from n3, so it should stop after n3
	graph := &goraff.Graph{}
	err := g.Go(graph)
	assert.NoError(err)

	assert.Len(graph.NodeByName("action1"), 1)
	assert.Equal("action1", graph.NodeByName("action1")[0].Reader().GetStr("action1_key"))
	assert.Len(graph.NodeByName("action2"), 1)
	assert.Equal("action1 :: action2", graph.NodeByName("action2")[0].Reader().GetStr("action2_key"))
	assert.Len(graph.NodeByName("action3"), 1)
	assert.Equal("action2 :: action3", graph.NodeByName("action3")[0].Reader().GetStr("action3_key"))
	assert.Len(graph.NodeByName("action4"), 0) // Action 4 should not have run
}

func TestGraph_ConditionalEdges(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Scaff{}

	a1 := &actionMock{name: "action1"}
	n1 := g.AddBlock("action1", a1)
	a2 := &actionMock{name: "action2", lastName: "action1", expectNoRun: true, t: t}
	n2 := g.AddBlock("action2", a2)
	a3 := &actionMock{name: "action3", lastName: "action1"}
	n3 := g.AddBlock("action3", a3)

	g.SetEntrypoint(n1)
	// Both n2 and n3 should follow n1, but only n3 should match the condition
	g.AddJoin(n1, n2, goraff.FollowIfKeyMatches(n1, "action1_key", "should not match"))
	g.AddJoin(n1, n3, goraff.FollowIfKeyMatches(n1, "action1_key", "action1"))

	graph := &goraff.Graph{}
	err := g.Go(graph)
	assert.NoError(err)

	assert.Equal("action1", graph.FirstNodeByName(n1).Reader().GetStr("action1_key"))
	assert.Nil(graph.NodeByID(n2)) // Action 2 should not have run
	assert.Equal("action1 :: action3", graph.FirstNodeByName(n3).Reader().GetStr("action3_key"))
}

func TestGraph_AddEdge_Node1NotFound(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Scaff{}
	err := g.AddJoin("node1", "node2", nil)
	assert.Error(err)
	assert.Equal("block not found: node1", err.Error())
}

func TestGraph_AddEdge_Node2NotFound(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Scaff{}
	a1 := &actionMock{name: "action1"}
	n1 := g.AddBlock("action1", a1)
	err := g.AddJoin(n1, "node2", nil)
	assert.Error(err)
	assert.Equal("block not found: node2", err.Error())
}

func TestGraph_FanOutNodes_Parallel(t *testing.T) {
	// In this test we are checking tha we can fan out from a node
	// and, importantly, that the actions run in parallel
	// We will check parallelisation by delaying each action by a second.
	// The first runs immediately, the next three should run in parallel
	// and making sure the whole graph completes in around 2 seconds (not 4)
	assert := assert.New(t)
	g := &goraff.Scaff{}

	a1 := &actionMock{name: "action1", delay: 1 * time.Second}
	n1 := g.AddBlock("action1", a1)
	a2 := &actionMock{name: "action2", lastName: "action1", delay: 1 * time.Second}
	n2 := g.AddBlock("action2", a2)
	a3 := &actionMock{name: "action3", lastName: "action1", delay: 1 * time.Second}
	n3 := g.AddBlock("action3", a3)
	a4 := &actionMock{name: "action4", lastName: "action1", delay: 1 * time.Second}
	n4 := g.AddBlock("action4", a4)

	g.SetEntrypoint(n1)
	g.AddJoin(n1, n2, nil)
	g.AddJoin(n1, n3, nil)
	g.AddJoin(n1, n4, nil)

	start := time.Now()

	graph := &goraff.Graph{}
	err := g.Go(graph)
	assert.NoError(err)

	elapsed := time.Since(start)
	assert.True(elapsed < 2500*time.Millisecond, "Elapsed time should be less than 2.5 seconds (first node, parallel nodes, and a bit of leeway)")

	assert.Equal("action1", graph.FirstNodeByName(n1).Reader().GetStr("action1_key"))
	assert.Equal("action1 :: action2", graph.FirstNodeByName(n2).Reader().GetStr("action2_key"))
	assert.Equal("action1 :: action3", graph.FirstNodeByName(n3).Reader().GetStr("action3_key"))
	assert.Equal("action1 :: action4", graph.FirstNodeByName(n4).Reader().GetStr("action4_key"))
}

type mockFollowIfWantsDone struct {
	nodeIDs []string
	t       *testing.T
}

func (f *mockFollowIfWantsDone) Match(s *goraff.ReadableGraph) (bool, error) {
	assert := assert.New(f.t)
	for _, nodeID := range f.nodeIDs {
		st, err := s.FirstNodeByName(nodeID)
		assert.Nil(err)
		d := st.Done()
		fmt.Println(d)
		assert.NotNil(st)
		assert.True(st.Done())
	}
	return true, nil
}

func TestGraph_StateIsMarkedDoneBeforeTriggers(t *testing.T) {
	// The state should be marked done before the triggers are checked
	// Because some triggers may rely on the state being done
	assert := assert.New(t)
	g := &goraff.Scaff{}

	a1 := &actionMock{name: "action1"}
	n1 := g.AddBlock("action1", a1)
	a2 := &actionMock{name: "action2", lastName: "action1"}
	n2 := g.AddBlock("action2", a2)
	a3 := &actionMock{name: "action3", lastName: "action2"}
	n3 := g.AddBlock("action3", a3)

	g.SetEntrypoint(n1)
	g.AddJoin(n1, n2, nil)
	followIf := &mockFollowIfWantsDone{nodeIDs: []string{n2}, t: t}
	g.AddJoin(n2, n3, followIf)

	graph := &goraff.Graph{}
	err := g.Go(graph)
	assert.NoError(err)

	assert.Equal("action1", graph.FirstNodeByName(n1).Reader().GetStr("action1_key"))
	assert.Equal("action1 :: action2", graph.FirstNodeByName(n2).Reader().GetStr("action2_key"))
	assert.Equal("action2 :: action3", graph.FirstNodeByName(n3).Reader().GetStr("action3_key"))
}

func TestGraph_EntrypointNotSet(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Scaff{}

	graph := &goraff.Graph{}
	err := g.Go(graph)

	assert.Error(err)
	assert.Equal("error validating graph: entrypoint not set", err.Error())
}

type actionMockCheckReader struct {
	expectNilReader bool
	t               *testing.T
}

func (a *actionMockCheckReader) Do(s *goraff.Node, r *goraff.ReadableGraph, triggeringNS *goraff.ReadableNode) error {
	if a.expectNilReader {
		if triggeringNS != nil {
			a.t.Error("Expected nil reader but got a non-nil reader")
		}
		s.SetStr("check_reader_key", "reader is nil")
	} else {
		if triggeringNS == nil {
			a.t.Error("Expected non-nil reader but got a nil reader")
		}
		s.SetStr("check_reader_key", "reader is not nil")
	}
	return nil
}

func TestGraph_FlowMgr_ReaderPassing(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Scaff{}

	// Define an action mock that will check the triggeringNS reader for nil
	checkReaderAction1 := &actionMockCheckReader{
		expectNilReader: true,
		t:               t,
	}
	n1 := g.AddBlock("action1", checkReaderAction1)

	// Define another action mock that will be triggered by the first and expects a non-nil reader
	checkReaderAction2 := &actionMockCheckReader{
		expectNilReader: false,
		t:               t,
	}
	n2 := g.AddBlock("action2", checkReaderAction2)

	g.SetEntrypoint(n1)
	g.AddJoin(n1, n2, nil)

	graph := &goraff.Graph{}
	err := g.Go(graph)
	assert.NoError(err)

	assert.Equal("reader is nil", graph.FirstNodeByName(n1).Reader().GetStr("check_reader_key"))
	assert.Equal("reader is not nil", graph.FirstNodeByName(n2).Reader().GetStr("check_reader_key"))
}

func TestGraph_Go_ValidateEntrypoint(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Scaff{}

	a1 := &actionMock{name: "action1"}
	g.AddBlock("action1", a1)
	a2 := &actionMock{name: "action2", lastName: "action1"}
	g.AddBlock("action2", a2)

	// No entrypoint set
	graph := &goraff.Graph{}
	err := g.Go(graph)
	assert.Error(err)
	assert.Equal("error validating graph: entrypoint not set", err.Error())
}

func TestGraph_Go_ValidateUniqueBlockNames(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Scaff{}

	a1 := &actionMock{name: "action1"}
	n1 := g.AddBlock("action1", a1)
	a2 := &actionMock{name: "action1"}
	g.AddBlock("action1", a2)

	g.SetEntrypoint(n1)
	graph := &goraff.Graph{}
	err := g.Go(graph)
	assert.Error(err)
	assert.Equal("error validating graph: block name not unique: action1", err.Error())
}
