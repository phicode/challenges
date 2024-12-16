package main

// https://adventofcode.com/2024/day/16

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
	"github.com/phicode/challenges/lib/graphs"
	"github.com/phicode/challenges/lib/rowcol"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day16/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day16/input.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day16/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day16/input.txt")
	//lib.Profile(1, "day16.pprof", "Part 2", ProcessPart2, "aoc24/day16/input.txt")
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
	grid  rowcol.Grid[byte]
	graph *Graph
}

type Graph struct {
	Nodes map[NodeKey]*Node
	//Edges map[Edge]bool
}

type NodeKey struct {
	Facing rowcol.Direction
	P      rowcol.Pos
}

func (key NodeKey) String() string {
	return fmt.Sprintf("%v %c", key.P, key.Facing.PrintChar())
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[NodeKey]*Node),
	}
}

func (g *Graph) AddEdge(a, b NodeKey) {
	na := g.getNode(a)
	nb := g.getNode(b)
	na.Neigh = append(na.Neigh, b)
	nb.Neigh = append(nb.Neigh, a)
}

func (g *Graph) getNode(key NodeKey) *Node {
	n := g.Nodes[key]
	if n == nil {
		n = &Node{Key: key}
		g.Nodes[key] = n
	}
	return n
}

type Node struct {
	Key   NodeKey
	Neigh []NodeKey
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	grid := rowcol.NewByteGridFromStrings(lines)
	graph := NewGraph()
	for p, v := range grid.Iterator() {
		if v == '#' {
			continue
		}
		for _, facing := range rowcol.Directions {
			// add moves
			a := NodeKey{Facing: facing, P: p}
			next := p.AddDir(facing)
			if grid.IsValidPos(next) && grid.GetPos(next) != '#' {
				b := NodeKey{Facing: facing, P: next}
				graph.AddEdge(a, b)
			}

			// add turns
			left := NodeKey{Facing: facing.Left(), P: p}
			right := NodeKey{Facing: facing.Right(), P: p}
			graph.AddEdge(a, left)
			graph.AddEdge(a, right)
		}
	}
	//fmt.Println("graph nodes:", len(graph.Nodes))

	return Input{
		grid:  grid,
		graph: graph,
	}
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	startPos := input.grid.MustFindFirst(func(v byte) bool { return v == 'S' })
	endPos := input.grid.MustFindFirst(func(v byte) bool { return v == 'E' })
	var data []NodeKey = lib.MapKeys(input.graph.Nodes)
	var start = func(a NodeKey) bool { return a.P == startPos && a.Facing == rowcol.Right }
	var neigh = func(t NodeKey) []NodeKey {
		node := input.graph.Nodes[t]
		return node.Neigh
	}
	var cost func(NodeKey, NodeKey) int = score
	// map[NodeKey]*lib.Node[NodeKey]
	withCost := lib.DijkstraWithCost(data, start, neigh, cost)
	var minScore *lib.Node[NodeKey]
	for k, v := range withCost {
		if k.P == endPos {
			if minScore == nil || v.Distance < minScore.Distance {
				minScore = v
				//fmt.Println("Endpos Distance: ", v.Distance)
			}
		}
	}
	if lib.LogLevel >= lib.LogDebug {
		path := minScore.GetPath()
		g := input.grid.Copy()
		for i, n := range path {
			key := n.Value
			fmt.Printf("Position %d: %v, facing: %v\n", i, key.P, key.Facing)
			g.SetPos(key.P, key.Facing.PrintChar())
		}
		rowcol.PrintByteGrid(&g)
	}

	return minScore.Distance
}

func score(a, b NodeKey) int {
	assert.False(a.Facing == b.Facing.Reverse())

	if a.Facing == b.Facing {
		return 1
	}
	return 1000
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) int {
	startPos := input.grid.MustFindFirst(func(v byte) bool { return v == 'S' })
	endPos := input.grid.MustFindFirst(func(v byte) bool { return v == 'E' })
	var data []NodeKey = lib.MapKeys(input.graph.Nodes)
	var start = func(a NodeKey) bool { return a.P == startPos && a.Facing == rowcol.Right }
	var neigh = func(t NodeKey) []NodeKey {
		node := input.graph.Nodes[t]
		return node.Neigh
	}
	var cost func(NodeKey, NodeKey) int = score
	// map[NodeKey]*lib.Node[NodeKey]
	withCost := graphs.DijkstraAllWithCost(data, start, neigh, cost)
	var minScore *graphs.NodeAll[NodeKey]
	for k, v := range withCost {
		if k.P == endPos {
			if minScore == nil || v.Distance < minScore.Distance {
				minScore = v
				//fmt.Println("Endpos Distance: ", v.Distance)
			}
		}
	}
	positions := make(map[NodeKey]bool)
	follow(minScore, positions)
	ps := make(map[rowcol.Pos]bool)
	for pos := range positions {
		ps[pos.P] = true
	}
	//fmt.Println("nodekeys: ", len(positions))
	//fmt.Println("positions: ", len(ps))

	return len(ps)
}

func follow(node *graphs.NodeAll[NodeKey], positions map[NodeKey]bool) {
	if positions[node.Value] {
		return
	}
	positions[node.Value] = true
	//if len(node.Prev) > 1 {
	//	fmt.Printf("split @ %v\n", node.Value)
	//}
	for _, prev := range node.Prev {
		follow(prev, positions)
	}
}
