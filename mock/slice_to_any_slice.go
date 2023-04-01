package mock

// SliceToAnySlice is useful in mock objects while working with variadic arguments
func SliceToAnySlice[T any](sl []T) []any {
	r := make([]any, len(sl))
	for i, v := range sl {
		r[i] = any(v)
	}
	return r
}
