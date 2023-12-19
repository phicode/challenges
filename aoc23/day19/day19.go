package main

// https://adventofcode.com/2023/day/19

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 2

func main() {
	ProcessPart1("aoc23/day19/example.txt")
	ProcessPart1("aoc23/day19/input.txt")

	//ProcessPart2("aoc23/day19/example.txt")
	//ProcessPart2("aoc23/day19/input.txt")
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
	_ = lines

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
	Label  string
	Op     rune
	Value  int
	Target string // target workflow, Accept or Reject
}

func (c Condition) String() string {
	return fmt.Sprintf("%s%c%d:%s", c.Label, c.Op, c.Value, c.Target)
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

type Pair struct {
	Label string
	Value int
}
type Part []Pair

func (p Part) Get(label string) int {
	for _, pair := range p {
		if pair.Label == label {
			return pair.Value
		}
	}
	panic(fmt.Errorf("label %q not found in part %s", label, p))
}
func (p Part) String() string {
	var b bytes.Buffer
	b.WriteRune('{')
	for i, pair := range p {
		if i > 0 {
			b.WriteRune(',')
		}
		b.WriteString(pair.Label)
		b.WriteRune('=')
		b.WriteString(strconv.Itoa(pair.Value))
	}
	b.WriteRune('}')
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
		p = append(p, ParsePair(pair))
	}
	return p
}

func ParsePair(pair string) Pair {
	label, valueS := Split2(pair, '=')
	value, err := strconv.Atoi(valueS)
	if err != nil {
		panic(err)
	}
	return Pair{Label: label, Value: value}
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
		Label:  label,
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
	var x, m, a, s int
	for _, part := range parts {
		if AcceptPart(wfs, part) {
			if VERBOSE >= 2 {
				fmt.Println("accepted part:", part)
			}
			x += part.Get("x")
			m += part.Get("m")
			a += part.Get("a")
			s += part.Get("s")
		}
	}
	return x + m + a + s
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
		value := part.Get(cond.Label)
		if cond.Applies(value) {
			return cond.Target
		}
	}
	return w.Default
}
