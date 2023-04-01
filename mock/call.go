package mock

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/k0kubun/pp/v3"
)

type Call struct {
	t                  Testing
	method             string
	argsMatchers       []Matcher
	prevCalls          []*Call
	nextCalled         []*Call
	minCalls, maxCalls int
	actualCallsNum     int
	action             func([]interface{}) []interface{}
	origin             string
}

// CallHolder is implemented by Call and also by any types that embed Call (e.g. by generated typed Call)
type CallHolder interface {
	GetCall() *Call
}

var _ CallHolder = (*Call)(nil)

// GetCall implements CallHolder interface, just returns the receiver c
func (c *Call) GetCall() *Call {
	return c
}

const infCalls = 1e8 // close enough to infinity

func newCall(t Testing, method string, origin string, argsMatchers ...Matcher) *Call {
	t.Helper()
	return &Call{
		t:              t,
		method:         method,
		argsMatchers:   argsMatchers,
		prevCalls:      nil,
		minCalls:       1,
		maxCalls:       infCalls,
		actualCallsNum: 0,
		action:         func(i []interface{}) []interface{} { return nil },
		origin:         origin,
	}
}

// Times set how many times the call should be called.
// Use -1 for max to identify infinity.
func (c *Call) Times(min, max int) *Call {
	if max == -1 {
		max = infCalls
	}
	c.minCalls, c.maxCalls = min, max
	return c
}

func (c *Call) CalledTimes() int {
	return c.actualCallsNum
}

// Returns true if the minimum number of calls have been made.
func (c *Call) satisfied() bool {
	return c.actualCallsNum >= c.minCalls
}

func (c *Call) callingNext(next *Call) {
	c.nextCalled = append(c.nextCalled, next)
}

func (c *Call) call(args []interface{}) []interface{} {
	c.actualCallsNum++
	for _, prevCall := range c.prevCalls {
		prevCall.callingNext(c)
	}
	return c.action(args)
}

// After sets that the call c can be called only after the prevCall is called its min times.
//
// After the prevCall is called the call c cannot be called.
func (c *Call) After(prevCall CallHolder) *Call {
	c.prevCalls = append(c.prevCalls, prevCall.GetCall())
	return c
}

func InOrder(calls ...CallHolder) {
	for i := 1; i < len(calls); i++ {
		calls[i].GetCall().After(calls[i-1])
	}
}

// If yes, returns nil. If no, returns error with message explaining why it does not match.
func (c *Call) matches(args []interface{}) error {
	if len(args) != len(c.argsMatchers) {
		return fmt.Errorf("expected call %s expects %d arguments, got %d",
			c.origin, len(c.argsMatchers), len(args))
	}

	for i, matcher := range c.argsMatchers {
		if !matcher(args[i]) {
			pp := pp.New()
			pp.SetColoringEnabled(false)
			return fmt.Errorf(
				"expected call %s doesn't match %v argument, got:\n\t%v",
				c.origin, ordinal(i+1), strings.ReplaceAll(pp.Sprint(args[i]), "\n", "\n\t"),
			)
		}
	}

	for _, prevCall := range c.prevCalls {
		if !prevCall.satisfied() {
			return fmt.Errorf("expected call %s should be called after call %v",
				c.origin, prevCall.origin)
		}
	}

	if len(c.nextCalled) != 0 {
		return fmt.Errorf("some expected calls planned to be after expected call %s have already been called", c.origin)
	}

	if c.actualCallsNum >= c.maxCalls {
		return fmt.Errorf("expected call %s has already been called the max number of times", c.origin)
	}

	return nil
}

func (c *Call) Return(rets ...interface{}) *Call {
	c.action = func([]interface{}) []interface{} {
		return rets
	}
	return c
}

func (c *Call) RunReturn(f interface{}) *Call {
	c.t.Helper()
	v := reflect.ValueOf(f)
	if ft := v.Type(); len(c.argsMatchers) != ft.NumIn() {
		c.t.Fatalf("Wrong number of arguments in RunReturn func for call %v: got %d, want %d",
			c.origin, len(c.argsMatchers), ft.NumIn())
	}
	c.action = func(args []interface{}) []interface{} {
		c.t.Helper()
		vArgs := make([]reflect.Value, len(args))
		for i := 0; i < len(args); i++ {
			vArgs[i] = reflect.ValueOf(args[i])
		}
		vRets := v.Call(vArgs)
		rets := make([]interface{}, len(vRets))
		for i, ret := range vRets {
			rets[i] = ret.Interface()
		}
		return rets
	}
	return c
}
