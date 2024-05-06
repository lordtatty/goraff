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
	g := goraff.New()
	// must be an empty goraff.Graph{}
	assert.Equal(goraff.Graph{}, *g)
}

func TestNode_ID(t *testing.T) {
	assert := assert.New(t)
	n := goraff.Node{}
	// The node should have an ID which looks like a UUID
	assert.Regexp("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", n.ID())
	// Assert that the ID is unique
	n2 := goraff.Node{}
	assert.NotEqual(n.ID(), n2.ID())
}

func TestGraph_AddNode(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Graph{}
	a := &actionMock{name: "action1"}
	g.AddNode(a)
	// The graph should have one node
	assert.Equal(1, g.Len())
}

type actionMock struct {
	name        string
	lastName    string
	expectNoRun bool
	t           *testing.T
	delay       time.Duration
	err         error
}

func (a *actionMock) Do(s *goraff.NodeState, r *goraff.StateReadOnly, triggeringNodeID string) error {
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
		s.Set(key, a.name)
		return nil
	}
	ls := r.NodeState(triggeringNodeID)
	lastKey := fmt.Sprintf("%s_key", a.lastName)
	lastVal := ls.Get(lastKey)
	// split string on " :: " and take the last element
	parts := strings.Split(lastVal, " :: ")
	lastVal = parts[len(parts)-1]
	val := fmt.Sprintf("%s :: %s", lastVal, a.name)
	s.Set(key, val)
	return nil
}

func TestGraph_Go_NoEdges(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Graph{}
	a1 := &actionMock{name: "action1"}
	a2 := &actionMock{name: "action2", lastName: "action1", expectNoRun: true, t: t}

	n1 := g.AddNode(a1)
	n2 := g.AddNode(a2)
	g.SetEntrypoint(n1)
	g.Go()

	state := g.State()
	assert.Equal("action1", state.NodeState(n1).Get("action1_key"))
	assert.Nil(state.NodeState(n2)) // Action 2 should not have run
}

func TestGraph_NodeHasError(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Graph{}
	a1 := &actionMock{name: "action1"}
	a2 := &actionMock{name: "action2", lastName: "action1", err: fmt.Errorf("error"), t: t}

	n1 := g.AddNode(a1)
	n2 := g.AddNode(a2)

	g.AddEdge(n1, n2, nil)

	g.SetEntrypoint(n1)
	err := g.Go()
	assert.Error(err)
	assert.Equal("error running node: error", err.Error())
}

func TestGraph_Go_WithEdges(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Graph{}

	a1 := &actionMock{name: "action1"}
	n1 := g.AddNode(a1)
	a2 := &actionMock{name: "action2", lastName: "action1"}
	n2 := g.AddNode(a2)
	a3 := &actionMock{name: "action3", lastName: "action2"}
	n3 := g.AddNode(a3)
	a4 := &actionMock{name: "action4", expectNoRun: true, t: t}
	g.AddNode(a4) // thi should not run

	g.SetEntrypoint(n1)
	// with no condition we always follow the edge
	g.AddEdge(n1, n2, nil)
	g.AddEdge(n2, n3, nil)
	// No edge from n3, so it should stop after n3
	g.Go()

	state := g.State()
	assert.Equal("action1", state.NodeState(n1).Get("action1_key"))
	assert.Equal("action1 :: action2", state.NodeState(n2).Get("action2_key"))
	assert.Equal("action2 :: action3", state.NodeState(n3).Get("action3_key"))
	assert.Equal("", state.NodeState(n2).Get("action4_key")) // Action 4 should not have run
}

func TestGraph_ConditionalEdges(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Graph{}

	a1 := &actionMock{name: "action1"}
	n1 := g.AddNode(a1)
	a2 := &actionMock{name: "action2", lastName: "action1", expectNoRun: true, t: t}
	n2 := g.AddNode(a2)
	a3 := &actionMock{name: "action3", lastName: "action1"}
	n3 := g.AddNode(a3)

	g.SetEntrypoint(n1)
	// Both n2 and n3 should follow n1, but only n3 should match the condition
	g.AddEdge(n1, n2, goraff.FollowIfKeyMatches(n1, "action1_key", "should not match"))
	g.AddEdge(n1, n3, goraff.FollowIfKeyMatches(n1, "action1_key", "action1"))

	g.Go()

	state := g.State()
	assert.Equal("action1", state.NodeState(n1).Get("action1_key"))
	assert.Nil(state.NodeState(n2)) // Action 2 should not have run
	assert.Equal("action1 :: action3", state.NodeState(n3).Get("action3_key"))
}

func TestGraph_AddEdge_Node1NotFound(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Graph{}
	err := g.AddEdge("node1", "node2", nil)
	assert.Error(err)
	assert.Equal("node not found: node1", err.Error())
}

func TestGraph_AddEdge_Node2NotFound(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Graph{}
	a1 := &actionMock{name: "action1"}
	n1 := g.AddNode(a1)
	err := g.AddEdge(n1, "node2", nil)
	assert.Error(err)
	assert.Equal("node not found: node2", err.Error())
}

func TestGraph_FanOutNodes_Parallel(t *testing.T) {
	// In this test we are checking tha we can fan out from a node
	// and, importantly, that the actions run in parallel
	// We will check parallelisation by delaying each action by a second.
	// The first runs immediately, the next three should run in parallel
	// and making sure the whole graph completes in around 2 seconds (not 4)
	assert := assert.New(t)
	g := &goraff.Graph{}

	a1 := &actionMock{name: "action1", delay: 1 * time.Second}
	n1 := g.AddNode(a1)
	a2 := &actionMock{name: "action2", lastName: "action1", delay: 1 * time.Second}
	n2 := g.AddNode(a2)
	a3 := &actionMock{name: "action3", lastName: "action1", delay: 1 * time.Second}
	n3 := g.AddNode(a3)
	a4 := &actionMock{name: "action4", lastName: "action1", delay: 1 * time.Second}
	n4 := g.AddNode(a4)

	g.SetEntrypoint(n1)
	g.AddEdge(n1, n2, nil)
	g.AddEdge(n1, n3, nil)
	g.AddEdge(n1, n4, nil)

	start := time.Now()
	g.Go()
	elapsed := time.Since(start)
	assert.True(elapsed < 2500*time.Millisecond, "Elapsed time should be less than 2.5 seconds (first node, parallel nodes, and a bit of leeway)")

	state := g.State()
	assert.Equal("action1", state.NodeState(n1).Get("action1_key"))
	assert.Equal("action1 :: action2", state.NodeState(n2).Get("action2_key"))
	assert.Equal("action1 :: action3", state.NodeState(n3).Get("action3_key"))
	assert.Equal("action1 :: action4", state.NodeState(n4).Get("action4_key"))
}

type mockFollowIfWantsDone struct {
	nodeIDs []string
	t       *testing.T
}

func (f *mockFollowIfWantsDone) Match(s *goraff.StateReadOnly) bool {
	assert := assert.New(f.t)
	for _, nodeID := range f.nodeIDs {
		st := s.NodeState(nodeID)
		d := st.Done()
		fmt.Println(d)
		assert.NotNil(st)
		assert.True(st.Done())
	}
	return true
}

func TestGraph_StateIsMarkedDoneBeforeTriggers(t *testing.T) {
	// The state should be marked done before the triggers are checked
	// Because some triggers may rely on the state being done
	assert := assert.New(t)
	g := &goraff.Graph{}

	a1 := &actionMock{name: "action1"}
	n1 := g.AddNode(a1)
	a2 := &actionMock{name: "action2", lastName: "action1"}
	n2 := g.AddNode(a2)
	a3 := &actionMock{name: "action3", lastName: "action2"}
	n3 := g.AddNode(a3)

	g.SetEntrypoint(n1)
	g.AddEdge(n1, n2, nil)
	followIf := &mockFollowIfWantsDone{nodeIDs: []string{n2}, t: t}
	g.AddEdge(n2, n3, followIf)

	g.Go()

	state := g.State()
	assert.Equal("action1", state.NodeState(n1).Get("action1_key"))
	assert.Equal("action1 :: action2", state.NodeState(n2).Get("action2_key"))
	assert.Equal("action2 :: action3", state.NodeState(n3).Get("action3_key"))
}

func TestGraph_EntrypointNotSet(t *testing.T) {
	assert := assert.New(t)
	g := &goraff.Graph{}
	err := g.Go()
	assert.Error(err)
	assert.Equal("entrypoint not set", err.Error())
}
