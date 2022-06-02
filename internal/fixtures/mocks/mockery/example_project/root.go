// Code generated by mockigo. DO NOT EDIT.

package example_project

import foo "github.com/subtle-byte/mockigo/internal/fixtures/mockery/example_project/foo"
import match "github.com/subtle-byte/mockigo/match"
import mock "github.com/subtle-byte/mockigo/mock"

var _ = match.Any[int]

type Root struct {
	mock *mock.Mock
}

func NewRoot(t mock.Testing) *Root {
	t.Helper()
	return &Root{mock: mock.NewMock(t)}
}

type _Root_Expecter struct {
	mock *mock.Mock
}

func (_mock *Root) EXPECT() _Root_Expecter {
	 return _Root_Expecter{mock: _mock.mock}
}

type _Root_ReturnsFoo_Call struct {
	*mock.Call
}

func (_mock *Root) ReturnsFoo() (foo.Foo, error) {
	_mock.mock.T.Helper()
	_results := _mock.mock.Called("ReturnsFoo", )
	var _r0 foo.Foo
	if _got := _results.Get(0); _got != nil {
		_r0 = _got.(foo.Foo)
	}
	_r1 := _results.Error(1)
	return _r0, _r1
}

func (_expecter _Root_Expecter) ReturnsFoo() _Root_ReturnsFoo_Call {
	return _Root_ReturnsFoo_Call{Call: _expecter.mock.ExpectCall("ReturnsFoo", )}
}

func (_call _Root_ReturnsFoo_Call) Return(_r0 foo.Foo, _r1 error) _Root_ReturnsFoo_Call {
	_call.Call.Return(_r0, _r1)
	return _call
}

func (_call _Root_ReturnsFoo_Call) RunReturn(f func() (foo.Foo, error)) _Root_ReturnsFoo_Call {
	_call.Call.RunReturn(f)
	return _call
}

type _Root_TakesBaz_Call struct {
	*mock.Call
}

func (_mock *Root) TakesBaz(_a0 *foo.Baz) () {
	_mock.mock.T.Helper()
	_mock.mock.Called("TakesBaz", _a0)
}

func (_expecter _Root_Expecter) TakesBaz(_a0 match.Arg[*foo.Baz]) _Root_TakesBaz_Call {
	return _Root_TakesBaz_Call{Call: _expecter.mock.ExpectCall("TakesBaz", _a0.Arg)}
}

func (_call _Root_TakesBaz_Call) Return() _Root_TakesBaz_Call {
	_call.Call.Return()
	return _call
}

func (_call _Root_TakesBaz_Call) RunReturn(f func(*foo.Baz) ()) _Root_TakesBaz_Call {
	_call.Call.RunReturn(f)
	return _call
}