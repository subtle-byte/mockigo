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
	not := Not(eq)
	matchers := ArgsToMatchers([]Arg[int]{
		{Matcher: eq.Matcher},
		{Matcher: any.Matcher},
		{Matcher: matcherBy.Matcher},
		{Matcher: not.Matcher},
	})
	require.Len(t, matchers, 4)
	require.Equal(t, true, matchers[0](7))
	require.Equal(t, false, matchers[0](8))
	require.Equal(t, true, matchers[1](7))
	require.Equal(t, true, matchers[2](10))
	require.Equal(t, false, matchers[2](2))
	require.Equal(t, false, matchers[2]("hello"))
	require.Equal(t, true, matchers[3](6))
	require.Equal(t, false, matchers[3](7))
	require.Equal(t, true, matchers[3](9))
}
