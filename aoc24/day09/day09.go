package main

// https://adventofcode.com/2024/day/9

import (
	"flag"
	"fmt"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
)

func main() {
	flag.Parse()
	lib.Timed("Part 1", ProcessPart1, "aoc24/day09/example.txt")
	lib.Timed("Part 1", ProcessPart1, "aoc24/day09/input.txt")

	lib.Timed("Part 2", ProcessPart2, "aoc24/day09/example.txt")
	lib.Timed("Part 2", ProcessPart2, "aoc24/day09/input.txt")

	//lib.Profile(1, "day09.pprof", "Part 2", ProcessPart2, "aoc24/dayXX/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	disk := ReadAndParseInput(name)
	result := SolvePart1(disk)
	fmt.Println("Result:", result)
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	blocks := ParseBlocks(lines[0])
	Rearrange(blocks)
	fileIds := blocks.ToFileIdsArray()
	result := Checksum(fileIds)
	fmt.Println("Result:", result)
}

////////////////////////////////////////////////////////////

type Disk struct {
	Size    int
	FileIds []int
}

func ReadAndParseInput(name string) *Disk {
	lines := lib.ReadLines(name)
	assert.True(len(lines) == 1)
	return ParseInput(lines[0])
}
func ParseInput(line string) *Disk {
	diskSize := CountAll(line)
	fileIds := make([]int, diskSize)
	for i := 0; i < diskSize; i++ {
		fileIds[i] = -1
	}
	nextFileId := 0
	idx := 0
	for i, r := range line {
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
	return &Disk{Size: diskSize, FileIds: fileIds}
}

////////////////////////////////////////////////////////////

func SolvePart1(disk *Disk) int {
	fileIds := disk.FileIds
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

////////////////////////////////////////////////////////////

func SolvePart2(disk *Disk) int {
	return 0
}
