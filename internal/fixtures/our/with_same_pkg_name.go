package fixtures

import (
	h "html/template"
	t "text/template"
)

type SomeInterface interface {
	Foo(i t.Template) h.Template
}
