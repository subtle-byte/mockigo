// Code generated by mockigo. DO NOT EDIT.

//go:generate mockigo 

package comment

import match "github.com/subtle-byte/mockigo/match"
import mock "github.com/subtle-byte/mockigo/mock"

var _ = match.Any[int]

type IfaceWithBuildTagInCommentMock struct {
	mock *mock.Mock
}

func NewIfaceWithBuildTagInCommentMock(t mock.Testing) *IfaceWithBuildTagInCommentMock {
	t.Helper()
	return &IfaceWithBuildTagInCommentMock{mock: mock.NewMock(t)}
}

type _IfaceWithBuildTagInCommentMock_Expecter struct {
	mock *mock.Mock
}

func (_mock *IfaceWithBuildTagInCommentMock) EXPECT() _IfaceWithBuildTagInCommentMock_Expecter {
	 return _IfaceWithBuildTagInCommentMock_Expecter{mock: _mock.mock}
}

type _IfaceWithBuildTagInCommentMock_Sprintf_Call struct {
	*mock.Call
}

func (_mock *IfaceWithBuildTagInCommentMock) Sprintf(format string, a ...interface{}) (string) {
	_mock.mock.T.Helper()
	_args := []any{format, mock.SliceToAnySlice(a)}
	_results := _mock.mock.Called("Sprintf", _args...)
	_r0 := _results.Get(0).(string)
	return _r0
}

func (_expecter _IfaceWithBuildTagInCommentMock_Expecter) Sprintf(format match.Arg[string], a ...match.Arg[interface{}]) _IfaceWithBuildTagInCommentMock_Sprintf_Call {
	_args := append([]mock.Matcher{format.Matcher}, match.ArgsToMatchers(a)...)
	return _IfaceWithBuildTagInCommentMock_Sprintf_Call{Call: _expecter.mock.ExpectCall("Sprintf", _args...)}
}

func (_call _IfaceWithBuildTagInCommentMock_Sprintf_Call) Return(_r0 string) _IfaceWithBuildTagInCommentMock_Sprintf_Call {
	_call.Call.Return(_r0)
	return _call
}

func (_call _IfaceWithBuildTagInCommentMock_Sprintf_Call) RunReturn(f func(format string, a ...interface{}) (string)) _IfaceWithBuildTagInCommentMock_Sprintf_Call {
	_call.Call.RunReturn(f)
	return _call
}
