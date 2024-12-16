package lib

import (
	"math"
	"slices"
)

func Dijkstra[T comparable](data []T, start func(a T) bool, neigh func(t T) []T) map[T]*Node[T] {
	costOne := func(_, _ T) int { return 1 }
	return DijkstraWithCost(data, start, neigh, costOne)
}

func DijkstraWithCost[T comparable](data []T, start func(a T) bool, neigh func(t T) []T, cost func(T, T) int) map[T]*Node[T] {
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

	q := NewHeapWithUpdater[*Node[T]](nodes, nodeIsLess[T], updateNodeIndex[T])

	for q.Len() > 0 {
		u := q.Pop()
		u.visited = true

		for _, v := range neigh(u.Value) {
			v := nodeByValue[v]
			if v.visited {
				continue
			}
			alt := u.Distance + cost(u.Value, v.Value)
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

func (k *Node[T]) GetPath() []*Node[T] {
	var path []*Node[T]
	cur := k
	for cur != nil {
		path = append(path, cur)
		cur = cur.Prev
	}
	slices.Reverse(path)
	return path
}

func nodeIsLess[T any](a, b *Node[T]) bool {
	return a.Distance < b.Distance
}
func updateNodeIndex[T any](a *Node[T], i int) {
	a.idx = i
}
