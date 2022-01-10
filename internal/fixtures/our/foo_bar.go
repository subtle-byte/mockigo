package fixtures

import "time"

type FooBar interface {
	Foo(a int) int
	Bar(t time.Duration)
}

type BarFoo interface {
	Foo(int, string) int
	Bar(int) string
}
