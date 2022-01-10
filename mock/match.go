package mock

import (
	"reflect"
)

type Matcher func(x interface{}) bool

// Any matcher

type anyMatcher struct{}

func (anyMatcher) Matches(interface{}) bool {
	return true
}

func Any() Matcher { return anyMatcher{}.Matches }

// Eq matcher

type eqMatcher struct {
	x interface{}
}

func (e eqMatcher) Matches(x interface{}) bool {
	return reflect.DeepEqual(e.x, x)

}

func Eq(x interface{}) Matcher { return eqMatcher{x}.Matches }

// Not matcher

type notMatcher struct {
	m Matcher
}

func (n notMatcher) Matches(x interface{}) bool {
	return !n.m(x)
}

func Not(x Matcher) Matcher {
	return notMatcher{x}.Matches
}
