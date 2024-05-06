// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// LLMClient is an autogenerated mock type for the LLMClient type
type LLMClient struct {
	mock.Mock
}

type LLMClient_Expecter struct {
	mock *mock.Mock
}

func (_m *LLMClient) EXPECT() *LLMClient_Expecter {
	return &LLMClient_Expecter{mock: &_m.Mock}
}

// Chat provides a mock function with given fields: systemMsg, userMsg, stream
func (_m *LLMClient) Chat(systemMsg string, userMsg string, stream chan string) (string, error) {
	ret := _m.Called(systemMsg, userMsg, stream)

	if len(ret) == 0 {
		panic("no return value specified for Chat")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, chan string) (string, error)); ok {
		return rf(systemMsg, userMsg, stream)
	}
	if rf, ok := ret.Get(0).(func(string, string, chan string) string); ok {
		r0 = rf(systemMsg, userMsg, stream)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string, chan string) error); ok {
		r1 = rf(systemMsg, userMsg, stream)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LLMClient_Chat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Chat'
type LLMClient_Chat_Call struct {
	*mock.Call
}

// Chat is a helper method to define mock.On call
//   - systemMsg string
//   - userMsg string
//   - stream chan string
func (_e *LLMClient_Expecter) Chat(systemMsg interface{}, userMsg interface{}, stream interface{}) *LLMClient_Chat_Call {
	return &LLMClient_Chat_Call{Call: _e.mock.On("Chat", systemMsg, userMsg, stream)}
}

func (_c *LLMClient_Chat_Call) Run(run func(systemMsg string, userMsg string, stream chan string)) *LLMClient_Chat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(chan string))
	})
	return _c
}

func (_c *LLMClient_Chat_Call) Return(_a0 string, _a1 error) *LLMClient_Chat_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *LLMClient_Chat_Call) RunAndReturn(run func(string, string, chan string) (string, error)) *LLMClient_Chat_Call {
	_c.Call.Return(run)
	return _c
}

// NewLLMClient creates a new instance of LLMClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLLMClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *LLMClient {
	mock := &LLMClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}