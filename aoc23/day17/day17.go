package main

// https://adventofcode.com/2023/day/17

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/rowcol"
)

var VERBOSE = 0

func main() {
	ProcessPart1("aoc23/day17/example.txt") // 102
	ProcessPart1("aoc23/day17/input.txt")   // 886

	ProcessPart2("aoc23/day17/example.txt")  // 94
	ProcessPart2("aoc23/day17/example2.txt") // 71
	ProcessPart2("aoc23/day17/input.txt")    // 1055
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	Process(name, 1, 3)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	Process(name, 4, 10)
}

func Process(name string, minMove, maxMove int) {
	lines := lib.ReadLines(name)
	g := ParseInput(lines)
	t0 := time.Now()
	nodes := BuildNodes(g, minMove, maxMove)
	Follow(g, nodes)
	fmt.Println(time.Since(t0))
	fmt.Println()
}

////////////////////////////////////////////////////////////

type Grid struct {
	rowcol.Grid[int]
}

func ParseInput(lines []string) *Grid {
	rows := len(lines)
	cols := len(lines[0])
	g := rowcol.NewGrid[int](rows, cols)
	for rowIdx, l := range lines {
		row := make([]int, cols)
		for i, v := range l {
			if v < '0' || v > '9' {
				panic("invalid input")
			}
			row[i] = int(v - '0')
		}
		g.SetRow(rowIdx, row)
	}
	return &Grid{g}
}

type Direction struct {
	R, C int
}

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

type Pos struct {
	R, C int
}

func (p Pos) Add(dir Direction) Pos {
	return Pos{
		R: p.R + dir.R,
		C: p.C + dir.C,
	}
}

func (p Pos) IsZero() bool   { return p.R == 0 && p.C == 0 }
func (p Pos) String() string { return fmt.Sprintf("(%d,%d)", p.C, p.R) }

var (
	Left  = Direction{R: 0, C: -1}
	Right = Direction{R: 0, C: +1}
	Up    = Direction{R: -1, C: 0}
	Down  = Direction{R: +1, C: 0}
	Stand = Direction{R: 0, C: 0}
)

func (d Direction) Left() Direction {
	return Direction{-d.C, d.R}
}
func (d Direction) Right() Direction {
	return Direction{d.C, d.R}
}
func (d Direction) Reverse() Direction {
	return Direction{-d.R, -d.C}
}

var Directions = []Direction{Left, Right, Up, Down}

type NodeKey struct {
	Pos Pos
	Dir Direction // how we got here
}

type Node struct {
	Visited   bool
	Neighbors []NodeKey
	// value: cost
	NeighborCost map[NodeKey]int
}

func (n *Node) AddNeighbor(pos NodeKey, cost int) {
	if slices.Contains(n.Neighbors, pos) {
		panic("double add")
	}
	n.Neighbors = append(n.Neighbors, pos)
	n.NeighborCost[pos] = cost
}

func BuildNodes(g *Grid, minMove, maxMove int) map[NodeKey]*Node {
	nodes := make(map[NodeKey]*Node)
	start := NodeKey{Pos{0, 0}, Stand}
	nodes[start] = NewNode()
	TestMovesAndFollow(g, nodes, start, minMove, maxMove)
	if VERBOSE >= 1 {
		fmt.Println("part2 nodes:", len(nodes))
	}
	return nodes
}

func NewNode() *Node {
	return &Node{NeighborCost: make(map[NodeKey]int)}
}

func TestMovesAndFollow(g *Grid, nodes map[NodeKey]*Node, vertex NodeKey, minMove, maxMove int) {
	currentNode := nodes[vertex]
	if currentNode == nil {
		panic(fmt.Errorf("node not found: %v", vertex))
	}
	if currentNode.Visited {
		//fmt.Println("node already visited, aborting")
		return
	}
	currentNode.Visited = true
	for _, dir := range Directions {
		if dir == vertex.Dir || dir == vertex.Dir.Reverse() {
			continue
		}
		for i := minMove; i <= maxMove; i++ {
			targetVertex := NodeKey{Pos{vertex.Pos.R + i*dir.R, vertex.Pos.C + i*dir.C}, dir}
			if g.IsValidPosition(targetVertex.Pos.R, targetVertex.Pos.C) {
				node, found := nodes[targetVertex]
				if !found {
					node = NewNode()
					nodes[targetVertex] = node
				}
				TestMovesAndFollow(g, nodes, targetVertex, minMove, maxMove)
				currentNode.AddNeighbor(targetVertex, CalcCost(g, targetVertex.Pos, dir, i))
			}
		}
	}
}

func CalcCost(g *Grid, pos Pos, dir Direction, l int) int {
	var cost int
	reverse := dir.Reverse()
	for i := 0; i < l; i++ {
		cost += g.Get(pos.R, pos.C)
		pos = pos.Add(reverse)
	}
	return cost
}

func Follow(g *Grid, nodes map[NodeKey]*Node) {
	nodeKeys := lib.MapKeys(nodes)
	start := func(k NodeKey) bool { return k.Pos.IsZero() }
	neighbors := func(k NodeKey) []NodeKey {
		return nodes[k].Neighbors
	}
	cost := func(u NodeKey, v NodeKey) int {
		c := nodes[u].NeighborCost[v]
		if c <= 0 {
			panic("invalid state")
		}
		return c
	}
	paths := lib.DijkstraWithCost(nodeKeys, start, neighbors, cost)
	heatLoss := VisitPart2EndNodes(g, nodes, paths)
	fmt.Println("heat loss:", heatLoss)
}

func VisitPart2EndNodes(g *Grid, nodes map[NodeKey]*Node, paths map[NodeKey]*lib.Node[NodeKey]) int {
	rows, cols := g.Rows(), g.Columns()
	minHeatLoss := math.MaxInt
	for k := range nodes {
		if k.Pos.R == rows-1 && k.Pos.C == cols-1 { // end node
			node := paths[k]
			if VERBOSE >= 2 {
				fmt.Println("End node found", node.Distance)
			}
			if node.Distance < minHeatLoss {
				minHeatLoss = node.Distance
			}
			if VERBOSE >= 2 {
				DescribePath(node, nodes)
			}
		}
	}
	return minHeatLoss
}

func DescribePath(end *lib.Node[NodeKey], nodes map[NodeKey]*Node) {
	if end.Prev == nil {
		fmt.Println("found unvisited end node")
		return
	}
	var path []*lib.Node[NodeKey]
	path = append(path, end)
	current := end
	for !current.Value.Pos.IsZero() { // start node
		current = current.Prev
		path = append(path, current)
	}
	slices.Reverse(path)
	fmt.Println("path:")
	var totalCost int
	for i := 1; i < len(path); i++ {
		prev, current := path[i-1], path[i]
		edgeCost := nodes[prev.Value].NeighborCost[current.Value]
		totalCost += edgeCost
		fmt.Println(prev.Value, " -> ", current.Value, ", cost:", edgeCost, ", sum:", totalCost)
	}
	fmt.Println()
}
