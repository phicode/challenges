package lib

func Map[S, T any](xs []S, f func(S) T) []T {
	ts := make([]T, len(xs))
	for i, x := range xs {
		ts[i] = f(x)
	}
	return ts
}
