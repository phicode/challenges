package main

// https://adventofcode.com/2023/day/10

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 2

func main() {
	// Distance: 4
	ProcessPart1("aoc23/day10/example.txt")
	// Distance: 8
	ProcessPart1("aoc23/day10/example2.txt")

	ProcessPart1("aoc23/day10/input.txt")

	// Enclosed: 4
	ProcessPart2("aoc23/day10/example1part2.txt")
	// Enclosed: 4
	ProcessPart2("aoc23/day10/example2part2.txt")
	// Enclosed: 8
	ProcessPart2("aoc23/day10/example3part2.txt")
	// Enclosed: 10
	ProcessPart2("aoc23/day10/example4part2.txt")
	ProcessPart2("aoc23/day10/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	g := ParseGrid(lines)
	dist, _ := SolvePart1(g)
	fmt.Println("Distance:", dist)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	g := ParseGrid(lines)
	dist, cleaned := SolvePart1(g)
	fmt.Println("Distance:", dist)
	enclosed := SolvePart2(cleaned)
	fmt.Println("Enclosed:", enclosed)
	cleaned.Print()

	fmt.Println()
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

type Tile byte

func (t Tile) CanMove(direction Pos) bool {
	if direction == DirectionNorth { // north / up
		return t == Vertical || t == BendNE || t == BendNW
	}
	if direction == DirectionSouth { // south / down
		return t == Vertical || t == BendSE || t == BendSW
	}
	if direction == DirectionEast { // east / right
		return t == Horizontal || t == BendNE || t == BendSE
	}
	if direction == DirectionWest { // west / left
		return t == Horizontal || t == BendNW || t == BendSW
	}
	return false
}

func (t Tile) MoveA() Pos {
	switch t {
	case Vertical:
		return DirectionNorth
	case Horizontal:
		return DirectionEast
	case BendNE:
		return DirectionNorth
	case BendNW:
		return DirectionNorth
	case BendSW:
		return DirectionSouth
	case BendSE:
		return DirectionSouth
	case Ground:
		return Pos{}
	case Start:
		return Pos{}
	}
	panic("invalid tile")
}
func (t Tile) MoveB() Pos {
	switch t {
	case Vertical:
		return DirectionSouth
	case Horizontal:
		return DirectionWest
	case BendNE:
		return DirectionEast
	case BendNW:
		return DirectionWest
	case BendSW:
		return DirectionWest
	case BendSE:
		return DirectionEast
	case Ground:
		return Pos{}
	case Start:
		return Pos{}
	}
	panic("invalid tile")
}

func (t Tile) MatchesMove(a, b Pos) bool {
	return (a == t.MoveA() && b == t.MoveB()) ||
		(a == t.MoveB() && b == t.MoveA())
}

const (
	// | is a vertical pipe connecting north and south.
	Vertical Tile = '|'
	// - is a horizontal pipe connecting east and west.
	Horizontal Tile = '-'
	// L is a 90-degree bend connecting north and east.
	BendNE Tile = 'L'
	// J is a 90-degree bend connecting north and west.
	BendNW Tile = 'J'
	// 7 is a 90-degree bend connecting south and west.
	BendSW Tile = '7'
	// F is a 90-degree bend connecting south and east.
	BendSE Tile = 'F'
	// . is ground; there is no pipe in this tile.
	Ground Tile = '.'
	// S is the starting position of the animal; there is a pipe on this tile, but your sketch doesn't show what shape the pipe has.
	Start Tile = 'S'
)

var Tiles = []Tile{
	Vertical, Horizontal, BendNE, BendNW, BendSW, BendSE, Ground, Start,
}

type Pos struct {
	X, Y int
}

func (p Pos) Add(b Pos) Pos    { return Pos{p.X + b.X, p.Y + b.Y} }
func (p Pos) Sub(b Pos) Pos    { return Pos{p.X - b.X, p.Y - b.Y} }
func (p Pos) Reverse() Pos     { return Pos{-p.X, -p.Y} }
func (p Pos) NormalRight() Pos { return Pos{p.Y, -p.X} }
func (p Pos) NormalLeft() Pos  { return Pos{-p.Y, p.X} }

var (
	DirectionNorth = Pos{0, -1}
	DirectionSouth = Pos{0, 1}
	DirectionEast  = Pos{1, 0}
	DirectionWest  = Pos{-1, 0}
)

var Directions = []Pos{
	DirectionNorth,
	DirectionSouth,
	DirectionEast,
	DirectionWest,
}

type Grid struct {
	Start Pos
	tiles [][]Tile
}

func ParseGrid(lines []string) *Grid {
	g := Grid{}
	g.tiles = make([][]Tile, len(lines))
	for i, line := range lines {
		g.tiles[i] = []Tile(line)
	}
	return &g
}

func (g *Grid) FindStart() Pos {
	for y, row := range g.tiles {
		for x, col := range row {
			if col == Start {
				return Pos{x, y}
			}
		}
	}
	panic("start position not found")
}

func (g *Grid) FindStartParameters(start Pos) (Tile, []Pos) {
	var connections []Pos
	for _, dir := range Directions {
		to := start.Add(dir)
		if g.CanMove(to, dir.Reverse()) {
			connections = append(connections, to)
		}
	}
	if len(connections) != 2 || connections[0] == connections[1] {
		panic("invalid state")
	}
	return TranslateTile(start, connections[0], connections[1]), connections
}

func (g *Grid) Rows() int    { return len(g.tiles) }
func (g *Grid) Columns() int { return len(g.tiles[0]) }
func (g *Grid) IsValid(p Pos) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < g.Columns() && p.Y < g.Rows()
}

func (g *Grid) CanMove(p, direction Pos) bool {
	if !g.IsValid(p) {
		return false
	}
	tile := g.tiles[p.Y][p.X]
	return tile.CanMove(direction)
}

func (g *Grid) CopyDimensions() *Grid {
	w, h := g.Columns(), g.Rows()
	clone := Grid{
		Start: g.Start,
		tiles: make([][]Tile, h),
	}
	for y := 0; y < h; y++ {
		clone.tiles[y] = make([]Tile, w)
		for x := 0; x < w; x++ {
			clone.tiles[y][x] = '.'
		}
	}
	return &clone
}

func (g *Grid) Set(p Pos, tile Tile) { g.tiles[p.Y][p.X] = tile }
func (g *Grid) Get(p Pos) Tile       { return g.tiles[p.Y][p.X] }

func (g *Grid) Print() {
	for _, row := range g.tiles {
		for _, t := range row {
			fmt.Printf("%c", t)
		}
		fmt.Println()
	}
}

func (g *Grid) CountMarked() (int, int) {
	w, h := g.Columns(), g.Rows()
	var ground, marked int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v := g.tiles[y][x]
			if v == '.' {
				ground++
			}
			if v == 'O' {
				marked++
			}
		}
	}
	return ground, marked
}

////////////////////////////////////////////////////////////
// Part 1

func SolvePart1(g *Grid) (int, *Grid) {
	g.Start = g.FindStart()
	cleaned := g.CopyDimensions()
	startTile, connections := g.FindStartParameters(g.Start)
	cleaned.Set(g.Start, startTile)
	cleaned.Set(connections[0], g.Get(connections[0]))
	cleaned.Set(connections[1], g.Get(connections[1]))
	distance := 1
	previous := []Pos{g.Start, g.Start}
	for connections[0] != connections[1] {
		previous, connections = Advance(g, previous, connections)
		cleaned.Set(connections[0], g.Get(connections[0]))
		cleaned.Set(connections[1], g.Get(connections[1]))
		distance++
	}
	return distance, cleaned
}

func Advance(g *Grid, previous, connections []Pos) ([]Pos, []Pos) {
	var rv [2]Pos
	rv[0] = AdvanceOne(g, previous[0], connections[0])
	rv[1] = AdvanceOne(g, previous[1], connections[1])
	return connections, rv[:]
}

func AdvanceOne(g *Grid, prev, pos Pos) Pos {
	for _, dir := range Directions {
		if g.CanMove(pos, dir) && pos.Add(dir) != prev {
			return pos.Add(dir)
		}
	}
	panic("no movement possibility found")
}

func TranslateTile(a, b, c Pos) Tile {
	moveA := b.Sub(a)
	moveB := c.Sub(a)
	for _, tile := range Tiles {
		if tile.MatchesMove(moveA, moveB) {
			return tile
		}
	}
	panic("could not translate tile movement")
}

////////////////////////////////////////////////////////////
// Part 2

func SolvePart2(g *Grid) int {
	//markoutside(g, Pos{0, 0})

	_, connections := g.FindStartParameters(g.Start)
	previous := g.Start
	current := connections[0]
	for current != g.Start {
		a, b := OutsidePositions(previous, current)
		markoutside(g, a)
		markoutside(g, b)
		next := AdvanceOne(g, previous, current)
		previous, current = current, next
	}

	a, b := g.CountMarked()
	return min(a, b)
}

func OutsidePositions(from Pos, to Pos) (Pos, Pos) {
	normal := to.Sub(from).NormalLeft()
	return from.Add(normal), to.Add(normal)
}

func markoutside(g *Grid, pos Pos) {
	if g.Get(pos) != '.' {
		return
	}
	g.Set(pos, 'O')
	for _, dir := range Directions {
		to := pos.Add(dir)
		if g.IsValid(to) && g.Get(to) == '.' {
			markoutside(g, to)
		}
	}
}

//func Intersections(g *Grid, start Pos, end Pos) int {
//	intr := 0
//	for start != end {
//		if g.Get(start) != Ground {
//			intr++
//		}
//		start = start.Add(DirectionEast)
//	}
//	return intr
//}
