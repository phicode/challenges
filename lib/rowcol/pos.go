package rowcol

import "fmt"

type Pos struct {
	Row, Col int
}

func (p Pos) AddDir(dir Direction) Pos {
	return Pos{Row: p.Row + dir.Row, Col: p.Col + dir.Col}
}
func (p Pos) Add(b Pos) Pos         { return Pos{Row: p.Row + b.Row, Col: p.Col + b.Col} }
func (p Pos) AddS(row, col int) Pos { return Pos{Row: p.Row + row, Col: p.Col + col} }
func (p Pos) Sub(b Pos) Pos         { return Pos{Row: p.Row - b.Row, Col: p.Col - b.Col} }
func (p Pos) Reverse() Pos          { return Pos{Row: -p.Row, Col: -p.Col} }
func (p Pos) Right() Pos            { return Pos{Row: p.Col, Col: -p.Row} }
func (p Pos) Left() Pos             { return Pos{Row: -p.Col, Col: p.Row} }
func (p Pos) IsZero() bool          { return p.Row == 0 && p.Col == 0 }
func (p Pos) Abs() Pos              { return Pos{Row: abs(p.Row), Col: abs(p.Col)} }

// Position in Col, Row (analog to X, Y)
func (p Pos) String() string { return fmt.Sprintf("(%d,%d)", p.Col, p.Row) }

func (p Pos) Less(o Pos) bool {
	if p.Row < o.Row {
		return true
	}
	if p.Row == o.Row {
		return p.Col < o.Col
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Direction Pos

func (d Direction) String() string {
	switch d {
	case Right:
		return "Right"
	case Left:
		return "Left"
	case Up:
		return "Up"
	case Down:
		return "Down"
	default:
		return "Unknown"
	}
}

var (
	Left  = Direction{Row: 0, Col: -1}
	Right = Direction{Row: 0, Col: +1}
	Up    = Direction{Row: -1, Col: 0}
	Down  = Direction{Row: +1, Col: 0}
	Stand = Direction{Row: 0, Col: 0}
)

func (d Direction) Right() Direction { return Direction{Row: d.Col, Col: -d.Row} }
func (d Direction) Left() Direction  { return Direction{Row: -d.Col, Col: d.Row} }
func (d Direction) Reverse() Direction {
	return Direction{Row: -d.Row, Col: -d.Col}
}

func (d Direction) MulS(s int) Pos {
	return Pos{Row: d.Row * s, Col: d.Col * s}
}

var Directions = []Direction{Left, Right, Up, Down}
