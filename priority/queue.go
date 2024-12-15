// Package priority implements a generic priority queue; that is a queue who's items are assigned a priority
// and are popped off the queue in (descending) order of this priority.
package priority

import (
	"errors"
)

// Element holds an element in the priority queue along with it's priority.
type Element[T any] struct {
	Item     T   // The stored item itself
	Priority int // The priority of the element, highest gets popped first
}

// Queue is a generic priority queue.
type Queue[T any] struct {
	container []Element[T] // Underlying slice
}

// New builds and returns a new, empty priority Queue.
//
// If you already have a list of items you wish to transform into a priority queue,
// consider using [From] or [FromFunc] as they are more performant than constructing
// and empty queue and filling it in a loop.
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// From builds and returns a priority Queue from an already established []Element.
//
// This is more performant than creating a new empty Queue and using Push, but requires
// the caller to construct the slice of Element themselves.
func From[T any](elements []Element[T]) *Queue[T] {
	queue := &Queue[T]{
		container: make([]Element[T], 0, len(elements)), // Create a new slice so we own the data internally
	}

	queue.container = append(queue.container, elements...)

	// Heapify the container
	queue.init()

	return queue
}

// FromFunc builds and returns a priority Queue from an already established list of items, taking
// a closure that is called to calculate the priority of each item.
//
// This is more performant than creating a new empty Queue and using Push, but requires
// the caller to construct the slice of items themselves.
//
// It is slightly less performant than calling [From] as the priority is calculated
// on the fly.
//
// Take care implementing the priorityFunc as it is called for every element in items, it should
// be as performant as possible.
func FromFunc[T any](items []T, priorityFunc func(item T) int) *Queue[T] {
	container := make([]Element[T], 0, len(items))
	for _, item := range items {
		container = append(container, Element[T]{Item: item, Priority: priorityFunc(item)})
	}

	queue := &Queue[T]{
		container: container,
	}
	// Heapify the container
	queue.init()

	return queue
}

// Push adds an item and it's priority to the queue.
//
// If you already have an existing slice of items, consider converting it to []Element
// and using [From] as it is more performant.
func (q *Queue[T]) Push(item T, priority int) {
	q.container = append(q.container, Element[T]{Item: item, Priority: priority})
	q.siftUp(len(q.container) - 1)
}

// Pop removes and returns the element with the highest priority.
func (q *Queue[T]) Pop() (T, error) {
	if len(q.container) == 0 {
		var zero T
		return zero, errors.New("pop from empty priority queue")
	}

	// Swap the first (highest priority) and last element
	n := len(q.container) - 1
	q.swap(0, n)

	// Return the last element (now the highest priority) and trim the queue
	elem := q.container[n]
	q.container = (q.container)[:n]

	// Update heap order
	q.siftDown(0, n)

	return elem.Item, nil
}

// Size returns the number of elements currently in the queue.
func (q *Queue[T]) Size() int {
	return len(q.container)
}

// Empty returns whether the queue is empty.
func (q *Queue[T]) Empty() bool {
	return len(q.container) == 0
}

// init heapifies the underlying container, establishing the heap invariants required by
// the other methods. It is only used when creating a priority queue with [From].
func (q *Queue[T]) init() {
	n := len(q.container)
	for i := n/2 - 1; i >= 0; i-- { //nolint: mnd
		q.siftDown(i, n)
	}
}

// siftUp moves an item (by index) up the heap until it's in the correct position
// for it's priority.
func (q *Queue[T]) siftUp(index int) {
	for {
		parent := (index - 1) / 2 //nolint: mnd
		if parent == index || q.container[index].Priority < q.container[parent].Priority {
			break
		}
		q.swap(parent, index)
		index = parent
	}
}

// siftDown moves an item (by index) down the heap until it's in the correct position.
func (q *Queue[T]) siftDown(index, length int) bool {
	i := index
	for {
		leftChild := 2*i + 1 //nolint: mnd
		if leftChild >= length || leftChild < 0 {
			break
		}
		toSwap := leftChild
		if rightChild := leftChild + 1; rightChild < length && q.less(rightChild, leftChild) {
			toSwap = rightChild
		}
		if !q.less(toSwap, i) {
			break
		}
		q.swap(i, toSwap)
		i = toSwap
	}
	return i > index
}

// swap swaps two elements in the heap by index.
func (q *Queue[T]) swap(i, j int) {
	q.container[i], q.container[j] = q.container[j], q.container[i]
}

// less reports whether element i should come before element j in priority order.
func (q *Queue[T]) less(i, j int) bool {
	return q.container[i].Priority > q.container[j].Priority
}
