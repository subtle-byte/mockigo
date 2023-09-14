package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListFormatsNamedNotVariadic(t *testing.T) {
	formats := TupleFormatter{
		rawNames:   []string{"a", "_", "bc"},
		types:      []string{"int", "rune", "[]*string"},
		isVariadic: false,
	}.Format("prefix_", map[string]struct{}{})
	assert.Equal(t, TupleFormats{
		Args:                     "a, prefix_1, bc",
		ArgsMatchers:             "a.Matcher, prefix_1.Matcher, bc.Matcher",
		RawParams:                "a int, _ rune, bc []*string",
		NamedParams:              "a int, prefix_1 rune, bc []*string",
		NamedArgedParams:         "a match.Arg[int], prefix_1 match.Arg[rune], bc match.Arg[[]*string]",
		VariadicArgs:             "",
		VariadicArgsEval:         "",
		VariadicArgsMatchersEval: "",
	}, formats)
}

func TestListFormatsNotNamedVariadic(t *testing.T) {
	formats := TupleFormatter{
		rawNames:   []string{"", "", ""},
		types:      []string{"int", "rune", "...*string"},
		isVariadic: true,
	}.Format("prefix_", map[string]struct{}{})
	assert.Equal(t, TupleFormats{
		Args:                     "prefix_0, prefix_1, prefix_2...",
		ArgsMatchers:             "prefix_0.Matcher, prefix_1.Matcher, match.ArgsToInterfaces(prefix_2)...",
		RawParams:                "int, rune, ...*string",
		NamedParams:              "prefix_0 int, prefix_1 rune, prefix_2 ...*string",
		NamedArgedParams:         "prefix_0 match.Arg[int], prefix_1 match.Arg[rune], prefix_2 ...match.Arg[*string]",
		VariadicArgs:             "_args...",
		VariadicArgsEval:         "\t_args := []any{prefix_0, prefix_1, mock.SliceToAnySlice(prefix_2)}\n",
		VariadicArgsMatchersEval: "\t_args := append([]mock.Matcher{prefix_0.Matcher, prefix_1.Matcher}, match.ArgsToMatchers(prefix_2)...)\n",
	}, formats)
}
