package slice

import "slices"

type (
	Identifier[E any, ID comparable] func(E) ID
	Converter[T any, R any]          func(T) (R, error)
	Less[T any]                      func(t1 T, t2 T) bool // return t1 <= t2
)

// 内建方法：slices.Contains slices.Delete slices.Repeat slices.Concat slices.Max slices.Min ...

func IdentifierSelf[E comparable](e E) E {
	return e
}

// Delete 简化 slices.Delete
func Delete[E comparable](sl []E, spec E) []E {
	return DeleteByID(sl, spec, IdentifierSelf)
}

// DeleteByID 根据标识符删除
func DeleteByID[E any, ID comparable](sl []E, spec E, identifier Identifier[E, ID]) []E {
	return slices.DeleteFunc(sl, func(e E) bool { return identifier(e) == identifier(spec) })
}

func DeleteDuplicate[T comparable](sl []T) []T {
	return DeleteDuplicateByID(sl, IdentifierSelf)
}

// DeleteDuplicateByID 根据标识符去重
func DeleteDuplicateByID[E any, ID comparable](sl []E, identifier Identifier[E, ID]) []E {
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

func Filter[T any](slice []T, filter func(T) bool) []T {
	var result = make([]T, 0, len(slice))
	for _, v := range slice {
		if filter(v) {
			result = append(result, v)
		}
	}
	return result
}

func Group[E comparable](sl []E) map[E][]E {
	return GroupByID(sl, IdentifierSelf)
}

func GroupByID[E any, ID comparable](sl []E, identifier Identifier[E, ID]) map[ID][]E {
	return GroupByFunc(sl, identifier, func(e E, tmp []E) []E { return append(tmp, e) })
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

func Find[T comparable](sl []T, e T) (T, bool) {
	return FindByID(sl, e, IdentifierSelf)
}

func FindByID[E any, ID comparable](sl []E, id ID, identifier Identifier[E, ID]) (E, bool) {
	return FindByFunc(sl, func(e E) bool { return identifier(e) == id })
}

// FindByFunc 根据函数查找, 没有则返回默认值
func FindByFunc[T any](sl []T, fn func(T) bool) (T, bool) {
	var t T
	indexFunc := slices.IndexFunc(sl, fn)
	if indexFunc == -1 {
		return t, false
	}
	return sl[indexFunc], true
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

func ToMap[E any, ID comparable](sl []E, id Identifier[E, ID]) map[ID]E {
	return GroupByFunc(sl, id, func(e E, _ E) E { return e })
}

func Make[E any](e ...E) []E {
	return e
}

// Max 自定义less，内建 slices.Max 仅支持 cmp.Order 类型
func Max[E any](less Less[E], sl ...E) (max E) {
	return Min(func(e E, e2 E) bool { return less(e2, e) }, sl...)
}

// Min 同理 Max
func Min[E any](less Less[E], sl ...E) (min E) {
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

func IsEmpty[T any](sl []T) bool {
	return len(sl) == 0
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
