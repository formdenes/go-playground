package helpers

// Map returns new array by fn
func Map[T, U any](els []T, fn func(T) U) []U {
	res := make([]U, len(els))
	for i, el := range els {
		res[i] = fn(el)
	}

	return res
}

func Reduce[T any, U any](els []T, fn func(U, T, int, []T) U, initVal U) U {
	if len(els) == 0 {
		var zero U
		return zero
	}
	res := initVal
	for i, el := range els {
		res = fn(res, el, i, els)
	}
	return res
}

func Contains[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
	var res []T
	for _, item := range slice {
		if predicate(item) {
			res = append(res, item)
		}
	}
	return res
}

func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}
