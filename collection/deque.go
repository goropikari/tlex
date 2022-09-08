package collection

import "container/list"

type Deque[T any] struct {
	l *list.List
}

func NewDeque[T any]() *Deque[T] {
	return &Deque[T]{
		l: list.New(),
	}
}

func (q *Deque[T]) PushBack(x T) {
	q.l.PushBack(x)
}

func (q *Deque[T]) PushFront(x T) {
	q.l.PushFront(x)
}

func (q *Deque[T]) Front() T {
	e := q.l.Front()
	return e.Value.(T)
}

func (q *Deque[T]) Back() T {
	e := q.l.Back()
	return e.Value.(T)
}

func (q *Deque[T]) PopFront() {
	e := q.l.Front()
	q.l.Remove(e)
}

func (q *Deque[T]) PopBack() {
	e := q.l.Back()
	q.l.Remove(e)
}

func (q *Deque[T]) Size() int {
	return q.l.Len()
}
