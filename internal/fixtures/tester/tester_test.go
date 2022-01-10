package tester_test

import (
	html_template "html/template"
	"strconv"
	"testing"
	text_template "text/template"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mocks "github.com/subtle-byte/mockigo/internal/fixtures/mocks/our"
	"github.com/subtle-byte/mockigo/match"
	"github.com/subtle-byte/mockigo/mock"
)

func TestSimpleInterface(t *testing.T) {
	simple := mocks.NewSimpleInterface(t)

	mock.InOrder(1, 1,
		simple.EXPECT().
			Bar(match.MatchedBy(func(i int) bool {
				return i > 5
			})).
			RunReturn(func(i int) int {
				return i - 5
			}),
		simple.EXPECT().
			Bar(match.Any[int]()).
			Return(0),
	)
	r1 := simple.Bar(8)
	r2 := simple.Bar(2)
	require.Equal(t, 3, r1)
	require.Equal(t, 0, r2)
}

func TestSomeInterface(t *testing.T) {
	someInterface := mocks.NewSomeInterface(t)

	var textTemplate text_template.Template = text_template.Template{}
	var htmlTemplate html_template.Template = html_template.Template{}

	someInterface.EXPECT().
		Foo(match.Eq(textTemplate)).
		Return(htmlTemplate)

	var returnedHtmlTemplate html_template.Template = someInterface.Foo(textTemplate)

	require.Equal(t, htmlTemplate, returnedHtmlTemplate)
}

func TestFooBar(t *testing.T) {
	fb := mocks.NewFooBar(t)
	mock.InOrder(1, -1,
		fb.EXPECT().Foo(match.Eq(7)).Return(8),
		fb.EXPECT().Bar(match.Eq(time.Second)),
	)
	r := fb.Foo(7)
	assert.Equal(t, 8, r)
	fb.Bar(time.Second)
}

func TestBarFoo(t *testing.T) {
	m := mocks.NewBarFoo(t)
	mock.InOrder(1, -1,
		m.EXPECT().Foo(match.Any[int](), match.Eq("hello")).Return(45),
		m.EXPECT().Bar(match.Any[int]()).RunReturn(func(n int) string {
			return strconv.Itoa(n)
		}),
	)
	fooRet := m.Foo(100, "hello") // == 45
	barRet := m.Bar(200)          // == "200"

	assert.Equal(t, 45, fooRet)
	assert.Equal(t, "200", barRet)
}
