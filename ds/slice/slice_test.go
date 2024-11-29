package slice

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
)

func TestRemoveSliceDuplicate(t *testing.T) {
	type args[E any, ID comparable] struct {
		sl         []E
		identifier Identifier[E, ID]
	}
	type testCase[E any, ID comparable] struct {
		name string
		args args[E, ID]
		want []E
	}
	tests := []testCase[string, string]{
		{
			name: "test1",
			args: args[string, string]{sl: []string{"a", "b", "a"}, identifier: IdentifierSelf[string]},
			want: []string{"a", "b"},
		},
		{
			name: "test2",
			args: args[string, string]{sl: []string{"a", "b", "c"}, identifier: IdentifierSelf[string]},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicateByID(tt.args.sl, tt.args.identifier); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicateByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceContains(t *testing.T) {
	type args[T comparable] struct {
		s []T
		e T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "test1",
			args: args[int]{s: []int{1, 2, 3}, e: 1},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceConvert(t *testing.T) {
	type args[T any, R any] struct {
		slice     []T
		converter Converter[T, R]
	}
	type testCase[T any, R any] struct {
		name    string
		args    args[T, R]
		want    []R
		wantErr bool
	}
	tests := []testCase[string, int]{
		{
			name: "test1",
			args: args[string, int]{slice: []string{"1", "2", "3"}, converter: func(s string) (int, error) {
				tmp, _ := strconv.ParseInt(s, 10, 64)
				return int(tmp), nil
			}},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.args.slice, tt.args.converter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceDifferenceSet(t *testing.T) {
	type args[E any, ID comparable] struct {
		s1         []E
		s2         []E
		identifier Identifier[E, ID]
	}
	type testCase[E any, ID comparable] struct {
		name string
		args args[E, ID]
		want []E
	}
	tests := []testCase[string, string]{
		{
			name: "test1",
			args: args[string, string]{s1: []string{"a", "b", "c"}, s2: []string{"a", "b"}, identifier: IdentifierSelf[string]},
			want: []string{"c"},
		},
		{
			name: "test2",
			args: args[string, string]{s1: []string{"a", "b", "c"}, s2: []string{"a", "b", "c"}, identifier: IdentifierSelf[string]},
			want: []string{},
		},
		{
			name: "test3",
			args: args[string, string]{s1: []string{"a", "b", "c"}, s2: []string{"d", "e", "f"}, identifier: IdentifierSelf[string]},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DifferenceSet(tt.args.s1, tt.args.s2, tt.args.identifier); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DifferenceSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceDuplicateN(t *testing.T) {
	type args[E any] struct {
		e E
		n int64
	}
	type testCase[E any] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[int64]{
		{
			name: "test1",
			args: args[int64]{e: 1, n: 3},
			want: []int64{1, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DuplicateN(tt.args.e, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DuplicateN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceFindValByFunc(t *testing.T) {
	type args[T any] struct {
		sl []T
		fn func(T) bool
	}
	type testCase[T any] struct {
		name  string
		args  args[T]
		want  T
		want1 bool
	}
	tests := []testCase[int]{
		{
			name: "test1",
			args: args[int]{sl: []int{1, 2, 3}, fn: func(i int) bool {
				return i == 1
			}},
			want:  1,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FindValByFunc(tt.args.sl, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindValByFunc() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindValByFunc() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSliceGroupByID(t *testing.T) {
	type args[E any, ID comparable] struct {
		sl         []E
		identifier Identifier[E, ID]
	}
	type testCase[E any, ID comparable] struct {
		name string
		args args[E, ID]
		want map[ID][]E
	}
	tests := []testCase[int, int]{
		{
			name: "test1",
			args: args[int, int]{sl: []int{1, 2, 3, 3}, identifier: IdentifierSelf[int]},
			want: map[int][]int{1: {1}, 2: {2}, 3: {3, 3}},
		},
		{
			name: "test2",
			args: args[int, int]{sl: []int{5764858630966808144, 5764859556309959697}, identifier: IdentifierSelf[int]},
			want: map[int][]int{5764858630966808144: {5764858630966808144}, 5764859556309959697: {5764859556309959697}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupByID(tt.args.sl, tt.args.identifier); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceIntersect(t *testing.T) {
	type args[E any, ID comparable] struct {
		s1         []E
		s2         []E
		identifier Identifier[E, ID]
	}
	type testCase[E any, ID comparable] struct {
		name string
		args args[E, ID]
		want []E
	}
	tests := []testCase[int, int]{
		{
			name: "test1",
			args: args[int, int]{s1: []int{1, 2, 3}, s2: []int{1, 2}, identifier: IdentifierSelf[int]},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intersect(tt.args.s1, tt.args.s2, tt.args.identifier); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceIsEmpty(t *testing.T) {
	type args[T any] struct {
		sl []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "test1",
			args: args[int]{sl: []int{}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.sl); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceRemove(t *testing.T) {
	type args[E comparable] struct {
		sl   []E
		spec E
	}
	type testCase[E comparable] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[string]{
		{
			name: "test1",
			args: args[string]{sl: []string{"a", "b", "c"}, spec: "a"},
			want: []string{"b", "c"},
		},
		{
			name: "test2",
			args: args[string]{sl: []string{"a", "b", "c", "b"}, spec: "b"},
			want: []string{"a", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Remove(tt.args.sl, tt.args.spec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}

	type str struct {
	}

	tests2 := []testCase[*str]{
		{
			name: "test1",
			args: args[*str]{sl: []*str{&str{}, &str{}, nil}, spec: nil},
			want: []*str{&str{}, &str{}},
		},
	}

	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if got := Remove(tt.args.sl, tt.args.spec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceRemoveByFunc(t *testing.T) {
	type args[T any] struct {
		sl     []T
		remove func(T) bool
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[string]{
		{
			name: "test1",
			args: args[string]{sl: []string{"a", "b", "c"}, remove: func(s string) bool {
				return s == "a"
			}},
			want: []string{"b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveByFunc(tt.args.sl, tt.args.remove); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveByFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

type book struct {
	ID        int64
	Serials   string
	Publisher string
}

func TestSliceGroupByFunc(t *testing.T) {

	type args[E any, ID comparable, R any] struct {
		sl         []E
		identifier Identifier[E, ID]
		groupFunc  func(E, R) R
	}
	type testCase[E any, ID comparable, R any] struct {
		name string
		args args[E, ID, R]
		want map[ID]R
	}
	tests := []testCase[book, string, map[string][]int64]{
		{
			name: "test1",
			args: args[book, string, map[string][]int64]{
				sl: []book{
					{ID: 1, Serials: "se1", Publisher: "a"},
					{ID: 2, Serials: "se1", Publisher: "a"},
					{ID: 3, Serials: "se1", Publisher: "b"},
					{ID: 4, Serials: "se1", Publisher: "c"},
					{ID: 5, Serials: "se2", Publisher: "c"},
					{ID: 6, Serials: "se2", Publisher: "c"},
				},
				identifier: func(book book) string {
					return book.Publisher
				},
				groupFunc: func(e book, m map[string][]int64) map[string][]int64 {
					if m == nil {
						return map[string][]int64{e.Serials: {e.ID}}
					}
					m[e.Serials] = append(m[e.Serials], e.ID)
					return m
				},
			},
			want: map[string]map[string][]int64{
				"a": {
					"se1": {1, 2},
				},
				"b": {
					"se1": {3},
				},
				"c": {
					"se1": {4},
					"se2": {5, 6},
				},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupByFunc(tt.args.sl, tt.args.identifier, tt.args.groupFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupByFunc() = %v, want %v", JsonString(got), JsonString(tt.want))
			}
		})
	}
}

// json string
func JsonString(v any) string {

	b, _ := json.Marshal(v)
	return string(b)
}

func TestSliceContainsByFunc(t *testing.T) {
	type args[T any] struct {
		s            []T
		containsFunc func(T) bool
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[book]{
		{
			name: "test1",
			args: args[book]{
				s: []book{
					{ID: 1, Serials: "se1", Publisher: "a"},
					{ID: 2, Serials: "se2", Publisher: "a"},
				},
				containsFunc: func(b book) bool {
					return b.ID == 1 && b.Serials == "se1" && b.Publisher == "a"
				}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsByFunc(tt.args.s, tt.args.containsFunc); got != tt.want {
				t.Errorf("ContainsByFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceMerge(t *testing.T) {
	type args[T any] struct {
		sls [][]T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[string]{
		{
			name: "test1",
			args: args[string]{sls: [][]string{{"a", "b"}, {"c", "d"}}},
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "test2",
			args: args[string]{sls: [][]string{{"a", "b"}, {"c", "d"}, {"e", "f"}}},
			want: []string{"a", "b", "c", "d", "e", "f"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Merge(tt.args.sls...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceFilter(t *testing.T) {
	type args[T any] struct {
		slice  []T
		filter func(T) bool
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "test1",
			args: args[int]{
				slice:  []int{1, 2, 3},
				filter: func(i int) bool { return false },
			},
			want: []int{},
		},
		{
			name: "test2",
			args: args[int]{
				slice:  []int{1, 2, 3},
				filter: func(i int) bool { return true },
			},
			want: []int{1, 2, 3},
		},
		{
			name: "test3",
			args: args[int]{
				slice:  []int{1, 2, 3},
				filter: func(i int) bool { return i%2 == 0 },
			},
			want: []int{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.slice, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceRemoveDuplicateByID(t *testing.T) {
	type args[T comparable] struct {
		sl []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[string]{
		{
			name: "test1",
			args: args[string]{sl: []string{"a", "b", "a"}},
			want: []string{"a", "b"},
		},
		{
			name: "test2",
			args: args[string]{sl: []string{"a", "b", "c"}},
			want: []string{"a", "b", "c"},
		},
		{
			name: "test3",
			args: args[string]{sl: []string{}},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicate(tt.args.sl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceForEach(t *testing.T) {
	type args[T any] struct {
		slice   []T
		forEach func(T)
	}
	type testCase[T any] struct {
		name string
		args args[T]
	}
	tests := []testCase[string]{
		{
			name: "test1",
			args: args[string]{slice: []string{"a", "b", "c"}, forEach: func(s string) {
				t.Log(s)
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ForEach(tt.args.slice, tt.args.forEach)
		})
	}
}

func TestSlice2Map(t *testing.T) {
	type args[E any, ID comparable] struct {
		sl []E
		id Identifier[E, ID]
	}
	type testCase[E any, ID comparable] struct {
		name string
		args args[E, ID]
		want map[ID]E
	}
	tests := []testCase[string, int64]{
		{
			name: "11",
			args: args[string, int64]{
				sl: []string{"1", "2", "3"},
				id: func(s string) int64 {
					i, _ := strconv.ParseInt(s, 10, 64)
					return i
				},
			},
			want: map[int64]string{
				1: "1",
				2: "2",
				3: "3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToMap(tt.args.sl, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceGroupByFunc1(t *testing.T) {
	type args[E any, ID comparable, R any] struct {
		sl         []E
		identifier Identifier[E, ID]
		groupFunc  func(E, R) R
	}
	type testCase[E any, ID comparable, R any] struct {
		name string
		args args[E, ID, R]
		want map[ID]R
	}
	tests := []testCase[book, int64, book]{
		{
			name: "111",
			args: args[book, int64, book]{
				sl: []book{
					{1, "s1", "p1"},
					{2, "s2", "p2"},
				},
				identifier: func(b book) int64 {
					return b.ID
				},
				groupFunc: func(b book, _ book) book {
					return b
				},
			},
			want: map[int64]book{
				1: {1, "s1", "p1"},
				2: {2, "s2", "p2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupByFunc(tt.args.sl, tt.args.identifier, tt.args.groupFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupByFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceMake(t *testing.T) {
	type args[E any] struct {
		e []E
	}
	type testCase[E any] struct {
		name string
		args args[E]
		want []E
	}
	tests := []testCase[string]{
		{
			name: "test1",
			args: args[string]{e: []string{"a", "b", "c"}},
			want: []string{"a", "b", "c"},
		},
		{
			name: "test2",
			args: args[string]{e: []string{}},
			want: []string{},
		},
		{
			name: "test3",
			args: args[string]{e: nil},
			want: nil,
		},
		{
			name: "test4",
			args: args[string]{e: []string{"a"}},
			want: []string{"a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Make(tt.args.e...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Make() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceMax(t *testing.T) {
	type args[E any] struct {
		sl   []E
		less func(E, E) bool
	}
	type testCase[E any] struct {
		name    string
		args    args[E]
		wantMax E
	}
	tests := []testCase[int]{
		{
			name: "test1",
			args: args[int]{sl: []int{1, 2, 3}, less: func(i, i2 int) bool {
				return i < i2
			}},
			wantMax: 3,
		},
		{
			name: "test2",
			args: args[int]{sl: []int{1, 2, 3}, less: func(i, i2 int) bool {
				return i > i2
			}},
			wantMax: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMax := Max(tt.args.sl, tt.args.less); !reflect.DeepEqual(gotMax, tt.wantMax) {
				t.Errorf("Max() = %v, want %v", gotMax, tt.wantMax)
			}
		})
	}
}

func TestSliceMin(t *testing.T) {
	type args[E any] struct {
		sl   []E
		less func(E, E) bool
	}
	type testCase[E any] struct {
		name    string
		args    args[E]
		wantMin E
	}
	tests := []testCase[int]{
		{
			name: "test1",
			args: args[int]{sl: []int{1, 2, 3}, less: func(i, i2 int) bool {
				return i < i2
			}},
			wantMin: 1,
		},
		{
			name: "test2",
			args: args[int]{sl: []int{1, 2, 3}, less: func(i, i2 int) bool {
				return i > i2
			}},
			wantMin: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMin := Min(tt.args.sl, tt.args.less); !reflect.DeepEqual(gotMin, tt.wantMin) {
				t.Errorf("Min() = %v, want %v", gotMin, tt.wantMin)
			}
		})
	}
}
