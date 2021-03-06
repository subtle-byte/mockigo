// Code generated by mockigo. DO NOT EDIT.

package fixtures

import match "github.com/subtle-byte/mockigo/match"
import mock "github.com/subtle-byte/mockigo/mock"
import some_interface "github.com/subtle-byte/mockigo/internal/fixtures/our/some_interface"

var _ = match.Any[int]

type GenericFunc[Y some_interface.SomeInterface] struct {
	mock *mock.Mock
}

func NewGenericFunc[Y some_interface.SomeInterface](t mock.Testing) *GenericFunc[Y] {
	t.Helper()
	return &GenericFunc[Y]{mock: mock.NewMock(t)}
}

type _GenericFunc_Expecter[Y some_interface.SomeInterface] struct {
	mock *mock.Mock
}

func (_mock *GenericFunc[Y]) EXPECT() _GenericFunc_Expecter[Y] {
	 return _GenericFunc_Expecter[Y]{mock: _mock.mock}
}

type _GenericFunc_Execute_Call[Y some_interface.SomeInterface] struct {
	*mock.Call
}

func (_mock *GenericFunc[Y]) Execute() (Y) {
	_mock.mock.T.Helper()
	_results := _mock.mock.Called("Execute", )
	_r0 := _results.Get(0).(Y)
	return _r0
}

func (_expecter _GenericFunc_Expecter[Y]) Execute() _GenericFunc_Execute_Call[Y] {
	return _GenericFunc_Execute_Call[Y]{Call: _expecter.mock.ExpectCall("Execute", )}
}

func (_call _GenericFunc_Execute_Call[Y]) Return(_r0 Y) _GenericFunc_Execute_Call[Y] {
	_call.Call.Return(_r0)
	return _call
}

func (_call _GenericFunc_Execute_Call[Y]) RunReturn(f func() (Y)) _GenericFunc_Execute_Call[Y] {
	_call.Call.RunReturn(f)
	return _call
}
