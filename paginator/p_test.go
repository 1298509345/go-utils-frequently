package paginator

import (
	"testing"
)

type SliceSource[T any] struct {
	data []T
}

func (s *SliceSource[T]) GetChunk(offset, limit int) ([]T, error) {
	end := min(offset+limit, len(s.data))
	return s.data[offset:end], nil
}

func (s *SliceSource[T]) Total() int {
	return len(s.data)
}

func TestPag(t *testing.T) {
	src1 := &SliceSource[int]{data: []int{1, 2, 3}}
	src2 := &SliceSource[int]{data: []int{4, 5, 6, 7, 8, 9, 10}}
	src3 := &SliceSource[int]{data: []int{11, 12, 13}}

	paginator := NewPaginator[int]([]DataSource[int]{src1, src2, src3}, 3)

	// 分页验证
	for i := 0; true; i++ {
		page, _ := paginator.GetPage(i)
		if len(page) == 0 {
			break
		}
		t.Log(page)
	}
}
