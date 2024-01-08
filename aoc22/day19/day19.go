package main

// https://adventofcode.com/2022/day/19

import (
	"fmt"
	"strings"

	"git.bind.ch/phil/challenges/lib"
)

var VERBOSE = 1

func main() {
	ProcessPart1("aoc22/day19/example.txt")
	ProcessPart1("aoc22/day19/input.txt")

	//ProcessPart2("aoc22/day19/example.txt")
	//ProcessPart2("aoc22/day19/input.txt")
}

func ProcessPart1(name string) {
	fmt.Println("Part 1 input:", name)
	lines := lib.ReadLines(name)
	blueprints := ParseBlueprints(lines)
	var sum int
	for _, bp := range blueprints {
		geodes := SolvePart1(bp)
		quality := geodes * bp.Num
		sum += quality
		fmt.Println("Blueprint", bp.Num, ", geodes", geodes, ", quality:", quality)
		fmt.Println("combinations:", combinations)
	}
	fmt.Println("Sum:", sum)

	fmt.Println()
}

func ProcessPart2(name string) {
	fmt.Println("Part 2 input:", name)
	lines := lib.ReadLines(name)
	_ = lines

	fmt.Println()
}

////////////////////////////////////////////////////////////

type Recipe struct {
	Produce Resource
	Cost    Resource
}

type Resource struct {
	Ore      int
	Clay     int
	Obsidian int
	Geode    int
}

func (r *Resource) Set(name string, value int) {
	switch name {
	case "ore":
		r.Ore = value
	case "clay":
		r.Clay = value
	case "obsidian":
		r.Obsidian = value
	case "geode":
		r.Geode = value
	default:
		panic(fmt.Errorf("invalid resource: %q", name))
	}
}

func (r *Resource) CanBuild(recipe Recipe) bool {
	return r.Ore >= recipe.Cost.Ore &&
		r.Clay >= recipe.Cost.Clay &&
		r.Obsidian >= recipe.Cost.Obsidian &&
		r.Geode >= recipe.Cost.Geode
}

func (r *Resource) MinutesUntil(s State, recipe Recipe) (bool, int) {
	oreOk, tOre := MinuteUntil(recipe.Cost.Ore, s.Resources.Ore, s.Robots.Ore)
	clayOk, tClay := MinuteUntil(recipe.Cost.Clay, s.Resources.Clay, s.Robots.Clay)
	obsOk, tObs := MinuteUntil(recipe.Cost.Obsidian, s.Resources.Obsidian, s.Robots.Obsidian)
	geoOk, tGeo := MinuteUntil(recipe.Cost.Geode, s.Resources.Geode, s.Robots.Geode)

	return oreOk && clayOk && obsOk && geoOk, max(tOre, tClay, tObs, tGeo)
}

func MinuteUntil(cost, current, income int) (bool, int) {
	if income == 0 && cost > current {
		return false, 0
	}
	if cost <= current {
		return true, 0
	}
	required := cost - current
	t := required / income
	if t*income < required {
		t++
	}
	return true, t
}

func (r *Resource) Sub(b Resource) {
	r.Ore -= b.Ore
	r.Clay -= b.Clay
	r.Obsidian -= b.Obsidian
	r.Geode -= b.Geode
}
func (r *Resource) Add(b Resource) {
	r.Ore += b.Ore
	r.Clay += b.Clay
	r.Obsidian += b.Obsidian
	r.Geode += b.Geode
}

type Blueprint struct {
	Num     int
	Recipes []Recipe
}

func ParseBlueprints(lines []string) []*Blueprint {
	var rv []*Blueprint
	var current *Blueprint
	for _, l := range lines {
		if strings.HasPrefix(l, "Blueprint") {
			current = &Blueprint{}
			// add the "empty" recipe
			//current.Recipes = append(current.Recipes, Recipe{})

			if n, err := fmt.Sscanf(l, "Blueprint %d:", &current.Num); n != 1 || err != nil {
				panic(fmt.Errorf("invalid input: %q", l))
			}
			rv = append(rv, current)
			continue
		}
		if l == "" {
			continue
		}
		l = strings.TrimSpace(l)
		var recipe Recipe
		// Production
		var res string
		if n, err := fmt.Sscanf(l, "Each %s robot costs", &res); n != 1 || err != nil {
			panic(fmt.Errorf("invalid input: %q", l))
		}
		recipe.Produce.Set(res, 1)
		// Cost
		produceAndCost := strings.Split(l, "costs")
		if len(produceAndCost) != 2 {
			panic("invalid input")
		}
		costs := strings.Split(produceAndCost[1], "and")
		for _, cost := range costs {
			cost = strings.TrimSpace(cost)
			var amount int
			if n, err := fmt.Sscanf(cost, "%d %s", &amount, &res); n != 2 || err != nil {
				panic(fmt.Errorf("invalid input: %q", cost))
			}
			res, _ = strings.CutSuffix(res, ".")
			recipe.Cost.Set(res, amount)
		}
		current.Recipes = append(current.Recipes, recipe)
	}
	return rv
}

////////////////////////////////////////////////////////////
// Part 1

type State struct {
	Robots    Resource
	Resources Resource
}

func (s State) RunMinute(bp *Blueprint, mem Memory, minute, end int) int {
	key := Key{State: s, Minute: minute}
	if memoization {
		if v, found := mem[key]; found {
			panic("asdf")
			return v
		}
	}
	// production phase
	var best int
	for _, recipe := range bp.Recipes {
		canbuild, minutesUntil := s.Resources.MinutesUntil(s, recipe)
		if !canbuild {
			continue
		}
		if minutesUntil == 0 {
			geodes := s.RunMinuteProduction(bp, mem, minute, end, recipe)
			best = max(best, geodes)
		} else {
			if minute+minutesUntil > end {
				continue
			}
			state := s.SkipMinutes(minutesUntil)
			geodes := state.RunMinuteProduction(bp, mem, minute+minutesUntil, end, recipe)
			best = max(best, geodes)
		}

	}
	if memoization {
		mem[key] = best
	}
	//fmt.Printf("Minute %d/%d geodes: %d\n", minute, end, geodes)
	return best
}

var combinations int

func (s State) RunMinuteProduction(bp *Blueprint, mem Memory, minute int, end int, recipe Recipe) int {
	s.Resources.Sub(recipe.Cost)
	s.Resources.Add(s.Robots)
	s.Robots.Add(recipe.Produce)
	combinations++
	if minute == end {
		return s.Resources.Geode
	}
	return s.RunMinute(bp, mem, minute+1, end)
}

func (s State) SkipMinutes(n int) State {
	for i := 0; i < n; i++ {
		s.Resources.Add(s.Robots)
	}
	return s
}

func SolvePart1(bp *Blueprint) int {
	var s State
	s.Robots.Ore = 1
	mem := make(Memory)
	v := s.RunMinute(bp, mem, 1, 24)
	fmt.Println("memo size:", len(mem))
	return v
}

const memoization = false

type Key struct {
	//RecipeIdx int
	State  State
	Minute int
}

type Memory map[Key]int

////////////////////////////////////////////////////////////
// Part 2
