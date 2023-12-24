package main

// https://adventofcode.com/2023/day/24

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc23/day24/example.txt", 7, 27)                          // 2
	ProcessPart1("aoc23/day24/input.txt", 200000000000000, 400000000000000) // 18184
	//
	//ProcessPart2("aoc23/day24/example.txt")
	//ProcessPart2("aoc23/day24/input.txt")
}

func ProcessPart1(name string, low, high float64) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	hailstones := ParseHailstones(lines)
	for _, h := range hailstones {
		fmt.Println(h)
	}
	SolvePart1(hailstones, low, high)
	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

func log(v int, msg string) {
	if v <= VERBOSE {
		fmt.Println(msg)
	}
}

////////////////////////////////////////////////////////////

type Ray struct {
	P, V V3
}

type V3 struct {
	X, Y, Z float64
}

func (a V3) String() string    { return fmt.Sprintf("(%.3f, %.3f, %.3f)", a.X, a.Y, a.Z) }
func (a V3) Add(b V3) V3       { return V3{a.X + b.X, a.Y + b.Y, a.Z + b.Z} }
func (a V3) Sub(b V3) V3       { return V3{a.X - b.X, a.Y - b.Y, a.Z - b.Z} }
func (a V3) Div(b V3) V3       { return V3{a.X / b.X, a.Y / b.Y, a.Z / b.Z} }
func (a V3) Mul(b V3) V3       { return V3{a.X * b.X, a.Y * b.Y, a.Z * b.Z} }
func (v V3) MulS(s float64) V3 { return V3{v.X * s, v.Y * s, v.Y * s} }

func ParseHailstones(lines []string) []Ray {
	var rv []Ray
	for _, l := range lines {
		rv = append(rv, ParseHailstone(l))
	}
	return rv
}

// px py pz @ vx vy vz
func ParseHailstone(l string) Ray {
	var hs Ray
	n, err := fmt.Sscanf(l, "%f, %f, %f @ %f, %f, %f", &hs.P.X, &hs.P.Y, &hs.P.Z, &hs.V.X, &hs.V.Y, &hs.V.Z)
	if n != 6 || err != nil {
		panic(fmt.Errorf("parse error, n=%d, err=%w", n, err))
	}
	return hs
}

func SolvePart1(hailstones []Ray, low, high float64) {
	n := 0
	for i, a := range hailstones {
		for _, b := range hailstones[i+1:] {
			p, ok := a.RaysIntersect(b)
			if ok {
				if a.Future(p) && b.Future(p) {
					if p.X >= low && p.X <= high && p.Y >= low && p.Y <= high {
						fmt.Printf("%v and %v collide INSIDE at %v\n", a.P, b.P, p)
						n++
					} else {
						fmt.Printf("%v and %v collide OUTSIDE at %v\n", a.P, b.P, p)
					}
				} else {
					fmt.Printf("%v and %v collide IN THE PAST at %v\n", a.P, b.P, p)
				}

			}
		}
	}
	fmt.Println("Collisions:", n)
}

func (a Ray) RaysIntersect(b Ray) (point V3, intersect bool) {
	var denom = (b.V.Y)*(a.V.X) - (b.V.X)*(a.V.Y)
	if denom == 0 {
		// parallel
		return V3{}, false
	}
	var num = (b.V.X)*(a.P.Y-b.P.Y) - (b.V.Y)*(a.P.X-b.P.X)
	// fraction of distance between a.{A.X,A.Y} => a.{B.X,B.Y} between 0 and 1
	var f = num / denom
	var p = V3{
		a.P.X + f*(a.V.X),
		a.P.Y + f*(a.V.Y),
		0,
	}
	return p, true
}

func (a Ray) Future(q V3) bool {
	// q = p + x * v
	// x  = (q - p) / v
	// if x > 0: position q is in the future
	x := q.Sub(a.P).Div(a.V)
	return x.X > 0 && x.Y > 0
}
