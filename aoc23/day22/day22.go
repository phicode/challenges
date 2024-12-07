package main

// https://adventofcode.com/2023/day/22

import (
	"fmt"

	"github.com/phicode/challenges/lib"
)

var VERBOSE = 2

func main() {
	ProcessPart1("aoc23/day22/example.txt")

	VERBOSE = 1
	ProcessPart1("aoc23/day22/input.txt") // 443

	ProcessPart2("aoc23/day22/example.txt") // 7
	ProcessPart2("aoc23/day22/input.txt")   // 69915
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	bricks := ParseBricks(lines)
	SolvePart1(bricks)
	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	bricks := ParseBricks(lines)
	SolvePart2(bricks)
	fmt.Println()
}

////////////////////////////////////////////////////////////

func ParseBricks(lines []string) []Brick {
	var bricks []Brick
	for i, l := range lines {
		b := ParseBrick(l)
		b.N = i
		bricks = append(bricks, b)
	}
	return bricks
}

func ParseBrick(l string) Brick {
	var b Brick
	n, err := fmt.Sscanf(l, "%d,%d,%d~%d,%d,%d",
		&b.S.X, &b.S.Y, &b.S.Z,
		&b.E.X, &b.E.Y, &b.E.Z)
	if n != 6 || err != nil {
		panic("parse error")
	}
	if b.S.X > b.E.X || b.S.Y > b.E.Y || b.S.Z > b.E.Z {
		panic("end position lower than start position")
	}
	return b
}

type V3 struct {
	X, Y, Z int
}

// V2L is a vector of 2 points, labeled
type V2L struct {
	N    int
	A, B int
}

type Brick struct {
	N    int
	S, E V3
}

type Space struct {
	bs []Brick
}

func (b Brick) String() string {
	return fmt.Sprintf("%d: %d,%d,%d~%d,%d,%d (v:%d)", b.N, b.S.X, b.S.Y, b.S.Z, b.E.X, b.E.Y, b.E.Z, b.Volume())
}
func (b Brick) Volume() int {
	return (b.E.X - b.S.X + 1) * (b.E.Y - b.S.Y + 1) * (b.E.Z - b.S.Z + 1)
}

func (b Brick) Contains(p V3) bool {
	return b.S.X >= p.X && b.S.Y >= p.Y && b.S.Z >= p.Z &&
		b.E.X <= p.X && b.E.Y <= p.Y && b.E.Z <= p.Z
}

func (b Brick) OnTopOf(o Brick) bool {
	if b.S.Z-1 != o.E.Z {
		return false
	}
	return b.XYIntersect(o)
}

func (b Brick) XYIntersect(o Brick) bool {
	return IntervalIntersect(b.S.X, b.E.X, o.S.X, o.E.X) &&
		IntervalIntersect(b.S.Y, b.E.Y, o.S.Y, o.E.Y)
}

func IntervalIntersect(alow, ahigh, blow, bhigh int) bool {
	return !(blow > ahigh || alow > bhigh)
}

// CanRest determines if the Brick b can rest in Space s without falling.
func (s Space) CanRest(b Brick) bool {
	if b.S.Z <= 1 {
		return true // brick is supported by the ground
	}
	for _, test := range s.bs {
		if b.N != test.N && b.OnTopOf(test) {
			return true
		}
	}
	return false
}

// finds the bricks that support Brick b
func (s Space) FindSupports(b Brick) []Brick {
	if b.S.Z <= 1 {
		return []Brick{} // supported by the ground
	}
	var rv []Brick
	for _, test := range s.bs {
		if test.N != b.N && b.OnTopOf(test) {
			rv = append(rv, test)
		}
	}
	return rv
}

func (s Space) Find(b Brick) []Brick {
	if b.S.Z <= 1 {
		return []Brick{} // supported by the ground
	}
	var rv []Brick
	for _, test := range s.bs {
		if test.N != b.N && b.OnTopOf(test) {
			rv = append(rv, test)
		}
	}
	return rv
}

// IsSupporter determines if Brick b is supporting any other bricks
func (s Space) IsSupporter(b Brick) bool {
	for _, test := range s.bs {
		if test.N != b.N && test.OnTopOf(b) {
			return true
		}
	}
	return false
}

func SolvePart1(bricks []Brick) {
	s := FillSpace2(bricks)
	r := FindRemovable(s)
	fmt.Println("removable:", r)
}

func SolvePart2(bricks []Brick) {
	s := FillSpace2(bricks)
	d := s.Dependencies()
	mem := make(map[int]int)
	sum := 0
	for _, b := range s.bs {
		chain := s.ChainReaction(mem, b, d)
		if chain > 0 {
			fmt.Println("removing brick", b, "causes chain reaction:", chain)
		}
		sum += chain
	}
	fmt.Println("sum of chain reactions:", sum)
}

func FindRemovable(s Space) int {
	isRequired := make(map[int]bool)
	for _, b := range s.bs {
		supports := s.FindSupports(b)
		if len(supports) == 1 {
			fmt.Println("brick", b.N, "depends on:", supports[0])
			isRequired[supports[0].N] = true
		}
	}

	var r int
	for _, b := range s.bs {
		if !isRequired[b.N] {
			r++
		}
	}
	return r
}
func (s Space) ChainReaction(mem map[int]int, b Brick, d *Dependencies) int {
	if value, found := mem[b.N]; found {
		return value
	}

	removed := make(map[int]bool)
	removed[b.N] = true
	supports := d.Supports[b.N]
	for len(supports) > 0 {
		var nextsupports []int
		for _, sup := range supports {
			if ChainReacts(sup, removed, d) {
				removed[sup] = true
				nextsupports = append(nextsupports, d.Supports[sup]...)
			}

		}
		supports = nextsupports
	}
	n := len(removed) - 1 // do not count the starting brick
	fmt.Println("chain reaction of", b, "removes", n, "bricks")
	mem[b.N] = n
	return n
}

func ChainReacts(id int, removed map[int]bool, d *Dependencies) bool {
	dependsOn := d.DependsOn[id]
	for _, dep := range dependsOn {
		if _, found := removed[dep]; !found {
			// one of this brick's supports has not been removed: brick does not chain react
			return false
		}
	}
	return true
}

func (s Space) ValidateAll() {
	for _, b := range s.bs {
		s.ValidatePointsEmpty(b, b.AllPoints())
	}
}

func FillSpace2(bricks []Brick) Space {
	s := Space{bricks}

	for {
		var changes int
		for i, b := range bricks {
			//supports := s.FindSupports(b)
			//if len(supports) == 0 {
			if !s.CanRest(b) {
				if VERBOSE >= 2 {
					fmt.Println("brick falls one space", b)
				}
				b.S.Z--
				b.E.Z--
				s.bs[i] = b
				changes++
				//s.ValidateAll()
			}
		}
		if changes == 0 {
			break
		}
	}
	s.ValidateAll()
	if VERBOSE >= 2 {
		fmt.Println("XZ")
		fmt.Println("----------------------------------------")
		s.DrawXZ()
		fmt.Println("----------------------------------------")
		fmt.Println()
		fmt.Println("YZ")
		fmt.Println("----------------------------------------")
		s.DrawYZ()
		fmt.Println("----------------------------------------")
	}
	return s
}

func (b Brick) AllPoints() []V3 {
	ps := make([]V3, 0, b.Volume())
	for x := b.S.X; x <= b.E.X; x++ {
		for y := b.S.Y; y <= b.E.Y; y++ {
			for z := b.S.Z; z <= b.E.Z; z++ {
				ps = append(ps, V3{x, y, z})
			}
		}
	}
	return ps
}

func (b Brick) XZPoints() []V2L {
	var ps []V2L
	for x := b.S.X; x <= b.E.X; x++ {
		for z := b.S.Z; z <= b.E.Z; z++ {
			ps = append(ps, V2L{b.N, x, z})
		}
	}
	return ps
}
func (b Brick) YZPoints() []V2L {
	var ps []V2L
	for y := b.S.Y; y <= b.E.Y; y++ {
		for z := b.S.Z; z <= b.E.Z; z++ {
			ps = append(ps, V2L{b.N, y, z})
		}
	}
	return ps
}

func (s Space) DrawXZ() {
	var ps []V2L
	for _, b := range s.bs {
		ps = append(ps, b.XZPoints()...)
	}
	Draw(ps)
}
func (s Space) DrawYZ() {
	var ps []V2L
	for _, b := range s.bs {
		ps = append(ps, b.YZPoints()...)
	}
	Draw(ps)
}

func Draw(ps []V2L) {
	// A is the X dimension of the display coordinate system.
	// B is the Y dimension of the display coordinate system.
	xmax := lib.Reduce(ps, MaxA, 0)
	ymax := lib.Reduce(ps, MaxB, 0)

	for y := ymax; y > 0; y-- {
		for x := 0; x <= xmax; x++ {
			p, n := Contains(ps, x, y)
			if n == 0 {
				fmt.Print(".")
			} else if n == 1 {
				fmt.Print(translateBrick(p.N))
			} else {
				fmt.Print("?")
			}
		}
		fmt.Println()
	}
}

func translateBrick(n int) string {
	if n >= 26 {
		return "+"
		//panic("only supports the example input data")
	}
	b := byte('A' + n)
	return string(b)
}

func Contains(ps []V2L, a int, b int) (V2L, int) {
	var sample V2L
	var n int
	for _, p := range ps {
		if p.A == a && p.B == b {
			sample = p
			n++
		}
	}
	return sample, n
}

func MaxA(t V2L, acc int) int { return max(t.A, acc) }
func MaxB(t V2L, acc int) int { return max(t.B, acc) }

func (s Space) ValidatePointsEmpty(placed Brick, points []V3) {
	//fmt.Println("validating", len(points), "points")
	for _, b := range s.bs {
		if b.N == placed.N {
			continue
		}
		for _, p := range points {
			if b.Contains(p) {
				panic(fmt.Errorf("brick %v contains point %v of %v", b, p, placed))
			}
		}
	}
}

func (s Space) Dependencies() *Dependencies {
	d := NewDeps()
	for _, b := range s.bs {
		supports := s.FindSupports(b)
		// b depends on all supports
		for _, sup := range supports {
			d.Add(b.N, sup.N)
		}
	}
	return d
}

type Dependencies struct {
	// key depends on values
	// to-node -> from-node
	DependsOn map[int][]int
	// key supports values: from-node -> to-node
	Supports map[int][]int
}

func NewDeps() *Dependencies {
	return &Dependencies{
		DependsOn: make(map[int][]int),
		Supports:  make(map[int][]int),
	}
}

func (d *Dependencies) Add(from, to int) {
	add(d.DependsOn, from, to)
	add(d.Supports, to, from)
}

func add(m map[int][]int, key, value int) {
	xs := m[key]
	xs = append(xs, value)
	m[key] = xs
}
