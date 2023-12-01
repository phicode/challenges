package main

// https://adventofcode.com/2022/day/3

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	Process("aoc22/day03/example.txt")
	Process("aoc22/day03/input.txt")
}

func Process(name string) {
	fmt.Println("input:", name)
	rs := ReadInput(name)

	sum := 0
	for i, r := range rs {
		conflicts := Conflicts(r.A, r.B)
		_ = i
		//fmt.Println("conflicts in rucksack", i, r.A, r.B, ":", conflicts)
		for _, c := range conflicts {
			sum += int(c)
		}
	}
	fmt.Println("total conflict priority sum:", sum)

	groups(rs)

	fmt.Println()
}

func ReadInput(name string) []Rucksack {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	var rs []Rucksack

	for s.Scan() {
		line := s.Bytes()
		rs = append(rs, ParseRucksack(line))
	}

	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return rs
}

func ParseRucksack(line []byte) Rucksack {
	if len(line)%2 != 0 {
		panic(fmt.Errorf("invalid input: %q", line))
	}
	for _, x := range line {
		valid := (x >= 'a' && x <= 'z') || (x >= 'A' && x <= 'Z')
		if !valid {
			panic(fmt.Errorf("invalid input: %q", line))
		}
	}
	r := Rucksack{
		A: toPriorities(line[:len(line)/2]),
		B: toPriorities(line[len(line)/2:]),
	}
	return r
}

func toPriorities(bytes []byte) []byte {
	prios := make([]byte, len(bytes))
	for i, v := range bytes {
		prios[i] = valueToPrio[v]
	}
	return prios
}

// only 1 item type per rucksack
type Rucksack struct {
	A []byte
	B []byte
}

var (
	valueToPrio = map[byte]byte{}
)

func init() {
	for x := 'a'; x <= 'z'; x++ {
		valueToPrio[byte(x)] = byte(x - 'a' + 1)
	}
	for x := 'A'; x <= 'Z'; x++ {
		valueToPrio[byte(x)] = byte(x - 'A' + 27)
	}
}

func Conflicts(as, bs []byte) []byte {
	var ina [53]bool
	var conflicts []byte
	for _, a := range as {
		ina[a] = true
	}
	for _, b := range bs {
		if ina[b] {
			// reset, so that it the item type is not reported twice
			ina[b] = false
			conflicts = append(conflicts, b)
		}
	}
	return conflicts
}

func groups(rs []Rucksack) {
	numGroups := len(rs) / 3
	sum := 0
	for i := 0; i < numGroups; i++ {
		first := i * 3
		common := blub(rs[first : first+3])
		fmt.Println("group", i, "has in common", common)
		if len(common) != 1 {
			panic("unexpected")
		}
		sum += int(common[0])
	}
	fmt.Println("prio sum", sum)
}

func blub(rs []Rucksack) []byte {
	var inall [53]int
	for _, r := range rs {
		var inthis [53]bool
		for _, v := range r.A {
			if !inthis[v] {
				inthis[v] = true
				inall[v]++
			}
		}
		for _, v := range r.B {
			if !inthis[v] {
				inthis[v] = true
				inall[v]++
			}
		}
	}
	var out []byte
	for prio, num := range inall {
		if num == len(rs) {
			out = append(out, byte(prio))
		}
	}
	return out
}
