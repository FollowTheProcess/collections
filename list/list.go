// Package list implements a doubly linked list.
//
// A doubly-linked list holds a pair of pointers to the first and last element
// of the list (often referred to as the head and the tail; respectively).
//
// The list can be traversed both forwards and backwards and offers cheap (O(1)) insertion
// and removal from arbitrary indexes.
package list

import (
	"errors"
	"iter"
)

// Node is a single Node in the list.
type Node[T any] struct {
	item T        // The actual item stored in the list
	prev *Node[T] // The previous element in the list
	next *Node[T] // The next element in the list
}

// NewNode returns a new list [Node].
func NewNode[T any](item T) *Node[T] {
	return &Node[T]{item: item}
}

// Item returns the item stored in the [Node].
func (n *Node[T]) Item() T {
	return n.item
}

// List is a doubly-linked list.
type List[T any] struct {
	first *Node[T] // The first element in the list
	last  *Node[T] // The last element in the list
	len   int      // The number of elements in the list
}

// New returns a new [List].
func New[T any]() *List[T] {
	return &List[T]{}
}

// Append adds an item to the end (tail) of the list, returning the list [Node] it was inserted into.
// It may be retrieved afterwards with l.Last().
func (l *List[T]) Append(item T) *Node[T] {
	node := NewNode(item)
	if l.last != nil {
		// List has items in it, insert after last
		l.insertAfter(l.last, node)
	} else {
		// Empty list
		l.Prepend(item)
	}

	return node
}

// Prepend adds an item to the front (head) of the list, returning the list [Node] it was inserted into.
// It may be retrieved afterwards with l.First().
func (l *List[T]) Prepend(item T) *Node[T] {
	node := NewNode(item)

	if l.first != nil {
		// List has items in it, insert before first
		l.insertBefore(l.first, node)
	} else {
		// Empty list
		l.first = node
		l.last = node
		l.len++
	}

	return node
}

// First returns a pointer to the node at the start (head) of the list, leaving
// the list unmodified.
//
// If the list is empty an error is returned.
func (l *List[T]) First() (*Node[T], error) {
	if l.first == nil {
		return nil, errors.New("First() called on empty list")
	}
	return l.first, nil
}

// Last returns a pointer to the node at the end (tail) of the list, leaving
// the list unmodified.
//
// If the list is empty an error is returned.
func (l *List[T]) Last() (*Node[T], error) {
	if l.last == nil {
		return nil, errors.New("Last() called on empty list")
	}
	return l.last, nil
}

// Len returns the number of elements in the list.
func (l *List[T]) Len() int {
	return l.len
}

// Pop removes the last node from the list and returns it.
//
// If the list is empty, Pop() returns an error.
func (l *List[T]) Pop() (*Node[T], error) {
	if l.last == nil {
		return nil, errors.New("Pop() called on empty list")
	}

	return l.remove(l.last), nil
}

// PopFirst removes the first node from the list and returns it.
//
// If the list is empty, PopFirst() returns an error.
func (l *List[T]) PopFirst() (*Node[T], error) {
	if l.first == nil {
		return nil, errors.New("PopFirst() called on empty list")
	}

	return l.remove(l.first), nil
}

// Remove removes a [Node] from the list, returning it after removal.
func (l *List[T]) Remove(node *Node[T]) *Node[T] {
	return l.remove(node)
}

// All returns an iterator over the items in the list, in order.
func (l *List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for elem := l.first; elem != nil; elem = elem.next {
			if !yield(elem.item) {
				return
			}
		}
	}
}

// Backwards returns an iterator over the items in the list, in reverse order.
func (l *List[T]) Backwards() iter.Seq[T] {
	return func(yield func(T) bool) {
		for elem := l.last; elem != nil; elem = elem.prev {
			if !yield(elem.item) {
				return
			}
		}
	}
}

// insertAfter inserts a new node after an existing one.
func (l *List[T]) insertAfter(after, node *Node[T]) {
	node.prev = after

	if after.next != nil {
		// We're inserting in the middle of the list somewhere
		node.next = after.next
		after.next.prev = node
	} else {
		// We're inserting at the end of the list
		node.next = nil
		l.last = node
	}

	after.next = node
	l.len++
}

// insertBefore inserts a new element before an existing one.
func (l *List[T]) insertBefore(before, node *Node[T]) {
	node.next = before

	if before.prev != nil {
		// We're inserting in the middle of the list somewhere
		node.prev = before.prev
		before.prev.next = node
	} else {
		// We're inserting at the start of the list
		node.prev = nil
		l.first = node
	}

	before.prev = node
	l.len++
}

// remove removes a node from the list, returning the one it just removed.
func (l *List[T]) remove(node *Node[T]) *Node[T] {
	if node.prev != nil {
		// Removing from somewhere in the middle of the list
		node.prev.next = node.next
	} else {
		// Removing the first element
		l.first = node.next
	}

	if node.next != nil {
		// Removing from somewhere in the middle of the list
		node.next.prev = node.prev
	} else {
		// Removing the last element
		l.last = node.prev
	}

	// Avoid memory leaks
	node.next = nil
	node.prev = nil

	l.len--

	return node
}
