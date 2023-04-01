// Code generated by mockigo. DO NOT EDIT.

//go:generate mockigo 

package example_project

import foo "github.com/subtle-byte/mockigo/internal/fixtures/mockery/example_project/foo"
import match "github.com/subtle-byte/mockigo/match"
import mock "github.com/subtle-byte/mockigo/mock"

var _ = match.Any[int]

type RootMock struct {
	mock *mock.Mock
}

func NewRootMock(t mock.Testing) *RootMock {
	t.Helper()
	return &RootMock{mock: mock.NewMock(t)}
}

type _RootMock_Expecter struct {
	mock *mock.Mock
}

func (_mock *RootMock) EXPECT() _RootMock_Expecter {
	 return _RootMock_Expecter{mock: _mock.mock}
}

type _RootMock_ReturnsFoo_Call struct {
	*mock.Call
}

func (_mock *RootMock) ReturnsFoo() (foo.Foo, error) {
	_mock.mock.T.Helper()
	_results := _mock.mock.Called("ReturnsFoo", )
	var _r0 foo.Foo
	if _got := _results.Get(0); _got != nil {
		_r0 = _got.(foo.Foo)
	}
	_r1 := _results.Error(1)
	return _r0, _r1
}

func (_expecter _RootMock_Expecter) ReturnsFoo() _RootMock_ReturnsFoo_Call {
	return _RootMock_ReturnsFoo_Call{Call: _expecter.mock.ExpectCall("ReturnsFoo", )}
}

func (_call _RootMock_ReturnsFoo_Call) Return(_r0 foo.Foo, _r1 error) _RootMock_ReturnsFoo_Call {
	_call.Call.Return(_r0, _r1)
	return _call
}

func (_call _RootMock_ReturnsFoo_Call) RunReturn(f func() (foo.Foo, error)) _RootMock_ReturnsFoo_Call {
	_call.Call.RunReturn(f)
	return _call
}

type _RootMock_TakesBaz_Call struct {
	*mock.Call
}

func (_mock *RootMock) TakesBaz(_a0 *foo.Baz) () {
	_mock.mock.T.Helper()
	_mock.mock.Called("TakesBaz", _a0)
}

func (_expecter _RootMock_Expecter) TakesBaz(_a0 match.Arg[*foo.Baz]) _RootMock_TakesBaz_Call {
	return _RootMock_TakesBaz_Call{Call: _expecter.mock.ExpectCall("TakesBaz", _a0.Matcher)}
}

func (_call _RootMock_TakesBaz_Call) Return() _RootMock_TakesBaz_Call {
	_call.Call.Return()
	return _call
}

func (_call _RootMock_TakesBaz_Call) RunReturn(f func(*foo.Baz) ()) _RootMock_TakesBaz_Call {
	_call.Call.RunReturn(f)
	return _call
}