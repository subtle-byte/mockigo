package mockery

import (
	"net/http"

	my_http "github.com/subtle-byte/mockigo/internal/fixtures/mockery/http"
)

// Example is an example
type Example interface {
	A() http.Flusher
	B(fixtureshttp string) my_http.MyStruct
}
