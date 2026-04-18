// Package queue implements a FIFO queue generic over any type.
//
// The queue is not safe for concurrent access across goroutines, the caller is responsible for
// synchronising concurrent access.
package queue // import "go.followtheprocess.codes/collections/queue"

import (
	"errors"
	"fmt"
	"iter"
)

// Queue is a FIFO queue generic over any type.
type Queue[T any] struct {
	container []T // Ring buffer backing array (len == capacity)
	head      int // Index of the next item to pop
	tail      int // Index of the next slot to push into
	size      int // Number of items currently in the queue
}

// New constructs and returns a new Queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// WithCapacity constructs and returns a new Queue with the given capacity.
//
// This can be a useful performance improvement when the expected maximum size of the queue is
// known ahead of time as it eliminates the need for reallocation.
func WithCapacity[T any](capacity int) *Queue[T] {
	return &Queue[T]{container: make([]T, capacity)}
}

// From builds a [Queue] from an existing slice of items, pushing items
// into the queue in the order of the slice.
//
// The queue will be preallocated the size of len(items).
func From[T any](items []T) *Queue[T] {
	queue := WithCapacity[T](len(items))
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
	if len(q.container) == 0 {
		q.container = make([]T, 1)
	}

	if q.size == len(q.container) {
		q.grow()
	}

	q.container[q.tail] = item
	q.tail = (q.tail + 1) % len(q.container)
	q.size++
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
	if q.size == 0 {
		var none T

		return none, errors.New("pop from empty queue")
	}

	item := q.container[q.head]
	var zero T
	q.container[q.head] = zero
	q.head = (q.head + 1) % len(q.container)
	q.size--

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
	return q.size
}

// IsEmpty returns whether or not the queue is empty.
//
//	s := queue.New[string]()
//	s.IsEmpty() // true
//	s.Push("hello")
//	s.IsEmpty() // false
func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}

// All returns the an iterator over the queue in FIFO order.
//
//	q := queue.New[string]()
//	q.Push("hello")
//	q.Push("there")
//	qlices.Collect(s.All()) // [hello there]
func (q *Queue[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range q.size {
			if !yield(q.container[(q.head+i)%len(q.container)]) {
				return
			}
		}
	}
}

// String satisfies the [fmt.Stringer] interface and allows a Queue to be printed.
func (q *Queue[T]) String() string {
	items := make([]T, q.size)
	for i := range q.size {
		items[i] = q.container[(q.head+i)%len(q.container)]
	}

	return fmt.Sprintf("%v", items)
}

// grow doubles the capacity of the ring buffer, copying items into logical order.
func (q *Queue[T]) grow() {
	newContainer := make([]T, len(q.container)*2)
	n := copy(newContainer, q.container[q.head:])
	copy(newContainer[n:], q.container[:q.head])
	q.head = 0
	q.tail = q.size
	q.container = newContainer
}
