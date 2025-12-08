package main

// https://adventofcode.com/2022/day/18

import (
	"fmt"

	"github.com/phicode/challenges/lib"
)

var VERBOSE = 1

func main() {
	// 64
	ProcessPart1("aoc22/day18/example.txt")
	ProcessPart1("aoc22/day18/input.txt")

	ProcessPart2("aoc22/day18/example.txt")
	ProcessPart2("aoc22/day18/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	ps := ParsePositions(lines)
	area := FindSurfaceArea(ps)
	fmt.Println("Area:", area)
	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	ps := ParsePositions(lines)
	f := NewFiller(ps)
	f.Run()
	area := f.FindSurfaceArea()
	fmt.Println("Area:", area)
	fmt.Println()
}

////////////////////////////////////////////////////////////

type Position struct {
	X, Y, Z int
}

func (p Position) String() string {
	return fmt.Sprintf("%d,%d,%d", p.X, p.Y, p.Z)
}

func ParsePositions(lines []string) []Position {
	var rv []Position
	for _, l := range lines {
		var p Position
		n, err := fmt.Sscanf(l, "%d,%d,%d", &p.X, &p.Y, &p.Z)
		if n != 3 || err != nil {
			panic(fmt.Errorf("invalid input; n=%d, err=%w", n, err))
		}
		rv = append(rv, p)
	}
	return rv
}

////////////////////////////////////////////////////////////
// PART 1

func FindSurfaceArea(ps []Position) int {
	m := make(map[Position]struct{})
	for _, p := range ps {
		m[p] = struct{}{}
	}
	var area int
	for _, p := range ps {
		area += Exposed(m, p, 1, 0, 0)
		area += Exposed(m, p, -1, 0, 0)
		area += Exposed(m, p, 0, 1, 0)
		area += Exposed(m, p, 0, -1, 0)
		area += Exposed(m, p, 0, 0, 1)
		area += Exposed(m, p, 0, 0, -1)
	}
	return area
}

func Exposed(m map[Position]struct{}, p Position, addX, addY, addZ int) int {
	test := Position{p.X + addX, p.Y + addY, p.Z + addZ}
	if _, found := m[test]; found {
		return 0
	}
	return 1
}

////////////////////////////////////////////////////////////
// PART 2

type State byte

const (
	NotVisited State = iota
	Voxel
	Outside
)

type Filler struct {
	ps   []Position
	cube *Cube[State]
	Max  Position
	//toVisit *list.List
}

func NewFiller(ps []Position) *Filler {
	_max := lib.Reduce(ps, Position{}, MaxPos)
	c := NewCube[State](_max)
	for _, p := range ps {
		c.Set(p, Voxel) // contains a voxel
	}
	return &Filler{
		ps:   ps,
		cube: c,
		Max:  _max,
	}
}

func (f *Filler) Run() {
	f.Follow(f.Max)
	var n int
	for _, x := range f.cube.vs {
		if x == NotVisited {
			n++
		}
	}
	fmt.Println("not visited:", n)
}

func (f *Filler) Follow(p Position) {
	shouldFollow := p.X >= 0 && p.Y >= 0 && p.Z >= 0 &&
		p.X <= f.Max.X && p.Y <= f.Max.Y && p.Z <= f.Max.Z &&
		f.cube.Get(p) == NotVisited
	if !shouldFollow {
		return
	}
	f.cube.Set(p, Outside)

	// follow all possible neighbors
	f.FollowNeigh(p, 1, 0, 0)
	f.FollowNeigh(p, -1, 0, 0)
	f.FollowNeigh(p, 0, 1, 0)
	f.FollowNeigh(p, 0, -1, 0)
	f.FollowNeigh(p, 0, 0, 1)
	f.FollowNeigh(p, 0, 0, -1)
}

func (f *Filler) FollowNeigh(p Position, x, y, z int) {
	p = Position{p.X + x, p.Y + y, p.Z + z}
	f.Follow(p)
}

type Cube[T any] struct {
	// dimensions
	d  Position
	vs []T
}

func NewCube[T any](dimensions Position) *Cube[T] {
	dimensions.X++
	dimensions.Y++
	dimensions.Z++
	size := dimensions.X * dimensions.Y * dimensions.Z
	return &Cube[T]{
		d:  dimensions,
		vs: make([]T, size),
	}
}
func (c *Cube[T]) ValidPosition(p Position) bool {
	return p.X >= 0 && p.Y >= 0 && p.Z >= 0 &&
		p.X < c.d.X && p.Y < c.d.Y && p.Z < c.d.Z
}
func (c *Cube[T]) posToIndex(p Position) int {
	if !c.ValidPosition(p) {
		panic(fmt.Errorf("invalid position=%v for cube of dimension=%v", p, c.d))
	}
	return p.X + (p.Y * c.d.X) + (p.Z * c.d.X * c.d.Y)
}

func (c *Cube[T]) Set(p Position, v T) { c.vs[c.posToIndex(p)] = v }
func (c *Cube[T]) Get(p Position) T    { return c.vs[c.posToIndex(p)] }

func MaxPos(a, b Position) Position {
	return Position{
		max(a.X, b.X),
		max(a.Y, b.Y),
		max(a.Z, b.Z),
	}
}

func (f *Filler) FindSurfaceArea() int {

	var area int
	for _, p := range f.ps {
		area += f.Exposed(p, 1, 0, 0)
		area += f.Exposed(p, -1, 0, 0)
		area += f.Exposed(p, 0, 1, 0)
		area += f.Exposed(p, 0, -1, 0)
		area += f.Exposed(p, 0, 0, 1)
		area += f.Exposed(p, 0, 0, -1)
	}
	return area
}

func (f *Filler) Exposed(p Position, addX, addY, addZ int) int {
	test := Position{p.X + addX, p.Y + addY, p.Z + addZ}
	if !f.cube.ValidPosition(test) {
		// outside of tracked cube
		// = voxel that touches the borders
		// = exposed area
		return 1
	}
	if f.cube.Get(test) == Outside {
		return 1
	}
	return 0
}
