package heap

import (
	"container/heap"
)

// New 内建heap实现不限制size，且必须传入一个完整实现 heap.Interface 的实例，使用较繁琐（Less和类型绑定，无法实现带泛型通用方法）
func New[T any](size int, le func(t1, t2 T) bool) *Heap[T] {
	return &Heap[T]{
		maxSize: size,
		innerHeap: &inner[T]{
			vals: make([]T, 0, size+1),
			less: le,
		},
	}
}

type Heap[T any] struct {
	innerHeap *inner[T]
	maxSize   int
}

func (h *Heap[T]) Push(v T) {
	heap.Push(h.innerHeap, v)
	if h.innerHeap.Len() > h.maxSize {
		heap.Pop(h.innerHeap)
	}
}

func (h *Heap[T]) Pop() (T, bool) {
	if h.innerHeap.Len() == 0 {
		var zeroVal T
		return zeroVal, false
	}
	return heap.Pop(h.innerHeap).(T), true
}

func (h *Heap[T]) Peek() (T, bool) {
	if h.innerHeap.Len() == 0 {
		var zeroVal T
		return zeroVal, false
	}
	return h.innerHeap.vals[0], true
}

func (h *Heap[T]) Clear() {
	h.innerHeap.vals = h.innerHeap.vals[:0]
}

type inner[T any] struct {
	vals []T
	less func(i, j T) bool // return t1 <= t2 小的靠近堆顶
}

func (h *inner[T]) Len() int {
	return len(h.vals)
}

func (h *inner[T]) Less(i, j int) bool {
	return h.less(h.vals[i], h.vals[j])
}

func (h *inner[T]) Swap(i, j int) {
	h.vals[i], h.vals[j] = h.vals[j], h.vals[i]
}

func (h *inner[T]) Push(t any) {
	tt, ok := t.(T)
	if !ok {
		panic("inner heap push: type mismatch")
	}
	h.vals = append(h.vals, tt)
}

func (h *inner[T]) Pop() any {
	if h.Len() == 0 {
		return nil
	}
	t := h.vals[h.Len()-1]
	h.vals = h.vals[:h.Len()-1]
	return t
}

//func (h *Heap[T]) up(i int) {
//	for i > 0 {
//		parent := (i - 1) / 2
//		if h.less(h.Vals[parent], h.Vals[i]) {
//			break
//		}
//		h.Vals[parent], h.Vals[i] = h.Vals[i], h.Vals[parent]
//		i = parent
//	}
//}
//
//func (h *Heap[T]) down(i int) {
//	for {
//		left := 2*i + 1
//		right := 2*i + 2
//		smallest := i
//		if left < len(h.Vals) && h.less(h.Vals[left], h.Vals[smallest]) {
//			smallest = left
//		}
//		if right < len(h.Vals) && h.less(h.Vals[right], h.Vals[smallest]) {
//			smallest = right
//		}
//		if smallest == i {
//			break
//		}
//		h.Vals[i], h.Vals[smallest] = h.Vals[smallest], h.Vals[i]
//		i = smallest
//	}
//}
