package main

// https://adventofcode.com/2023/day/20

import (
	"container/list"
	"fmt"
	"strings"

	"git.bind.ch/phil/challenges/lib"
	"git.bind.ch/phil/challenges/lib/math"
)

var VERBOSE = 2

func main() {
	ProcessPart1("aoc23/day20/example.txt", 1)
	ProcessPart1("aoc23/day20/example2.txt", 4)

	VERBOSE = 1
	ProcessPart1("aoc23/day20/example.txt", 1000)  // 32000000
	ProcessPart1("aoc23/day20/example2.txt", 1000) // 11687500
	ProcessPart1("aoc23/day20/input.txt", 1000)    // 807069600

	ProcessPart2("aoc23/day20/input.txt")
}

func ProcessPart1(name string, pulses int) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	net := ParseNetwork(lines)

	for i := 0; i < pulses; i++ {
		net.Button(false)
		if VERBOSE >= 2 {
			fmt.Println()
		}
	}
	fmt.Println("high signals:", net.NumSignalHigh)
	fmt.Println("low signals:", net.NumSignalLow)
	fmt.Println("product:", net.NumSignalLow*net.NumSignalHigh)
	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	net := ParseNetwork(lines)

	if len(net.Untyped) != 1 {
		panic("invalid number of untyped outputs")
	}
	utName := net.Untyped[0]
	m := net.FindFlowContributors(utName)
	fmt.Println(len(m), "nodes contribute to", utName)

	contributors := net.FindImmediateContributors(utName)
	fmt.Println("watching contributors", contributors)

	// first index: contributor
	//  second index:
	contributorsHigh := make([][]int, len(contributors))

	for buttonPress := 1; buttonPress < 100_000_000; buttonPress++ {
		ok, mb := net.Button(true)
		if !ok {
			fmt.Println("aborting after", buttonPress, "button presses")
			return
		}

		for i, contrib := range contributors {
			if mb.DidSendHigh(contrib) {
				fmt.Println("found high signal of contributor", i, "@", buttonPress)
				contributorsHigh[i] = append(contributorsHigh[i], buttonPress)
			}
			allGotData := true
			for _, ch := range contributorsHigh {
				if len(ch) < 3 {
					allGotData = false
				}
			}
			if allGotData {
				fmt.Println("got enough data on button press", buttonPress)
				lcm := AnalyseRepetitions(contributorsHigh)
				fmt.Println("repetition after:", lcm)
				return
			}
		}
	}
}

func AnalyseRepetitions(high [][]int) int {
	var reps []int
	for i, rep := range high {
		diff1 := rep[1] - rep[0]
		diff2 := rep[2] - rep[1]
		fmt.Println("reps", rep[0], rep[1], rep[2], "=", diff1, diff2)
		if diff1 != diff2 {
			panic(fmt.Errorf("bad luck, series % does not repeat", i))
		}
		reps = append(reps, diff1)
	}
	return math.LcmN(reps)
}

////////////////////////////////////////////////////////////

type Network struct {
	Broadcaster  []string
	Conjunctions []*Conjunction
	// key: target node
	// value: all nodes that are an input to the target node
	Inputs        map[string][]string
	All           map[string]Receiver
	NumSignalLow  int
	NumSignalHigh int
	Untyped       []string
}

func (n *Network) AddFlipFlop(name string, outputs []string) {
	f := &FlipFlop{
		Name:    name,
		State:   Low,
		Outputs: outputs,
	}
	n.All[name] = f
}

func (n *Network) AddConjunction(name string, outputs []string) {
	c := &Conjunction{
		Name:    name,
		Outputs: outputs,
	}
	n.All[name] = c
	n.Conjunctions = append(n.Conjunctions, c)
}
func (n *Network) AddUntyped(name string) {
	u := &Untyped{}
	n.All[name] = u
	n.Untyped = append(n.Untyped, name)
}

func (n *Network) AddInput(src, dst string) {
	if n.Inputs == nil {
		n.Inputs = make(map[string][]string)
	}
	inputs := n.Inputs[dst]
	inputs = append(inputs, src)
	n.Inputs[dst] = inputs
}

func (n *Network) AddInputs(src string, dsts []string) {
	for _, dst := range dsts {
		n.AddInput(src, dst)
	}
}

func ParseNetwork(lines []string) *Network {
	var n = Network{
		All: make(map[string]Receiver),
	}
	for _, l := range lines {
		bc, gotBc := strings.CutPrefix(l, "broadcaster -> ")
		if gotBc {
			n.Broadcaster = strings.Split(bc, ", ")
			continue
		}
		parts := strings.Split(l, " -> ")
		if len(parts) != 2 {
			panic("invalid input")
		}
		src := parts[0]
		t := src[0]
		src = src[1:]
		dsts := strings.Split(parts[1], ", ")
		if t == '%' {
			n.AddFlipFlop(src, dsts)
		} else if t == '&' {
			n.AddConjunction(src, dsts)
		} else {
			panic("invalid input")
		}
		n.AddInputs(src, dsts)
	}
	// set conjunction inputs
	for _, c := range n.Conjunctions {
		c.SetInputs(n.Inputs[c.Name])
	}
	// find untyped outputs
	for name, _ := range n.Inputs {
		if _, found := n.All[name]; !found {
			n.AddUntyped(name)
			if VERBOSE >= 1 {
				fmt.Println("found untyped part:", name)
			}
		}
	}
	return &n
}

type Pulse bool

const (
	Low  Pulse = false
	High Pulse = true
)

type FlipFlop struct {
	Name    string
	State   Pulse
	Outputs []string
}

func (f *FlipFlop) Receive(s Signal, mailbox *Mailbox) bool {
	if s.Pulse == Low {
		f.State = !f.State
		for _, out := range f.Outputs {
			mailbox.Post(Signal{f.Name, out, f.State})
		}
	}
	return true
}

type Conjunction struct {
	Name       string
	Inputs     []Pulse
	InputNames []string
	Outputs    []string
}

func (c *Conjunction) SetInputs(names []string) {
	c.Inputs = make([]Pulse, len(names))
	c.InputNames = names
}

func (c *Conjunction) Receive(s Signal, mailbox *Mailbox) bool {
	idx := -1
	for i, in := range c.InputNames {
		if s.Src == in {
			idx = i
			break
		}
	}
	if idx == -1 {
		panic("sender not found")
	}
	c.Inputs[idx] = s.Pulse
	sendPulse := Low
	for _, input := range c.Inputs {
		if input == Low {
			sendPulse = High
			break
		}
	}

	for _, out := range c.Outputs {
		mailbox.Post(Signal{c.Name, out, sendPulse})
	}
	return true
}

type Untyped struct {
}

func (u *Untyped) Receive(s Signal, mailbox *Mailbox) bool {
	// Low signal = false = abort
	// High signal = true = continue
	return bool(s.Pulse)
}

type Receiver interface {
	Receive(signal Signal, mailbox *Mailbox) bool
}

type Mailbox struct {
	queue   list.List
	Signals []Signal
}

func (m *Mailbox) Post(s Signal) {
	m.queue.PushBack(s)
	m.Signals = append(m.Signals, s)
}

func (m *Mailbox) Poll() Signal {
	if m.queue.Len() == 0 {
		return Signal{}
	}
	front := m.queue.Front()
	m.queue.Remove(front)
	return front.Value.(Signal)
}

func (m *Mailbox) DidSendHigh(contrib string) bool {
	for _, s := range m.Signals {
		if s.Src == contrib {
			return s.Pulse == High
		}
	}
	return false
}

type Signal struct {
	Src   string
	Dest  string
	Pulse Pulse
}

func (s Signal) String() string {
	sig := "low"
	if s.Pulse {
		sig = "high"
	}
	return fmt.Sprintf("%s -%s-> %s", s.Src, sig, s.Dest)
}

func (n *Network) Button(handleAbort bool) (bool, *Mailbox) {
	mb := &Mailbox{}

	// signal from button to broadcaster
	n.NumSignalLow++
	if VERBOSE >= 2 {
		fmt.Println(Signal{Src: "button", Dest: "broadcaster", Pulse: Low})
	}

	for _, dest := range n.Broadcaster {
		mb.Post(Signal{Src: "broadcaster", Dest: dest, Pulse: Low})
	}

	for mb.queue.Len() > 0 {
		signal := mb.Poll()
		if VERBOSE >= 2 {
			fmt.Println(signal)
		}
		if signal.Pulse == High {
			n.NumSignalHigh++
		} else {
			n.NumSignalLow++
		}
		receiver := n.All[signal.Dest]
		if receiver == nil {
			panic(fmt.Errorf("receiver not found: %s", signal.Dest))
		}
		ok := receiver.Receive(signal, mb)
		if handleAbort && !ok {
			return false, mb
		}
	}
	return true, mb
}

func (n *Network) FindFlowContributors(s string) map[string]struct{} {
	m := make(map[string]struct{})
	inputs := n.Inputs[s]
	n.AddAllContributors(m, inputs)
	return m
}

func (n *Network) AddAllContributors(m map[string]struct{}, inputs []string) {
	for _, in := range inputs {
		if _, found := m[in]; found {
			continue
		}
		m[in] = struct{}{}
		n.AddAllContributors(m, n.Inputs[in])
	}
}

func (n *Network) FindImmediateContributors(name string) []string {
	contributors := n.Inputs[name]
	for len(contributors) == 1 {
		name = contributors[0]
		if !n.IsConjunction(name) {
			panic(fmt.Errorf("flipflop %q leads into untyped output", name))
		}
		contributors = n.Inputs[name]
	}
	for _, contrib := range contributors {
		if !n.IsConjunction(contrib) {
			panic(fmt.Errorf("flipflop %s leads into node %s", contrib, name))
		}
	}
	return contributors
}

func (n *Network) IsConjunction(name string) bool {
	for _, c := range n.Conjunctions {
		if c.Name == name {
			return true
		}
	}
	return false
}
