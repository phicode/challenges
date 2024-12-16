package main

// https://adventofcode.com/2024/day/12

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	//lib.Timed("Part 1", ProcessPart1, "aoc24/day12/example.txt")
	//lib.Timed("Part 1", ProcessPart1, "aoc24/day12/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day12/example.txt")
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day12/input.txt")

	//lib.Profile(1, "day12.pprof", "Part 2", ProcessPart2, "aoc24/day12/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart1(input)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	input := ReadAndParseInput(name)
	result := SolvePart2(input)
	fmt.Println("Result:", result)
}

////////////////////////////////////////////////////////////

type Input struct {
	grid rowcol.Grid[byte]
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	return Input{rowcol.NewByteGridFromStrings(lines)}
}

////////////////////////////////////////////////////////////

type State struct {
	grid    rowcol.Grid[byte]
	visited rowcol.Grid[bool]
	scratch rowcol.Grid[byte]
}

func (s *State) follow(p rowcol.Pos) (byte, []rowcol.Pos) {
	plant := s.grid.GetPos(p)
	s.visited.SetPos(p, true)
	var ps []rowcol.Pos
	ps = append(ps, p)
	return plant, s.followPlantDirs(plant, ps, p)
}

func (s *State) followPlant(plant byte, ps []rowcol.Pos, p rowcol.Pos) []rowcol.Pos {
	if !s.grid.IsValidPos(p) || s.grid.GetPos(p) != plant {
		return ps
	}
	if s.visited.GetPos(p) {
		return ps
	}
	s.visited.SetPos(p, true)
	ps = append(ps, p)
	return s.followPlantDirs(plant, ps, p)
}

func (s *State) followPlantDirs(plant byte, ps []rowcol.Pos, p rowcol.Pos) []rowcol.Pos {
	ps = s.followPlant(plant, ps, p.AddDir(rowcol.Up))
	ps = s.followPlant(plant, ps, p.AddDir(rowcol.Down))
	ps = s.followPlant(plant, ps, p.AddDir(rowcol.Left))
	ps = s.followPlant(plant, ps, p.AddDir(rowcol.Right))
	return ps
}

func SolvePart1(input Input) int {
	s := State{grid: input.grid, visited: rowcol.NewGrid[bool](input.grid.Size())}
	total := 0
	for pos := range s.visited.PosIterator() {
		if s.visited.GetPos(pos) {
			continue
		}
		plant, ps := s.follow(pos)
		peri := perimeter(ps, s.grid, plant)
		//fmt.Printf("Plant: %c ; area: %d, perimeter: %d\n", plant, len(ps), peri)
		total += peri * len(ps)
	}
	return total
}

func perimeter(ps []rowcol.Pos, grid rowcol.Grid[byte], plant byte) int {
	peri := 0
	for _, p := range ps {
		peri += perimeterContribution(p, grid, plant)
	}
	return peri
}

func perimeterContribution(pos rowcol.Pos, grid rowcol.Grid[byte], plant byte) int {
	sameNeighbors := 0
	for _, dir := range rowcol.Directions {
		test := pos.AddDir(dir)
		if grid.IsValidPos(test) && grid.GetPos(test) == plant {
			sameNeighbors++
		}
	}
	// 4 same neighbors => plant is fully surounded => no perimeter contribution
	// 3 plant has 1 neighbors on 3 sided
	return 4 - sameNeighbors
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	s := State{grid: input.grid, visited: rowcol.NewGrid[bool](input.grid.Size())}
	total := 0
	for pos := range s.visited.PosIterator() {
		if s.visited.GetPos(pos) {
			continue
		}
		plant, area, edges := s.followp2(pos)
		sides := edges //- 1
		fmt.Printf("Plant: %c ; area: %d, sides: %d\n", plant, area, sides)
		total += area * sides
	}
	return total
}

func (s *State) followp2(p rowcol.Pos) (byte, int, int) {
	plant := s.grid.GetPos(p)
	if plant == 'C' {
		fmt.Println("debug")
	}
	s.visited.SetPos(p, true)
	a, e := s.followPlantDirsp2(plant, p)
	a++
	e += NumEdges(s.grid, p)
	return plant, a, e
}

func (s *State) followPlantp2(plant byte, p rowcol.Pos) (int, int) {
	if !s.grid.IsValidPos(p) || s.grid.GetPos(p) != plant {
		return 0, 0
	}
	if s.visited.GetPos(p) {
		return 0, 0
	}
	s.visited.SetPos(p, true)
	area, edges := s.followPlantDirsp2(plant, p)
	area++
	edges += NumEdges(s.grid, p)
	return area, edges
}

func (s *State) followPlantDirsp2(plant byte, p rowcol.Pos) (int, int) {
	a1, e1 := s.followPlantp2(plant, p.AddDir(rowcol.Up))
	a2, e2 := s.followPlantp2(plant, p.AddDir(rowcol.Down))
	a3, e3 := s.followPlantp2(plant, p.AddDir(rowcol.Left))
	a4, e4 := s.followPlantp2(plant, p.AddDir(rowcol.Right))
	return a1 + a2 + a3 + a4, e1 + e2 + e3 + e4
}

////////////////////////////////////////////////////////////

type Poly struct {
	// len(coordinates) = area
	coordinates []rowcol.Pos
	vertices    []rowcol.Pos
	sides       int
	plant       byte
}

func (p Poly) Area() int { return len(p.coordinates) }
func SolvePart2old(input Input) int {
	// build polygons
	// count number of sides of the polygon
	// find enclaves

	s := State{
		grid:    input.grid,
		visited: rowcol.NewGrid[bool](input.grid.Size()),
		scratch: rowcol.NewGrid[byte](input.grid.Size()),
	}

	var all []*Poly
	for pos := range s.visited.PosIterator() {
		if s.visited.GetPos(pos) {
			continue
		}
		plant, ps := s.follow(pos)
		vecs := buildPolygon(s, ps, plant)
		sides := countsides(vecs)
		fmt.Printf("%c vector length: %d, sides: %d\n", plant, len(vecs), sides)

		all = append(all, &Poly{
			coordinates: ps,
			vertices:    vecs,
			sides:       sides,
			plant:       plant,
		})
	}

	for _, poly := range all {
		if enclosed, p := isEnclosed(poly, s); enclosed {
			other := findShape(p, all)
			other.sides += poly.sides
		}
	}

	total := 0
	for _, poly := range all {
		total += poly.Area() * poly.sides
	}
	return total
}

func findShape(p rowcol.Pos, all []*Poly) *Poly {
	for _, poly := range all {
		for _, c := range poly.coordinates {
			if p == c {
				return poly
			}
		}
	}
	panic("not found")
}

func isEnclosed(p *Poly, s State) (bool, rowcol.Pos) {
	var neighbor byte
	var neighpos rowcol.Pos
	for _, c := range p.coordinates {
		for _, dir := range rowcol.Directions {
			pos := c.AddDir(dir)
			if !s.grid.IsValidPos(pos) {
				return false, rowcol.Pos{} // at the edge of the map -> not enclosed
			}
			v := s.grid.GetPos(pos)
			if v == p.plant {
				continue
			}
			if neighbor == 0 {
				neighbor = v
				neighpos = pos
				continue
			}
			if neighbor != v {
				return false, rowcol.Pos{} // different neighbors
			}
		}
	}
	assert.True(neighbor != 0)
	fmt.Printf("%c is enclosed by %c\n", p.plant, neighbor)
	return true, neighpos
}

func countsides(vecs []rowcol.Pos) int {
	a, b := vecs[0], vecs[1]
	dir := b.Sub(a)
	sides := 1
	for i := 1; i < len(vecs)-1; i++ {
		a, b = vecs[i], vecs[i+1]
		newdir := b.Sub(a)
		if newdir != dir {
			dir = newdir
			sides++
		}
	}
	return sides
}

func buildPolygon(s State, ps []rowcol.Pos, plant byte) []rowcol.Pos {
	grid := s.scratch
	grid.Reset(0)
	for _, p := range ps {
		grid.SetPos(p, plant)
	}

	_min := rowcol.MinPos(ps)
	dir := rowcol.Right
	var vecs []rowcol.Pos
	vecs = append(vecs, _min)

	cur := _min
	for {
		if valueToRight(grid, cur, dir.Left(), plant) {
			dir = dir.Left()
			cur = cur.AddDir(dir)
			vecs = append(vecs, cur)
		} else if valueToRight(grid, cur, dir, plant) {
			cur = cur.AddDir(dir)
			vecs = append(vecs, cur)
		} else if valueToRight(grid, cur, dir.Right(), plant) {
			dir = dir.Right()
			cur = cur.AddDir(dir)
			vecs = append(vecs, cur)
		} else {
			panic("invalid state")
		}
		if cur == _min {
			break
		}
	}
	return vecs
}

func valueToRight(grid rowcol.Grid[byte], p rowcol.Pos, dir rowcol.Direction, v byte) bool {
	topleft := TopLeftCorner(p, dir)
	return grid.IsValidPos(topleft) &&
		grid.GetPos(topleft) == v
}

// given a position p and a direction dir, returns the index of the cell which lies to the "right" the movement vector
func TopLeftCorner(p rowcol.Pos, dir rowcol.Direction) rowcol.Pos {
	if dir == rowcol.Left {
		return p.AddS(-1, -1)
	}
	if dir == rowcol.Right {
		return p
	}
	if dir == rowcol.Up {
		return p.AddS(-1, 0)
	}
	if dir == rowcol.Down {
		return p.AddS(0, -1)
	}
	panic("invalid direction")
}

// see notes
// a is the diagonal neighbor, b and c are the horizontal and vertial neighbors
func IsEdge(a, b, c bool) bool {
	return !((!a && !b && c) ||
		(!a && b && !c) ||
		(a && b && c))
}

func NumEdges(grid rowcol.Grid[byte], p rowcol.Pos) int {
	// abc
	// dXe
	// fgh
	v := grid.GetPos(p)
	pa := p.AddDir(rowcol.UpLeft)
	pb := p.AddDir(rowcol.Up)
	pc := p.AddDir(rowcol.UpRight)
	pd := p.AddDir(rowcol.Left)
	pe := p.AddDir(rowcol.Right)
	pf := p.AddDir(rowcol.DownLeft)
	pg := p.AddDir(rowcol.Down)
	ph := p.AddDir(rowcol.DownRight)
	edges := 0
	a := IsSame(grid, pa, v)
	b := IsSame(grid, pb, v)
	c := IsSame(grid, pc, v)
	d := IsSame(grid, pd, v)
	e := IsSame(grid, pe, v)
	f := IsSame(grid, pf, v)
	g := IsSame(grid, pg, v)
	h := IsSame(grid, ph, v)
	// top left quadrant
	//if grid.IsValidPos(pa) {
	if IsEdge(a, b, d) {
		edges++
	}
	//} else {
	//	edges++
	//}
	// top right quadrant
	//if grid.IsValidPos(pc) {
	if IsEdge(c, b, e) {
		edges++
	}
	//} else {
	//	edges++
	//}
	// top left quadrant
	//if grid.IsValidPos(pf) {
	if IsEdge(f, d, g) {
		edges++
	}
	//} else {
	//	edges++
	//}
	// top left quadrant
	//if grid.IsValidPos(ph) {
	if IsEdge(h, e, g) {
		edges++
	}
	//} else {
	//	edges++
	//}
	return edges
}

func IsSame(g rowcol.Grid[byte], p rowcol.Pos, v byte) bool {
	return g.IsValidPos(p) && g.GetPos(p) == v
}
