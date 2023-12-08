package lib

func Filter[T any](xs []T, filter func(T) bool) []T {
	var rv []T
	for _, x := range xs {
		if filter(x) {
			rv = append(rv, x)
		}
	}
	return rv
}

func All[T any](xs []T, predicate func(T) bool) bool {
	for _, x := range xs {
		if !predicate(x) {
			return false
		}
	}
	return true
}

func Count[T any](xs []T, predicate func(T) bool) int {
	var c int
	for _, x := range xs {
		if predicate(x) {
			c++
		}
	}
	return c
}
