package main

// https://adventofcode.com/2022/day/1

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
)

func main() {
	Process("aoc22/day01/example.txt")
	Process("aoc22/day01/input.txt")
}

func Process(name string) {
	fmt.Println("Input:", name)
	elfs := ReadInput(name)
	fmt.Printf("#elfs: %d\n", len(elfs))
	var maxElf *Elf
	top3 := &TopX{N: 3}
	for i, elf := range elfs {
		top3.Push(elf)
		if i == 0 {
			maxElf = elf
		} else {
			if elf.Sum > maxElf.Sum {
				maxElf = elf
			}
		}
	}
	fmt.Printf("max-elf: %d\n", maxElf.Sum)
	fmt.Printf("top 3 elf sum: %d\n", top3.Sum)
	fmt.Println()
}

type TopX struct {
	N   int
	Sum int
	heap.Interface
	TopElfs MinHeapElfs
}

type MinHeapElfs []*Elf

var _ heap.Interface = (*MinHeapElfs)(nil)

func (h MinHeapElfs) Len() int           { return len(h) }
func (h MinHeapElfs) Less(i, j int) bool { return h[i].Sum < h[j].Sum }
func (h MinHeapElfs) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeapElfs) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Elf))
}

func (h *MinHeapElfs) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (x *TopX) Push(elf *Elf) {
	if len(x.TopElfs) < x.N {
		x.Sum += elf.Sum
		heap.Push(&x.TopElfs, elf)
		return
	}
	top := heap.Pop(&x.TopElfs).(*Elf)
	if top.Sum > elf.Sum {
		heap.Push(&x.TopElfs, top)
	} else {
		x.Sum -= top.Sum
		x.Sum += elf.Sum
		heap.Push(&x.TopElfs, elf)
	}
}

type Elf struct {
	Sum   int
	Foods []int
}

func ReadInput(name string) []*Elf {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	elfs := make([]*Elf, 0, 128)
	elf := &Elf{}
	elfs = append(elfs, elf)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			// next elf
			elf = &Elf{}
			elfs = append(elfs, elf)
			continue
		}
		food, err := strconv.Atoi(line)
		if err != nil {
			panic(fmt.Errorf("invalid number: %q", line))
		}
		elf.Foods = append(elf.Foods, food)
		elf.Sum += food
	}
	if err := s.Err(); err != nil {
		panic(s.Err())
	}
	return elfs
}
