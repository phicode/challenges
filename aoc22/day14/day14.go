package main

// https://adventofcode.com/2022/day/14

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

const EXTEND_AMOUNT = 500

func main() {
	Process("aoc22/day14/example.txt", false)
	Process("aoc22/day14/input.txt", false)
	Process("aoc22/day14/example.txt", true)
	Process("aoc22/day14/input.txt", true)
}

func Process(name string, ext bool) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	programs := ParsePrograms(lines)
	w := World{programs: programs}
	w.Build(ext)
	w.Print()
	sands := 0
	for {
		if !w.OneSandRound() {
			break
		}
		sands++
		//w.Print()
	}
	w.Print()
	fmt.Println("sands:", sands)

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

func ParsePrograms(lines []string) [][]Step {
	var progs [][]Step
	for _, line := range lines {
		progs = append(progs, ParseSteps(line))
	}
	return progs
}

func ParseSteps(line string) []Step {
	parts := strings.Split(line, " -> ")
	var steps []Step
	for _, part := range parts {
		var s Step
		n, err := fmt.Sscanf(part, "%d,%d", &s.X, &s.Y)
		if n != 2 || err != nil {
			panic(fmt.Errorf("failed to parse step, %w", err))
		}
		steps = append(steps, s)
	}
	return steps
}

type Step struct {
	X, Y int
}

type World struct {
	programs       [][]Step
	ShiftX, ShiftY int
	Cols, Rows     int
	data           []byte
}

func (w *World) Build(ext bool) {
	w.ScaleAll(ext)
	w.data = bytes.Repeat([]byte("."), w.Cols*w.Rows)
	for _, steps := range w.programs {
		w.RunSteps(steps)
	}
	if ext {
		for x := 0; x < w.Cols; x++ {
			w.data[w.idx(x, w.Rows-1)] = '#'
		}
	}
}

func (w *World) ScaleAll(ext bool) {
	minx, miny := 500, 0 // position of sand source
	for _, steps := range w.programs {
		for _, step := range steps {
			minx = min(minx, step.X)
			miny = min(miny, step.Y)
		}
	}
	if ext {
		minx -= EXTEND_AMOUNT
	}

	w.ShiftX, w.ShiftY = minx, miny
	var maxx, maxy int
	for _, steps := range w.programs {
		for i, _ := range steps {
			steps[i].X -= minx
			steps[i].Y -= miny
			maxx = max(maxx, steps[i].X)
			maxy = max(maxy, steps[i].Y)
		}
	}
	if ext {
		maxx += EXTEND_AMOUNT
		maxy += 2
	}
	w.Cols = maxx + 1
	w.Rows = maxy + 1
}

func (w *World) RunSteps(steps []Step) {
	for i := 0; i < len(steps)-1; i++ {
		a := steps[i]
		b := steps[i+1]
		if a.X == b.X { // vertical line
			y1 := min(a.Y, b.Y)
			y2 := max(a.Y, b.Y)
			for y := y1; y <= y2; y++ {
				w.data[w.idx(a.X, y)] = '#'
			}
		} else if a.Y == b.Y { // horizontal line
			x1 := min(a.X, b.X)
			x2 := max(a.X, b.X)
			for x := x1; x <= x2; x++ {
				w.data[w.idx(x, a.Y)] = '#'
			}
		} else {
			panic("diagonal line")
		}
	}
	w.data[w.idx(500-w.ShiftX, 0-w.ShiftY)] = '+'
}

func (w *World) idx(x, y int) int {
	return w.Cols*y + x
}

func (w *World) Print() {
	for i := 0; i < w.Rows; i++ {
		fmt.Println(string(w.data[i*w.Cols : (i+1)*w.Cols]))
	}
	fmt.Println()
}

func (w *World) OneSandRound() bool {
	x, y := 500-w.ShiftX, 0-w.ShiftY
	// move down while the space below is air
	for {
		if y+1 >= w.Rows {
			return false
		}
		if w.get(x, y+1) == '.' {
			y++
			continue
			//w.Print()
		}
		if x-1 < 0 {
			return false
		}
		// diagonal left
		if w.get(x-1, y+1) == '.' {
			x--
			y++
			//w.Print()
			continue // go back to first loop condition
		}
		if x+1 >= w.Cols {
			return false
		}
		// diagonal right
		if w.get(x+1, y+1) == '.' {
			x++
			y++
			//w.Print()
			continue // go back to first loop condition
		}
		break
	}
	if v := w.get(x, y); v == 'o' {
		return false
	}
	w.data[w.idx(x, y)] = 'o'
	return true
}

func (w *World) get(x, y int) byte {
	return w.data[w.idx(x, y)]
}
