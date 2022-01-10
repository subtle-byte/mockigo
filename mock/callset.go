package mock

import (
	"bytes"
	"errors"
	"fmt"
)

type callSet struct {
	calls map[callSetKey][]*Call
}

type callSetKey struct {
	method string
}

func newCallSet() *callSet {
	return &callSet{make(map[callSetKey][]*Call)}
}

// Add adds a new expected call.
func (cs callSet) Add(call *Call) {
	key := callSetKey{call.method}
	m := cs.calls
	m[key] = append(m[key], call)
}

func (cs callSet) FindMatch(method string, args []interface{}) (*Call, error) {
	key := callSetKey{method}
	methodCalls := cs.calls[key]
	var callsErrors bytes.Buffer
	for _, call := range methodCalls {
		err := call.matches(args)
		if err != nil {
			_, _ = fmt.Fprintf(&callsErrors, "\n%v", err)
		} else {
			return call, nil
		}
	}
	if len(methodCalls) == 0 {
		_, _ = fmt.Fprintf(&callsErrors, " there are no expected calls of the method %q", method)
	}
	return nil, errors.New(callsErrors.String())
}

func (cs callSet) Unsatisfied() []*Call {
	unsatisfied := []*Call(nil)
	for _, methodCalls := range cs.calls {
		for _, call := range methodCalls {
			if !call.satisfied() {
				unsatisfied = append(unsatisfied, call)
			}
		}
	}
	return unsatisfied
}
