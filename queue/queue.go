// Package queue implements a FIFO queue generic over any type.
//
// The queue is not safe for concurrent access across goroutines, the caller is responsible for
// synchronising concurrent access.
package queue

import (
	"errors"
	"fmt"
	"iter"
)

// Queue is a FIFO queue generic over any type.
//
// A Queue should be instantiated by the New function and not directly.
type Queue[T any] struct {
	container []T // Underlying slice
}

// New constructs and returns a new Queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// From builds a [Queue] from an existing slice of items, pushing items
// into the queue in the order of the slice.
//
// The queue will be preallocated the size of len(items).
func From[T any](items []T) *Queue[T] {
	queue := &Queue[T]{container: make([]T, 0, len(items))}
	for _, item := range items {
		queue.Push(item)
	}

	return queue
}

// Collect builds a [Queue] from an iterator of items, pushing items
// into the queue in the order of iteration.
func Collect[T any](items iter.Seq[T]) *Queue[T] {
	queue := New[T]()
	for item := range items {
		queue.Push(item)
	}

	return queue
}

// Push adds an item to the back of the queue.
//
//	q := queue.New[string]()
//	q.Push("hello")
func (q *Queue[T]) Push(item T) {
	q.container = append(q.container, item)
}

// Pop removes an item from the front of the queue, if the queue
// is empty, an error will be returned.
//
//	q := queue.New[string]()
//	q.Push("hello")
//	q.Push("there")
//	item, _ := q.Pop()
//	fmt.Println(item) // "hello"
func (q *Queue[T]) Pop() (T, error) {
	l := len(q.container)
	if l == 0 {
		var none T
		return none, errors.New("pop from empty queue")
	}
	item := (q.container)[0]
	q.container = (q.container)[1:]

	return item, nil
}

// Size returns the number of items in the queue.
//
//	s := queue.New[string]()
//	s.Size() // 0
//	s.Push("hello")
//	s.Push("there")
//	s.Size() // 2
func (q *Queue[T]) Size() int {
	return len(q.container)
}

// Empty returns whether or not the queue is empty.
//
//	s := queue.New[string]()
//	s.Empty() // true
//	s.Push("hello")
//	s.Empty() // false
func (q *Queue[T]) Empty() bool {
	return len(q.container) == 0
}

// Items returns the an iterator over the queue in FIFO order.
//
//	q := queue.New[string]()
//	q.Push("hello")
//	q.Push("there")
//	qlices.Collect(s.Items()) // [hello there]
func (q *Queue[T]) Items() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range q.container {
			if !yield(item) {
				return
			}
		}
	}
}

// String satisfies the [fmt.Stringer] interface and allows a Queue to be printed.
func (q *Queue[T]) String() string {
	return fmt.Sprintf("%v", q.container)
}
