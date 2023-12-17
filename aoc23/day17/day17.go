package main

// https://adventofcode.com/2023/day/17

import (
	"fmt"
	"math"
	"slices"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 0

func main() {
	ProcessPart1("aoc23/day17/example.txt")
	ProcessPart1("aoc23/day17/input.txt")

	ProcessPart2("aoc23/day17/example.txt")  // 94
	ProcessPart2("aoc23/day17/example2.txt") // 71
	ProcessPart2("aoc23/day17/input.txt")    // 1055
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
	g := ParseInput(lines)
	nodes := BuildPart2Nodes(g)

	//if VERBOSE >= 2 {
	//	DebugPart2Node(Pos{0, 0}, nodes)
	//}

	FollowPart2(g, nodes)

	fmt.Println()
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
	nodesMap := BuildPart1Nodes(g)
	if VERBOSE >= 1 {
		fmt.Println("nodes:", len(nodesMap))
	}

	if VERBOSE >= 3 {
		DebugNode(Pos{}, nodesMap)
		DebugNode(Pos{1, 1}, nodesMap)
		DebugNode(Pos{}, nodesMap)
	}

	nodes := lib.MapKeys(nodesMap)
	isStart := func(h HistoryNodeKey) bool { return h.P.C == 0 && h.P.R == 0 }
	neighbors := func(h HistoryNodeKey) []HistoryNodeKey {
		return nodesMap[h].Neighbors
	}
	cost := func(u, v HistoryNodeKey) int { return nodesMap[v].HeatLoss } // cost is stored on the target vertex
	dijkstra := lib.DijkstraWithCost[HistoryNodeKey](nodes, isStart, neighbors, cost)
	heatLoss := AccumulateHeatLoss(g, nodesMap, dijkstra)
	fmt.Println("heat loss:", heatLoss)
}

func AccumulateHeatLoss(g *Grid, nodes map[HistoryNodeKey]*HistoryNode, paths map[HistoryNodeKey]*lib.Node[HistoryNodeKey]) int {
	rows, cols := g.Rows(), g.Columns()
	minHeatLoss := math.MaxInt
	var minEnd *lib.Node[HistoryNodeKey]
	for k := range nodes {
		if k.P.R == rows-1 && k.P.C == cols-1 { // end node
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

func BuildPart1Nodes(g *Grid) map[HistoryNodeKey]*HistoryNode {
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
	if !slices.Contains(node.Neighbors, key) {
		node.Neighbors = append(node.Neighbors, key)
	}
}

func DebugNode(p Pos, nodes map[HistoryNodeKey]*HistoryNode) {
	fmt.Println("debugging:", p)
	for k, v := range nodes {
		if k.P == p {
			fmt.Println(k.Moves, "neighbors:", len(v.Neighbors), "heat loss:", v.HeatLoss)
		}
	}
	fmt.Println()
}

type Part2Node struct {
	Visited   bool
	Neighbors []Part2Key
	// value: cost
	NeighborCost map[Part2Key]int
}

func (n *Part2Node) AddNeighbor(pos Part2Key, cost int) {
	if slices.Contains(n.Neighbors, pos) {
		panic("double add")
	}
	n.Neighbors = append(n.Neighbors, pos)
	n.NeighborCost[pos] = cost
}

func BuildPart2Nodes(g *Grid) map[Part2Key]*Part2Node {
	nodes := make(map[Part2Key]*Part2Node)
	start := Part2Key{Pos{0, 0}, Stand}
	nodes[start] = NewPart2Node()
	TestMovesAndFollowPart2(g, nodes, start)
	if VERBOSE >= 1 {
		fmt.Println("part2 nodes:", len(nodes))
	}
	return nodes
}

func NewPart2Node() *Part2Node {
	return &Part2Node{NeighborCost: make(map[Part2Key]int)}
}

func TestMovesAndFollowPart2(g *Grid, nodes map[Part2Key]*Part2Node, vertex Part2Key) {
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
		for i := 4; i <= 10; i++ {
			targetVertex := Part2Key{Pos{vertex.Pos.R + i*dir.R, vertex.Pos.C + i*dir.C}, dir}
			if g.IsValidPosition(targetVertex.Pos.R, targetVertex.Pos.C) {
				node, found := nodes[targetVertex]
				if !found {
					node = NewPart2Node()
					nodes[targetVertex] = node
				}
				TestMovesAndFollowPart2(g, nodes, targetVertex)
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

func FollowPart2(g *Grid, nodes map[Part2Key]*Part2Node) {
	nodeKeys := lib.MapKeys(nodes)
	start := func(k Part2Key) bool { return k.Pos.IsZero() }
	neighbors := func(k Part2Key) []Part2Key {
		return nodes[k].Neighbors
	}
	cost := func(u Part2Key, v Part2Key) int {
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

func VisitPart2EndNodes(g *Grid, nodes map[Part2Key]*Part2Node, paths map[Part2Key]*lib.Node[Part2Key]) int {
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
				DescribePart2Path(node, nodes)
			}
		}
	}
	return minHeatLoss
}

func DescribePart2Path(end *lib.Node[Part2Key], nodes map[Part2Key]*Part2Node) {
	if end.Prev == nil {
		fmt.Println("found unvisited end node")
		return
	}
	var path []*lib.Node[Part2Key]
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

type Part2Key struct {
	Pos Pos
	Dir Direction // how we got here
}
