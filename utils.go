package go_utils_frequently

func PtrOf[T any](t T) *T {
	return &t
}

// If 三目运算
func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfShortCircuit[T any](condition bool, trueVal func() T, falseVal func() T) T {
	if condition {
		return trueVal()
	}
	return falseVal()
}
