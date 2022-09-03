package mock

import (
	"reflect"
)

type Matcher func(x interface{}) bool

func Any() Matcher {
	return func(interface{}) bool {
		return true
	}
}

func Eq(x interface{}) Matcher {
	return func(y interface{}) bool {
		return reflect.DeepEqual(x, y)
	}
}

func Not(m Matcher) Matcher {
	return func(x interface{}) bool {
		return !m(x)
	}
}
