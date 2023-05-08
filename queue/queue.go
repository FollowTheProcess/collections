// Package queue implements a FIFO queue generic over any type.
//
// The queue is implemented using an internal slice pointer, and is not thread safe.
package queue

import (
	"errors"
	"fmt"
)

// ErrPopFromEmptyQueue is returned when Pop() is called on an empty queue.
var ErrPopFromEmptyQueue = errors.New("pop from empty queue")

// Queue is a FIFO queue generic over any type.
//
// A Queue should be instantiated by the New function and not directly,
// doing so will result in a nil pointer dereference.
type Queue[T any] struct {
	container *[]T // Underlying slice
}

// New constructs and returns a new stack.
func New[T any]() Queue[T] {
	container := make([]T, 0) // Initialise the underlying slice
	return Queue[T]{container: &container}
}

// Push adds an item to the back of the queue.
//
//	q := queue.New[string]()
//	q.Push("hello")
func (q Queue[T]) Push(item T) {
	*q.container = append(*q.container, item)
}

// Pop removes an item from the front of the queue, if the queue
// is empty, ErrPopFromEmptyQueue will be returned.
//
//	q := queue.New[string]()
//	q.Push("hello")
//	q.Push("there")
//	item, _ := q.Pop()
//	fmt.Println(item) // "hello"
func (q Queue[T]) Pop() (T, error) {
	l := len(*q.container)
	if l == 0 {
		var none T
		return none, ErrPopFromEmptyQueue
	}
	item := (*q.container)[0]
	*q.container = (*q.container)[1:]

	return item, nil
}

// Length returns the number of items in the queue.
//
//	s := queue.New[string]()
//	s.Length() // 0
//	s.Push("hello")
//	s.Push("there")
//	s.Length() // 2
func (q Queue[T]) Length() int {
	return len(*q.container)
}

// IsEmpty returns whether or not the queue is empty.
//
//	s := queue.New[string]()
//	s.IsEmpty() // true
//	s.Push("hello")
//	s.IsEmpty() // false
func (q Queue[T]) IsEmpty() bool {
	return len(*q.container) == 0
}

// Items returns the items in the queue as a new slice (copy).
//
//	q := queue.New[string]()
//	q.Push("hello")
//	q.Push("there")
//	q.Items() // [hello there]
func (q Queue[T]) Items() []T {
	return append([]T{}, *q.container...)
}

// String satisfies the stringer interface and allows a stack to be printed.
func (q Queue[T]) String() string {
	return fmt.Sprintf("%v", *q.container)
}
