package match

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	eq := Eq(7)
	any := Any[int]()
	f := func(i int) bool { return i > 5 }
	matcherBy := MatchedBy(f)
	matchers := ArgsToMatchers([]Arg[int]{
		{Arg: eq.Arg},
		{Arg: any.Arg},
		{Arg: matcherBy.Arg},
	})
	require.Len(t, matchers, 3)
	require.Equal(t, true, matchers[0](7))
	require.Equal(t, false, matchers[0](8))
	require.Equal(t, true, matchers[1](7))
	require.Equal(t, true, matchers[2](10))
	require.Equal(t, false, matchers[2](2))
	require.Equal(t, false, matchers[2]("hello"))
}
