package main

// https://adventofcode.com/2022/day/16

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"slices"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

const MEMOIZATION = true

func main() {
	fmt.Println("==== Part 1")
	Process("aoc22/day16/example.txt") // 1651
	Process("aoc22/day16/input.txt")   // 1724

	//Process2("aoc22/day16/example.txt", 30) // 2031

	fmt.Println("==== Part 2")
	ProcessStep2V2("aoc22/day16/example.txt", 26) // 1707
	//ProcessStep2V2("aoc22/day16/input.txt", 26)   // ??

	//fmt.Println("==== Part 2")
	Process2("aoc22/day16/example.txt", 26) // 1707
	//Process2("aoc22/day16/input.txt", 26)   // ??
}

func Process(name string) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	valves := ParseValves(lines)

	mem := newmemory()
	aa := valves["AA"]
	best := aa.BestValue(mem, 30)
	fmt.Println("best value:", best)
	fmt.Println("memory size:", len(mem.best))

	fmt.Println()
}

func ProcessStep2V2(name string, rem int) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	valves := ParseValves(lines)

	mem := newmemory()
	aa := valves["AA"]
	var s StateV2
	s.A = aa
	s.B = aa
	s.Rem = rem
	best := s.BestValueStep2(mem)
	fmt.Println("best value:", best)
	fmt.Println("memory size:", len(mem.best))

	fmt.Println()
}

func Process2(name string, rounds int) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	valves := ParseValves(lines)

	state := NewState(rounds, valves)
	best := state.run(0)

	fmt.Println("best value:", best)
	fmt.Println("memory size:", len(state.mem.best))

	//for k, v := range state.mem.best {
	//	fmt.Println(k, "=", v)
	//}

	//solver := NewSolver(valves)
	//solver.Rounds(rounds)

	fmt.Println()
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

////////////////////////////////////////////////////////////

type Valve struct {
	Idx     int
	Name    string
	Flow    int
	Tunnels []*Valve
}

type Valves map[string]*Valve

func ParseValves(lines []string) Valves {
	vs := make(Valves)
	for i, line := range lines {
		name, flow, tunnels := ParseValve(line)
		//fmt.Printf("valve %q (%d) tunnels: %d\n", name, flow, len(tunnels))

		v := vs.Get(name)
		v.Idx = i
		v.Flow = flow
		for _, tun := range tunnels {
			v.Tunnels = append(v.Tunnels, vs.Get(tun))
		}
	}
	return vs
}

func (vs Valves) Get(n string) *Valve {
	if v, ok := vs[n]; ok {
		return v
	}
	v := &Valve{Name: n}
	vs[n] = v
	return v
}

func ParseValve(line string) (string, int, []string) {
	var name string
	var flow int
	var tunnels string
	n, err := fmt.Sscanf(line, "Valve %s has flow rate=%d;", &name, &flow)
	if n != 2 || err != nil {
		panic(fmt.Errorf("invalid line, n=%d, err=%v - %q", n, err, line))
	}
	// 9th space
	spc := 0
	for i, v := range line {
		if v == ' ' {
			spc++
		}
		if spc == 9 {
			tunnels = line[i+1:]
			break
		}
	}
	if len(tunnels) == 0 {
		panic(fmt.Errorf("invalid input: %q", line))
	}

	return name, flow, strings.Split(tunnels, ", ")
}

func (v *Valve) ValueAtTime(rem int) int {
	// water begins to flow the next round -> subtract one more from remaining time
	return max(0, (rem-1)*v.Flow)
}

func (v *Valve) BestNeighbor(mem *memory, rem int) int {
	var best int
	for i, t := range v.Tunnels {
		val := t.BestValue(mem, rem)
		if i == 0 || val > best {
			best = val
		}
	}
	return best
}

func (v *Valve) BestValue(mem *memory, rem int) int {
	key := NewKey(v, nil, rem, mem.openValves)
	best, found := mem.get(key)
	if found {
		return best
	}
	best = v.CalcBestValue(mem, rem)
	mem.put(key, best)
	return best
}

func (v *Valve) CalcBestValue(mem *memory, rem int) int {
	if rem <= 1 {
		return 0
	}

	if v.Flow == 0 && rem > 1 {
		// time remaining -1 for moving to neighbor
		return v.BestNeighbor(mem, rem-1)
	}

	// combinations:
	//   A: activate current valve, then search neighbors
	//   B: search neighbors, activate current valve
	//   C: search neighbors if valve is already open

	// C
	if mem.openValves.IsOpen(v.Idx) {
		return v.BestNeighbor(mem, rem-1)
	}

	// A
	mem.openValves.SetOpen(v.Idx)
	valueA := v.ValueAtTime(rem)
	bestA := v.BestNeighbor(mem, rem-2)
	totalA := valueA + bestA
	mem.openValves.SetClosed(v.Idx)

	// B
	totalB := v.BestNeighbor(mem, rem-1)

	return max(totalA, totalB)
}

func (v *Valve) Combinations(ov OpenValves, sp *ShortestPaths) []Combination {
	var cs []Combination

	if v.Flow > 0 && !ov.IsOpen(v.Idx) {
		cs = append(cs, Combination{
			Start: v,
			Open:  true,
			End:   v,
		})
	}

	//for _, t := range v.Tunnels {
	//	cs = append(cs, Combination{
	//		Start: v,
	//		Open:  false,
	//		End:   t,
	//	})
	//}

	steps := NextSteps(v, sp, ov)
	for _, step := range steps {
		cs = append(cs, Combination{
			Start:  v,
			Open:   false,
			End:    step.Next,
			Target: step.Target,
		})
	}

	return cs
}

type Key struct {
	OV OpenValves

	// Valve a index: upper 8 bits
	// Valve b index: lower 8 bits
	Positions int16

	//NameA string
	//NameB string
	Rem int16
}

func (k Key) String() string {
	posA := k.Positions & 0xff
	posB := k.Positions >> 8

	return fmt.Sprintf("pos (%d/%d)(r:%d),ov:%v", posA, posB, k.Rem, k.OV)
}

type memory struct {
	best       map[Key]int
	openValves OpenValves
}

func newmemory() *memory {
	return &memory{best: make(map[Key]int)}
}

func (m *memory) get(k Key) (int, bool) {
	a, found := m.best[k]
	return a, found
}
func (m *memory) put(k Key, v int) {
	m.best[k] = v
}

type OpenValves uint64

func (o *OpenValves) SetOpen(i int)       { *o |= 1 << i }
func (o *OpenValves) IsOpen(i int) bool   { return *o&(1<<i) != 0 }
func (o *OpenValves) IsClosed(i int) bool { return *o&(1<<i) == 0 }
func (o *OpenValves) SetClosed(i int)     { *o = *o & ^(1 << i) }
func (o *OpenValves) Invert() OpenValves  { return ^(*o) }

func (o *OpenValves) CountClosed(allOpen OpenValves) int {
	todo := *o ^ allOpen
	return bits.OnesCount64(uint64(todo))
}

type Combination struct {
	Start  *Valve
	Open   bool
	End    *Valve
	Target *Valve
}

type Round struct {
	Rem    int
	valveA *Valve
	valveB *Valve
	Value  int
	OV     OpenValves
}

type State struct {
	vs      Valves
	sp      *ShortestPaths
	rounds  []Round
	mem     *memory
	AllOpen OpenValves
}

func NewState(rounds int, vs Valves) *State {
	s := &State{
		vs:     vs,
		sp:     NewShortestPaths(vs),
		rounds: make([]Round, rounds),
		mem:    newmemory(),
	}
	for i := 0; i < rounds; i++ {
		s.rounds[i].Rem = rounds - i - 1 // effects apply next round
	}
	s.rounds[0].valveA = vs["AA"]
	s.rounds[0].valveB = vs["AA"]
	for _, v := range vs {
		if v.Flow > 0 {
			s.AllOpen.SetOpen(v.Idx)
		}
	}
	return s
}

func (s *State) run(i int) int {
	// abort condition
	if i >= len(s.rounds) {
		return 0
	}
	r := &s.rounds[i]
	// copy initial state from previous round
	if i > 0 {
		prev := s.rounds[i-1]
		r.OV = prev.OV
		r.valveA = prev.valveA
		r.valveB = prev.valveB
	}
	// no more valves to open -> no more flow to add
	if r.OV == s.AllOpen {
		return 0
	}
	if MEMOIZATION {
		key := NewKey(r.valveA, r.valveB, r.Rem, r.OV)
		// memoization
		if v, ok := s.mem.get(key); ok {
			return v
		}
	}
	numOpen := r.OV.CountClosed(s.AllOpen)

	csA := r.valveA.Combinations(r.OV, s.sp)
	csB := r.valveB.Combinations(r.OV, s.sp)
	best := 0

	for _, cA := range csA {
		for _, cB := range csB {
			// combination where both open the same valve
			if sameOpen(cA, cB) {

				if numOpen == 1 {
					r.apply(true, cA)
					best = max(best, r.Value)
					r.undo(true, cA)
				}
				// else: ignore combination
			} else {
				//if numOpen > 1 && (cA.Target == cB.Target) {
				//	continue
				//}
				r.apply(true, cA)
				r.apply(false, cB)

				best = max(best, r.Value+s.run(i+1))

				r.undo(false, cB)
				r.undo(true, cA)
			}
		}
	}

	if MEMOIZATION {
		key := NewKey(r.valveA, r.valveB, r.Rem, r.OV)
		//if best > 0 {
		s.mem.put(key, best)
		//}
		if l := len(s.mem.best); l > 0 && l%1_000_000 == 0 {
			fmt.Println("memo size:", len(s.mem.best))
		}
	}
	return best
}

func sameOpen(a, b Combination) bool {
	return a.Open && b.Open && a.Start == b.Start
}

func NewKey(a, b *Valve, rem int, ov OpenValves) Key {
	var positions int
	if b == nil {
		positions = a.Idx
	} else {
		if a.Name > b.Name {
			a, b = b, a
		}
		positions = a.Idx<<8 | b.Idx
	}
	return Key{
		OV:        ov,               // 8 bytes
		Positions: int16(positions), // 2 bytes
		//NameA: a.Name,  // 8 bytes
		//NameB: b.Name,  // 8 bytes
		Rem: int16(rem), //2 bytes
	}
}

func (r *Round) apply(a bool, c Combination) {
	if c.Open {
		r.Value += r.Rem * c.Start.Flow
		r.OV.SetOpen(c.Start.Idx)
		return
	}
	if a {
		r.valveA = c.End
	} else {
		r.valveB = c.End
	}
}
func (r *Round) undo(a bool, c Combination) {
	if c.Open {
		r.Value -= r.Rem * c.Start.Flow
		r.OV.SetClosed(c.Start.Idx)
		return
	}
	if a {
		r.valveA = c.Start
	} else {
		r.valveB = c.Start
	}
}

type Solver struct {
	Valves       Valves
	ToOpenValves map[string]*Valve
}

func NewSolver(valves Valves) *Solver {
	s := &Solver{
		Valves:       valves,
		ToOpenValves: make(map[string]*Valve),
	}
	for _, v := range valves {
		if v.Flow > 0 {
			s.ToOpenValves[v.Name] = v
		}
	}
	return s
}

func (s *Solver) Rounds(minutes int) {
	s.BuildValvePermutations()

	pos := s.Valves.Get("AA")
	for i := 1; i <= minutes; i++ {
		fmt.Printf("== Minute %d ==\n", i)
		distances := s.getDistances(pos)
		_ = distances
		for _, toOpen := range s.ToOpenValves {
			dist := distances[toOpen.Name]
			rem := minutes - i - dist.Distance
			if rem < 0 {
				panic("negative remaining time")
			}
			pressure := toOpen.Flow * rem
			fmt.Printf("can open %q in %d minutes for %d pressure\n", toOpen.Name, dist.Distance, pressure)
		}
		break
	}
}

// next steps
// - calculate shortest path to all valves
// - calculate which valve to open next
// - plot path
func (s *Solver) getDistances(pos *Valve) map[string]*lib.Node[string] {
	var names []string
	for _, v := range s.Valves {
		names = append(names, v.Name)
	}
	isStart := func(name string) bool {
		return name == pos.Name
	}
	neighbors := func(name string) []string {
		v := s.Valves[name]
		var neighs []string
		for _, neigh := range v.Tunnels {
			neighs = append(neighs, neigh.Name)
		}
		return neighs
	}
	dijkstra := lib.Dijkstra(names, isStart, neighbors)

	fmt.Println(len(dijkstra))
	return dijkstra
}

func (s *Solver) BuildValvePermutations() [][]int {
	var toOpen []*Valve
	for _, v := range s.ToOpenValves {
		toOpen = append(toOpen, v)
	}
	fmt.Printf("building %d permutations to open valves\n", len(toOpen)*len(toOpen))
	return nil
}

type PermutationIterator struct {
	valves []*Valve
}

type RoundState struct {
	Pressure int
	Total    int
	Position string
	Target   string
	OV       OpenValves
}

// finds the next steps to take, given
//   - the current position
//   - all valves
//   - the list of currently open valves
func NextSteps(pos *Valve, sp *ShortestPaths, ov OpenValves) []Step {
	var rv []Step
	for _, v := range sp.vs {
		if v != pos && v.Flow > 0 && !ov.IsOpen(v.Idx) {
			s := Step{Target: v, Next: ShortestPathTo(sp, pos, v)}
			rv = append(rv, s)
		}
	}
	return rv
}

type Step struct {
	// where we want to go to
	Target *Valve

	// the next valve to navigate in order to reach Target in the least steps possible
	Next *Valve
}

type ShortestPaths struct {
	vs Valves
	// key: source node
	// value:
	//   key: target node
	//   value: lib.Node of target node
	shortest map[string]map[string]*lib.Node[string]
}

func NewShortestPaths(vs Valves) *ShortestPaths {
	var names []string
	for _, v := range vs {
		names = append(names, v.Name)
	}
	neighbors := func(name string) []string {
		v := vs[name]
		var neighs []string
		for _, neigh := range v.Tunnels {
			neighs = append(neighs, neigh.Name)
		}
		return neighs
	}
	sp := &ShortestPaths{
		vs:       vs,
		shortest: make(map[string]map[string]*lib.Node[string]),
	}
	for _, v := range vs {
		isStart := func(n string) bool { return n == v.Name }
		sp.shortest[v.Name] = lib.Dijkstra(names, isStart, neighbors)
	}
	return sp
}

func ShortestPathTo(sp *ShortestPaths, from *Valve, to *Valve) *Valve {
	graph := sp.shortest[from.Name]
	toNode := graph[to.Name]
	if toNode == nil || toNode.Prev == nil {
		panic("invalid state")
	}
	for toNode.Prev.Value != from.Name {
		if toNode.Prev == nil {
			panic("invalid state")
		}
		toNode = toNode.Prev
	}
	next := sp.vs.Get(toNode.Value)
	if !slices.Contains(from.Tunnels, next) {
		panic("invalid state")
	}
	return next
}

////////////////////////////////////////////////////////////
// V2

func (s *StateV2) BestValueStep2(mem *memory) int {
	if s.Rem <= 1 {
		return 0
	}
	key := NewKey(s.A, s.B, s.Rem, s.OV)
	//key := Key{Rem: int16(s.Rem), OV: s.OV}
	if v, ok := mem.get(key); ok {
		return v
	}
	var best int
	for i := 0; i <= len(s.A.Tunnels); i++ {
		for j := 0; j <= len(s.B.Tunnels); j++ {
			var next StateV2
			ov := s.OV
			var flow int // flow added by this step
			if i < len(s.A.Tunnels) {
				// move to another node
				next.A = s.A.Tunnels[i]
			} else {
				// stay at current node and open vale (if possible)
				next.A = s.A
				flow += ov.Open(s.A, s.Rem)
			}
			if j < len(s.B.Tunnels) {
				// move to another node
				next.B = s.B.Tunnels[j]
			} else {
				// stay at current node and open vale (if possible)
				next.B = s.B
				flow += ov.Open(s.B, s.Rem)
			}
			next.OV = ov
			next.Rem = s.Rem - 1

			result := flow + next.BestValueStep2(mem)
			best = max(best, result)
		}
	}
	mem.put(key, best)
	return best
}

type StateV2 struct {
	A   *Valve
	B   *Valve
	Rem int
	OV  OpenValves
}

func (ov *OpenValves) Open(v *Valve, rem int) int {
	if !ov.IsOpen(v.Idx) {
		ov.SetOpen(v.Idx)
		return v.ValueAtTime(rem)
	}
	return 0
}
