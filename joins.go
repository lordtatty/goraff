package goraff

import "fmt"

// Condition is a condition that must be met for a join to be taken
type FollowIf interface {
	Match(s *ReadableGraph) (bool, error)
}

// Manage joins
type Joins struct {
	joins  map[string][]*Join
	Blocks *Blocks
	errs   []error
}

func (j *Joins) trackErr(err error) error {
	j.errs = append(j.errs, err)
	return err
}

func (j *Joins) Validate() error {
	if j.errs != nil {
		return fmt.Errorf("joins have errors: %v", j.errs)
	}
	return nil
}

type ErrBlockNotFound struct {
	ID string
}

func (e ErrBlockNotFound) Error() string {
	return "block not found: " + e.ID
}

func (j *Joins) Add(fromName, toName string, condition FollowIf) error {
	if j.Blocks == nil {
		return fmt.Errorf("joins must be associated with a Blocks struct")
	}
	from := j.Blocks.Get(fromName)
	if from == nil {
		return j.trackErr(ErrBlockNotFound{
			ID: fromName,
		})
	}
	to := j.Blocks.Get(toName)
	if to == nil {
		return j.trackErr(ErrBlockNotFound{
			ID: toName,
		})
	}
	if j.joins == nil {
		j.joins = make(map[string][]*Join)
	}
	e := &Join{From: from, To: to, Condition: condition}
	if _, ok := j.joins[fromName]; !ok {
		j.joins[fromName] = []*Join{}
	}
	j.joins[fromName] = append(j.joins[fromName], e)
	return nil
}

func (j *Joins) Get(from string) ([]*Join, bool) {
	if _, ok := j.joins[from]; !ok {
		return nil, false
	}
	return j.joins[from], true
}

// Join connects two blocks in a scaff
type Join struct {
	From      *Block
	To        *Block
	Condition FollowIf
}

func (e *Join) TriggersMet(s *ReadableGraph) (bool, error) {
	if e.Condition == nil {
		// without a conditon, we always trigger the join
		return true, nil
	}
	return e.Condition.Match(s)
}

type followIfKeyMatchesName struct {
	Name  string
	Key   string
	Value string
}

func (e *followIfKeyMatchesName) Match(s *ReadableGraph) (bool, error) {
	n, err := s.FirstNodeByName(e.Name)
	if err != nil {
		return false, fmt.Errorf("error getting node state: %w", err)
	}
	return n.GetStr(e.Key) == e.Value, nil
}

func FollowIfKeyMatches(nodeID, key, value string) FollowIf {
	return &followIfKeyMatchesName{Name: nodeID, Key: key, Value: value}
}

type followIfNodesCompleted struct {
	NodeIDs []string
}

func (e *followIfNodesCompleted) Match(s *ReadableGraph) (bool, error) {
	for _, nodeID := range e.NodeIDs {
		st, err := s.FirstNodeByName(nodeID)
		if err != nil {
			return false, fmt.Errorf("error getting node state: %w", err)
		}
		if st == nil {
			return false, nil
		}
		if !st.Done() {
			return false, nil
		}
	}
	return true, nil
}

func FollowIfNodesCompleted(nodeIDs ...string) FollowIf {
	return &followIfNodesCompleted{NodeIDs: nodeIDs}
}
