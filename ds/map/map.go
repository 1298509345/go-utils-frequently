package _map

import (
	"maps"
	"slices"
)

// 内建方法 maps.Copy, maps.Clone ...

func Keys[K comparable, V any](m map[K]V) (ret []K) {
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

func Vals[K comparable, V any](m map[K]V) (ret []V) {
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

func Merge[K comparable, V any](a, b map[K]V) map[K]V {
	return MergeByFunc(a, b, func(_, v V) V { return v })
}

func MergeByFunc[K comparable, V any](a, b map[K]V, merge func(old V, new V) V) map[K]V {
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

func OneKey[K comparable, V any](m map[K]V) (ret K) {
	for k := range m {
		return k
	}
	return
}

func OneVal[K comparable, V any](m map[K]V) (ret V) {
	for _, v := range m {
		return v
	}
	return
}

func IsEmpty[K comparable, V any](m map[K]V) bool {
	return len(m) == 0
}

func Delete[K comparable, V any](m map[K]V, ks ...K) map[K]V {
	maps.DeleteFunc(m, func(k K, _ V) bool { return slices.Contains(ks, k) })
	return m
}
