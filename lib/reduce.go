package lib

func Reduce[T any, A any](ts []T, acc A, reducer func(t T, acc A) A) A {
	for _, t := range ts {
		acc = reducer(t, acc)
	}
	return acc
}

// ReduceT reduced values of the same type
func ReduceT[T any](ts []T, acc T, reducer func(t T, acc T) T) T {
	for _, t := range ts {
		acc = reducer(t, acc)
	}
	return acc
}
