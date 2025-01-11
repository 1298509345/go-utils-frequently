package _map

import (
	"maps"
	"slices"
)

// 内建方法 maps.Copy, maps.Clone ...

func MapKeys[K comparable, V any](m map[K]V) (ret []K) {
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

func MapVals[K comparable, V any](m map[K]V) (ret []V) {
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

func MapMerge[K comparable, V any](a, b map[K]V) map[K]V {
	return MapMergeByFunc(a, b, func(_, v V) V { return v })
}

func MapMergeByFunc[K comparable, V any](a, b map[K]V, merge func(old V, new V) V) map[K]V {
	ret := make(map[K]V, len(a))
	for k, v := range a {
		ret[k] = v
	}
	for k, v := range b {
		if old, ok := ret[k]; ok {
			ret[k] = merge(old, v)
		} else {
			ret[k] = v
		}
	}
	return ret
}

func MapOneKey[K comparable, V any](m map[K]V) (ret K) {
	for k := range m {
		return k
	}
	return
}

func MapOneVal[K comparable, V any](m map[K]V) (ret V) {
	for _, v := range m {
		return v
	}
	return
}

func MapIsEmpty[K comparable, V any](m map[K]V) bool {
	return len(m) == 0
}

func MapDelete[K comparable, V any](m map[K]V, ks ...K) map[K]V {
	maps.DeleteFunc(m, func(k K, _ V) bool { return slices.Contains(ks, k) })
	return m
}
