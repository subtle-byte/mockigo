package mock

import (
	"fmt"
	"runtime"
	"strings"
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
	getCaller := mock.getCaller
	if getCaller == nil {
		getCaller = runtime.Caller
	}
	_, file, line, ok := getCaller(skip + 1)
	if ok && strings.HasSuffix(file, "libexec/src/testing/testing.go") {
		_, file, line, ok = getCaller(skip)
	}
	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return "unknown file"
}

func (mock *Mock) ExpectCall(method string, args ...Matcher) *Call {
	mock.T.Helper()
	call := newCall(mock, method, mock.callerInfo(2), args...)
	mock.mu.Lock()
	defer mock.mu.Unlock()
	mock.expectedCalls.Add(call)
	return call
}

type Rets struct {
	rets []interface{}
	call *Call
}

func (r Rets) Len() int {
	return len(r.rets)
}

func (r Rets) Get(i int) interface{} {
	if i >= r.Len() {
		r.call.mock.T.Fatalf("Call %v does not have return value at index %v", r.call.origin, i)
	}
	return r.rets[i]
}

func (r Rets) Error(i int) error {
	ret := r.Get(i)
	if ret == nil {
		return nil
	}
	e, ok := ret.(error)
	if !ok {
		r.call.mock.T.Fatalf("Call %v does not have return value of the error type at index %v", r.call.origin, i)
	}
	return e
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

		expected.actualCallsNum++
		call = expected
	}()

	var rets = call.action(args)
	return Rets{
		rets: rets,
		call: call,
	}
}
