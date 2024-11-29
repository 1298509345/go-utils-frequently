package batchprocessor

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

const (
	defaultBatchSize        = 100
	defaultConcurrencyLimit = 1
	maxConcurrencyLimit     = 20
)

type (
	Processor[T any] func(context.Context, []T) error
	Fetcher[T any]   func(ctx context.Context, page int, pageSize int) ([]T, error)

	BatchProcessor[T any] struct {
		ProcFunc         Processor[T]
		BatchSize        int
		ConcurrencyLimit int
	}
	Option[T any] func(*BatchProcessor[T])
)

func New[T any](options ...Option[T]) *BatchProcessor[T] {
	bp := &BatchProcessor[T]{}
	for _, option := range options {
		option(bp)
	}
	bp.init()

	return bp
}

func WithBatchSize[T any](batchSize int) Option[T] {
	return func(bp *BatchProcessor[T]) {
		bp.BatchSize = batchSize
	}
}

func WithProcessor[T any](proc Processor[T]) Option[T] {
	return func(bp *BatchProcessor[T]) {
		bp.ProcFunc = proc
	}
}

func WithConcurrencyLimit[T any](limit int) Option[T] {
	return func(bp *BatchProcessor[T]) {
		bp.ConcurrencyLimit = limit
	}
}

func (bp *BatchProcessor[T]) init() {
	if bp.BatchSize == 0 {
		bp.BatchSize = defaultBatchSize
	}
	if bp.ConcurrencyLimit <= 0 || bp.ConcurrencyLimit > maxConcurrencyLimit {
		bp.ConcurrencyLimit = defaultConcurrencyLimit
	}
}

func (bp *BatchProcessor[T]) Process(ctx context.Context, data []T) error {
	bp.init()

	eg := errgroup.Group{}
	eg.SetLimit(bp.ConcurrencyLimit)

	for start := 0; start < len(data); start += bp.BatchSize {
		end := start + bp.BatchSize
		if end > len(data) {
			end = len(data)
		}
		startCopy, endCopy := start, end
		eg.Go(func() error {
			if err := bp.ProcFunc(ctx, data[startCopy:endCopy]); err != nil {
				return fmt.Errorf("error processing batch from index %d to %d: %w", startCopy, endCopy, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

type batchInfo[T any] struct {
	batch []T
	page  int
	err   error
}

func (bp *BatchProcessor[T]) ProcessFetcher(ctx context.Context, fetcher Fetcher[T], startPage int) error {
	if fetcher == nil {
		return fmt.Errorf("no fetcher provided")
	}
	bp.init()

	var page = 1
	if startPage > page {
		page = startPage
	}

	var (
		eg      = errgroup.Group{}
		batches = make(chan batchInfo[T], bp.ConcurrencyLimit)
	)
	eg.SetLimit(bp.ConcurrencyLimit)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				batches <- batchInfo[T]{err: fmt.Errorf("panic:%v", err)}
			}
		}()
		defer close(batches)

		for ; ; page++ {
			oneBatch, err := fetcher(ctx, page, bp.BatchSize)
			if len(oneBatch) == 0 {
				break
			}
			batches <- batchInfo[T]{batch: oneBatch, page: page, err: err}
		}
	}()

	for oneBatch := range batches {
		curBatch := batchInfo[T]{
			err:   oneBatch.err,
			page:  oneBatch.page,
			batch: make([]T, len(oneBatch.batch)),
		}
		copy(curBatch.batch, oneBatch.batch)
		eg.Go(func() error {
			if curBatch.err != nil {
				return fmt.Errorf("error fetching curBatch: %w", curBatch.err)
			}
			if err := bp.ProcFunc(ctx, curBatch.batch); err != nil {
				return fmt.Errorf("error processing curBatch: %w, page: %v", err, curBatch.page)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
