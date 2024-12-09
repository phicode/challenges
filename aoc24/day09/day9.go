package main

// https://adventofcode.com/2024/day/9

import (
	"flag"
	"fmt"
	"github.com/phicode/challenges/lib/assert"
	"strconv"
	"strings"

	"github.com/phicode/challenges/lib"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day09/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day09/input.txt")
	//
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day09/example.txt")
	//lib.Timed("Part 2", ProcessPart2, "aoc24/day09/input.txt")

	//lib.Profile(1, "day09.pprof", "Part 2", ProcessPart2, "aoc24/dayXX/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	//input := ReadAndParseInput(name)
	//result := SolvePart1(input)
	result := SolvePart1(name)
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
	Start, End *Block
}

func ReadAndParseInput(name string) *Input {
	lines := lib.ReadLines(name)
	assert.True(len(lines) == 1)
	return ParseInput(lines[0])
}

func ParseInput(line string) *Input {
	input := []rune(line)
	start := new(Block)
	end := start
	nextFileId := 0
	for i := 0; i < len(input); i++ {
		num := int(input[i] - '0')
		if i == 0 {
			start.FileId = nextFileId
			nextFileId++
			start.Len = num
			continue
		}
		current := new(Block)
		current.Prev = end
		end.Next = current
		end = current
		if i%2 == 1 {
			current.FileId = -1
		} else {
			current.FileId = nextFileId
			nextFileId++
		}
		current.Len = num
	}
	return &Input{start, end}
}

type Block struct {
	FileId     int // -1 = free
	Len        int
	Next, Prev *Block
}

func (in *Input) String() string {
	sb := strings.Builder{}
	cur := in.Start
	for cur != nil {
		if cur.FileId == -1 {
			sb.WriteString(strings.Repeat(".", cur.Len))
		} else {
			idxStr := strconv.Itoa(cur.FileId)
			sb.WriteString(strings.Repeat(idxStr, cur.Len))
		}
		cur = cur.Next
	}
	return sb.String()
}

func (in *Input) Remove(f *Block) {
	prev, next := f.Remove()
	if next != nil {
		in.End = next.End()
	} else {
		in.End = prev.End()
	}
}

func (in *Input) Split(f *Block) (*Block, *Block) {
	p, q := f.Split()
	if in.Start == f {
		in.Start = p
	}
	if in.End == f {
		in.End = q
	}
	return p, q
}

////////////////////////////////////////////////////////////

func (b *Block) NextFree() *Block {
	cur := b
	for cur != nil {
		if cur.FileId == -1 {
			return cur
		}
		cur = cur.Next
	}
	return nil
}

func (b *Block) PrevFile() *Block {
	cur := b
	for cur != nil {
		if cur.FileId != -1 {
			return cur
		}
		cur = cur.Prev
	}
	return nil
}

func (b *Block) Remove() (*Block, *Block) {
	prev := b.Prev
	next := b.Next
	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}
	return prev, next
}

func (b *Block) End() *Block {
	cur := b
	for cur.Next != nil {
		cur = cur.Next
	}
	return cur
}

func (b *Block) Split() (*Block, *Block) {
	p, q := new(Block), new(Block)
	prev, next := b.Prev, b.Next
	p.Next = q
	q.Prev = p
	if prev != nil {
		prev.Next = p
		p.Prev = prev
	}
	if next != nil {
		next.Prev = q
		q.Next = next
	}
	b.Prev, b.Next = nil, nil
	return p, q
}

func SolvePart1(name string) int {
	lines := lib.ReadLines(name)
	return SolvePart1Naive(lines[0])
}
func SolvePart1Fancy(input *Input) int {
	free := input.Start.NextFree()
	for free != nil {
		assert.True(input.Start != nil)
		assert.True(input.End != nil)
		f := input.End.PrevFile()
		assert.True(f != nil)
		if f.Len >= free.Len { // free block is fully "taken over"
			nextFree := free.NextFree()
			free.FileId = f.FileId
			f.Len -= free.Len
			if f.Len == 0 {
				input.Remove(f)
			}
			free = nextFree
			continue
		}

		// free block has more space than the file
		// free.Len > f.Len
		a, b := input.Split(free)
		a.FileId = f.FileId
		a.Len = f.Len
		b.FileId = -1
		b.Len = free.Len - f.Len
		input.Remove(f)
		free = b
	}
	fmt.Println(input.String())
	return 0
}

////////////////////////////////////////////////////////////

func SolvePart2(input *Input) int {
	return 0
}

type Naive struct {
}

func SolvePart1Naive(instr string) int {
	diskSize := CountAll(instr)
	fileIds := make([]int, diskSize)
	for i := 0; i < diskSize; i++ {
		fileIds[i] = -1
	}
	nextFileId := 0
	idx := 0
	for i, r := range instr {
		n := int(r - '0')
		if i%2 == 1 {
			idx += n
			continue
		}
		for j := 0; j < n; j++ {
			fileIds[idx] = nextFileId
			idx++
		}
		nextFileId++
	}

	// swap fileIds
	freeIdx := NextFreeIdx(fileIds, 0)
	fileIdx := PrevFileIdx(fileIds, len(fileIds)-1)
	for freeIdx < fileIdx {
		fileIds[freeIdx] = fileIds[fileIdx]
		fileIds[fileIdx] = -1
		freeIdx = NextFreeIdx(fileIds, freeIdx+1)
		fileIdx = PrevFileIdx(fileIds, fileIdx-1)
	}
	return Checksum(fileIds)
}

func Checksum(ids []int) int {
	total := 0
	for i, id := range ids {
		if id != -1 {
			total += i * id
		}
	}
	return total
}

func NextFreeIdx(ids []int, i int) int {
	for ids[i] != -1 {
		i++
	}
	return i
}
func PrevFileIdx(ids []int, i int) int {
	for ids[i] == -1 {
		i--
	}
	return i
}

func CountAll(instr string) int {
	total := 0
	for _, r := range instr {
		assert.True(r >= '0' && r <= '9')
		total += int(r - '0')
	}
	return total
}
