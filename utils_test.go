package go_utils_frequently

import (
	"testing"
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
}

func TestIf(t *testing.T) {
	var a int
	t.Log(If(a > 1, any(1), 2))
	t.Log(IfShortCircuit(a > 1, func() int { return 1 }, func() int { return 2 }))
}
