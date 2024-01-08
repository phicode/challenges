package main

// https://adventofcode.com/2023/day/24

import (
	"fmt"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 2

func main() {
	ProcessPart1("aoc23/day24/example.txt", 7, 27) // 2
	VERBOSE = 1
	ProcessPart1("aoc23/day24/input.txt", 200000000000000, 400000000000000) // 18184

	VERBOSE = 2
	ProcessPart2("aoc23/day24/example.txt")
	VERBOSE = 1
	ProcessPart2("aoc23/day24/input.txt")
}

func ProcessPart1(name string, low, high float64) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	hailstones := ParseHailstones(lines)
	if VERBOSE >= 2 {
		for _, h := range hailstones {
			fmt.Println(h)
		}
	}
	SolvePart1(hailstones, low, high)
	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	raysI := ParseHailstonesI(lines)
	//rays := ParseHailstones(lines)
	if VERBOSE >= 2 {
		for _, h := range raysI {
			fmt.Println(h)
		}
	}

	a, b := raysI[0], raysI[1]
	findEqualVec(a, b, raysI[2:])
	//test(a, b, c, 5, 3)
	//test(b, a, c, 3, 5)

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
type RayI struct {
	P, V V3I
}

type V3 struct {
	X, Y, Z float64
}
type V3I struct {
	X, Y, Z int
}

func (a V3) String() string    { return fmt.Sprintf("(%.3f, %.3f, %.3f)", a.X, a.Y, a.Z) }
func (r Ray) String() string   { return fmt.Sprintf("%v @ %v", r.P, r.V) }
func (a V3) Add(b V3) V3       { return V3{a.X + b.X, a.Y + b.Y, a.Z + b.Z} }
func (a V3) Sub(b V3) V3       { return V3{a.X - b.X, a.Y - b.Y, a.Z - b.Z} }
func (a V3) Div(b V3) V3       { return V3{a.X / b.X, a.Y / b.Y, a.Z / b.Z} }
func (a V3) Mul(b V3) V3       { return V3{a.X * b.X, a.Y * b.Y, a.Z * b.Z} }
func (a V3) MulS(s float64) V3 { return V3{a.X * s, a.Y * s, a.Z * s} }

func (a V3I) String() string  { return fmt.Sprintf("(%d, %d, %d)", a.X, a.Y, a.Z) }
func (r RayI) String() string { return fmt.Sprintf("%v @ %v", r.P, r.V) }
func (a V3I) Add(b V3I) V3I   { return V3I{a.X + b.X, a.Y + b.Y, a.Z + b.Z} }
func (a V3I) Sub(b V3I) V3I   { return V3I{a.X - b.X, a.Y - b.Y, a.Z - b.Z} }
func (a V3I) Div(b V3I) V3I   { return V3I{a.X / b.X, a.Y / b.Y, a.Z / b.Z} }
func (a V3I) Mul(b V3I) V3I   { return V3I{a.X * b.X, a.Y * b.Y, a.Z * b.Z} }
func (a V3I) DivS(s int) V3I  { return V3I{a.X / s, a.Y / s, a.Z / s} }
func (a V3I) MulS(s int) V3I  { return V3I{a.X * s, a.Y * s, a.Z * s} }

func ParseHailstones(lines []string) []Ray {
	var rv []Ray
	for _, l := range lines {
		rv = append(rv, ParseHailstone(l))
	}
	return rv
}
func ParseHailstonesI(lines []string) []RayI {
	var rv []RayI
	for _, l := range lines {
		rv = append(rv, ParseHailstoneI(l))
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
func ParseHailstoneI(l string) RayI {
	var hs RayI
	n, err := fmt.Sscanf(l, "%d, %d, %d @ %d, %d, %d", &hs.P.X, &hs.P.Y, &hs.P.Z, &hs.V.X, &hs.V.Y, &hs.V.Z)
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
						if VERBOSE >= 2 {
							fmt.Printf("%v and %v collide INSIDE at %v\n", a.P, b.P, p)
						}
						n++
					} else {
						if VERBOSE >= 2 {
							fmt.Printf("%v and %v collide OUTSIDE at %v\n", a.P, b.P, p)
						}
					}
				} else {
					if VERBOSE >= 2 {
						fmt.Printf("%v and %v collide IN THE PAST at %v\n", a.P, b.P, p)
					}
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

func findEqualVec(a RayI, b RayI, others []RayI) {
outer:
	for t1 := 0; t1 < 1_000_000; t1++ {
		fmt.Println("t1", t1)
		for t2 := t1 + 1; t2 < 1_000_000; t2++ {
			candidate, ok, abort := test(a, b, t1, t2)
			if abort {
				//fmt.Println("aborting")
				continue outer
			}
			if !ok {
				continue
			}

			if times, ok := allmatch(others, candidate, t1, t2); ok {
				fmt.Println("found rock:", candidate, "t1:", t1, "t2:", t2, "times:", times)
				return
			}

			candidate, ok, abort = test(a, b, t2, t1)
			if abort {
				//fmt.Println("aborting")
				continue outer
			}
			if !ok {
				continue
			}

			if times, ok := allmatch(others, candidate, t1, t2); ok {
				fmt.Println("found rock:", candidate, "t1:", t1, "t2:", t2, "times:", times)
				return
			}
		}
	}
}

//func findEqualVecV2(a RayI, b RayI, others []RayI) {
//	for tdiff:= 1; tdiff < 100000; tdiff++ {
//			candidate, ok := testWithTimeDiff(a, b, tdiff)
//			if !ok {
//				continue
//			}
//
//			if times, ok := allmatch(others, candidate, t1, t2); ok {
//				fmt.Println("found rock:", candidate, "t1:", t1, "t2:", t2, "times:", times)
//			}
//		}
//	}
//}

func allmatch(others []RayI, candidate RayI, t1 int, t2 int) ([]int, bool) {
	var times []int
	for _, o := range others {
		if t, ok := Intersects3d(candidate, o); ok {
			times = append(times, t)
			fmt.Printf("t1=%v, t2=%d -> %v - also intersects %v at t=%d\n", t1, t2, candidate, o, t)
		} else {
			return nil, false
		}
	}
	return times, true
}

func test(a RayI, b RayI, t1 int, t2 int) (RayI, bool, bool) {
	// where a is at time t1
	a_t1 := a.P.Add(a.V.MulS(t1))
	b_t2 := b.P.Add(b.V.MulS(t2))
	vec := b_t2.Sub(a_t1)
	diffT := t2 - t1
	if !DividesEven(vec, diffT) {
		return RayI{}, false, Abort(vec, diffT)
	}

	// rock velocity
	rv := vec.DivS(diffT)
	// rock position at t0
	rp_t0 := a_t1.Sub(rv.MulS(t1))

	rock := RayI{rp_t0, rv}
	return rock, true, false
}

//func testWithTimeDiff(a RayI, b RayI, tdiff int) (RayI, bool) {
//	// where a is at time t1
//	a_t1 := a.P.Add(a.V.MulS(t1))
//	b_t2 := b.P.Add(b.V.MulS(t2))
//	vec := b_t2.Sub(a_t1)
//	diffT := t2 - t1
//	if !DividesEven(vec, diffT) {
//		return RayI{}, false
//	}
//
//	// rock velocity
//	rv := vec.DivS(diffT)
//	// rock position at t0
//	rp_t0 := a_t1.Sub(rv.MulS(t1))
//
//	rock := RayI{rp_t0, rv}
//	return rock, true
//}

func DividesEven(vec V3I, t int) bool {
	return vec.X%t == 0 &&
		vec.Y%t == 0 &&
		vec.Z%t == 0
}

func Abort(vec V3I, t int) bool {
	// these numbers can never divide with a zero modulo
	return abs(vec.X) < t || abs(vec.Y) < t || abs(vec.Z) < t
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func Intersects3d(r, h RayI) (int, bool) {
	// r.P + x*r.V   = h.P + x*h.V
	// x*r.V - x*h.V = h.P - r.P
	// x(r.V - h.V) = h.P - r.P
	// x = (h.P - r.P) / (r.V - h.V)

	p := h.P.Sub(r.P)
	q := r.V.Sub(h.V)
	if q.X == 0 || q.Y == 0 || q.Z == 0 {
		return 0, false
	}
	x := p.Div(q)

	if x.X == x.Y && x.X == x.Z {
		return x.X, true
	}
	return 0, false
}
