package main

// https://adventofcode.com/2022/day/6

import (
	"bufio"
	"fmt"
	"os"

	"github.com/phicode/challenges/lib/assets"
)

func main() {
	Process("aoc22/day06/example.txt")
	Process("aoc22/day06/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	for i, line := range lines {
		sop4 := FindSOP([]byte(line), 4)
		sop14 := FindSOP([]byte(line), 14)
		fmt.Println("line", i, "start-of-packet-4", sop4)
		fmt.Println("line", i, "start-of-packet-14", sop14)
	}

	fmt.Println()
}

func ReadInput(name string) []string {
	f, err := os.Open(assets.MustFind(name))
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

type Accumulator struct {
	L       int       // search string length
	n       int       // position in stream
	buf     []byte    // sliding window under consideration
	acc     [256]byte // which byte occurs how often
	dups    [256]byte // which bytes have duplicates
	numdups int       // number of bytes with duplicates
}

func NewAcc(l int) *Accumulator {
	return &Accumulator{
		L:   l,
		n:   -1,
		buf: make([]byte, l),
	}
}

func (a *Accumulator) Push(x byte) int {
	a.n++
	pos := a.n % a.L // insert position in sliding buffer

	if a.n >= a.L {
		// remove existing character from accumulator
		prev := a.buf[pos]
		a.acc[prev]--
		if a.acc[prev] == 1 { // duplication removed
			a.dups[prev] = 0
			a.numdups--
		}
	}

	a.buf[pos] = x
	a.acc[x]++
	if a.acc[x] == 2 { // new duplicate character found
		a.dups[x] = 1
		a.numdups++
	}

	if a.n >= a.L && a.numdups == 0 {
		return a.n + 1
	}
	return -1
}

func FindSOP(xs []byte, l int) int {
	a := NewAcc(l)
	for _, x := range xs {
		idx := a.Push(x)
		if idx != -1 {
			return idx
		}
	}
	return -1
	//l := len(xs)
	//if l < 4 {
	//	return -1
	//}
	//var a, b, c, d byte
	//a, b, c, d = 0, xs[0], xs[1], xs[2]
	//for i := 3; i < l; i++ {
	//	a, b, c, d = b, c, d, xs[i]
	//	if a != b && a != c && a != d && b != c && b != d && c != d {
	//		return i + 1
	//	}
	//}
	//return -1
}
