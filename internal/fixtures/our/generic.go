package fixtures

import someinterface "github.com/subtle-byte/mockigo/internal/fixtures/our/some_interface"

type GenericInterface[T any, B someinterface.SomeInterface, G ~int | float32] interface {
	SomeMethod(B) T
}

type GenericFunc[Y someinterface.SomeInterface] func() Y

type GenericComparable[T comparable] func() T
