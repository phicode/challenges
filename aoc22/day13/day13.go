package main

// https://adventofcode.com/2022/day/13

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	ProcessPart1("aoc22/day13/example.txt")
	ProcessPart1("aoc22/day13/input.txt")
	ProcessProcessPart2("aoc22/day13/example.txt")
	ProcessProcessPart2("aoc22/day13/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	ps := Programs(lines)
	idxSum := 0
	for i, p := range ps {
		ok := p.Compare()
		if ok {
			idxSum += i + 1
		}
		fmt.Printf("== Pair %d == right=%t\n", i+1, ok)
	}

	fmt.Printf("index sum: %d\n", idxSum)

	fmt.Println()
}

func ProcessProcessPart2(name string) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	lines = append(lines, "[[2]]", "[[6]]")
	var nodes []Node
	for _, line := range lines {
		nodes = append(nodes, Parse(line))
	}
	sort.Sort(NodeSorter(nodes))
	prod := 1
	for i, node := range nodes {
		s := node.String()
		fmt.Println(node)
		if s == "[[2]]" || s == "[[6]]" {
			prod *= (i + 1)
		}
	}
	fmt.Println("product:", prod)
}

type NodeSorter []Node

func (n NodeSorter) Len() int           { return len(n) }
func (n NodeSorter) Less(i, j int) bool { return n[i].CompareTo(n[j]) == 1 }
func (n NodeSorter) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

var _ sort.Interface = (NodeSorter)(nil)

func Programs(lines []string) []Program {
	var left Node
	var ps []Program
	for _, line := range lines {
		if line == "" {
			continue
		}
		n := Parse(line)
		if left == nil {
			left = n
			continue
		}
		ps = append(ps, Program{left, n})
		left = nil
	}
	return ps
}

type Program struct {
	left  Node
	right Node
}

func (p Program) Compare() bool {
	return p.left.CompareTo(p.right) == 1
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
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return lines
}

////////////////////////////////////////////////////////////

func Parse(line string) Node {
	tokens := tokenize(line)
	var root = &List{}
	var stack []*List
	current := root
	for i, t := range tokens {
		if i == 0 {
			if t.Type != '[' {
				panic("first element is not a start list")
			}
			continue
		}
		switch t.Type {
		case '[':
			stack = append(stack, current)
			current = current.AddList()
		case ']':
			if len(stack) == 0 {
				if i < len(tokens)-1 {
					panic("invalid end of list")
				}
				continue
			}
			current = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		case 'i':
			current.AddValue(t.Value)
		default:
			panic("unknown token")
		}
	}
	if len(stack) > 0 {
		panic("stack not empty")
	}
	if s := root.String(); s != line {
		panic(fmt.Errorf("%q != %q", line, s))
	}
	return root
}

func tokenize(line string) []Token {
	var ts []Token
	var hasvalue bool
	var value int
	for _, v := range line {

		if v == '[' || v == ']' {
			if hasvalue {
				ts = append(ts, Token{Type: 'i', Value: value})
				hasvalue = false
				value = 0
			}
			ts = append(ts, Token{Type: byte(v)})
			continue
		}
		if v == ',' {
			if hasvalue {
				ts = append(ts, Token{Type: 'i', Value: value})
				hasvalue = false
				value = 0
			}
			continue
		}
		if v < '0' || v > '9' {
			panic("invalid input")
		}
		hasvalue = true
		value *= 10
		value += int(v - '0')
	}
	return ts
}

type Node interface {
	IsList() bool
	AppendTo(b *bytes.Buffer)

	// first return value
	// 1: right order
	// 0: continue
	// -1: not right order
	CompareTo(Node) int
	fmt.Stringer
}

type List struct {
	Content []Node
}

type Value int

var _ Node = (*List)(nil)
var _ Node = (*Value)(nil)

func (l *List) IsList() bool { return true }

func (l *List) AddList() *List {
	nested := &List{}
	l.Content = append(l.Content, nested)
	return nested
}

func (l *List) AddValue(value int) {
	v := Value(value)
	l.Content = append(l.Content, &v)
}

func (l *List) String() string {
	b := &bytes.Buffer{}
	l.AppendTo(b)
	return b.String()
}
func (l *Value) String() string {
	b := &bytes.Buffer{}
	l.AppendTo(b)
	return b.String()
}
func (l *List) AppendTo(b *bytes.Buffer) {
	b.WriteByte('[')
	for i, n := range l.Content {
		if i > 0 {
			b.WriteByte(',')
		}
		n.AppendTo(b)
	}
	b.WriteByte(']')
}

func (l *Value) IsList() bool { return false }

func (l *Value) AppendTo(b *bytes.Buffer) {
	b.WriteString(strconv.Itoa(int(*l)))
}

type Token struct {
	Type  byte // [, ], i
	Value int
}

func (l *List) CompareTo(node Node) int {
	var lb *List
	if node.IsList() {
		lb = node.(*List)
	} else { // value
		lb = &List{Content: []Node{node}}
	}
	lena, lenb := len(l.Content), len(lb.Content)
	minlen := min(lena, lenb)
	for i := 0; i < minlen; i++ {
		cmp := l.Content[i].CompareTo(lb.Content[i])
		if cmp != 0 {
			return cmp
		}
	}
	if lena == lenb {
		return 0
	}
	if lena < lenb {
		return 1
	}
	return -1
}

func (l *Value) CompareTo(node Node) int {
	if node.IsList() {
		// compare integer to list
		la := &List{Content: []Node{l}}
		return la.CompareTo(node)
	}

	nb := node.(*Value)
	if *l < *nb {
		return 1
	}
	if *l > *nb {
		return -1
	}
	return 0
}
