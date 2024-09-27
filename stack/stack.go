// Package stack implements a LIFO stack generic over any type.
//
// The stack is not safe for concurrent access across goroutines, the caller is responsible for
// synchronising concurrent access.
package stack

import "errors"

// Stack is a LIFO stack.
type Stack[T any] struct {
	top *element[T] // The stack top pointer
}

// element is an element stored in the stack.
type element[T any] struct {
	item T           // The actual item
	next *element[T] // A pointer to the next element, or nil if there is none
}

// New returns a new [Stack].
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push pushes a new item onto the top of the stack.
func (s *Stack[T]) Push(item T) {
	elem := &element[T]{
		item: item,
		next: s.top, // Store the old top as the next one down
	}

	// Make this one the new top
	s.top = elem
}

// Pop pops an item off the top of the stack.
//
// If the stack is empty, an error will be returned.
func (s *Stack[T]) Pop() (T, error) {
	if s.top == nil {
		return *new(T), errors.New("pop from empty stack")
	}

	// Get the thing off the stack top
	toReturn := s.top

	// Move the top down one
	s.top = toReturn.next
	toReturn.next = nil

	return toReturn.item, nil
}

// Empty reports whether the stack is empty.
func (s *Stack[T]) Empty() bool {
	return s.top == nil
}
