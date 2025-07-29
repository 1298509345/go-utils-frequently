package heap

import (
	"math/rand/v2"
	"sort"
	"testing"
)

func TestHeap_Add(t *testing.T) {
	// 创建堆
	mh := New(5, func(t1, t2 int) bool {
		return t1 <= t2 // 3 5 7 10 9
		// return t1 > t2 // 9 7 5 3 1
	})

	// 添加元素
	mh.Push(5)
	mh.Push(3)
	mh.Push(7)
	mh.Push(1)
	mh.Push(9)
	t.Log(mh.innerHeap.vals)
	mh.Push(10)
	t.Log(mh.innerHeap.vals)
}

func TestHeap(t *testing.T) {
	var randInput, randInput2 []int

	for i := 0; i < 20; i++ {
		randInput = append(randInput, rand.IntN(200))
	}
	for i := 0; i < 20; i++ {
		randInput2 = append(randInput2, rand.IntN(200))
	}

	t.Log(randInput)
	t.Log(randInput2)

	h1 := New(10, func(t1, t2 int) bool {
		return t1 <= t2
	})

	for _, v := range randInput {
		h1.Push(v)
	}
	t.Log(h1.innerHeap.vals)
	sort.IntSlice(randInput).Sort()
	t.Log(randInput)

	h2 := New(10, func(t1, t2 int) bool {
		return t1 > t2
	})

	for _, v := range randInput2 {
		h2.Push(v)
	}
	t.Log(h2.innerHeap.vals)

	sort.IntSlice(randInput2).Sort()
	t.Log(randInput2)
}

func findKthLargest(nums []int, k int) int {
	newh := New(k, func(t1, t2 int) bool {
		return t1 <= t2
	})

	for _, v := range nums {
		newh.Push(v)
	}
	pop, _ := newh.Pop()
	return pop
}
