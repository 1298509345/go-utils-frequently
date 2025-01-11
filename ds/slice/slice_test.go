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
			if got := DeleteDuplicateByID(tt.args.sl, tt.args.identifier); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteDuplicateByID() = %v, want %v", got, tt.want)
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
			got, got1 := FindByFunc(tt.args.sl, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByFunc() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindByFunc() got1 = %v, want %v", got1, tt.want1)
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
			if got := Delete(tt.args.sl, tt.args.spec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
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
			if got := Delete(tt.args.sl, tt.args.spec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
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
			if got := DeleteDuplicate(tt.args.sl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteDuplicate() = %v, want %v", got, tt.want)
			}
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
			if gotMax := Max(tt.args.less, tt.args.sl...); !reflect.DeepEqual(gotMax, tt.wantMax) {
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
			if gotMin := Min(tt.args.less, tt.args.sl...); !reflect.DeepEqual(gotMin, tt.wantMin) {
				t.Errorf("Min() = %v, want %v", gotMin, tt.wantMin)
			}
		})
	}
}

func TestRemoveByID(t *testing.T) {
	type args[E any, ID comparable] struct {
		sl         []E
		spec       E
		identifier Identifier[E, ID]
	}
	type testCase[E any, ID comparable] struct {
		name string
		args args[E, ID]
		want []E
	}
	tests := []testCase[book, int64]{
		{
			name: "test1",
			args: args[book, int64]{
				sl: []book{
					{ID: 1, Serials: "se1", Publisher: "a"},
					{ID: 2, Serials: "se1", Publisher: "a"},
				},
				spec:       book{ID: 1, Serials: "se1", Publisher: "a"},
				identifier: func(b book) int64 { return b.ID },
			},
			want: []book{{ID: 2, Serials: "se1", Publisher: "a"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteByID(tt.args.sl, tt.args.spec, tt.args.identifier); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup(t *testing.T) {
	type args[E comparable] struct {
		sl []E
	}
	type testCase[E comparable] struct {
		name string
		args args[E]
		want map[E][]E
	}
	tests := []testCase[string]{
		{
			name: "test1",
			args: args[string]{sl: []string{"a", "b", "a"}},
			want: map[string][]string{"a": {"a", "a"}, "b": {"b"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Group(tt.args.sl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	type args[E any, ID comparable] struct {
		sl         []E
		id         ID
		identifier Identifier[E, ID]
	}
	type testCase[E any, ID comparable] struct {
		name  string
		args  args[E, ID]
		want  E
		want1 bool
	}
	tests := []testCase[book, int64]{
		{
			name: "test1",
			args: args[book, int64]{
				sl: []book{
					{ID: 1, Serials: "se1", Publisher: "a"},
					{ID: 2, Serials: "se1", Publisher: "a"},
				},
				id:         1,
				identifier: func(b book) int64 { return b.ID },
			},
			want: book{ID: 1, Serials: "se1", Publisher: "a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FindByID(tt.args.sl, tt.args.id, tt.args.identifier)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindByID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFind(t *testing.T) {
	type args[T comparable] struct {
		sl []T
		e  T
	}
	type testCase[T comparable] struct {
		name  string
		args  args[T]
		want  T
		want1 bool
	}
	tests := []testCase[int]{
		{
			name:  "test1",
			args:  args[int]{sl: []int{1, 2, 3}, e: 1},
			want:  1,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Find(tt.args.sl, tt.args.e)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Find() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
