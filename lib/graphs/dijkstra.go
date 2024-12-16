package graphs

import (
	"fmt"
	"math"

	"github.com/phicode/challenges/lib"
	"github.com/phicode/challenges/lib/assert"
)

func DijkstraAll[T comparable](data []T, start func(a T) bool, neigh func(t T) []T) map[T]*NodeAll[T] {
	costOne := func(_, _ T) int { return 1 }
	return DijkstraAllWithCost(data, start, neigh, costOne)
}

func DijkstraAllWithCost[T comparable](
	data []T,
	start func(a T) bool,
	neigh func(t T) []T,
	cost func(T, T) int,
) map[T]*NodeAll[T] {

	nodes := make([]*NodeAll[T], len(data))
	nodeByKey := make(map[T]*NodeAll[T])
	nStart := 0
	for i, d := range data {
		node := &NodeAll[T]{Value: d, Distance: math.MaxInt, idx: i}
		nodes[i] = node
		nodeByKey[d] = node
		if start(d) {
			nodes[i].Distance = 0
			nStart++
		}
	}
	assert.True(nStart > 0)

	q := lib.NewHeapWithUpdater[*NodeAll[T]](nodes, nodeIsLess[T], updateNodeIndex[T])

	for q.Len() > 0 {
		u := q.Pop()
		u.visited = true

		for _, neighKey := range neigh(u.Value) {
			v := nodeByKey[neighKey]
			if v.visited {
				continue
			}
			alt := u.Distance + cost(u.Value, v.Value)
			assert.False(alt <= 0)
			if alt < v.Distance {
				v.Prev = []*NodeAll[T]{u}
				v.Distance = alt
				// changing the distance requires the heap to be fixed
				q.Fix(v.idx)
			} else if alt == v.Distance {
				v.Prev = append(v.Prev, u) // track all nodes with equal distance
			}
		}
	}

	return nodeByKey
}

type NodeAll[T any] struct {
	Value    T
	Distance int
	Prev     []*NodeAll[T]
	idx      int
	visited  bool
}

func (n *NodeAll[T]) String() string {
	return fmt.Sprintf("%v", n.Value)
}

//
//func (k *NodeAll[T]) GetPaths() [][]*NodeAll[T] {
//	var path []*NodeAll[T]
//	cur := k
//	for cur != nil {
//		path = append(path, cur)
//		assert.True(len(cur.Prev) > 0)
//		if len(cur.Prev) == 0 {
//			cur = cur.Prev[0]
//		}
//
//	}
//	slices.Reverse(path)
//	return path
//}

func nodeIsLess[T any](a, b *NodeAll[T]) bool {
	return a.Distance < b.Distance
}
func updateNodeIndex[T any](a *NodeAll[T], i int) {
	a.idx = i
}
