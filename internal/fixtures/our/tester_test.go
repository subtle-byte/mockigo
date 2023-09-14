package fixtures

import (
	html_template "html/template"
	"strconv"
	"testing"
	text_template "text/template"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/subtle-byte/mockigo/match"
	"github.com/subtle-byte/mockigo/mock"
)

func TestSimpleInterface(t *testing.T) {
	simple := NewSimpleInterfaceMock(t)

	mock.InOrder(
		simple.EXPECT().
			Bar(match.MatchedBy(func(i int) bool {
				return i > 5
			})).
			RunReturn(func(i int) int {
				return i - 5
			}).
			Times(1, 1),
		simple.EXPECT().
			Bar(match.Any[int]()).
			Return(0),
	)
	require.Equal(t, 3, simple.Bar(8))
	require.Equal(t, 0, simple.Bar(8))
	require.Equal(t, 0, simple.Bar(2))
}

func TestSomeInterface(t *testing.T) {
	someInterface := NewSomeInterfaceMock(t)

	var textTemplate text_template.Template = text_template.Template{}
	var htmlTemplate html_template.Template = html_template.Template{}

	someInterface.EXPECT().
		Foo(match.Eq(textTemplate)).
		Return(htmlTemplate)

	var returnedHtmlTemplate html_template.Template = someInterface.Foo(textTemplate)

	require.Equal(t, htmlTemplate, returnedHtmlTemplate)
}

func TestFooBar(t *testing.T) {
	fb := NewFooBarMock(t)
	mock.InOrder(
		fb.EXPECT().Foo(match.Eq(7)).Return(8),
		fb.EXPECT().Bar(match.Eq(time.Second)),
	)
	assert.Equal(t, 8, fb.Foo(7))
	fb.Bar(time.Second)
}

func TestBarFoo(t *testing.T) {
	m := NewBarFooMock(t)
	mock.InOrder(
		m.EXPECT().Foo(match.Any[int](), match.Eq("hello")).Return(45),
		m.EXPECT().Bar(match.Any[int]()).RunReturn(func(n int) string {
			return strconv.Itoa(n)
		}),
	)
	assert.Equal(t, 45, m.Foo(100, "hello"))
	assert.Equal(t, "200", m.Bar(200))
}
