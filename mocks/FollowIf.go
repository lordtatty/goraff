// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import (
	goraff "github.com/lordtatty/goraff"
	mock "github.com/stretchr/testify/mock"
)

// FollowIf is an autogenerated mock type for the FollowIf type
type FollowIf struct {
	mock.Mock
}

type FollowIf_Expecter struct {
	mock *mock.Mock
}

func (_m *FollowIf) EXPECT() *FollowIf_Expecter {
	return &FollowIf_Expecter{mock: &_m.Mock}
}

// Match provides a mock function with given fields: s
func (_m *FollowIf) Match(s *goraff.StateReadOnly) bool {
	ret := _m.Called(s)

	if len(ret) == 0 {
		panic("no return value specified for Match")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*goraff.StateReadOnly) bool); ok {
		r0 = rf(s)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FollowIf_Match_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Match'
type FollowIf_Match_Call struct {
	*mock.Call
}

// Match is a helper method to define mock.On call
//   - s *goraff.StateReadOnly
func (_e *FollowIf_Expecter) Match(s interface{}) *FollowIf_Match_Call {
	return &FollowIf_Match_Call{Call: _e.mock.On("Match", s)}
}

func (_c *FollowIf_Match_Call) Run(run func(s *goraff.StateReadOnly)) *FollowIf_Match_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*goraff.StateReadOnly))
	})
	return _c
}

func (_c *FollowIf_Match_Call) Return(_a0 bool) *FollowIf_Match_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *FollowIf_Match_Call) RunAndReturn(run func(*goraff.StateReadOnly) bool) *FollowIf_Match_Call {
	_c.Call.Return(run)
	return _c
}

// NewFollowIf creates a new instance of FollowIf. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFollowIf(t interface {
	mock.TestingT
	Cleanup(func())
}) *FollowIf {
	mock := &FollowIf{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
