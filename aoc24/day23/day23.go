package main

// https://adventofcode.com/2024/day/23

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day23/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day23/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day23/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day23/input.txt")

	//lib.Profile(1, "day23.pprof", "Part 2", ProcessPart2, "aoc24/day23/input.txt")
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
	Nodes map[string]*Node
}

func (in Input) addEdge(a, b string) {
	na := in.getNode(a)
	nb := in.getNode(b)
	na.Neigh = append(na.Neigh, nb)
	nb.Neigh = append(nb.Neigh, na)
}

func (in Input) getNode(name string) *Node {
	n, ok := in.Nodes[name]
	if !ok {
		n = &Node{Name: name}
		in.Nodes[name] = n
	}
	return n
}

type Node struct {
	Name  string
	Neigh []*Node
}

func (n *Node) HasNeighbor(c *Node) bool {
	for _, neigh := range n.Neigh {
		if neigh.Name == c.Name {
			return true
		}
	}
	return false
}

func ReadAndParseInput(name string) Input {
	lines := lib.ReadLines(name)
	return ParseInput(lines)
}

func ParseInput(lines []string) Input {
	in := Input{make(map[string]*Node)}
	for _, line := range lines {
		parts := strings.Split(line, "-")
		assert.True(len(parts) == 2)
		in.addEdge(parts[0], parts[1])
	}
	return in
}

////////////////////////////////////////////////////////////

func SolvePart1(input Input) int {
	groups := ExtractGroupsOf3(input)
	withT := 0
	for _, group := range groups {
		if group.HasT {
			withT++
		}
	}
	return withT
}

func ExtractGroupsOf3(input Input) map[string]*Group {
	groups := make(map[string]*Group)
	for _, a := range input.Nodes {
		for i, b := range a.Neigh {
			for j, c := range a.Neigh {
				if i != j {
					if b.HasNeighbor(c) {
						ga, gb, gc := Sort3(a, b, c)
						hasT := ga.Name[0] == 't' || gb.Name[0] == 't' || gc.Name[0] == 't'
						name := GroupName(ga, gb, gc)
						if groups[name] == nil {
							groups[name] = &Group{
								Nodes: []*Node{a, b, c},
								HasT:  hasT,
							}
						}
					}
				}
			}
		}
	}
	return groups
}

type Group struct {
	Nodes []*Node
	HasT  bool
}

func GroupName(a, b, c *Node) string {
	return fmt.Sprintf("%s,%s,%s", a.Name, b.Name, c.Name)
}

func Sort3(a, b, c *Node) (*Node, *Node, *Node) {
	if b.Name < a.Name {
		a, b = b, a
	}
	if c.Name < a.Name {
		return c, a, b
	}
	if c.Name > b.Name {
		return a, b, c
	}
	return a, c, b
}

////////////////////////////////////////////////////////////

func SolvePart2(input Input) string {
	name := ""
	groups := ExtractGroupsOf3(input)
	for _, group := range groups {
		ExpandGroup(group, input.Nodes)
		n := group.Name()
		if len(n) > len(name) {
			name = n
		}
	}
	return name
}

func ExpandGroup(group *Group, nodes map[string]*Node) {
	for _, n := range nodes {
		group.Growgroup(n)
	}
}

func (g *Group) Growgroup(n *Node) {
	// node n must be neighbor with all current group members
	for _, groupMember := range g.Nodes {
		if groupMember.Name == n.Name {
			return
		}
		if !groupMember.HasNeighbor(n) {
			return
		}
	}
	g.Nodes = append(g.Nodes, n)
}

func (g *Group) Name() string {
	nodeNames := make([]string, len(g.Nodes))
	for i, node := range g.Nodes {
		nodeNames[i] = node.Name
	}
	sort.Strings(nodeNames)
	return strings.Join(nodeNames, ",")
}
