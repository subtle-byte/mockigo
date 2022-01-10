package util

func MapSlice[I any, R any](in []I, f func(I) R) []R {
	res := make([]R, len(in))
	for i, v := range in {
		res[i] = f(v)
	}
	return res
}

func MapSliceWithIndex[I any, R any](in []I, f func(i int, item I) R) []R {
	res := make([]R, len(in))
	for i, v := range in {
		res[i] = f(i, v)
	}
	return res
}
