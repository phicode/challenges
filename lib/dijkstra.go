package lib

import (
	"math"
)

func Dijkstra[T comparable](data []T, start func(a T) bool, neigh func(t T) []T) map[T]*Node[T] {
	nodes := make([]*Node[T], len(data))
	nodeByValue := make(map[T]*Node[T])
	for i, d := range data {
		node := &Node[T]{Value: d, Distance: math.MaxInt, idx: i}
		nodes[i] = node
		nodeByValue[d] = node
		if start(d) {
			nodes[i].Distance = 0
		}
	}

	q := NewHeapWithUpdater[*Node[T]](nodes, nodeless, nodeupdater)

	for q.Len() > 0 {
		u := q.Pop()
		u.visited = true

		for _, v := range neigh(u.Value) {
			v := nodeByValue[v]
			if v.visited {
				continue
			}
			alt := u.Distance + 1
			if alt < v.Distance {
				v.Prev = u
				v.Distance = alt
				// changing the distance requires the heap to be fixed
				q.Fix(v.idx)
			}
		}
	}

	return nodeByValue
}

type Node[T any] struct {
	Value    T
	Distance int
	Prev     *Node[T]
	idx      int
	visited  bool
}

func nodeless[T any](a, b *Node[T]) bool {
	return a.Distance < b.Distance
}
func nodeupdater[T any](a *Node[T], i int) {
	a.idx = i
}
