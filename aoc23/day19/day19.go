package main

// https://adventofcode.com/2023/day/19

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day19/example.txt") // 19114
	ProcessPart1("aoc23/day19/input.txt")   // 319295

	ProcessPart2("aoc23/day19/example.txt") // 167409079868000
	ProcessPart2("aoc23/day19/input.txt")   // 110807725108076
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	wfs, parts := ParseInput(lines)
	if VERBOSE >= 2 {
		for _, wf := range wfs {
			fmt.Println(wf)
		}
		for _, part := range parts {
			fmt.Println(part)
		}
	}
	sum := SolvePart1(wfs, parts)
	fmt.Println("sum:", sum)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	wfs, _ := ParseInput(lines)
	//t := time.Now()
	SolvePart2(wfs)
	//fmt.Println(time.Now().Sub(t))

	fmt.Println()
}

////////////////////////////////////////////////////////////

const (
	Accept = "A"
	Reject = "R"
)

type Workflows map[string]*Workflow

type Workflow struct {
	Name       string
	Default    string
	Conditions []Condition
}

type Condition struct {
	Index  int
	Op     rune
	Value  int
	Target string // target workflow, Accept or Reject
}

func (c Condition) String() string {
	return fmt.Sprintf("%s%c%d:%s", IndexToLabel[c.Index], c.Op, c.Value, c.Target)
}

func (c Condition) Applies(value int) bool {
	switch c.Op {
	case '>':
		return value > c.Value
	case '<':
		return value < c.Value
	default:
		panic("unknown operation")
	}
}

var LabelToIndex = map[string]int{
	"x": 0,
	"m": 1,
	"a": 2,
	"s": 3,
}
var IndexToLabel = map[int]string{
	0: "x",
	1: "m",
	2: "a",
	3: "s",
}

// values: x, m, a, s
type Part [4]int

func (p Part) String() string {
	var b bytes.Buffer
	b.WriteString("{x=")
	b.WriteString(strconv.Itoa(p[0]))
	b.WriteString(",m=")
	b.WriteString(strconv.Itoa(p[1]))
	b.WriteString(",a=")
	b.WriteString(strconv.Itoa(p[2]))
	b.WriteString(",s=")
	b.WriteString(strconv.Itoa(p[3]))
	b.WriteString("}")
	return b.String()
}

func ParseInput(lines []string) (Workflows, []Part) {
	var workflows = make(Workflows)
	var parts []Part

	parseWF := true
	for _, l := range lines {
		if l == "" {
			parseWF = false
			continue
		}
		if parseWF {
			wf := ParseWorkflow(l)
			workflows[wf.Name] = &wf
		} else {
			parts = append(parts, ParsePart(l))
		}
	}
	return workflows, parts
}

func ParsePart(line string) Part {
	line = line[1 : len(line)-1] // strip start and end {}
	pairs := strings.Split(line, ",")

	var p Part
	for _, pair := range pairs {
		i, v := ParsePair(pair)
		p[i] = v
	}
	return p
}

func ParsePair(pair string) (int, int) {
	label, valueS := Split2(pair, '=')
	value, err := strconv.Atoi(valueS)
	if err != nil {
		panic(err)
	}
	return LabelToIndex[label], value
}

func ParseWorkflow(l string) Workflow {
	var wf Workflow
	name, conditions := Split2(l, '{')
	wf.Name = name
	if conditions[len(conditions)-1] != '}' {
		panic("invalid input")
	}
	conditions = conditions[:len(conditions)-1]
	splitConditions := strings.Split(conditions, ",")
	for _, cond := range splitConditions {
		if strings.ContainsRune(cond, ':') {
			wf.Conditions = append(wf.Conditions, ParseCondition(cond))
		} else {
			// no colon -> default
			if wf.Default != "" {
				panic("invalid input:multiple defaults")
			}
			wf.Default = cond
		}
	}
	if wf.Default == "" {
		panic("invalid input: no default")
	}
	return wf
}

func ParseCondition(cond string) Condition {
	test, target := Split2(cond, ':')
	op, opIdx := FindOp(test)
	label := test[:opIdx]
	valueStr := test[opIdx+1:]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(err)
	}
	return Condition{
		Index:  LabelToIndex[label],
		Op:     op,
		Value:  value,
		Target: target,
	}
}

func FindOp(condition string) (rune, int) {
	for i, x := range condition {
		if x == '>' || x == '<' {
			return x, i
		}
	}
	panic(fmt.Errorf("no operation found in: %q", condition))
}

func Split2(s string, x rune) (string, string) {
	idx := strings.IndexRune(s, x)
	if idx < 1 || idx >= len(s)-1 {
		panic("invalid input")
	}
	return s[:idx], s[idx+1:]
}

////////////////////////////////////////////////////////////
// Part 1

func SolvePart1(wfs Workflows, parts []Part) int {
	var sum int
	for _, part := range parts {
		if AcceptPart(wfs, part) {
			if VERBOSE >= 2 {
				fmt.Println("accepted part:", part)
			}
			sum += part[0]
			sum += part[1]
			sum += part[2]
			sum += part[3]
		}
	}
	return sum
}

func AcceptPart(wfs Workflows, part Part) bool {
	current := wfs["in"]
	for {
		target := current.ApplyConditions(part)
		if target == Accept {
			return true
		}
		if target == Reject {
			return false
		}
		current = wfs[target]
	}
}

func (w *Workflow) ApplyConditions(part Part) string {
	for _, cond := range w.Conditions {
		value := part[cond.Index]
		if cond.Applies(value) {
			return cond.Target
		}
	}
	return w.Default
}

////////////////////////////////////////////////////////////
// Part 2

type Range struct {
	Start, End int // start inclusive, end exclusive
}

// Split Range r, so that x is Part of the lower Range
func (r Range) Split(x int) (Range, Range) {
	if x < r.Start || x >= r.End-1 {
		panic("invalid state")
	}
	return Range{r.Start, x + 1}, Range{x + 1, r.End}
}

func (r Range) Value() int {
	return r.End - r.Start
}

func (c Condition) AppliesFull(r Range) bool {
	return c.Applies(r.Start) && c.Applies(r.End-1)
}
func (c Condition) AppliesPartial(r Range) bool {
	a := c.Applies(r.Start)
	b := c.Applies(r.End - 1)
	return a != b
}

// Split returns the Range which adheres to the condition and the range that does not apply the condition
func (c Condition) Split(r Range) (Range, Range) {
	switch c.Op {
	case '>':
		// b contains the values that apply (value > )
		// the value itself is part of the not applying range
		a, b := r.Split(c.Value)
		return b, a
	case '<':
		// a contains the values that apply (value > )
		// the value itself is part of the not applying range
		a, b := r.Split(c.Value - 1)
		return a, b
	default:
		panic("invalid state")
	}
}

type Part2Solver struct {
	//remaining []RangePart
	Workflows Workflows
	Accepted  int
	Rejected  int
}

type RangePart struct {
	Values [4]Range
	Target string
}

func SolvePart2(wfs Workflows) {
	var s Part2Solver
	s.Workflows = wfs
	start := RangePart{
		Values: [4]Range{
			{1, 4001},
			{1, 4001},
			{1, 4001},
			{1, 4001},
		},
		Target: "in",
	}
	s.Solve(start)
	fmt.Println("Accepted:", s.Accepted)
}

func (s *Part2Solver) Solve(next RangePart) {
	current := s.Workflows[next.Target]
	if current == nil {
		panic(fmt.Sprintf("target not found: %q", next.Target))
	}

	for _, cond := range current.Conditions {
		a, b := next.ApplyCondition(cond)
		if VERBOSE >= 2 {
			fmt.Println("a:", a)
			fmt.Println("b:", b)
		}

		// a is the range that has transitioned to the next workflow
		s.Follow(a)

		// b is the range where the condition did not apply, test the next condition
		if b.IsZero() {
			return
		}
		if b.Target != next.Target {
			panic("invalid state")
		}
		next = b
	}
	next = next.WithTarget(current.Default)
	s.Follow(next)
}

func (s *Part2Solver) Follow(a RangePart) {
	if a.IsZero() {
		return
	}
	if a.Target == Accept {
		v := a.Value()
		s.Accepted += v
		if VERBOSE >= 2 {
			fmt.Println("ACCEPTING RANGE", v, "total:", s.Accepted)
		}
		return
	}
	if a.Target == Reject {
		v := a.Value()
		s.Rejected += v
		if VERBOSE >= 2 {
			fmt.Println("REJECTING RANGE", v, "total:", s.Rejected)
		}
		return
	}
	s.Solve(a)
}

// returns the sub-range that applies to the condition and the remaining range
func (r RangePart) ApplyCondition(c Condition) (RangePart, RangePart) {
	interval := r.Values[c.Index]
	if c.AppliesFull(interval) {
		return r.WithTarget(c.Target), RangePart{} // no split
	}
	if c.AppliesPartial(interval) {
		iA, iB := c.Split(interval)
		match := r.Replace(c.Index, iA).WithTarget(c.Target)
		remaining := r.Replace(c.Index, iB)
		return match, remaining
	}
	return RangePart{}, r
}

func (r RangePart) WithTarget(target string) RangePart {
	r.Target = target
	return r
}

func (r RangePart) Replace(index int, v Range) RangePart {
	r.Values[index] = v
	return r
}

func (r RangePart) Value() int {
	product := 1
	for _, r := range r.Values {
		product *= r.Value()
	}
	return product
}

func (r RangePart) IsZero() bool {
	return r.Target == ""
}
