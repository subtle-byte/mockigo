package mock

import (
	"fmt"
	"runtime"
	"strings"
)

// t is implementation of Testing for the test cases where fail is expected
type t struct {
	failed                bool
	cleanupFunc           func()
	cleaningUp            bool
	cleanupSetterFuncName string
	helperFuncNames       map[string]bool
	file                  string
	line                  int
}

var _ Testing = (*t)(nil)

func runWithT(f func(t *t)) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("cannot get caller")
	}
	t := &t{
		helperFuncNames: make(map[string]bool),
		file:            file,
		line:            line,
	}
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); !ok || msg != "test-failed" {
				panic(r)
			}
		}
		t.cleaningUp = true
		t.cleanupFunc()
	}()
	f(t)
}

func newMockWithT(t *t) *Mock {
	t.Helper()
	m := NewMock(t)
	m.getCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return t.caller(skip)
	}
	return m
}

func (t *t) callerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		panic("cannot get caller")
	}
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	return frame.Function
}
func (t *t) caller(skip int) (pc uintptr, file string, line int, ok bool) {
	pc, file, line, ok = runtime.Caller(skip + 1)
	if !ok {
		panic("cannot get caller")
	}

	if file == t.file {
		file = "runWithT"
		line -= t.line // to the line number be relative to the runWithT line
	} else {
		// Remove abs path prefix
		file = file[strings.LastIndex(file, "mockigo"):]
	}

	return
}
func (t *t) notHelperCallerName(skip int) string {
	if t.cleaningUp {
		return t.cleanupSetterFuncName
	}

	for t.helperFuncNames[t.callerName(skip+1)] {
		skip++
	}
	_, file, line, _ := t.caller(skip + 1)
	return fmt.Sprintf("%s:%d", file, line)
}

// FailNow is used by require.Equal
func (t *t) FailNow() {
	t.failed = true
	panic("test-failed")
}
func (t *t) Errorf(format string, args ...any) {
	t.failed = true
	fmt.Print(t.notHelperCallerName(1)+": ", fmt.Sprintf(format, args...))
}
func (t *t) Fatalf(format string, args ...any) {
	t.Helper()
	t.Errorf(format, args...)
	t.FailNow()
}
func (t *t) Helper() {
	t.helperFuncNames[t.callerName(1)] = true
}
func (t *t) Cleanup(f func()) {
	t.cleanupFunc = f
	t.cleanupSetterFuncName = t.notHelperCallerName(1)
}
func (t *t) Failed() bool { return t.failed }
