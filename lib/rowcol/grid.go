package rowcol

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

func (g *Grid[T]) Get(row, col int) T    { return g.Data[row][col] }
func (g *Grid[T]) Set(row, col int, v T) { g.Data[row][col] = v }
func (g *Grid[T]) IsValidPosition(row, col int) bool {
	return row >= 0 && col >= 0 && row < g.Rows() && col < g.Columns()
}

func (g *Grid[T]) Copy() Grid[T] {
	c := NewGrid[T](g.Rows(), g.Columns())
	for i, row := range g.Data {
		copy(c.Data[i], row)
	}
	return c
}

func (g *Grid[T]) Map(mapper func(input T) T) {
	rows, cols := g.Rows(), g.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g.Data[r][c] = mapper(g.Data[r][c])
		}
	}
}

func (g *Grid[T]) Count(cmp func(T) bool) int {
	rows, cols := g.Rows(), g.Columns()
	count := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if cmp(g.Data[r][c]) {
				count++
			}
		}
	}
	return count
}
func (g *Grid[T]) Visit(fn func(T)) {
	rows, cols := g.Rows(), g.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			fn(g.Data[r][c])
		}
	}
}
func (g *Grid[T]) Reduce(acc T, fn func(T, T) T) T {
	rows, cols := g.Rows(), g.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			acc = fn(acc, g.Data[r][c])
		}
	}
	return acc
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
