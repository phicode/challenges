package main

// https://adventofcode.com/2023/day/17

import (
	"fmt"
	"math"
	"slices"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day17/example.txt")
	ProcessPart1("aoc23/day17/input.txt")
	//
	//ProcessPart2("aoc23/day17/example.txt")
	//ProcessPart2("aoc23/day17/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	g := ParseInput(lines)
	Follow(g)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

type Grid struct {
	lib.Grid[int]
}

func ParseInput(lines []string) *Grid {
	rows := len(lines)
	cols := len(lines[0])
	g := lib.NewGrid[int](rows, cols)
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

type HistoryNodeKey struct {
	P     Pos
	Moves [3]Direction
}

func (h *HistoryNodeKey) CanMove(d Direction) bool {
	a, b, c := h.Moves[0], h.Moves[1], h.Moves[2]
	if a == d && a == b && b == c {
		return false
	}
	return true
}
func (h *HistoryNodeKey) CanMoveGrid(g *Grid, dir Direction) bool {
	t := h.P.Add(dir)
	return h.Moves[2].Reverse() != dir && g.IsValidPosition(t.R, t.C) && h.CanMove(dir)
}

func (h *HistoryNodeKey) ApplyMove(dir Direction) HistoryNodeKey {
	hnk := HistoryNodeKey{
		P: h.P.Add(dir),
	}
	hnk.Moves[0] = h.Moves[1]
	hnk.Moves[1] = h.Moves[2]
	hnk.Moves[2] = dir
	return hnk
}

type HistoryNode struct {
	HistoryNodeKey

	Neighbors []HistoryNodeKey
	HeatLoss  int
}

func Follow(g *Grid) {
	nodesMap := BuildNodes(g)
	if VERBOSE >= 1 {
		fmt.Println("nodes:", len(nodesMap))
	}

	//DebugNode(Pos{}, nodesMap)
	//DebugNode(Pos{1, 1}, nodesMap)
	//DebugNode(Pos{}, nodesMap)

	nodes := lib.MapKeys(nodesMap)
	isStart := func(h HistoryNodeKey) bool { return h.P.C == 0 && h.P.R == 0 }
	neighbors := func(h HistoryNodeKey) []HistoryNodeKey {
		return nodesMap[h].Neighbors
	}
	cost := func(h HistoryNodeKey) int { return nodesMap[h].HeatLoss }
	dijkstra := lib.DijkstraWithCost[HistoryNodeKey](nodes, isStart, neighbors, cost)
	fmt.Println(AccumulateHeatLoss(g, nodesMap, dijkstra))
}

func AccumulateHeatLoss(g *Grid, nodes map[HistoryNodeKey]*HistoryNode, paths map[HistoryNodeKey]*lib.Node[HistoryNodeKey]) int {
	rows, cols := g.Rows(), g.Columns()
	minHeatLoss := math.MaxInt
	var minEnd *lib.Node[HistoryNodeKey]
	for k, _ := range nodes {
		if k.P.R == rows-1 && k.P.C == cols-1 {
			node := paths[k]
			if VERBOSE >= 2 {
				fmt.Println("End node found", node.Distance)
			}
			if node.Distance < minHeatLoss {
				minHeatLoss = node.Distance
				minEnd = node
			}
		}
	}
	if VERBOSE >= 2 {
		DescribePath(minEnd)
	}
	return minHeatLoss
}

func DescribePath(end *lib.Node[HistoryNodeKey]) {
	var path []*lib.Node[HistoryNodeKey]
	path = append(path, end)
	current := end
	for !current.Value.P.IsZero() { // start node
		current = current.Prev
		path = append(path, current)
	}
	slices.Reverse(path)
	fmt.Println("path:")
	for _, p := range path {
		key := p.Value
		fmt.Println(key.P)
	}
}

func BuildNodes(g *Grid) map[HistoryNodeKey]*HistoryNode {
	start := &HistoryNode{
		HistoryNodeKey: HistoryNodeKey{P: Pos{0, 0}},
		HeatLoss:       g.Get(0, 0),
	}
	nodes := make(map[HistoryNodeKey]*HistoryNode)
	nodes[start.HistoryNodeKey] = start
	TestMovesAndFollow(g, nodes, start)
	return nodes
}

func TestMovesAndFollow(g *Grid, nodes map[HistoryNodeKey]*HistoryNode, node *HistoryNode) {
	for _, dir := range Directions {
		if node.CanMoveGrid(g, dir) {
			ApplyMoveAndFollow(g, nodes, node, dir)
		}
	}
}

func ApplyMoveAndFollow(g *Grid, nodes map[HistoryNodeKey]*HistoryNode, node *HistoryNode, dir Direction) {
	key := node.ApplyMove(dir)
	if _, found := nodes[key]; !found {
		heatLoss := g.Get(key.P.R, key.P.C)
		newNode := &HistoryNode{
			HistoryNodeKey: key,
			HeatLoss:       heatLoss,
		}
		nodes[key] = newNode
		TestMovesAndFollow(g, nodes, newNode)
	}
	if slices.Contains(node.Neighbors, key) {
		//fmt.Println("double")
	} else {
		node.Neighbors = append(node.Neighbors, key)
	}
}

func DebugNode(p Pos, nodes map[HistoryNodeKey]*HistoryNode) {
	fmt.Println("debugging:", p)
	for k, v := range nodes {
		if k.P == p {
			fmt.Println(k.Moves, "neighbors:", len(v.Neighbors), "heatloss:", v.HeatLoss)
		}
	}
	fmt.Println()
}
