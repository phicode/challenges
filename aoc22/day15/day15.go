package main

// https://adventofcode.com/2022/day/15

import (
	"bufio"
	"fmt"
	"os"

	"github.com/phicode/challenges/lib/assets"
)

func main() {
	ProcessPart1("aoc22/day15/example.txt", 10)
	ProcessPart1("aoc22/day15/input.txt", 2000000)
	ProcessPart2("aoc22/day15/example.txt", 20)
	ProcessPart2("aoc22/day15/input.txt", 4000000)
}

func ProcessPart1(name string, row int) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	sensors := ParseLines(lines)

	m := AccumulateLine(sensors, row)
	n := CountSharp(m)
	fmt.Println("no beacon in row", row, ":", n)

	fmt.Println()
}
func ProcessPart2(name string, size int) {
	fmt.Println("input:", name)
	lines := ReadInput(name)
	sensors := ParseLines(lines)

	//for y := 0; y <= size; y++ {
	//fmt.Printf("\r%d", y)
	s := Space{0, 0, size + 1, size + 1}
	space, b := BSP(sensors, s)
	fmt.Println("found solution", b)
	fmt.Println("found space", space)
	fmt.Println("tuning frequency", space.X*4_000_000+space.Y)

	fmt.Println()
}

func BSP(sensors []Sensor, space Space) (Space, bool) {
	parts := space.Split()
	for _, part := range parts {
		//if part.IsOneCell() {
		//	fmt.Println("bottom of stack", part)
		//}
		//if part.Contains(14, 11) {
		//	fmt.Println("this one should not match", part)
		//}
		if part.FullyCoveredAny(sensors) {
			//fmt.Println("fully covered", part)
			continue
		}
		if part.IsOneCell() {
			return part, true
		}
		s, ok := BSP(sensors, part)
		if ok {
			return s, ok
		}
	}
	return Space{}, false
}

func (s *Space) Split() [2]Space {
	if s.W < 2 && s.H < 2 {
		panic("too small")
	}
	if s.W >= s.H {
		wh := s.W / 2
		return [2]Space{
			{s.X, s.Y, wh, s.H},            // left
			{s.X + wh, s.Y, s.W - wh, s.H}, // right
		}
	}

	hh := s.H / 2
	return [2]Space{
		{s.X, s.Y, s.W, hh},            // top
		{s.X, s.Y + hh, s.W, s.H - hh}, // bottom
	}
}
func (s *Space) IsOneCell() bool {
	return s.W == 1 && s.H == 1
}

func (s *Space) FullyCoveredAny(sensors []Sensor) bool {
	for _, sensor := range sensors {
		if s.FullyCovered(sensor) {
			return true
		}
	}
	return false
}
func (s *Space) FullyCovered(sensor Sensor) bool {
	return sensor.ContainsXY(s.X, s.Y) &&
		sensor.ContainsXY(s.X+s.W-1, s.Y) &&
		sensor.ContainsXY(s.X+s.W-1, s.Y+s.H-1) &&
		sensor.ContainsXY(s.X, s.Y+s.H-1)
}

func (s *Space) Contains(x, y int) bool {
	return s.X <= x && x < s.X+s.W &&
		s.Y <= y && y < s.X+s.H
}

type Space struct {
	X, Y, W, H int
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

////////////////////////////////////////////////////////////

type Sensor struct {
	P      Pos
	Beacon Pos
	Reach  int
}

func (s Sensor) Contains(p Pos) bool {
	dist := s.P.DistanceTo(p)
	return dist <= s.Reach
}
func (s Sensor) ContainsXY(x, y int) bool {
	return s.Contains(Pos{x, y})
}

type Pos struct {
	X, Y int
}

func (p Pos) Add(x int, y int) Pos {
	return Pos{p.X + x, p.Y + y}
}

func (p Pos) DistanceTo(b Pos) int {
	return abs(p.X-b.X) + abs(p.Y-b.Y)
}

func ParseLines(lines []string) []Sensor {
	var sensors []Sensor
	for _, line := range lines {
		var sensor Sensor
		n, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
			&sensor.P.X, &sensor.P.Y, &sensor.Beacon.X, &sensor.Beacon.Y)
		if n != 4 || err != nil {
			panic("scan failed")
		}
		sensor.Reach = sensor.P.DistanceTo(sensor.Beacon)
		sensors = append(sensors, sensor)
	}
	return sensors
}

func AccumulateLine(sensors []Sensor, y int) map[int]byte {
	xs := make(map[int]byte)
	for _, sensor := range sensors {
		if sensor.Beacon.Y == y {
			xs[sensor.Beacon.X] = 'B'
		}
		if sensor.P.Y == y {
			xs[sensor.P.X] = 'S'
		}
	}
	for _, sensor := range sensors {
		reach := sensor.Reach
		distToRow := abs(sensor.P.Y - y)
		if distToRow > reach {
			continue
		}
		xMove := abs(reach - distToRow)
		for relx := -xMove; relx <= xMove; relx++ {
			absx := sensor.P.X + relx
			if _, ok := xs[absx]; !ok { // no value recorded yet
				xs[absx] = '#'
			}
		}
	}
	return xs
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func CountSharp(m map[int]byte) int {
	c := 0
	for _, v := range m {
		if v == '#' {
			c++
		}
	}
	return c
}
