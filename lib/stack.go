package lib

import "slices"

type Stack[T any] []T

func NewStack[T any]() Stack[T] {
	return make([]T, 0, 32)
}

func (s *Stack[T]) Push(x T) {
	*s = append(*s, x)
}

func (s *Stack[T]) Reverse() {
	slices.Reverse(*s)
}

func (s *Stack[T]) Pop() T {
	last := len(*s) - 1
	r := (*s)[last]
	*s = (*s)[:last]
	return r
}

func (s *Stack[T]) Peek() T {
	l := len(*s)
	if l == 0 {
		var t T
		return t // default T
	}
	return (*s)[l-1]
}
func (s *Stack[T]) Len() int { return len(*s) }
