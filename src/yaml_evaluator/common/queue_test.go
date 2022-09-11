package common

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateQueue(t *testing.T) {
    que := CreateQueue[int]()
    assert.True(t, que.IsEmpty())
}

func TestFront(t *testing.T) {
    que := CreateQueue[int]()
    fst := 5
    que.Enqueue(fst)
    assert.Equal(t, 1, que.Size())
    e, r := que.Front()
    assert.Equal(t, fst, e)
    assert.True(t, r)
    assert.False(t, que.IsEmpty())
}

func TestDequeue(t *testing.T) {
    que := CreateQueue[int]()
    fst := 5
    que.Enqueue(fst)
    assert.Equal(t, 1, que.Size())
    e, r := que.Dequeue()
    assert.Equal(t, fst, e)
    assert.True(t, r)
    assert.True(t, que.IsEmpty())
}

func TestSequentialWithQueue(t *testing.T) {
    que := CreateQueue[int]()
    fst, snd, trd := 5, 7, 11
    que.Enqueue(fst)
    que.Enqueue(snd)
    que.Enqueue(trd)
    assert.Equal(t, 3, que.Size())
    expectedOrder := []int{fst, snd, trd}
    count := 0
    for e, r := que.Dequeue(); r; e, r = que.Dequeue() {
        assert.True(t, r)
        assert.Equal(t, expectedOrder[count], e)
        count++
        assert.Equal(t, 3, que.Size() + count)
    }
    assert.Equal(t, 3, count)
    assert.True(t, que.IsEmpty())
}

