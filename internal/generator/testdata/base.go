package data

import "time"

type Nilable interface{}

type SimpleInterface interface {
	Foo(time.Time) (Nilable, error)
}

type GenericInterface[T any] interface {
	Foo(a int, b ...T) T
}

type Func func()

type NotInterface struct{}

var Var int

type NotExported interface{}

type Filtered interface{}
