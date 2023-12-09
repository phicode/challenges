package lib

func Reduce[T any, A any](ts []T, reducer func(t T, acc A) A, acc A) A {
	for _, t := range ts {
		acc = reducer(t, acc)
	}
	return acc
}
