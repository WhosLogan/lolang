package data

type Stack[T any] []any

func (st *Stack[T]) Len() int {
	return len(*st)
}

func (st *Stack[T]) Pop() T {
	v := (*st)[len(*st)-1]
	*st = (*st)[:len(*st)-1]
	return v
}

func (st *Stack[T]) Push(v T) {
	*st = append(*st, v)
}

func (st *Stack[T]) Empty() bool {
	return len(*st) == 0
}
