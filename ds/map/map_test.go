package _map

import (
	"reflect"
	"testing"
)

func TestMapKeys(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name    string
		args    args[K, V]
		wantRet []K
	}
	tests := []testCase[int, int]{
		{
			name:    "test1",
			args:    args[int, int]{m: map[int]int{1: 2}},
			wantRet: []int{1},
		},
		{
			name:    "test2",
			args:    args[int, int]{m: map[int]int{1: 2, 3: 4}},
			wantRet: []int{1, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := Keys(tt.args.m); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("Keys() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestMapMerge(t *testing.T) {
	type args[K comparable, V any] struct {
		a map[K]V
		b map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[int, string]{
		{
			name: "test1",
			args: args[int, string]{a: map[int]string{1: "2"}, b: map[int]string{3: "4"}},
			want: map[int]string{1: "2", 3: "4"},
		},
		{
			name: "test2",
			args: args[int, string]{a: map[int]string{1: "2"}, b: map[int]string{1: "4", 3: "4"}},
			want: map[int]string{1: "4", 3: "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Merge(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}

	type V struct {
		A int
		B string
	}

	tests2 := []testCase[int, *V]{
		{
			name: "test3",
			args: args[int, *V]{a: map[int]*V{1: {A: 1, B: "2"}}, b: map[int]*V{3: {A: 3, B: "4"}}},
			want: map[int]*V{1: {A: 1, B: "2"}, 3: {A: 3, B: "4"}},
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			got := Merge(tt.args.a, tt.args.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapOneKey(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name    string
		args    args[K, V]
		wantRet K
	}
	tests := []testCase[int, int]{
		{
			name:    "test1",
			args:    args[int, int]{m: map[int]int{1: 2}},
			wantRet: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := OneKey(tt.args.m); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("OneKey() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestMapOneVal(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name    string
		args    args[K, V]
		wantRet V
	}
	tests := []testCase[string, string]{
		{
			name:    "test1",
			args:    args[string, string]{m: map[string]string{"a": "b"}},
			wantRet: "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := OneVal(tt.args.m); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("MapOneVal() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestMapVals(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name    string
		args    args[K, V]
		wantRet []V
	}
	tests := []testCase[int64, int64]{
		{
			name:    "test1",
			args:    args[int64, int64]{m: map[int64]int64{1: 2}},
			wantRet: []int64{2},
		},
		{
			name:    "test2",
			args:    args[int64, int64]{m: map[int64]int64{1: 2, 3: 4}},
			wantRet: []int64{2, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := Vals(tt.args.m); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("Vals() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestMapMergeByFunc(t *testing.T) {
	type args[K comparable, V any] struct {
		a     map[K]V
		b     map[K]V
		merge func(V, V) V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[string, []int]{
		{
			name: "test1",
			args: args[string, []int]{
				a:     map[string][]int{"a": {1, 2}},
				b:     map[string][]int{"a": {3, 4}},
				merge: func(v, v2 []int) []int { return append(v, v2...) },
			},
			want: map[string][]int{"a": {1, 2, 3, 4}},
		},
		{
			name: "test2",
			args: args[string, []int]{
				a:     map[string][]int{"a": {1, 2}},
				b:     map[string][]int{"b": {3, 4}},
				merge: func(v, v2 []int) []int { return append(v, v2...) },
			},
			want: map[string][]int{"a": {1, 2}, "b": {3, 4}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeByFunc(tt.args.a, tt.args.b, tt.args.merge); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeByFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapIsEmpty(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want bool
	}
	tests := []testCase[int64, string]{
		{
			name: "test1",
			args: args[int64, string]{m: map[int64]string{}},
			want: true,
		},
		{
			name: "test2",
			args: args[int64, string]{m: map[int64]string{1: "2"}},
			want: false,
		},
		{
			name: "test3",
			args: args[int64, string]{m: nil},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.m); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapRemove(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
		k []K
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[int, string]{
		{
			name: "1",
			args: args[int, string]{
				m: map[int]string{1: "1", 2: "2", 3: "3"},
				k: []int{1, 2},
			},
			want: map[int]string{3: "3"},
		},
		{
			name: "2",
			args: args[int, string]{
				m: map[int]string{1: "1", 2: "2", 3: "3"},
				k: []int{1, 2, 4},
			},
			want: map[int]string{3: "3"},
		},
		{
			name: "3",
			args: args[int, string]{
				m: map[int]string{1: "1", 2: "2", 3: "3"},
				k: []int{1, 2, 3},
			},
			want: map[int]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Delete(tt.args.m, tt.args.k...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
