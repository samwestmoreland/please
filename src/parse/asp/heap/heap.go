package heap

import "arena"

func Append[T any](a *arena.Arena, s []T, ss ...T) []T {
	targetSize := len(s) + len(ss)
	if targetSize <= cap(s) || a == nil {
		return append(s, ss...)
	}
	ret := arena.MakeSlice[T](a, 0, targetSize*2)
	return append(append(ret, s...), ss...)
}
