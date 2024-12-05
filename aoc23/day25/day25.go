package main

// https://adventofcode.com/2023/day/25

import (
	"fmt"
	"math"
	"strings"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/assert"
)

var VERBOSE = 2

func main() {
	ProcessPart1("aoc23/day25/example.txt")
	VERBOSE = 0
	//ProcessPart1("aoc23/day25/input.txt")
	//
	//ProcessPart2("aoc23/day25/example.txt")
	//ProcessPart2("aoc23/day25/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	g := ParseGraph(lines)
	if VERBOSE >= 2 {
		g.PrintDot()
	}
	if VERBOSE >= 1 {
		g.Print()
	}
	//g.TestLoops()

	//g.BronKerbosch()

	a, b := g.KernighanLin()
	fmt.Println("a:", lib.MapKeys(a))
	fmt.Println("b:", lib.MapKeys(b))
	fmt.Println("Kernighan-Lin Result:", len(a)*len(b))

	sa, sb := g.MinimumConnections()
	fmt.Println("sa:", lib.MapKeys(sa))
	fmt.Println("sb:", lib.MapKeys(sb))
	fmt.Println("MinCon Result:", len(sa)*len(sb))

	sa2, sb2 := g.MinimumConnectionsWithStart(a, b)
	fmt.Println("sa2:", lib.MapKeys(sa2))
	fmt.Println("sb2:", lib.MapKeys(sb2))
	fmt.Println("MinCon2 Result:", len(sa2)*len(sb2))

	//foo := NewFoo(g)
	//foo.Run()

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

////////////////////////////////////////////////////////////

type Edge struct {
	a, b string
}

func (e Edge) Neighbor(name string) string {
	if e.a == name {
		return e.b
	}
	return e.a
}

type Vertex struct {
	name  string
	edges []Edge
}

func (v Vertex) NeighborOf(b *Vertex) bool {
	for _, e := range v.edges {
		n := e.Neighbor(v.name)
		if n == b.name {
			return true
		}
	}
	return false
}

type Graph struct {
	Vs map[string]*Vertex
}

func (g *Graph) AddVertex(name string) {
	v, ok := g.Vs[name]
	if !ok {
		v = &Vertex{name: name}
		g.Vs[name] = v
	}
}

func (g *Graph) AddEdge(name string, edge string) {
	e := NewEdge(name, edge)
	a, b := g.Vs[name], g.Vs[edge]
	a.edges = append(a.edges, e)
	b.edges = append(b.edges, e)
}

func NewEdge(a, b string) Edge {
	if a < b {
		return Edge{a, b}
	}
	return Edge{b, a}
}

func ParseGraph(lines []string) *Graph {
	g := &Graph{Vs: make(map[string]*Vertex)}
	for _, l := range lines {
		nameAndEdges := strings.Split(l, ": ")
		name := nameAndEdges[0]
		g.AddVertex(name)
		edges := strings.Split(nameAndEdges[1], " ")
		for _, edge := range edges {
			g.AddVertex(edge)
			g.AddEdge(name, edge)
		}
	}
	return g
}

func (g *Graph) PrintDot() {
	edges := make(map[Edge]bool)
	fmt.Println("graph {")
	for _, v := range g.Vs {
		for _, e := range v.edges {
			if edges[e] {
				continue
			}
			edges[e] = true
			fmt.Printf("  %q -- %q\n", e.a, e.b)
		}
	}
	fmt.Println("}")
}

func (g *Graph) Print() {
	for _, v := range g.Vs {
		edges := lib.Map(v.edges, func(e Edge) string { return e.Neighbor(v.name) })
		fmt.Printf("%v: %v\n", v.name, strings.Join(edges, " "))
	}
}

type Foo struct {
	G         *Graph
	Visited   map[string]*Vertex
	GroupA    []*Vertex
	All       []*Vertex
	Triples   map[Triple]bool
	QuadNodes map[string]bool
}

func NewFoo(g *Graph) *Foo {
	return &Foo{
		G:         g,
		Visited:   make(map[string]*Vertex),
		All:       lib.MapValues(g.Vs),
		Triples:   make(map[Triple]bool),
		QuadNodes: make(map[string]bool),
	}
}

func (f *Foo) Run() {
	first := f.All[0]
	f.Test(first)
	groupB := f.G.Exclude(f.GroupA)
	namesA := lib.Map(f.GroupA, func(v *Vertex) string { return v.name })
	namesB := lib.Map(groupB, func(v *Vertex) string { return v.name })

	fmt.Printf("group a (%d): %v\n", len(namesA), namesA)
	fmt.Printf("group b (%d): %v\n", len(namesB), namesB)
	assert.True(len(namesA)+len(namesB) == len(f.G.Vs))

	fmt.Println("quad nodes:", lib.MapKeys(f.QuadNodes))

	edges := f.G.Connections(f.GroupA, groupB)
	fmt.Println("edges between group a and b:", edges)
	fmt.Println("Product:", len(namesA)*len(namesB))
}

func (g *Graph) TestLoops() {
	for _, v := range g.Vs {
		g.FindLoop(v)
		//l := g.FindBFS(v, []*Vertex{v}, 0, distances)
		//fmt.Println(v.name, "loops to self:", l)
	}
}

//func (g *Graph) FindBFS(search *Vertex, current *Vertex, edges []Edge, i int, distances map[Edge]int) int {
//	var next []*Vertex
//	for _, v := range current.edges {
//
//		for _, e := range v.edges {
//			neighbor := e.Neighbor(v.name)
//			d, found := distances[neighbor]
//			if found {
//				return d + i + 1
//			}
//			distances[neighbor] = i + 1
//
//			//if neighbor == search.name {
//			//	if i > 2 {
//			//		return i + 1, []string{v.name}
//			//	} else {
//			//		continue
//			//	}
//			//}
//			//
//			//if distances[neighbor] {
//			//	continue
//			//}
//
//			//if nextExclude[neighbor] {
//			//	// already added
//			//	continue
//			//}
//			//nextExclude[neighbor] = true
//			//visited[neighbor] = true
//			vert := g.Vs[neighbor]
//			nextVertices = append(nextVertices, vert)
//		}
//	}
//	return g.FindBFS(search, nextVertices, i+1, distances)
//	//last := path[len(path)-1]
//	//path = append(path, FindNodeTo(vs, last))
//	//return i, path
//}

func FindNodeTo(vs []*Vertex, name string) string {
	for _, v := range vs {
		for _, e := range v.edges {
			if e.Neighbor(v.name) == name {
				return v.name
			}
		}
	}
	panic("invalid state")
}

func (g *Graph) FindLoop(search *Vertex) {
	visited := make(map[Edge]int)

	shortest := math.MaxInt
	for _, e := range search.edges {
		n := e.Neighbor(search.name)
		vn := g.Vs[n]
		s := g.DFSFindLoop(search, vn, 1, visited)
		shortest = min(shortest, s)
	}

	fmt.Println(search.name, "loop:", shortest)
}

func (g *Graph) DFSFindLoop(search, current *Vertex, dist int, visited map[Edge]int) int {
	shortest := math.MaxInt
	for _, e := range current.edges {
		d, found := visited[e]
		if found {
			if d < dist {
				continue
			}
		}
		// edge has not yet been visited, or the current distance is shorter than the previous

		fmt.Println(e, ":", dist)
		visited[e] = dist

		n := e.Neighbor(current.name)
		vn := g.Vs[n]
		g.DFSFindLoop(search, vn, dist+1, visited)
	}
	return shortest
}

func (f *Foo) Test(v *Vertex) {
	if f.Visited[v.name] != nil {
		return
	}
	f.Visited[v.name] = v
	f.GroupA = append(f.GroupA, v)

	// test all edge combinations, if they are triangles
	for i, a := range v.edges {
		for _, b := range v.edges[i+1:] {
			if f.G.IsTriangle(v, a, b) {
				an := a.Neighbor(v.name)
				bn := b.Neighbor(v.name)
				t := NewTriple(v.name, an, bn)
				if !f.Triples[t] {
					f.Triples[t] = true
					fmt.Printf("found triangle: %q, %q, %q\n", v.name, an, bn)
				}
				f.Test(f.G.Vs[an])
				f.Test(f.G.Vs[bn])
			}

			// quad tests
			f.TestQuad(v, f.G.Vs[a.Neighbor(v.name)], f.G.Vs[b.Neighbor(v.name)])
		}
	}
}

func (f *Foo) TestQuad(v, a, b *Vertex) {
	// test if a and b have a common node, which is not v
	for _, ae := range a.edges {
		for _, be := range b.edges {
			aeName := ae.Neighbor(a.name)
			beName := be.Neighbor(b.name)
			if aeName == beName && aeName != v.name {
				f.QuadNodes[v.name] = true
				f.QuadNodes[a.name] = true
				f.QuadNodes[b.name] = true
				f.QuadNodes[aeName] = true
				fmt.Printf("found quad: %v, %v, %v, %v\n", v.name, a.name, b.name, aeName)
			}
		}
	}

}

type Triple struct {
	A, B, C string
}

func NewTriple(a, b, c string) Triple {
	a, b, c = sort(a, b, c)
	return Triple{A: a, B: b, C: c}
}

func sort(a, b, c string) (string, string, string) {
	if a > b {
		a, b = b, a
	}
	if c < a {
		return c, a, b
	}
	if b < c {
		return a, b, c
	}
	return a, c, b
}

func (g *Graph) IsTriangle(v *Vertex, a, b Edge) bool {
	an := a.Neighbor(v.name)
	bn := b.Neighbor(v.name)
	return g.AreConnected(an, bn)
}

func (g *Graph) AreConnected(an string, bn string) bool {
	a := g.Vs[an]
	for _, e := range a.edges {
		if e.Neighbor(an) == bn {
			return true
		}
	}
	return false
}

func (g *Graph) Connections(as, bs []*Vertex) []Edge {
	edges := make(map[Edge]bool)
	for _, a := range as {
		for _, b := range bs {
			for _, e := range b.edges {
				if e.Neighbor(b.name) == a.name {
					edges[e] = true
				}
			}
		}
	}
	return lib.MapKeys(edges)
}

func (g *Graph) Exclude(vs []*Vertex) []*Vertex {
	all := make(map[string]*Vertex)
	for _, v := range g.Vs {
		all[v.name] = v
	}
	for _, v := range vs {
		delete(all, v.name)
	}
	return lib.MapValues(all)
}

type VertexSet map[string]*Vertex

func NewSet(vs []*Vertex) VertexSet {
	rv := make(VertexSet)
	for _, v := range vs {
		rv[v.name] = v
	}
	return rv
}
func (set VertexSet) Copy() VertexSet {
	cpy := make(VertexSet)
	for k, v := range set {
		cpy[k] = v
	}
	return cpy
}
func (set VertexSet) UnionOne(v *Vertex) VertexSet {
	cpy := set.Copy()
	cpy[v.name] = v
	return cpy
}

func (set VertexSet) Intersection(b VertexSet) VertexSet {
	rv := make(VertexSet)
	for k, v := range set {
		if _, found := b[k]; found {
			rv[k] = v
		}
	}

	return rv
}
func (set VertexSet) DifferenceSet(b VertexSet) VertexSet {
	rv := set.Copy()
	for k, _ := range b {
		delete(rv, k)
	}
	return rv
}
func (set VertexSet) DifferenceOne(b *Vertex) VertexSet {
	rv := set.Copy()
	delete(rv, b.name)
	return rv
}
func (set VertexSet) RemoveInPlace(b *Vertex) {
	delete(set, b.name)
}

func (set VertexSet) Contains(neigh string) bool {
	_, found := set[neigh]
	return found
}

// find Clique's
func (g *Graph) BronKerbosch() {
	r := make(VertexSet)
	p := NewSet(lib.MapValues(g.Vs))
	x := make(VertexSet)
	g.BronKerbosch1(r, p, x)
}

func (g *Graph) BronKerbosch1(r, p, x VertexSet) {
	if len(p) == 0 && len(x) == 0 {
		fmt.Println("maximum clique:", r)
	}

	for _, v := range p {
		{
			neigh := g.NeighborSet(v)
			r_ := r.UnionOne(v)
			p_ := p.Intersection(neigh)
			x_ := x.Intersection(neigh)
			g.BronKerbosch1(r_, p_, x_)
		}
		p = p.DifferenceOne(v)
		x = x.UnionOne(v)
	}
}

func (g *Graph) NeighborSet(v *Vertex) VertexSet {
	rv := make(VertexSet)
	for _, e := range v.edges {
		neighName := e.Neighbor(v.name)
		neigh := g.Vs[neighName]
		rv[neighName] = neigh
	}
	return rv
}

func (graph *Graph) KernighanLin() (VertexSet, VertexSet) {
	vs := lib.MapValues(graph.Vs)
	n := len(vs)
	mid := n / 2
	A, B := NewSet(vs[:mid]), NewSet(vs[mid:])

	//costs := CalcCosts(A, B)

	gmax := math.MaxInt
	for gmax > 0 {
		fmt.Printf("sizes: %d, %d ; gmax=%d\n", len(A), len(B), gmax)
		var gv []int
		var av, bv []string

		A_ := A.Copy()
		B_ := B.Copy()
		costs := CalcCosts(A_, B_)

		for i := 0; i < mid; i++ { // TODO: does this work with graphs of odd number of nodes ?
			aname, bname, g := FindMaxG(A_, B_, costs)

			A_.RemoveInPlace(graph.Vs[aname])
			B_.RemoveInPlace(graph.Vs[bname])
			//A_ = A_.DifferenceOne(graph.Vs[aname])
			//B_ = B_.DifferenceOne(graph.Vs[bname])

			// TODO: improve by updating the removed nodes in the costs structure
			costs = CalcCosts(A_, B_)

			gv = append(gv, g)
			av = append(av, aname)
			bv = append(bv, bname)
		}

		var k int
		k, gmax = FindK(gv)

		if gmax > 0 {
			for i := 0; i <= k; i++ {
				fmt.Println("exchanging", av[i], bv[i])
				ExchangeInPlace(A, B, av[i], bv[i])
			}
		}
		fmt.Printf("sizes: %d, %d ; gmax=%d\n", len(A), len(B), gmax)
	}
	if len(A) != len(B) {
		KernighanLinUnevenVertices(A, B)
	}
	return A, B
}

func KernighanLinUnevenVertices(as VertexSet, bs VertexSet) {
	fmt.Println("uneven number of nodes:", len(as), "!=", len(bs))
	// kernighan-lin exchanges pairs until no more improvements can be achieved.
	// however, one node might still be in the wrong group.
	// Test every node in the bigger set if it is better places in the "smaller" group

	// b is the bigger set
	assert.True(len(as) < len(bs))

	var maxd int = math.MinInt
	var maxb *Vertex
	for _, b := range bs {
		//bs.RemoveInPlace(b)
		e := CalcCost(b, as)
		i := CalcCost(b, bs)
		d := e - i // bigger d = higher external connectivity = candiate for exchange
		if d > maxd {
			maxd = d
			maxb = b
		}
		//bs[b.name] = b
	}
	if maxd > 0 {
		fmt.Println("changing node", maxb.name, "from set B to set A")
		bs.RemoveInPlace(maxb)
		as[maxb.name] = maxb
	}
}

func (graph *Graph) MinimumConnections() (VertexSet, VertexSet) {
	vs := lib.MapValues(graph.Vs)
	n := len(vs)
	mid := n / 2
	A, B := NewSet(vs[:mid]), NewSet(vs[mid:])

	return graph.MinimumConnectionsWithStart(A, B)
}

func (graph *Graph) MinimumConnectionsWithStart(A, B VertexSet) (VertexSet, VertexSet) {
	var last *Vertex
	for {
		fmt.Println("-----")
		node := findMostUnevenVertex(A, B)
		if node == nil || last == node {
			break
		}
		last = node
		fmt.Println("A:", lib.MapKeys(A))
		fmt.Println("B:", lib.MapKeys(B))
		if A[node.name] != nil {
			fmt.Println("moving node", node.name, "from A to B")
			delete(A, node.name)
			B[node.name] = node
		} else {
			fmt.Println("moving node", node.name, "from B to A")
			delete(B, node.name)
			A[node.name] = node
		}
	}
	return A, B
}

func findMostUnevenVertex(as, bs VertexSet) *Vertex {
	var maxval int = math.MinInt
	var maxvert *Vertex

	for _, a := range as {
		e := CalcCost(a, bs) // external
		i := CalcCost(a, as) // internal
		d := e - i           // bigger d = higher external connectivity = candiate for exchange
		if d > maxval {
			maxval = d
			maxvert = a
		}
	}

	for _, b := range bs {
		e := CalcCost(b, as) // external
		i := CalcCost(b, bs) // internal
		d := e - i           // bigger d = higher external connectivity = candiate for exchange
		if d > maxval {
			maxval = d
			maxvert = b
		}
	}
	fmt.Println("maxval:", maxval)
	//if maxval > 0 {
	//fmt.Println("changing node", maxvert.name)
	return maxvert
	//}
	//return nil
}

func ExchangeInPlace(as VertexSet, bs VertexSet, a string, b string) {
	av := as[a]
	bv := bs[b]
	assert.True(av != nil && bv != nil)
	delete(as, a)
	delete(bs, b)
	as[b] = bv
	bs[a] = av
}

func FindK(gv []int) (int, int) {
	var k int
	var gmax = math.MinInt

	var sum int
	for i, g := range gv {
		sum += g
		if sum > gmax {
			gmax = sum
			k = i
		}
	}
	return k, gmax
}

func FindMaxG(a VertexSet, b VertexSet, costs *Costs) (string, string, int) {
	var maxg = math.MinInt
	var maxa string
	var maxb string

	for aname, avert := range a {
		for bname, bvert := range b {
			// find a from A and b from B, such that g = D[a] + D[b] − 2×c(a, b) is maximal
			c := 0 // 2*c(a, b) = 2*0 = 0
			if avert.NeighborOf(bvert) {
				c = 2 // 2*c(a, b) = 2*1 = 2
			}
			g := costs.Get(aname) + costs.Get(bname) - c
			if g > maxg {
				maxg = g
				maxa = aname
				maxb = bname
			}
		}
	}
	assert.True(maxa != "")
	assert.True(maxb != "")
	return maxa, maxb, maxg
}

func CalcCosts(as, bs VertexSet) *Costs {
	c := NewCosts()
	for _, a := range as {
		e := CalcCost(a, bs)
		i := CalcCost(a, as)
		c.e[a.name] = e
		c.i[a.name] = i
	}
	for _, b := range bs {
		e := CalcCost(b, as)
		i := CalcCost(b, bs)
		c.e[b.name] = e
		c.i[b.name] = i
	}
	return c
}

func CalcCost(a *Vertex, set VertexSet) int {
	// test if a has neighbors in set
	cost := 0
	for _, e := range a.edges {
		neigh := e.Neighbor(a.name)
		if set.Contains(neigh) {
			cost++
		}
	}
	return cost
}

type Costs struct {
	e map[string]int // external cost
	i map[string]int // internal cost
}

func (c Costs) Get(name string) int {
	return c.e[name] - c.i[name]
}

func NewCosts() *Costs {
	return &Costs{
		e: make(map[string]int),
		i: make(map[string]int),
	}
}
