package mock

type Rets struct {
	rets []interface{}
	call *Call
}

func (r Rets) Len() int {
	return len(r.rets)
}

func (r Rets) Get(i int) interface{} {
	r.call.t.Helper()
	if i >= r.Len() {
		r.call.t.Fatalf("Call %v does not have %v return value", r.call.origin, ordinal(i+1))
	}
	return r.rets[i]
}

func (r Rets) Error(i int) error {
	r.call.t.Helper()
	ret := r.Get(i)
	if ret == nil {
		return nil
	}
	e, ok := ret.(error)
	if !ok {
		r.call.t.Fatalf("Call %v does not have %v return value of the error type", r.call.origin, ordinal(i+1))
	}
	return e
}
