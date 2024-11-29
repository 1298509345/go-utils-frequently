package slice

type Identifier[E any, ID comparable] func(E) ID

func IdentifierSelf[E comparable](e E) E {
	return e
}

// Remove 根据指定值删除
func Remove[E comparable](sl []E, spec E) []E {
	var res = make([]E, 0, len(sl))
	for _, item := range sl {
		if item != spec {
			res = append(res, item)
		}
	}
	return res
}

func RemoveDuplicate[T comparable](sl []T) []T {
	var (
		m      = make(map[T]struct{})
		result = make([]T, 0, len(sl))
	)

	for _, item := range sl {
		if _, ok := m[item]; !ok {
			m[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// RemoveDuplicateByID 根据标识符去重
func RemoveDuplicateByID[E any, ID comparable](sl []E, identifier Identifier[E, ID]) []E {
	var (
		m      = make(map[ID]struct{})
		result = make([]E, 0, len(sl))
	)

	for _, item := range sl {
		id := identifier(item)
		if _, ok := m[id]; !ok {
			m[id] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// RemoveByFunc 根据函数删除
func RemoveByFunc[T any](sl []T, remove func(T) bool) []T {
	var result = make([]T, 0, len(sl))
	for _, item := range sl {
		if !remove(item) {
			result = append(result, item)
		}
	}
	return result
}

// Contains 判断切片是否包含某个元素
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// ContainsByFunc 根据函数判断切片是否包含某个元素
func ContainsByFunc[T any](s []T, containsFunc func(T) bool) bool {
	for _, v := range s {
		if containsFunc(v) {
			return true
		}
	}
	return false
}

func Merge[T any](sls ...[]T) []T {
	var sumLen int
	for _, sl := range sls {
		sumLen += len(sl)
	}

	var ret = make([]T, 0, sumLen)
	for _, sl := range sls {
		ret = append(ret, sl...)
	}
	return ret
}

// Intersect 求两个切片的交集
func Intersect[E any, ID comparable](s1, s2 []E, identifier Identifier[E, ID]) []E {
	m := make(map[ID]struct{})
	for _, item := range s1 {
		id := identifier(item)
		m[id] = struct{}{}
	}

	var result = make([]E, 0, len(s2))
	for _, item := range s2 {
		id := identifier(item)
		if _, ok := m[id]; ok {
			result = append(result, item)
		}
	}
	return result
}

// DifferenceSet slice差集 s1-s2
func DifferenceSet[E any, ID comparable](s1, s2 []E, identifier Identifier[E, ID]) []E {
	m := make(map[ID]struct{})
	for _, item := range s2 {
		id := identifier(item)
		m[id] = struct{}{}
	}

	var result = make([]E, 0, len(s1))
	for _, item := range s1 {
		id := identifier(item)
		if _, ok := m[id]; !ok {
			result = append(result, item)
		}
	}
	return result
}

func GroupByID[E any, ID comparable](sl []E, identifier Identifier[E, ID]) map[ID][]E {
	m := make(map[ID][]E)
	for _, item := range sl {
		id := identifier(item)
		m[id] = append(m[id], item)
	}
	return m
}

func GroupByFunc[E, R any, ID comparable](sl []E, identifier Identifier[E, ID], groupFunc func(E, R) R) map[ID]R {
	m := make(map[ID]R)
	for _, item := range sl {
		id := identifier(item)
		tmp := m[id]
		m[id] = groupFunc(item, tmp)
	}
	return m
}

// FindValByFunc 根据函数查找, 没有则返回默认值
func FindValByFunc[T any](sl []T, fn func(T) bool) (T, bool) {
	for _, v := range sl {
		if fn(v) {
			return v, true
		}
	}
	var t T
	return t, false
}

func IsEmpty[T any](sl []T) bool {
	return len(sl) == 0
}

func DuplicateN[E any](e E, n int64) []E {
	var ret = make([]E, 0, n)
	for i := int64(0); i < n; i++ {
		ret = append(ret, e)
	}
	return ret
}

type Converter[T any, R any] func(T) (R, error)

func Convert[T any, R any](slice []T, converter Converter[T, R]) ([]R, error) {
	result := make([]R, len(slice))
	for i, v := range slice {
		converted, err := converter(v)
		if err != nil {
			return nil, err
		}
		result[i] = converted
	}
	return result, nil
}

func Filter[T any](slice []T, filter func(T) bool) []T {
	var result = make([]T, 0, len(slice))
	for _, v := range slice {
		if filter(v) {
			result = append(result, v)
		}
	}
	return result
}

func ForEach[T any](slice []T, forEach func(T)) {
	for _, v := range slice {
		forEach(v)
	}
}

func ToMap[E any, ID comparable](sl []E, id Identifier[E, ID]) map[ID]E {
	return GroupByFunc(sl, id, func(e E, _ E) E {
		return e
	})
}

func Make[E any](e ...E) []E {
	return e
}

func Max[E any](sl []E, less func(E, E) bool) (max E) {
	if len(sl) == 0 {
		return
	}
	max = sl[0]
	for i := 1; i < len(sl); i++ {
		if less(max, sl[i]) {
			max = sl[i]
		}
	}
	return
}

func Min[E any](sl []E, less func(E, E) bool) (min E) {
	if len(sl) == 0 {
		return
	}
	min = sl[0]
	for i := 1; i < len(sl); i++ {
		if less(sl[i], min) {
			min = sl[i]
		}
	}
	return
}
