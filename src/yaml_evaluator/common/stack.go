package common

type Stack[T any] struct {
    data []T
    size int
}

func CreateStack[T any]() *Stack[T] {
    return &Stack[T]{}
}

func (stk *Stack[T]) Push(element T) {
    stk.data = append(stk.data, element)
    stk.size++
}

func (stk *Stack[T]) Top() (T, bool) {
    var e T
    if stk.IsEmpty() {
        return e, false
    }
    e = stk.data[stk.Size() - 1]
    return e, true
}

func (stk *Stack[T]) Pop() (T, bool) {
    e, r := stk.Top()
    if !r {
        return e, r
    }
    stk.data = stk.data[:stk.Size() - 1]
    stk.size--
    return e, r
}

func (stk *Stack[T]) Size() int {
    return stk.size
}

func (stk *Stack[T]) IsEmpty() bool {
    return stk.Size() == 0
}
