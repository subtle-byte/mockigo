package util

func MapSlice[I any, R any](in []I, f func(I) R) []R {
	res := make([]R, len(in))
	for i, v := range in {
		res[i] = f(v)
	}
	return res
}

func SliceToSet[T comparable](sl []T) map[T]struct{} {
	res := make(map[T]struct{}, len(sl))
	for _, v := range sl {
		res[v] = struct{}{}
	}
	return res
}
