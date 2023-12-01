package lib

import (
	"container/heap"
)

type IndexUpdater[T any] func(t T, i int)

type Heap[T any] struct {
	impl himpl[T]
}

type himpl[T any] struct {
	less    func(a, b T) bool
	updater IndexUpdater[T]
	data    []T
}

func NewHeap[T any](data []T, less func(a, b T) bool) *Heap[T] {
	return NewHeapWithUpdater(data, less, nil)
}
func NewHeapWithUpdater[T any](data []T, less func(a, b T) bool, updater IndexUpdater[T]) *Heap[T] {
	h := &Heap[T]{
		impl: himpl[T]{
			less:    less,
			data:    data,
			updater: updater,
		},
	}
	if updater != nil {
		for i, t := range data {
			updater(t, i)
		}
	}
	heap.Init(&h.impl)
	return h
}

var _ heap.Interface = (*himpl[any])(nil)

func (h *himpl[T]) Len() int           { return len(h.data) }
func (h *himpl[T]) Less(i, j int) bool { return h.less(h.data[i], h.data[j]) }
func (h *himpl[T]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	if h.updater != nil {
		h.updater(h.data[i], i)
		h.updater(h.data[j], j)
	}
}

func (h *himpl[T]) Push(x any) {
	h.data = append(h.data, x.(T))
}
func (h *himpl[T]) Pop() any {
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
