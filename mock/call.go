package mock

import (
	"fmt"
	"reflect"
)

type Call struct {
	t                  Testing
	method             string
	argsMatchers       []Matcher
	prevCalls          []prevCall
	minCalls, maxCalls int
	actualCallsNum     int
	action             func([]interface{}) []interface{}
	origin             string
}

type PrevCall interface {
	After(minTimes, maxTimes int, previousCall PrevCall) *Call
	CallsNum() int
	Origin() string
}

type prevCall struct {
	minCalls, maxCalls int
	call               PrevCall
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

func (c *Call) Times(min, max int) *Call {
	if max == -1 {
		max = infCalls
	}
	if min > max {
		max = min
	}
	c.minCalls, c.maxCalls = min, max
	return c
}

func (c *Call) CallsNum() int {
	return c.actualCallsNum
}

func (c *Call) Origin() string {
	return c.origin
}

// Returns true if the minimum number of calls have been made.
func (c *Call) satisfied() bool {
	return c.actualCallsNum >= c.minCalls
}

// Returns true if the maximum number of calls have been made.
func (c *Call) exhausted() bool {
	return c.actualCallsNum >= c.maxCalls
}

func (c *Call) After(minTimes, maxTimes int, previousCall PrevCall) *Call {
	c.t.Helper()
	if maxTimes == -1 {
		maxTimes = infCalls
	}
	c.prevCalls = append(c.prevCalls, prevCall{
		minCalls: minTimes,
		maxCalls: maxTimes,
		call:     previousCall,
	})
	return c
}

func InOrder(minTimes, maxTimes int, calls ...PrevCall) {
	for i := 1; i < len(calls); i++ {
		calls[i].After(minTimes, maxTimes, calls[i-1])
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
			return fmt.Errorf(
				"expected call %s doesn't match the argument %#v at index %d",
				c.origin, args[i], i,
			)
		}
	}

	for _, prevCall := range c.prevCalls {
		if !(prevCall.minCalls <= prevCall.call.CallsNum() && prevCall.call.CallsNum() <= prevCall.maxCalls) {
			return fmt.Errorf("expected call %s should be called after call %v",
				c.origin, prevCall.call.Origin())
		}
	}

	if c.exhausted() {
		return fmt.Errorf("expected call %s has already been called the max number of times", c.origin)
	}

	return nil
}

func (c *Call) Return(rets ...interface{}) *Call {
	c.t.Helper()
	c.action = func([]interface{}) []interface{} {
		return rets
	}
	return c
}

func (c *Call) RunReturn(f interface{}) *Call {
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
