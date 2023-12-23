package main

// https://adventofcode.com/2023/day/23

import (
	"bytes"
	"fmt"
	"time"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/rowcol"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day23/example.txt") // 84
	ProcessPart1("aoc23/day23/input.txt")   // 2278

	ProcessPart2("aoc23/day23/example.txt") // 154
	ProcessPart2("aoc23/day23/input.txt")   // 6734
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	g := ParseGraph(lines)
	SolvePart1(g)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)

	start := time.Now()
	grid := ParseGraph(lines)
	graph := ExtractGraph(grid)

	// https://en.wikipedia.org/wiki/Longest_path_problem
	graph.SolvePart2()

	fmt.Println("t:", time.Since(start))

	fmt.Println()
}

////////////////////////////////////////////////////////////

var Directions = map[byte]rowcol.Direction{
	'>': rowcol.Right,
	'<': rowcol.Left,
	'^': rowcol.Up,
	'v': rowcol.Down,
}

type Grid struct {
	rowcol.Grid[byte]

	Start, End rowcol.Pos
}

func ParseGraph(lines []string) *Grid {
	g := &Grid{Grid: rowcol.NewByteGridFromStrings(lines)}
	g.Start = g.FindStart()
	g.End = g.FindEnd()
	return g
}

func (g *Grid) FindStart() rowcol.Pos {
	col := bytes.IndexByte(g.Grid.Data[0], '.')
	if col == -1 {
		panic("start not found")
	}
	return rowcol.Pos{Row: 0, Col: col}
}

func (g *Grid) FindEnd() rowcol.Pos {
	row := g.Rows() - 1
	col := bytes.IndexByte(g.Grid.Data[row], '.')
	if col == -1 {
		panic("end not found")
	}
	return rowcol.Pos{Row: row, Col: col}
}

func SolvePart1(g *Grid) {
	visited := make(map[rowcol.Pos]bool)
	visited[g.Start] = true
	d := AllSteps(g, g.Start, visited)
	fmt.Println("Distance:", d)
}

func AllSteps(g *Grid, current rowcol.Pos, visited map[rowcol.Pos]bool) int {
	maxd := 0
	for _, dir := range rowcol.Directions {
		next := current.AddDir(dir)
		d := OneStep(g, next, visited)
		maxd = max(maxd, d)
	}
	return maxd
}

func OneStep(g *Grid, pos rowcol.Pos, visited map[rowcol.Pos]bool) int {
	if !g.IsValidPosition(pos.Row, pos.Col) {
		return 0
	}
	if visited[pos] {
		return 0
	}
	if pos == g.End {
		return 1
	}
	v := g.Get(pos.Row, pos.Col)
	if v == '#' {
		return 0
	}
	visited[pos] = true
	defer func() { delete(visited, pos) }()
	if v == '.' {
		return AllSteps(g, pos, visited) + 1
	}
	dir, ok := Directions[v]
	if !ok {
		panic(fmt.Errorf("invalid field: %c", v))
	}
	return OneStep(g, pos.AddDir(dir), visited) + 1
}

// each position that has two or more paths leading away is a vertex, all other positions connect these vertices
func ExtractGraph(g *Grid) *Graph {
	graph := NewGraph(g)
	graph.AddVertex(g.Start)
	graph.AddVertex(g.End)
	visited := make(map[rowcol.Pos]bool)
	graph.Follow(g.Start, g.Start, g.Start.AddDir(rowcol.Down), 1, visited)
	return graph
}

type Graph struct {
	G            *Grid
	Nodes        map[rowcol.Pos]*Node
	Combinations int
}

func NewGraph(g *Grid) *Graph {
	return &Graph{
		G:     g,
		Nodes: make(map[rowcol.Pos]*Node),
	}
}

func (g *Graph) Follow(start, previous, current rowcol.Pos, dist int, visited map[rowcol.Pos]bool) {
	if current == g.G.End {
		g.AddEdge(start, current, dist)
		return
	}

	if ok, next := g.G.IsOneWay(previous, current); ok {
		g.Follow(start, current, next, dist+1, visited)
	} else {
		g.AddVertex(current)

		g.AddEdge(start, current, dist)
		g.AddEdge(current, start, dist)

		if visited[current] {
			return
		}
		visited[current] = true
		g.FollowVertex(current, previous, visited)
	}
}

func (g *Graph) FollowVertex(current rowcol.Pos, previous rowcol.Pos, visited map[rowcol.Pos]bool) {
	for _, dir := range rowcol.Directions {
		next := current.AddDir(dir)
		if next != previous && g.G.IsValidMove(next) {
			g.Follow(current, current, next, 1, visited)
		}
	}
}

func (g *Graph) AddVertex(v rowcol.Pos) {
	_, ok := g.Nodes[v]
	if !ok {
		g.Nodes[v] = &Node{P: v}
	}
}

func (g *Graph) AddEdge(a, b rowcol.Pos, dist int) {
	if VERBOSE >= 2 {
		fmt.Printf("%q -> %q [weight=%d]\n", a, b, dist)
	}

	n, found := g.Nodes[a]
	if !found {
		n = &Node{P: a}
		g.Nodes[a] = n
	}
	for i, neigh := range n.Neighbors {
		if neigh == b {
			if dist != n.Distances[i] {
				panic("invalid state")
			}
			return
		}
	}

	n.Neighbors = append(n.Neighbors, b)
	n.Distances = append(n.Distances, dist)
}

// a position is a one-way if there are only 2 possible ways to go to
func (g *Grid) IsOneWay(previous, p rowcol.Pos) (bool, rowcol.Pos) {
	n := 0
	var next rowcol.Pos
	for _, dir := range rowcol.Directions {
		cand := p.AddDir(dir)
		if g.IsValidMove(cand) {
			n++
			if cand != previous {
				next = cand
			}
		}
	}
	if n > 2 {
		return false, rowcol.Pos{}
	}
	if n < 2 {
		panic(fmt.Errorf("invalid state at position: %v", p))
	}
	return true, next
}

func (g *Grid) IsValidMove(p rowcol.Pos) bool {
	return g.IsValidPosition(p.Row, p.Col) && g.Get(p.Row, p.Col) != '#'
}

type Node struct {
	P         rowcol.Pos
	Neighbors []rowcol.Pos
	Distances []int
}

func (g *Graph) SolvePart2() {
	visited := make(map[rowcol.Pos]bool)
	startNode := g.GetNode(g.G.Start)
	dist, ok := g.FollowP2(startNode, visited)
	if !ok {
		fmt.Println("End not found !?")
	}
	fmt.Println("combinations:", g.Combinations)
	fmt.Println("distance:", dist)
}

func (g *Graph) FollowP2(current *Node, visited map[rowcol.Pos]bool) (int, bool) {
	g.Combinations++
	if g.G.End == current.P {
		return 0, true
	}

	var m int
	valid := false
	visited[current.P] = true

	for i, neigh := range current.Neighbors {
		if visited[neigh] {
			continue
		}
		n := g.GetNode(neigh)
		dist, ok := g.FollowP2(n, visited)
		if ok {
			dist += current.Distances[i]
			m = max(m, dist)
			valid = true
		}
	}
	delete(visited, current.P)
	return m, valid
}

func (g *Graph) GetNode(p rowcol.Pos) *Node {
	n := g.Nodes[p]
	if n == nil {
		panic(fmt.Errorf("no node found for key: %v", p))
	}
	return n
}
