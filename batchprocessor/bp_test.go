package batchprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/1298509345/go-utils-frequently/optional"
	"math/rand/v2"
	"reflect"
	"strconv"
	"testing"
)

func TestBatchProcessor_Process(t *testing.T) {
	type args[T any] struct {
		ctx  context.Context
		data []T
	}
	type testCase[T any] struct {
		name    string
		bp      BatchProcessor[T]
		args    args[T]
		wantRes any
	}

	tests := []testCase[int64]{
		{
			name: "1",
			bp: BatchProcessor[int64]{
				ProcFunc: func(ctx context.Context, data []int64) error {
					for idx, d := range data {
						data[idx] = d + 1
					}
					return nil
				},
				BatchSize:        10,
				ConcurrencyLimit: 10,
			},
			args: args[int64]{
				ctx:  context.Background(),
				data: []int64{1, 2, 3, 4, 5, 6},
			},
			wantRes: []int64{2, 3, 4, 5, 6, 7},
		},
		{
			name: "2",
			bp: BatchProcessor[int64]{
				ProcFunc: func(ctx context.Context, data []int64) error {
					for idx, d := range data {
						data[idx] = d + 1
					}
					return nil
				},
				BatchSize: 2,
			},
			args: args[int64]{
				ctx:  context.Background(),
				data: make([]int64, 0),
			},
			wantRes: []int64{},
		},
		{
			name: "3",
			bp: BatchProcessor[int64]{
				ProcFunc: func(ctx context.Context, data []int64) error {
					for idx, d := range data {
						data[idx] = d + 1
					}
					return nil
				},
				BatchSize: 2,
			},
			args: args[int64]{
				ctx:  context.Background(),
				data: nil, // 这里赋值有个隐式转换，相当于 data: []int64(nil)
			},
			wantRes: []int64(nil), // 这里如果直接赋值 nil, reflect.DeepEqual会返回false
		},
	}

	tests2 := []testCase[string]{
		{
			name: "4",
			bp: BatchProcessor[string]{
				ProcFunc: func(ctx context.Context, data []string) error {
					for idx, d := range data {
						data[idx] = d + d
					}
					return nil
				},
				BatchSize:        2,
				ConcurrencyLimit: 5,
			},
			args: args[string]{
				ctx:  context.Background(),
				data: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
			},
			wantRes: []string{"00", "11", "22", "33", "44", "55", "66", "77", "88", "99"},
		},
	}

	type testCaseStruct struct {
		A int
		B string
	}

	tests3 := []testCase[*testCaseStruct]{
		{
			name: "5",
			bp: BatchProcessor[*testCaseStruct]{
				ProcFunc: func(ctx context.Context, data []*testCaseStruct) error {
					for idx, d := range data {
						data[idx].A = d.A + 1
						data[idx].B = d.B + d.B
					}
					return nil
				},
				BatchSize:        1,
				ConcurrencyLimit: 5,
			},
			args: args[*testCaseStruct]{
				ctx: context.Background(),
				data: []*testCaseStruct{
					{
						A: 1,
						B: "1",
					},
					{
						A: 2,
						B: "2",
					},
					{
						A: 3,
						B: "3",
					},
				},
			},
			wantRes: []*testCaseStruct{
				{
					A: 2,
					B: "11",
				},
				{
					A: 3,
					B: "22",
				},
				{
					A: 4,
					B: "33",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bp.Process(tt.args.ctx, tt.args.data); !reflect.DeepEqual(tt.args.data, tt.wantRes) {
				t.Errorf("Process() error = %v, args = %v, wantRes %v", err, tt.args.data, tt.wantRes)
			}
			str, _ := json.Marshal(tt.args.data)
			t.Log(string(str))
		})
	}

	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bp.Process(tt.args.ctx, tt.args.data); !reflect.DeepEqual(tt.args.data, tt.wantRes) {
				t.Errorf("Process() error = %v, args = %v, wantRes %v", err, tt.args.data, tt.wantRes)
			}
			t.Log(tt.args.data)
		})
	}

	for _, tt := range tests3 {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bp.Process(tt.args.ctx, tt.args.data); !reflect.DeepEqual(tt.args.data, tt.wantRes) {
				t.Errorf("Process() error = %v, args = %v, wantRes %v", err, tt.args.data, tt.wantRes)
			}
			str, _ := json.Marshal(tt.args.data)
			t.Log(string(str))
		})
	}
}

func TestBatchProcessor_ProcessWithFetcher(t *testing.T) {
	type args[T any] struct {
		ctx context.Context
		Fetcher[T]
		startPage int
	}
	type testCase[T any] struct {
		name    string
		bp      BatchProcessor[T]
		args    args[T]
		wantErr bool
	}

	// random []int64
	var randomInt64 = func() (ret []int64) {
		for i := 0; i < 10; i++ {
			ret = append(ret, rand.Int64N(100))
		}
		return ret
	}()
	var cpRandomInt64 = make([]int64, len(randomInt64))
	for i, v := range randomInt64 {
		cpRandomInt64[i] = v + 1
	}
	t.Log(randomInt64)

	tests := []testCase[int64]{
		{
			name: "2",
			bp: BatchProcessor[int64]{
				ProcFunc: func(ctx context.Context, data []int64) error {
					for idx, d := range data {
						data[idx] = d + 1
					}
					return nil
				},
				BatchSize: 5,
			},
			args: args[int64]{
				ctx: context.Background(),
				Fetcher: func(_ context.Context, page int, pageSize int) ([]int64, error) {
					offset := (page - 1) * pageSize
					if offset >= len(randomInt64) {
						return nil, nil
					}
					if offset+pageSize > len(randomInt64) {
						return randomInt64[offset:], nil
					}
					return randomInt64[offset : offset+pageSize], nil
				},
				startPage: 0,
			},
			wantErr: false,
		},
		{
			name: "3",
			bp: BatchProcessor[int64]{
				ProcFunc: func(ctx context.Context, data []int64) error {
					for idx, d := range data {
						data[idx] = d + 1
					}
					return nil
				},
				BatchSize: 5,
			},
			args: args[int64]{
				ctx: context.Background(),
				Fetcher: func(_ context.Context, page int, pageSize int) ([]int64, error) {
					offset := (page - 1) * pageSize
					if offset >= len(randomInt64) {
						return nil, nil
					}
					if offset+pageSize > len(randomInt64) {
						return randomInt64[offset:], nil
					}
					return randomInt64[offset : offset+pageSize], nil
				},
				startPage: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bp.ProcessFetcher(tt.args.ctx, tt.args.Fetcher, tt.args.startPage); (err != nil) != tt.wantErr {
				t.Errorf("ProcessWithFetcher() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(randomInt64)
			t.Log(cpRandomInt64)
			//t.Log(reflect.DeepEqual(randomInt64[tt.args.startPage*tt.bp.BatchSize:], cpRandomInt64[tt.args.startPage*tt.bp.BatchSize:]))
		})
	}
}

func TestBatchProcessor_ProcessFetcher(t *testing.T) {
	type args[T any] struct {
		ctx       context.Context
		fetcher   Fetcher[T]
		startPage int
	}
	type testCase[T any] struct {
		name    string
		bp      BatchProcessor[T]
		args    args[T]
		wantErr bool
	}
	var (
		batchCnt  int
		inputData = func() []string {
			var ret []string
			for i := 0; i < 15; i++ {
				ret = append(ret, strconv.FormatInt(int64(i), 10))
			}
			return ret
		}()
	)

	tests := []testCase[string]{
		{
			name: "1",
			bp: BatchProcessor[string]{
				ProcFunc: func(ctx context.Context, strings []string) error {
					batchCnt++
					formatInt := strconv.FormatInt(int64(batchCnt), 10)
					for idx, s := range strings {
						strings[idx] = fmt.Sprintf("data:%v,batch:%v", s, formatInt)
					}
					return nil
				},
				BatchSize: 5,
			},
			args: args[string]{
				ctx: context.Background(),
				fetcher: func(_ context.Context, page int, pageSize int) ([]string, error) {
					offset := (page - 1) * pageSize
					if offset >= len(inputData) {
						return nil, nil
					}
					if offset+pageSize > len(inputData) {
						return inputData[offset:], nil
					}
					return inputData[offset : offset+pageSize], nil
				},
				startPage: 0,
			},
		},
		{
			name: "2",
			bp: BatchProcessor[string]{
				ProcFunc: func(ctx context.Context, strings []string) error {
					batchCnt++
					formatInt := strconv.FormatInt(int64(batchCnt), 10)
					for idx, s := range strings {
						strings[idx] = fmt.Sprintf("data:%v,batch:%v", s, formatInt)
					}
					return nil
				},
				BatchSize: defaultBatchSize,
			},
			args: args[string]{
				ctx: context.Background(),
				fetcher: func(_ context.Context, page int, pageSize int) ([]string, error) {
					return nil, nil
				},
				startPage: 5,
			},
		},
		{
			name: "3",
			bp: BatchProcessor[string]{
				ProcFunc: func(ctx context.Context, strings []string) error {
					for idx, s := range strings {
						strings[idx] = fmt.Sprintf("<%v>", s)
					}
					t.Log(strings)
					return nil
				},
				BatchSize:        5,
				ConcurrencyLimit: 3,
			},
			args: args[string]{
				ctx: context.Background(),
				fetcher: func(_ context.Context, page int, pageSize int) ([]string, error) {
					offset := (page - 1) * pageSize
					if offset >= len(inputData) {
						return nil, nil
					}
					if offset+pageSize > len(inputData) {
						return inputData[offset:], nil
					}
					return inputData[offset : offset+pageSize], nil
				},
				startPage: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bp.ProcessFetcher(tt.args.ctx, tt.args.fetcher, tt.args.startPage); (err != nil) != tt.wantErr {
				t.Errorf("ProcessFetcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkFetcher(b *testing.B) {
	type args[T any] struct {
		ctx       context.Context
		fetcher   Fetcher[T]
		startPage int
	}
	type testCase[T any] struct {
		name    string
		bp      BatchProcessor[T]
		args    args[T]
		wantErr bool
	}
	var (
		inputData = func() []string {
			var ret []string
			for i := 0; i < 15; i++ {
				ret = append(ret, strconv.FormatInt(int64(i), 10))
			}
			return ret
		}()
	)

	tests := []testCase[string]{
		{
			name: "3",
			bp: BatchProcessor[string]{
				ProcFunc: func(ctx context.Context, strings []string) error {
					for idx, s := range strings {
						strings[idx] = fmt.Sprintf("<%v>", s)
					}
					b.Log(strings)
					return nil
				},
				BatchSize:        5,
				ConcurrencyLimit: 2,
			},
			args: args[string]{
				ctx: context.Background(),
				fetcher: func(_ context.Context, page int, pageSize int) ([]string, error) {
					ret := make([]string, pageSize)
					offset := (page - 1) * pageSize
					if offset >= len(inputData) {
						return nil, nil
					}
					if offset+pageSize > len(inputData) {
						copy(ret, inputData[offset:])
						return ret, nil
					}
					copy(ret, inputData[offset:offset+pageSize])
					return ret, nil
				},
				startPage: 0,
			},
		},
	}

	for i := 0; i < b.N; i++ {
		for _, tt := range tests {
			b.Run(tt.name, func(t *testing.B) {
				if err := tt.bp.ProcessFetcher(tt.args.ctx, tt.args.fetcher, tt.args.startPage); (err != nil) != tt.wantErr {
					t.Errorf("ProcessFetcher() error = %v, wantErr %v", err, tt.wantErr)
				}
				//time.Sleep(time.Second)
			})
		}
	}
}

func TestExit(t *testing.T) {
	var ch = make(chan struct{})
	close(ch)
	select {
	case <-ch:
		t.Log("1")
	default:
		t.Log("2")
	}

	var src = []string{"1", "2"}
	var dst = make([]string, 0, cap(src))
	copy(dst, src)
	t.Log(dst)

}

type cpInner struct {
	a int
}

type cpOuter struct {
	inner *cpInner
	b     string
	err   error
}

func (c *cpOuter) String() string {
	return fmt.Sprintf("inner:%v,b:%v", c.inner, &c.b)
}

func TestCopy(t *testing.T) {
	var a = []cpOuter{{inner: &cpInner{1}, b: "1", err: fmt.Errorf("hh")}}
	ch := make(chan *cpOuter)
	go func() {
		defer close(ch)
		ch <- &a[0]
		return
	}()

	var b *cpOuter
	for outer := range ch {
		b = outer
	}
	b.err = nil

	t.Log(a[0])
}

func TestNew(t *testing.T) {
	type args[T any] struct {
		options []optional.Op[BatchProcessor[T]]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want *BatchProcessor[T]
	}

	tests := []testCase[int]{
		{
			name: "1",
			args: args[int]{
				options: []optional.Op[BatchProcessor[int]]{
					WithBatchSize[int](11),
					WithConcurrencyLimit[int](9),
				},
			},
			want: &BatchProcessor[int]{
				BatchSize:        11,
				ConcurrencyLimit: 9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
