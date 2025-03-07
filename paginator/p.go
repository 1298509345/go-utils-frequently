package paginator

// DataSource 数据源通用接口
type DataSource[T any] interface {
	GetChunk(offset int, limit int) ([]T, error)
	Total() int
}

type Paginator[T any] struct {
	sources  []DataSource[T]
	pageSize int
}

func NewPaginator[T any](sources []DataSource[T], pageSize int) *Paginator[T] {
	cpSource := make([]DataSource[T], len(sources))
	copy(cpSource, sources)
	return &Paginator[T]{
		sources:  cpSource,
		pageSize: pageSize,
	}
}

func (p *Paginator[T]) GetPage(pageNum int) ([]T, error) {
	if pageNum <= 0 {
		return nil, nil
	}
	var (
		skippedSoFar = pageNum * p.pageSize
		result       = make([]T, 0, p.pageSize)
	)

	// 遍历所有数据源获取元素
	for _, src := range p.sources {
		if skippedSoFar >= src.Total() {
			skippedSoFar -= src.Total()
			continue
		}

		// 计算当前数据源的有效区间
		remaining := src.Total() - skippedSoFar
		take := min(remaining, p.pageSize-len(result))

		chunk, err := src.GetChunk(skippedSoFar, take)
		if err != nil {
			return nil, err
		}
		result = append(result, chunk...)

		skippedSoFar = 0 // 后续数据源从0开始
		if len(result) == p.pageSize {
			break
		}
	}

	return result, nil
}
