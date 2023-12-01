package main

// https://adventofcode.com/2022/day/12

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math"
	"os"
	"slices"
	"time"
)

func main() {
	Process("aoc22/day12/example.txt", 1)
	Process("aoc22/day12/example.txt", 2)
	Process("aoc22/day12/input.txt", 1)
	Process("aoc22/day12/input.txt", 2)
}

func Process(name string, part int) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	g := NewGrid(lines)

	path := g.Dijkstra(part)
	g.PrintPath(path)

	fmt.Println()
}

func NewGrid(lines []string) *Grid {
	rows := len(lines)
	cols := len(lines[0])
	data := make([]byte, rows*cols)
	for i, row := range lines {
		copy(data[i*cols:], row)
	}
	g := &Grid{
		Rows: rows,
		Cols: cols,
		data: data,
	}
	g.start = g.Find('S')
	g.end = g.Find('E')
	return g
}

func ReadInput(name string) []string {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	var lines []string
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return lines
}

type Grid struct {
	Rows  int
	Cols  int
	data  []byte
	start Pos
	end   Pos
}

func (g *Grid) idx(p Pos) int {
	return p.Y*g.Cols + p.X
}

func (g *Grid) Find(b byte) Pos {
	idx := bytes.IndexByte(g.data, b)
	if idx == -1 {
		panic("invalid grid")
	}
	y := idx / g.Cols
	x := idx - y*g.Cols
	return Pos{x, y}
}

type Pos struct {
	X, Y int
}

func (p Pos) Up() Pos    { return Pos{p.X, p.Y - 1} }
func (p Pos) Down() Pos  { return Pos{p.X, p.Y + 1} }
func (p Pos) Left() Pos  { return Pos{p.X - 1, p.Y} }
func (p Pos) Right() Pos { return Pos{p.X + 1, p.Y} }

func (p Pos) Direction(b Pos) byte {
	if p.Up() == b {
		return '^'
	}
	if p.Down() == b {
		return 'v'
	}
	if p.Left() == b {
		return '<'
	}
	if p.Right() == b {
		return '>'
	}
	panic("invalid direction")
}

func (g *Grid) CanMoveTo(from, to Pos) bool {
	if !g.Valid(to) {
		return false
	}
	a, b := g.Get(from), g.Get(to)
	return (a == 'S' && b == 'a') || // start state
		(a == 'z' && b == 'E') || //
		(a+1 == b) || // one step up
		(a >= b && b != 'E') // same elevation or lower
}

func (g *Grid) Valid(p Pos) bool {
	return p.X >= 0 && p.Y >= 0 && p.X <= g.Cols-1 && p.Y <= g.Rows-1
}
func (g *Grid) Get(p Pos) byte { return g.data[g.idx(p)] }

func (g *Grid) Path(path []Pos) string {
	var b bytes.Buffer

	visited, spath := g.buildMap(path)

	fmt.Fprintln(&b, "path with", len(path)-1, "steps")
	fmt.Fprintln(&b, "path:", spath)
	for y := 0; y < g.Rows; y++ {
		s := y * g.Cols
		e := s + g.Cols
		fmt.Fprintln(&b, string(visited[s:e]))
	}
	return b.String()
}

func (g *Grid) PrintPath(path []Pos) {
	fmt.Println(g.Path(path))
}

func (g *Grid) Neighbors(pos Pos) []Pos {
	var ps []Pos
	if to := pos.Up(); g.CanMoveTo(pos, to) {
		ps = append(ps, to)
	}
	if to := pos.Down(); g.CanMoveTo(pos, to) {
		ps = append(ps, to)
	}
	if to := pos.Left(); g.CanMoveTo(pos, to) {
		ps = append(ps, to)
	}
	if to := pos.Right(); g.CanMoveTo(pos, to) {
		ps = append(ps, to)
	}
	return ps
}

func (g *Grid) buildMap(path []Pos) ([]byte, string) {
	v := bytes.Repeat([]byte("."), len(g.data))
	hist := ""
	for _, p := range path {
		hist += string(g.data[g.idx(p)])
	}
	for i := 0; i < len(path)-1; i++ {
		a := path[i]
		b := path[i+1]
		dir := a.Direction(b)
		v[g.idx(a)] = dir
	}
	return v, hist
}

func (g *Grid) BuildVertices() []*Vertex {
	backing := make([]Vertex, g.Cols*g.Rows)
	vs := make([]*Vertex, g.Cols*g.Rows)
	for x := 0; x < g.Cols; x++ {
		for y := 0; y < g.Rows; y++ {
			p := Pos{x, y}
			idx := g.idx(p)
			backing[idx] = vertex(p)
			backing[idx].Idx = idx
			v := &backing[idx]
			vs[idx] = v
			var neighs []Pos
			if to := p.Up(); g.CanMoveTo(p, to) {
				neighs = append(neighs, to)
			}
			if to := p.Down(); g.CanMoveTo(p, to) {
				neighs = append(neighs, to)
			}
			if to := p.Left(); g.CanMoveTo(p, to) {
				neighs = append(neighs, to)
			}
			if to := p.Right(); g.CanMoveTo(p, to) {
				neighs = append(neighs, to)
			}
			var neighbors []*Vertex
			for _, neigh := range neighs {
				n := &backing[g.idx(neigh)]
				neighbors = append(neighbors, n)
			}
			v.neighbors = neighbors
		}
	}
	return vs
}

func (g *Grid) Dijkstra(part int) []Pos {
	vertices := g.BuildVertices()
	if part == 1 {
		// normal initializer: S is the start point
		vertices[g.idx(g.start)].distance = 0
	} else {
		// part2 initializer: any vertex at elevation 'a' is a possible starting point
		for _, v := range vertices {
			if g.data[v.Idx] == 'a' {
				v.distance = 0
			}
		}
	}
	end := vertices[g.idx(g.end)]

	fmt.Println("# vertices:", len(vertices))
	q := VertexHeap(vertices)
	heap.Init(&q)

	t0 := time.Now()
	for len(q) > 0 {
		u := heap.Pop(&q).(*Vertex)
		u.Visited = true
		// abort early when the end is found
		if u.P == g.end {
			break
		}

		for _, v := range u.neighbors {
			if v.Visited {
				continue
			}
			alt := u.distance + 1
			if alt < v.distance {
				v.distance = alt
				v.previous = u
				heap.Fix(&q, v.Idx)
			}
		}
	}
	fmt.Println("dijkstra", time.Now().Sub(t0))
	current := end

	var path []Pos
	for current != nil {
		path = append(path, current.P)
		current = current.previous
	}
	slices.Reverse(path)
	return path
}

type Vertex struct {
	Idx       int
	P         Pos
	Visited   bool
	distance  int
	previous  *Vertex
	neighbors []*Vertex
}

func vertex(p Pos) Vertex {
	return Vertex{
		P:         p,
		distance:  math.MaxInt,
		previous:  nil,
		neighbors: nil,
	}
}

type VertexHeap []*Vertex

var _ heap.Interface = (*VertexHeap)(nil)

func (h VertexHeap) Len() int           { return len(h) }
func (h VertexHeap) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h VertexHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[j].Idx = j
	h[i].Idx = i
}

func (h *VertexHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	vert := x.(*Vertex)
	vert.Idx = len(*h)
	*h = append(*h, vert)
}
func (h *VertexHeap) Pop() any {
	old := *h
	n := len(old)
	vert := old[n-1]
	*h = old[0 : n-1]
	vert.Idx = -1
	return vert
}
