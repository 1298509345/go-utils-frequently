package go_utils_frequently

import (
	"reflect"
	"testing"

	_map "github.com/1298509345/go-utils-frequently/ds/map"
	"github.com/1298509345/go-utils-frequently/ds/slice"
)

func TestPtrOf(t *testing.T) {
	//var a any
	t.Log(*PtrOf(1))
	t.Log(*PtrOf("12"))
	t.Log(*PtrOf(struct {
		A string
	}{"dsa"}))
	t.Log(*PtrOf(any(nil)))
	ts := *PtrOf[[]string](nil)
	t.Log(ts == nil)
	t.Log(OfPtr(PtrOf("hh")))
}

func TestIf(t *testing.T) {
	var a int
	t.Log(If(a > 1, any(1), 2))
	t.Log(IfShortCircuit(a > 1, func() int { return 1 }, func() int { return 2 }))
}

func TestName(t *testing.T) {
	var tm = map[int][]string{
		1: {"1", "1"},
		2: {"2", "2"},
	}
	t.Log(slice.Merge(_map.Vals(tm)...))
	t.Log(slice.Merge(_map.Vals(tm)))
}

func TestP(t *testing.T) {
	t.Log("start")
	c1 := &CompexStruct{
		Array: []int64{1, 2, 3},
		Map: map[string]int{
			"1": 1,
			"2": 2,
		},
		Next: &CompexStruct{
			Array: []int64{4, 5, 6},
			Map: map[string]int{
				"2": 2,
			},
		},
	}
	c2 := &CompexStruct{
		Array: []int64{1, 2, 3},
		Map: map[string]int{
			"2": 2,
			"1": 1,
		},
		Next: &CompexStruct{
			Array: []int64{4, 5, 6},
			Map: map[string]int{
				"2": 2,
			},
		},
	}
	d := reflect.DeepEqual(c1, c2)
	t.Log("deep equal", d)
}

type CompexStruct struct {
	Array []int64
	Map   map[string]int
	Next  *CompexStruct
}
