package main

// https://adventofcode.com/2022/day/9

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	Process("aoc22/day09/example1.txt", 1)
	Process("aoc22/day09/input.txt", 1)
	Process("aoc22/day09/example2.txt", 9)
	Process("aoc22/day09/input.txt", 9)
}

func Process(name string, tails int) {
	fmt.Println("input:", name)
	steps := ReadInput(name)
	b := NewBoard(tails)
	for _, step := range steps {
		b.apply(step)
	}
	fmt.Println("visited", len(b.Visited))
	b.RenderVisited()

	fmt.Println()
}

func ReadInput(name string) []Step {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	var steps []Step
	for s.Scan() {
		line := s.Text()
		var s Step
		n, err := fmt.Sscanf(line, "%c %d", &s.Dir, &s.Num)
		if n != 2 || err != nil {
			panic("invalid input")
		}
		steps = append(steps, s)
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return steps
}

type Step struct {
	Dir rune
	Num int
}

type Pos struct {
	X int
	Y int
}

type Board struct {
	Head    Pos
	Tails   []Pos
	Visited map[Pos]bool
}

func NewBoard(tails int) *Board {
	return &Board{
		Tails: make([]Pos, tails),
		Visited: map[Pos]bool{
			Pos{0, 0}: true, // starting position
		},
	}
}

func (b *Board) apply(s Step) {
	for i := 0; i < s.Num; i++ {
		switch s.Dir {
		case 'R':
			b.MoveHead(1, 0)
		case 'L':
			b.MoveHead(-1, 0)
		case 'U':
			b.MoveHead(0, 1)
		case 'D':
			b.MoveHead(0, -1)
		default:
			panic(fmt.Errorf("invalid direction: %c", s.Dir))
		}
		b.FollowTails()
	}
}

func (b *Board) FollowTails() {
	l := len(b.Tails)
	for i := 0; i < l; i++ {
		b.FollowTail(i)
	}
}
func (b *Board) FollowTail(i int) {
	var xdist, ydist int
	if i == 0 {
		xdist = b.Head.X - b.Tails[i].X
		ydist = b.Head.Y - b.Tails[i].Y
	} else {
		xdist = b.Tails[i-1].X - b.Tails[i].X
		ydist = b.Tails[i-1].Y - b.Tails[i].Y
	}
	//fmt.Printf("distance x=%d, y=%d\n", xdist, ydist)

	if abs(xdist) <= 1 && abs(ydist) <= 1 {
		// no moves required
		return
	}

	// horizontal
	if ydist == 0 {
		b.MoveTail(i, sign(xdist), 0)
		return
	}

	// vertical
	if xdist == 0 {
		b.MoveTail(i, 0, sign(ydist))
		return
	}

	// diagonal
	b.MoveTail(i, sign(xdist), sign(ydist))
}

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}
func sign(a int) int {
	if a < 0 {
		return -1
	}
	if a == 0 {
		return 0
	}
	return 1
}

func (b *Board) MoveHead(x int, y int) {
	b.Head.X += x
	b.Head.Y += y
	//fmt.Printf("move head (%d, %d) => (%d, %d)\n", x, y, b.Head.X, b.Head.Y)
}
func (b *Board) MoveTail(i int, x int, y int) {
	b.Tails[i].X += x
	b.Tails[i].Y += y
	//fmt.Printf("move tail (%d, %d) => (%d, %d)\n", x, y, b.Tail.X, b.Tail.Y)

	// visited tracks the path of the last tail
	if i == len(b.Tails)-1 {
		b.Visited[b.Tails[i]] = true
	}
}

func (b *Board) RenderVisited() {
	var minx, maxx, miny, maxy int
	first := true
	for pos, _ := range b.Visited {
		if first {
			minx, maxx = pos.X, pos.X
			miny, maxy = pos.Y, pos.Y
			first = false
		} else {
			minx = min(minx, pos.X)
			maxx = max(maxx, pos.X)
			miny = min(miny, pos.Y)
			maxy = max(maxy, pos.Y)
		}
	}
	xrange := maxx - minx + 1
	yrange := maxy - miny + 1
	field := make([]rune, xrange*yrange)
	for i, _ := range field {
		field[i] = '.'
	}
	for pos, _ := range b.Visited {
		x, y := pos.X-minx, pos.Y-miny
		field[y*xrange+x] = '#'
	}
	for y := yrange - 1; y >= 0; y-- {
		fmt.Println(string(field[y*xrange : (y+1)*xrange]))
	}
}
