package rowcol

import (
	"fmt"
	"iter"
)

type Grid[T any] struct {
	Data [][]T
}

//func (g *Grid[byte]) Print() {
//	for _, row := range g.Data {
//		fmt.Println(string(row))
//	}
//	fmt.Println()
//}

func (g *Grid[T]) Rows() int                  { return len(g.Data) }
func (g *Grid[T]) Columns() int               { return len(g.Data[0]) }
func (g *Grid[T]) Size() (rows int, cols int) { return g.Rows(), g.Columns() }

func (g *Grid[T]) SetRow(i int, data []T) {
	if len(data) != g.Columns() {
		panic(fmt.Errorf("invalid data size, got=%d, expected=%d", len(data), g.Columns()))
	}
	copy(g.Data[i], data)
}

func (g *Grid[T]) Get(row, col int) T    { return g.Data[row][col] }
func (g *Grid[T]) GetPos(p Pos) T        { return g.Data[p.Row][p.Col] }
func (g *Grid[T]) Set(row, col int, v T) { g.Data[row][col] = v }
func (g *Grid[T]) SetPos(p Pos, v T)     { g.Data[p.Row][p.Col] = v }
func (g *Grid[T]) IsValidPosition(row, col int) bool {
	return row >= 0 && col >= 0 && row < g.Rows() && col < g.Columns()
}
func (g *Grid[T]) IsValidPos(p Pos) bool { return g.IsValidPosition(p.Row, p.Col) }

func (g *Grid[T]) SetIfValid(row, col int, v T) {
	if g.IsValidPosition(row, col) {
		g.Data[row][col] = v
	}
}
func (g *Grid[T]) SetIfValidPos(p Pos, v T) { g.SetIfValid(p.Row, p.Col, v) }

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

func (g *Grid[T]) VisitWithPos(fn func(T, Pos)) {
	rows, cols := g.Rows(), g.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			fn(g.Data[r][c], Pos{Row: r, Col: c})
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
func Reduce[T any, E any](g *Grid[T], acc E, fn func(E, T) E) E {
	rows, cols := g.Rows(), g.Columns()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			acc = fn(acc, g.Data[r][c])
		}
	}
	return acc
}

func (g *Grid[T]) Find(pred func(T) bool) (Pos, bool) {
	rows, cols := g.Size()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if pred(g.Data[r][c]) {
				return Pos{r, c}, true
			}
		}
	}
	return Pos{}, false
}

func (g *Grid[T]) FindAll(pred func(T) bool) []Pos {
	var rv []Pos
	rows, cols := g.Size()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if pred(g.Data[r][c]) {
				rv = append(rv, Pos{r, c})
			}
		}
	}
	return rv
}

func (g *Grid[T]) Reset(v T) {
	rows, cols := g.Size()
	// first row
	for c := 0; c < cols; c++ {
		g.Data[0][c] = v
	}
	// all other rows
	for r := 1; r < rows; r++ {
		copy(g.Data[r], g.Data[0])
	}
}

func (g *Grid[T]) PosIterator() iter.Seq[Pos] {
	rows, cols := g.Size()
	return func(yield func(Pos) bool) {
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				if !yield(Pos{r, c}) {
					return
				}
			}
		}
	}
}

/*
func (g *Grid[T]) ResetByRows(data []T) {
	if len(data) != g.Columns() {
		panic(fmt.Errorf("invalid data size, got=%d, expected=%d", len(data), g.Columns()))
	}
	rows := g.Rows()
	for r := 0; r < rows; r++ {
		copy(g.Data[r], data)
	}
}
*/

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
