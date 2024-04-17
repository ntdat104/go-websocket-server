package utils

func ForEach[T any](arr []T, fn func(T, int)) {
	for i, v := range arr {
		fn(v, i)
	}
}

func Map[T any, K any](arr []T, fn func(T, int) K) []K {
	result := make([]K, len(arr), cap(arr))
	for i, v := range arr {
		result[i] = fn(v, i)
	}
	return result
}

func Filter[T any](arr []T, fn func(T, int) bool) []T {
	result := make([]T, len(arr), cap(arr))
	for i, v := range arr {
		if fn(v, i) {
			result = append(result, v)
		}
	}
	return result
}

func Reduce[T any, K any](arr []T, fn func(K, T) K, initialValue K) K {
	accumulator := initialValue
	for _, v := range arr {
		accumulator = fn(accumulator, v)
	}
	return accumulator
}

func IndexOf[T comparable](arr []T, value T) int {
	for i, v := range arr {
		if v == value {
			return i
		}
	}
	return -1
}

func ToMap[T any, K comparable, V any](arr []T, fn func(T, int) (K, V)) map[K]V {
	resultMap := make(map[K]V, cap(arr))
	for i, v := range arr {
		key, value := fn(v, i)
		resultMap[key] = value
	}
	return resultMap
}
