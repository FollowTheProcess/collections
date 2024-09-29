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

// element is a single element in the list.
type element[T any] struct {
	item T           // The actual item stored in the list
	prev *element[T] // The previous element in the list
	next *element[T] // The next element in the list
}

// List is a doubly-linked list.
type List[T any] struct {
	first *element[T] // The first element in the list
	last  *element[T] // The last element in the list
	len   int         // The number of elements in the list
}

// New returns a new [List].
func New[T any]() *List[T] {
	return &List[T]{}
}

// Append adds an item to the end (tail) of the list. It may be retrieved
// afterwards with l.Last().
func (l *List[T]) Append(item T) {
	if l.last != nil {
		// List has items in it, insert after last
		l.insertAfter(l.last, &element[T]{item: item})
	} else {
		// Empty list
		l.Prepend(item)
	}
}

// Prepend adds an item to the front (head) of the list. It may be retrieved
// afterwards with l.First().
func (l *List[T]) Prepend(item T) {
	elem := &element[T]{item: item}

	if l.first != nil {
		// List has items in it, insert before first
		l.insertBefore(l.first, elem)
	} else {
		// Empty list
		l.first = elem
		l.last = elem
		l.len++
	}
}

// First returns the item at the start (head) of the list, leaving
// the list unmodified.
func (l *List[T]) First() (item T, ok bool) {
	var zero T
	if l.first == nil {
		return zero, false
	}
	return l.first.item, true
}

// Last returns the item at the end (tail) of the list, leaving
// the list unmodified.
func (l *List[T]) Last() (item T, ok bool) {
	var zero T
	if l.last == nil {
		return zero, false
	}
	return l.last.item, true
}

// Len returns the number of elements in the list.
func (l *List[T]) Len() int {
	return l.len
}

// Pop removes the last element of the list and returns it.
//
// If the list is empty, Pop() returns an error.
func (l *List[T]) Pop() (T, error) {
	var zero T
	last := l.last
	if last == nil {
		return zero, errors.New("pop from empty list")
	}
	l.remove(last)

	return last.item, nil
}

// PopFirst removes the first element of the list and returns it.
//
// If the list is empty, PopFirst() returns an error.
func (l *List[T]) PopFirst() (T, error) {
	var zero T
	first := l.first
	if first == nil {
		return zero, errors.New("pop from empty list")
	}
	l.remove(first)

	return first.item, nil
}

// Items returns an iterator over the items in the list, in order.
func (l *List[T]) Items() iter.Seq[T] {
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

// TODO(@FollowTheProcess): We need an Insert API but I don't really want to export element
// needs some thinking

// insertAfter inserts a new element after an existing one.
func (l *List[T]) insertAfter(after *element[T], elem *element[T]) {
	elem.prev = after

	if after.next != nil {
		// We're inserting in the middle of the list somewhere
		elem.next = after.next
		after.next.prev = elem
	} else {
		// We're inserting at the end of the list
		elem.next = nil
		l.last = elem
	}

	after.next = elem
	l.len++
}

// insertBefore inserts a new element before an existing one.
func (l *List[T]) insertBefore(before *element[T], elem *element[T]) {
	elem.next = before

	if before.prev != nil {
		// We're inserting in the middle of the list somewhere
		elem.prev = before.prev
		before.prev.next = elem
	} else {
		// We're inserting at the start of the list
		elem.prev = nil
		l.first = elem
	}

	before.prev = elem
	l.len++
}

// remove removes an element from the list.
func (l *List[T]) remove(elem *element[T]) {
	if elem.prev != nil {
		// Removing from somewhere in the middle of the list
		elem.prev.next = elem.next
	} else {
		// Removing the first element
		l.first = elem.next
	}

	if elem.next != nil {
		// Removing from somewhere in the middle of the list
		elem.next.prev = elem.prev
	} else {
		// Removing the last element
		l.last = elem.prev
	}

	// Avoid memory leaks
	elem.next = nil
	elem.prev = nil

	l.len--
}
