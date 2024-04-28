// Package queue implements a FIFO queue generic over any type.
//
// The queue is not safe for concurrent access across goroutines, the caller is responsible for
// synchronising concurrent access.
package queue

import (
	"errors"
	"fmt"
)

// ErrPopFromEmptyQueue is returned when Pop() is called on an empty queue.
var ErrPopFromEmptyQueue = errors.New("pop from empty queue")

// Queue is a FIFO queue generic over any type.
//
// A Queue should be instantiated by the New function and not directly.
type Queue[T any] struct {
	container []T     // Underlying slice
	options   options // Options for the queue.
}

// New constructs and returns a new Queue.
func New[T any](options ...Option) *Queue[T] {
	queue := Queue[T]{}
	for _, option := range options {
		option(&queue.options)
	}
	container := make([]T, 0, queue.options.capacity)
	queue.container = container
	return &queue
}

// Push adds an item to the back of the queue.
//
//	q := queue.New[string]()
//	q.Push("hello")
func (q *Queue[T]) Push(item T) {
	q.container = append(q.container, item)
}

// Pop removes an item from the front of the queue, if the queue
// is empty, ErrPopFromEmptyQueue will be returned.
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
		return none, ErrPopFromEmptyQueue
	}
	item := (q.container)[0]
	q.container = (q.container)[1:]

	return item, nil
}

// Length returns the number of items in the queue.
//
//	s := queue.New[string]()
//	s.Length() // 0
//	s.Push("hello")
//	s.Push("there")
//	s.Length() // 2
func (q *Queue[T]) Length() int {
	return len(q.container)
}

// Cap returns the current capacity of the queue.
//
//	s := queue.New[string](queue.WithCapacity(10))
//	s.Cap() // 10
func (q *Queue[T]) Cap() int {
	return cap(q.container)
}

// IsEmpty returns whether or not the queue is empty.
//
//	s := queue.New[string]()
//	s.IsEmpty() // true
//	s.Push("hello")
//	s.IsEmpty() // false
func (q *Queue[T]) IsEmpty() bool {
	return len(q.container) == 0
}

// Items returns the items in the queue as a new slice (copy).
//
//	q := queue.New[string]()
//	q.Push("hello")
//	q.Push("there")
//	q.Items() // [hello there]
func (q *Queue[T]) Items() []T {
	return append([]T{}, q.container...)
}

// String satisfies the [fmt.Stringer] interface and allows a Queue to be printed.
func (q *Queue[T]) String() string {
	return fmt.Sprintf("%v", q.container)
}
