package lib

import (
	"container/heap"
)

type IndexUpdater[T any] func(t T, i int)

type Heap[T any] struct {
	impl heapImpl[T]
}

type heapImpl[T any] struct {
	less    func(a, b T) bool
	updater IndexUpdater[T]
	data    []T
}

func NewHeap[T any](data []T, less func(a, b T) bool) *Heap[T] {
	return NewHeapWithUpdater(data, less, nil)
}

func NewHeapWithUpdater[T any](data []T, less func(a, b T) bool, indexUpdater IndexUpdater[T]) *Heap[T] {
	h := &Heap[T]{
		impl: heapImpl[T]{
			less:    less,
			data:    data,
			updater: indexUpdater,
		},
	}
	if indexUpdater != nil {
		for i, t := range data {
			indexUpdater(t, i)
		}
	}
	heap.Init(&h.impl)
	return h
}

var _ heap.Interface = (*heapImpl[any])(nil)

func (h *heapImpl[T]) Len() int           { return len(h.data) }
func (h *heapImpl[T]) Less(i, j int) bool { return h.less(h.data[i], h.data[j]) }
func (h *heapImpl[T]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	if h.updater != nil {
		h.updater(h.data[i], i)
		h.updater(h.data[j], j)
	}
}

func (h *heapImpl[T]) Push(x any) {
	h.data = append(h.data, x.(T))
}
func (h *heapImpl[T]) Pop() any {
	l := len(h.data)
	v := h.data[l-1]
	if h.updater != nil {
		h.updater(v, -1)
	}
	h.data = h.data[:l-1]
	return v
}

func (h *Heap[T]) Push(v T) {
	heap.Push(&h.impl, v)
}
func (h *Heap[T]) Pop() T {
	return heap.Pop(&h.impl).(T)
}
func (h *Heap[T]) Fix(i int) {
	heap.Fix(&h.impl, i)
}

func (h *Heap[T]) Len() int {
	return h.impl.Len()
}
