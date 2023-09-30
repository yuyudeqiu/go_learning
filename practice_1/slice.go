package practice_1

import "fmt"

type List[T any] struct {
	list []T
}

func NewList[T any]() List[T] {
	return List[T]{}
}

func (l *List[T]) Print() {
	fmt.Println(l.list)
}

func (l *List[T]) Add(element T) bool {
	l.list = append(l.list, element)
	return true
}

func (l *List[T]) Delete(idx int) (T, bool) {
	tmp := l.list[idx]
	l.list = append(l.list[:idx], l.list[idx+1:]...)
	return tmp, true
}
