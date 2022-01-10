package test

import (
	"github.com/subtle-byte/mockigo/internal/fixtures/mockery/test"
)

type C int

type ImportsSameAsPackage interface {
	A() test.B
	B() KeyManager
	C(C)
}
