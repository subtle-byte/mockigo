package test

import (
	"github.com/subtle-byte/mockigo/internal/fixtures/mockery/http"
)

type HasConflictingNestedImports interface {
	RequesterNS
	Z() http.MyStruct
}
