package heap

// New 原生heap实现不限制size，且必须传入一个完整实现 heap.Interface 的实例，使用较繁琐（Less和类型绑定，无法实现带泛型通用方法）
func New[T any](size int, le func(t1, t2 T) bool) *Heap[T] {
	return &Heap[T]{
		Vals: make([]T, 0, size+1),
		Size: size,
		Less: le,
	}
}

type Heap[T any] struct {
	Vals []T
	Size int
	Less func(t1, t2 T) bool // return t1 <= t2 小的靠近堆顶
}

func (h *Heap[T]) Add(v T) {
	h.Vals = append(h.Vals, v)
	h.up(len(h.Vals) - 1)
	if len(h.Vals) > h.Size {
		h.Remove()
	}
}

func (h *Heap[T]) First() (T, bool) {
	var t T
	if len(h.Vals) > 0 {
		return h.Vals[0], true
	}
	return t, false
}

func (h *Heap[T]) up(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if h.Less(h.Vals[parent], h.Vals[i]) {
			break
		}
		h.Vals[parent], h.Vals[i] = h.Vals[i], h.Vals[parent]
		i = parent
	}
}

func (h *Heap[T]) Remove() (T, bool) {
	var t T
	if len(h.Vals) == 0 {
		return t, false
	}
	t = h.Vals[0]
	last := len(h.Vals) - 1
	h.Vals[0] = h.Vals[last]
	h.Vals = h.Vals[:last]
	if last > 0 {
		h.down(0)
	}
	return t, true
}

func (h *Heap[T]) down(i int) {
	for {
		left := 2*i + 1
		right := 2*i + 2
		smallest := i
		if left < len(h.Vals) && h.Less(h.Vals[left], h.Vals[smallest]) {
			smallest = left
		}
		if right < len(h.Vals) && h.Less(h.Vals[right], h.Vals[smallest]) {
			smallest = right
		}
		if smallest == i {
			break
		}
		h.Vals[i], h.Vals[smallest] = h.Vals[smallest], h.Vals[i]
		i = smallest
	}
}

func (h *Heap[T]) Clear() {
	h.Vals = h.Vals[:0]
}
