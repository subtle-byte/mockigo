// Code generated by mockigo. DO NOT EDIT.

package test

import match "github.com/subtle-byte/mockigo/match"
import mock "github.com/subtle-byte/mockigo/mock"

var _ = match.Any[int]

type RequesterReturnElided struct {
	mock *mock.Mock
}

func NewRequesterReturnElided(t mock.Testing) *RequesterReturnElided {
	t.Helper()
	return &RequesterReturnElided{mock: mock.NewMock(t)}
}

type _RequesterReturnElided_Expecter struct {
	mock *mock.Mock
}

func (_mock *RequesterReturnElided) EXPECT() _RequesterReturnElided_Expecter {
	 return _RequesterReturnElided_Expecter{mock: _mock.mock}
}

type _RequesterReturnElided_Get_Call struct {
	*mock.Call
}

func (_mock *RequesterReturnElided) Get(path string) (a int, b int, c int, err error) {
	_mock.mock.T.Helper()
	_results := _mock.mock.Called("Get", path)
	_r0 := _results.Get(0).(int)
	_r1 := _results.Get(1).(int)
	_r2 := _results.Get(2).(int)
	_r3 := _results.Error(3)
	return _r0, _r1, _r2, _r3
}

func (_expecter _RequesterReturnElided_Expecter) Get(path match.Arg[string]) _RequesterReturnElided_Get_Call {
	return _RequesterReturnElided_Get_Call{Call: _expecter.mock.ExpectCall("Get", path.Arg)}
}

func (_call _RequesterReturnElided_Get_Call) Return(a int, b int, c int, err error) _RequesterReturnElided_Get_Call {
	_call.Call.Return(a, b, c, err)
	return _call
}

func (_call _RequesterReturnElided_Get_Call) RunReturn(f func(path string) (a int, b int, c int, err error)) _RequesterReturnElided_Get_Call {
	_call.Call.RunReturn(f)
	return _call
}