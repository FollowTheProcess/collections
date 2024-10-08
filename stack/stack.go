// Package stack implements a LIFO stack generic over any type.
//
// The stack is not safe for concurrent access across goroutines, the caller is responsible for
// synchronising concurrent access.
package stack

import (
	"errors"
	"fmt"
	"iter"
)

// Stack is a LIFO stack generic over any type.
//
// A Stack should be instantiated by the New function and not directly.
type Stack[T any] struct {
	container []T // Underlying slice
}

// New constructs and returns a new stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// From builds a [Stack] from an existing slice of items, pushing items
// into the stack in the order of the slice.
//
// The stack will be preallocated the size of len(items).
func From[T any](items []T) *Stack[T] {
	stack := &Stack[T]{container: make([]T, 0, len(items))}
	for _, item := range items {
		stack.Push(item)
	}

	return stack
}

// Collect builds a [Stack] from an iterator of items, pushing items
// into the stack in the order of iteration.
func Collect[T any](items iter.Seq[T]) *Stack[T] {
	stack := New[T]()
	for item := range items {
		stack.Push(item)
	}

	return stack
}

// Push adds an item to the top of stack.
//
//	s := stack.New[string]()
//	s.Push("hello")
func (s *Stack[T]) Push(item T) {
	s.container = append(s.container, item)
}

// Pop removes an item from the top of the stack, if the stack
// is empty, an error will be returned.
//
//	s := stack.New[string]()
//	s.Push("hello")
//	s.Push("there")
//	s.Push("general")
//	s.Push("kenobi")
//
//	item, _ := s.Pop()
//	fmt.Println(item) // "kenobi"
func (s *Stack[T]) Pop() (T, error) {
	l := len(s.container)
	if l == 0 {
		var none T
		return none, errors.New("pop from empty stack")
	}
	item := s.container[l-1]
	s.container = s.container[:l-1]

	return item, nil
}

// Size returns the number of items in the stack.
//
//	s := stack.New[string]()
//	s.Push("hello")
//	s.Push("there")
//	s.Size() // 2
func (s *Stack[T]) Size() int {
	return len(s.container)
}

// Empty returns whether or not the stack is empty.
//
//	s := stack.New[string]()
//	s.Empty() // true
//	s.Push("hello")
//	s.Empty() // false
func (s *Stack[T]) Empty() bool {
	return len(s.container) == 0
}

// Items returns the an iterator over the stack in LIFO order.
//
//	s := stack.New[string]()
//	s.Push("hello")
//	s.Push("there")
//	slices.Collect(s.Items()) // [there hello]
func (s *Stack[T]) Items() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := len(s.container) - 1; i >= 0; i-- {
			if !yield(s.container[i]) {
				return
			}
		}
	}
}

// String satisfies the [fmt.Stringer] interface and allows a stack to be printed.
func (s *Stack[T]) String() string {
	return fmt.Sprintf("%v", s.container)
}
