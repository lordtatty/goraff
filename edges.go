package goraff

// Condition is a condition that must be met for an edge to be taken
type EdgeCondition interface {
	Match(s *StateReadOnly) bool
}

// Edge represents an edge in the graph
type Edge struct {
	From      *Node
	To        *Node
	Condition EdgeCondition
}

func (e *Edge) Match(s *StateReadOnly) bool {
	if e.Condition == nil {
		// without a conditon, we always follow the edge
		return true
	}
	return e.Condition.Match(s)
}

type edgeConditionKeyMatches struct {
	NodeID string
	Key    string
	Value  string
}

func (e *edgeConditionKeyMatches) Match(s *StateReadOnly) bool {
	return s.NodeState(e.NodeID).Get(e.Key) == e.Value
}

func FollowIfKeyMatches(nodeID, key, value string) EdgeCondition {
	return &edgeConditionKeyMatches{NodeID: nodeID, Key: key, Value: value}
}
