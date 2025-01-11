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
	mh.Add(5)
	mh.Add(3)
	mh.Add(7)
	mh.Add(1)
	mh.Add(9)
	t.Log(mh.Vals)
	mh.Add(10)
	t.Log(mh.Vals)
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
		h1.Add(v)
	}
	t.Log(h1.Vals)
	sort.IntSlice(randInput).Sort()
	t.Log(randInput)

	h2 := New(10, func(t1, t2 int) bool {
		return t1 > t2
	})

	for _, v := range randInput2 {
		h2.Add(v)
	}
	t.Log(h2.Vals)

	sort.IntSlice(randInput2).Sort()
	t.Log(randInput2)
}
