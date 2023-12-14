package lib

import "fmt"

type Grid[T any] struct {
	Data [][]T
}

func (g *Grid[T]) Rows() int    { return len(g.Data) }
func (g *Grid[T]) Columns() int { return len(g.Data[0]) }

func (g *Grid[T]) SetRow(i int, data []T) {
	if len(data) != g.Columns() {
		panic(fmt.Errorf("invalid data size, got=%d, expected=%d", len(data), g.Columns()))
	}
	copy(g.Data[i], data)
}

func NewByteGridFromStrings(xs []string) Grid[byte] {
	g := NewGrid[byte](len(xs), len(xs[0]))
	for i, x := range xs {
		g.SetRow(i, []byte(x))
	}
	return g
}

func NewGrid[T any](rows, columns int) Grid[T] {
	var g Grid[T]
	g.Data = make([][]T, rows)
	for r := 0; r < rows; r++ {
		g.Data[r] = make([]T, columns)
	}
	return g
}
