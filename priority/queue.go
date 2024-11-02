// Package priority implements a generic priority queue; that is a queue who's items are assigned a priority
// and are popped off the queue in (descending) order of this priority.
package priority

import "errors"

// element holds an element in the priority queue along with it's priority.
type element[T any] struct {
	item     T   // The stored item
	priority int // The priority of the element
}

// Queue is a generic priority queue.
type Queue[T any] struct {
	container []element[T] // Underlying slice
}

// New builds and returns a new, empty priority Queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Push adds an item and it's priority to the queue.
func (q *Queue[T]) Push(item T, priority int) {
	q.container = append(q.container, element[T]{item: item, priority: priority})
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

	return elem.item, nil
}

// Size returns the number of elements currently in the queue.
func (q *Queue[T]) Size() int {
	return len(q.container)
}

// Empty returns whether the queue is empty.
func (q *Queue[T]) Empty() bool {
	return len(q.container) == 0
}

// siftUp moves an item (by index) up the heap until it's in the correct position
// for it's priority.
func (q *Queue[T]) siftUp(index int) {
	for {
		parent := (index - 1) / 2 //nolint: mnd
		if parent == index || q.container[index].priority < q.container[parent].priority {
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
	return q.container[i].priority > q.container[j].priority
}
