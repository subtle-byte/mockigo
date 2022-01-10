package mock

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type T struct {
	failed  bool
	cleanup func()
}

func (t *T) Errorf(format string, args ...any) {
	fmt.Printf("ERROR: "+format, args...)
}
func (t *T) FailNow() {
	panic("test-failed")
}
func (t *T) Fatalf(format string, args ...any) {
	t.failed = true
	fmt.Printf("FATAL: "+format, args...)
	panic("test-failed")
}
func (t *T) Helper()          {}
func (t *T) Cleanup(f func()) { t.cleanup = f }
func (t *T) Failed() bool     { return t.failed }

func runTest(f func(t *T)) {
	t := &T{}
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); !ok || msg != "test-failed" {
				panic(r)
			}
		}
		t.cleanup()
	}()
	f(t)
}

func newMock(t Testing) *Mock {
	m := NewMock(t)
	m.getCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		_, file, line, ok = runtime.Caller(skip)
		// Remove abs path prefix
		file = file[strings.Index(file, "mockigo"):]

		// Making line number relative to the test func start
		_, _, lineOfCaller, _ := runtime.Caller(skip + 2)
		line -= lineOfCaller
		line++

		return
	}
	return m
}

func Example_simple() {
	runTest(func(t *T) {
		m := newMock(t)
		m.ExpectCall("Foo", Eq("hello"))
		m.Called("Foo", "hello")
		m.Called("Foo", "bye")
	})
	//Output:
	// FATAL: Unexpected call of method "Foo" because:
	// expected call mockigo/mock/mock_test.go:3 doesn't match the argument "bye" at index 0
}

func TestMock_withReturn(t *testing.T) {
	m := NewMock(t)
	m.ExpectCall("Foo", Eq("hello")).Return(43)
	ret := m.Called("Foo", "hello")
	ret0 := ret.Get(0).(int)
	require.Equal(t, 43, ret0)
}

func TestMock_withRunReturn(t *testing.T) {
	m := NewMock(t)
	m.ExpectCall("Foo", Eq("hello")).RunReturn(func(s string) int {
		return 43
	})
	ret := m.Called("Foo", "hello")
	ret0 := ret.Get(0).(int)
	require.Equal(t, 43, ret0)
}

func Example_withTimes() {
	runTest(func(t *T) {
		m := newMock(t)
		m.ExpectCall("BarBar").Times(0, 1)
		m.ExpectCall("Bar").Times(1, -1)
		m.Called("Bar")
		m.Called("Bar")
		m.Called("Bar")
		m.ExpectCall("Foo", Eq("hello")).Times(1, -100)
		m.Called("Foo", "hello")
		m.Called("Foo", "hello")
	})
	//Output:
	// FATAL: Unexpected call of method "Foo" because:
	// expected call mockigo/mock/mock_test.go:8 has already been called the max number of times
}

func Example_withAfter() {
	runTest(func(t *T) {
		m := newMock(t)
		fooCall := m.ExpectCall("Foo")
		m.ExpectCall("Bar").After(1, -1, fooCall)
		m.Called("Bar")
	})
	//Output:
	// FATAL: Unexpected call of method "Bar" because:
	// expected call mockigo/mock/mock_test.go:4 should be called after call mockigo/mock/mock_test.go:3
}

func Example_withInOrder() {
	runTest(func(t *T) {
		m := newMock(t)
		InOrder(1, -1,
			m.ExpectCall("Foo"),
			m.ExpectCall("Bar"),
		)
		m.Called("Bar")
	})
	//Output:
	// FATAL: Unexpected call of method "Bar" because:
	// expected call mockigo/mock/mock_test.go:5 should be called after call mockigo/mock/mock_test.go:4
}

func TestMock_withMatchers(t *testing.T) {
	m := NewMock(t)
	m.ExpectCall("Foo", Any(), Eq("hello"), Not(Eq(0)))
	m.Called("Foo", 100, "hello", 1)
}

func TestMock_withMany(t *testing.T) {
	m := NewMock(t)
	InOrder(1, 1,
		m.ExpectCall("Foo", Any(), Eq("hello")).Return(45),
		m.ExpectCall("Bar", Any()).RunReturn(func(n int) string {
			return strconv.Itoa(n)
		}),
	)
	fooRet := m.Called("Foo", 100, "hello").Get(0).(int) // == 45
	barRet := m.Called("Bar", 200).Get(0).(string)       // == "200"

	require.Equal(t, 45, fooRet)
	require.Equal(t, "200", barRet)
}

func Example_unusedMethods() {
	runTest(func(t *T) {
		m := newMock(t)
		m.ExpectCall("Foo")
	})
	//Output:
	// ERROR: clean up phase: missing call(s) to expected call mockigo/mock/mock_test.go:3
}

func Example_noReturn() {
	runTest(func(t *T) {
		m := newMock(t)
		m.ExpectCall("Foo")
		rets := m.Called("Foo")
		rets.Get(1)
	})
	//Output:
	// FATAL: Call mockigo/mock/mock_test.go:3 does not have return value at index 1
}

func Example_errReturn() {
	runTest(func(t *T) {
		m := newMock(t)
		m.ExpectCall("Foo").Return(os.ErrClosed, nil, 7)
		rets := m.Called("Foo")
		e := rets.Error(0)
		require.Equal(t, os.ErrClosed, e)
		e = rets.Error(1)
		require.True(t, e == nil)
		rets.Error(2)
	})
	//Output:
	// FATAL: Call mockigo/mock/mock_test.go:3 does not have return value of the error type at index 2
}

func Test_unknownFile(t *testing.T) {
	m := NewMock(t)
	m.getCaller = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return 0, "", 0, false
	}
	info := m.callerInfo(1)
	require.Equal(t, "unknown file", info)
}

func Example_wrongNumberOfArguments() {
	runTest(func(t *T) {
		m := newMock(t)
		m.ExpectCall("Foo")
		m.Called("Foo", 3, 4)
	})
	//Output:
	// FATAL: Unexpected call of method "Foo" because:
	// expected call mockigo/mock/mock_test.go:3 expects 0 arguments, got 2
}

func Example_wrongNumberOfRunReturnArguments() {
	runTest(func(t *T) {
		m := newMock(t)
		m.ExpectCall("Foo").RunReturn(func(int) {})
	})
	//Output:
	// FATAL: Wrong number of arguments in RunReturn func for call mockigo/mock/mock_test.go:3: got 0, want 1
}

func Example_unexpectedMethod() {
	runTest(func(t *T) {
		m := newMock(t)
		m.Called("Foo")
	})
	//Output:
	// FATAL: Unexpected call of method "Foo" because: there are no expected calls of the method "Foo"
}
