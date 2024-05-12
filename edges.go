package goraff

// Condition is a condition that must be met for an edge to be taken
type FollowIf interface {
	Match(s *GraphStateReader) bool
}

// Edge represents an edge in the graph
type Edge struct {
	From      *Node
	To        *Node
	Condition FollowIf
}

func (e *Edge) TriggersMet(s *GraphStateReader) bool {
	if e.Condition == nil {
		// without a conditon, we always follow the edge
		return true
	}
	return e.Condition.Match(s)
}

type followIfKeyMatchesName struct {
	Name  string
	Key   string
	Value string
}

func (e *followIfKeyMatchesName) Match(s *GraphStateReader) bool {
	return s.state.FirstNodeStateByName(e.Name).Reader().GetStr(e.Key) == e.Value
}

func FollowIfKeyMatches(nodeID, key, value string) FollowIf {
	return &followIfKeyMatchesName{Name: nodeID, Key: key, Value: value}
}

type followIfNodesCompleted struct {
	NodeIDs []string
}

func (e *followIfNodesCompleted) Match(s *GraphStateReader) bool {
	for _, nodeID := range e.NodeIDs {
		st := s.FirstNodeStateByName(nodeID)
		if st == nil {
			return false
		}
		if !st.Done() {
			return false
		}
	}
	return true
}

func FollowIfNodesCompleted(nodeIDs ...string) FollowIf {
	return &followIfNodesCompleted{NodeIDs: nodeIDs}
}
