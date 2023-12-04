package heap

import (
	"arena"
)

func Append[T any](a *arena.Arena, slice []T, values ...T) []T {
	targetSize := len(slice) + len(values)
	if targetSize <= cap(slice) {
		return append(slice, values...)
	}
	var newSlice []T
	if a == nil {
		newSlice = make([]T, 0, targetSize*2)
	} else {
		newSlice = arena.MakeSlice[T](a, 0, targetSize*2)
	}
	return append(append(newSlice, slice...), values...)
}

func New[T any](a *arena.Arena) *T {
	if a == nil {
		return new(T)
	}
	return arena.New[T](a)
}
