package main

// https://adventofcode.com/2023/day/8

import (
	"fmt"
	"regexp"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/math"
)

var DEBUG = 1

func main() {
	// Steps: 2
	ProcessStep1("aoc23/day08/example.txt")
	// Steps: 6
	ProcessStep1("aoc23/day08/example2.txt")

	ProcessStep1("aoc23/day08/input.txt")

	// Steps: 6
	ProcessStep2("aoc23/day08/step2example.txt")
	ProcessStep2("aoc23/day08/input.txt")
}

func ProcessStep1(name string) {
	fmt.Println("Step 1 input:", name)
	lines := lib.ReadLines(name)
	m := ParseMap(lines)
	steps := m.Traverse()
	fmt.Println("Steps:", steps)

	fmt.Println()
}

func ProcessStep2(name string) {
	fmt.Println("Step 2 input:", name)
	lines := lib.ReadLines(name)
	m := ParseMap(lines)
	steps := m.TraverseGhost()
	fmt.Println("Steps:", steps)

	fmt.Println()
}

func debug(v int, msg string) {
	if v <= DEBUG {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

var pattern = regexp.MustCompile(`^([0-9A-Z]+) = \(([0-9A-Z]+), ([0-9A-Z]+)\)$`)

func ParseMap(lines []string) *Map {
	m := NewMap()
	m.Directions = []rune(lines[0])
	for i := 2; i < len(lines); i++ {
		match := pattern.FindStringSubmatch(lines[i])
		if len(match) != 4 {
			panic(fmt.Errorf("invalid input: %q=%v", lines[i], match))
		}
		m.AddNode(match[1], match[2], match[3])
	}
	return m
}

type Map struct {
	Directions []rune
	Nodes      map[string]*Node
}

func NewMap() *Map {
	return &Map{
		Nodes: make(map[string]*Node),
	}
}

type Node struct {
	Name  string
	Left  *Node
	Right *Node

	GhostStart bool
	GhostEnd   bool
}

func (m *Map) AddNode(name, left, right string) {
	n := m.GetOrCreateNode(name)
	n.Left = m.GetOrCreateNode(left)
	n.Right = m.GetOrCreateNode(right)
}

func (m *Map) GetOrCreateNode(name string) *Node {
	n, found := m.Nodes[name]
	if !found {
		n = &Node{
			Name:       name,
			GhostStart: name[len(name)-1] == 'A',
			GhostEnd:   name[len(name)-1] == 'Z',
		}
		m.Nodes[name] = n
	}
	return n
}

func (m *Map) Traverse() int {
	steps := 0
	node := m.Nodes["AAA"]

	for node.Name != "ZZZ" {
		dir := m.Dir(steps)
		steps++
		if dir == 'L' {
			node = node.Left
		} else {
			node = node.Right
		}
	}
	return steps
}

func (m *Map) Dir(idx int) rune {
	return m.Directions[idx%len(m.Directions)]
}

func (m *Map) TraverseGhost() interface{} {
	steps := 0
	isStartNode := func(n *Node) bool { return n.GhostStart }
	isEndNode := func(n *Node) bool { return n.GhostEnd }
	nodes := lib.Filter(m.AllNodes(), isStartNode)
	fmt.Println("Ghosts:", len(nodes))

	endIndixes := make([]int, len(nodes))
	nEnds := 0

	for !lib.All(nodes, isEndNode) {
		dir := m.Dir(steps)
		steps++
		if dir == 'L' {
			for i, n := range nodes {
				nodes[i] = n.Left
			}
		} else {
			for i, n := range nodes {
				nodes[i] = n.Right
			}
		}

		for i, n := range nodes {
			if n.GhostEnd && endIndixes[i] == 0 {
				fmt.Println("found end index for ghost", i, ":", steps)
				endIndixes[i] = steps
				nEnds++
			}
		}
		if nEnds == len(nodes) {
			break
		}
	}
	return math.LcmN(endIndixes)
}

func (m *Map) AllNodes() []*Node {
	var rv []*Node
	for _, n := range m.Nodes {
		rv = append(rv, n)
	}
	return rv
}
