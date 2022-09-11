package common

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateStack(t *testing.T) {
    stk := CreateStack[int]()
    assert.True(t, stk.IsEmpty())
}

func TestTop(t *testing.T) {
    stk := CreateStack[int]()
    fst := 5
    stk.Push(fst)
    assert.Equal(t, 1, stk.Size())
    e, r := stk.Top()
    assert.Equal(t, fst, e)
    assert.True(t, r)
    assert.False(t, stk.IsEmpty())
}

func TestPop(t *testing.T) {
    stk := CreateStack[int]()
    fst := 5
    stk.Push(fst)
    assert.Equal(t, 1, stk.Size())
    e, r := stk.Pop()
    assert.Equal(t, fst, e)
    assert.True(t, r)
    assert.True(t, stk.IsEmpty())
}

func TestSequentialWithStack(t *testing.T) {
    stk := CreateStack[int]()
    fst, snd, trd := 5, 7, 11
    stk.Push(fst)
    stk.Push(snd)
    stk.Push(trd)
    assert.Equal(t, 3, stk.Size())
    expectedOrder := []int{trd, snd, fst}
    count := 0
    for e, r := stk.Pop(); r; e, r = stk.Pop() {
        assert.True(t, r)
        assert.Equal(t, expectedOrder[count], e)
        count++
    }
    assert.Equal(t, 3, count)
    assert.True(t, stk.IsEmpty())
}

