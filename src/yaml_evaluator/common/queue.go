package common

type Queue[T any] struct {
    data []T
    size int
}

func CreateQueue[T any]() *Queue[T] {
    return &Queue[T]{}
}

func (que *Queue[T]) Enqueue(element T) {
    que.data = append(que.data, element)
    que.size++
}

func (que *Queue[T]) Front() (T, bool) {
    var e T
    if que.IsEmpty() {
        return e, false
    }
    e = que.data[0]
    return e, true
}

func (que *Queue[T]) Dequeue() (T, bool) {
    e, r := que.Front()
    if !r {
        return e, r
    }
    que.data = que.data[1:]
    que.size--
    return e, r
}

func (que *Queue[T]) Size() int {
    return que.size
}

func (que *Queue[T]) IsEmpty() bool {
    return que.Size() == 0
}
