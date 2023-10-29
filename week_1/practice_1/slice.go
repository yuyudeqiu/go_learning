package practice_1

import "fmt"

type List[T any] struct {
	list []T
}

func NewList[T any]() List[T] {
	return List[T]{}
}

func (l *List[T]) Print() {
	fmt.Printf("len: %d\tcap: %d\t val: %v\n", len(l.list), cap(l.list), l.list)
}

func (l *List[T]) Add(element T) bool {
	l.list = append(l.list, element)
	return true
}

func (l *List[T]) Delete(idx int) (tmp T, ok bool) {
	if idx < 0 || idx >= len(l.list) {
		return tmp, false
	}
	tmp = l.list[idx]
	//l.list = append(l.list[:idx], l.list[idx+1:]...)
	copy(l.list[idx:], l.list[idx+1:])
	l.list = l.list[:len(l.list)-1]

	// 缩容
	if len(l.list)*2 <= cap(l.list) && len(l.list) >= 1024 {
		newList := make([]T, len(l.list))
		copy(newList, l.list[:len(l.list)])
		l.list = newList
	}

	return tmp, true
}
