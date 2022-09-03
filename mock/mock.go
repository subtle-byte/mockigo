package mock

import (
	"fmt"
	"runtime"
	"sync"
)

type Testing interface {
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Helper()
	Cleanup(func())
	Failed() bool
}

type Mock struct {
	T             Testing
	mu            sync.Mutex
	expectedCalls *callSet
	getCaller     func(skip int) (pc uintptr, file string, line int, ok bool)
}

func NewMock(t Testing) *Mock {
	t.Helper()
	mock := &Mock{
		T:             t,
		expectedCalls: newCallSet(),
		getCaller:     runtime.Caller,
	}
	t.Cleanup(func() {
		mock.T.Helper()
		mock.finish()
	})
	return mock
}

func (mock *Mock) finish() {
	if mock.T.Failed() {
		return
	}
	mock.T.Helper()
	mock.mu.Lock()
	defer mock.mu.Unlock()
	unsatisfied := mock.expectedCalls.Unsatisfied()
	for _, call := range unsatisfied {
		mock.T.Errorf("clean up phase: missing call(s) to expected call %v", call.origin)
	}
}

func (mock *Mock) callerInfo(skip int) string {
	_, file, line, ok := mock.getCaller(skip + 1)
	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return "unknown file"
}

func (mock *Mock) ExpectCall(method string, args ...Matcher) *Call {
	mock.T.Helper()
	call := newCall(mock.T, method, mock.callerInfo(2), args...)
	mock.mu.Lock()
	defer mock.mu.Unlock()
	mock.expectedCalls.Add(call)
	return call
}

func (mock *Mock) Called(method string, args ...interface{}) Rets {
	mock.T.Helper()

	var call *Call
	func() {
		mock.T.Helper()
		mock.mu.Lock()
		defer mock.mu.Unlock()

		expected, err := mock.expectedCalls.FindMatch(method, args)
		if err != nil {
			mock.T.Fatalf("Unexpected call of method %q because:%s", method, err)
		}

		call = expected
	}()

	var rets = call.call(args)
	return Rets{
		rets: rets,
		call: call,
	}
}
