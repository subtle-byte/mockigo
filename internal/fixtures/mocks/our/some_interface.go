// Code generated by mockigo. DO NOT EDIT.

package fixtures

import html_template "html/template"
import match "github.com/subtle-byte/mockigo/match"
import mock "github.com/subtle-byte/mockigo/mock"
import text_template "text/template"

var _ = match.Any[int]

type SomeInterface struct {
	mock *mock.Mock
}

func NewSomeInterface(t mock.Testing) *SomeInterface {
	t.Helper()
	return &SomeInterface{mock: mock.NewMock(t)}
}

type _SomeInterface_Expecter struct {
	mock *mock.Mock
}

func (_mock *SomeInterface) EXPECT() _SomeInterface_Expecter {
	 return _SomeInterface_Expecter{mock: _mock.mock}
}

type _SomeInterface_Foo_Call struct {
	*mock.Call
}

func (_mock *SomeInterface) Foo(i text_template.Template) (html_template.Template) {
	_mock.mock.T.Helper()
	_results := _mock.mock.Called("Foo", i)
	_r0 := _results.Get(0).(html_template.Template)
	return _r0
}

func (_expecter _SomeInterface_Expecter) Foo(i match.Arg[text_template.Template]) _SomeInterface_Foo_Call {
	return _SomeInterface_Foo_Call{Call: _expecter.mock.ExpectCall("Foo", i.Arg)}
}

func (_call _SomeInterface_Foo_Call) Return(_r0 html_template.Template) _SomeInterface_Foo_Call {
	_call.Call.Return(_r0)
	return _call
}

func (_call _SomeInterface_Foo_Call) RunReturn(f func(i text_template.Template) (html_template.Template)) _SomeInterface_Foo_Call {
	_call.Call.RunReturn(f)
	return _call
}