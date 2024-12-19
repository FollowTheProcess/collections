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
type Stack[T any] struct {
	container []T
}

// New constructs and returns a new stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{container: make([]T, 0)}
}

// WithCapacity constructs and returns a new stack with the given capacity.
//
// This can be a useful performance improvement when the expected maximum size of the stack is
// known ahead of time as it eliminates the need for reallocation.
func WithCapacity[T any](capacity int) *Stack[T] {
	return &Stack[T]{container: make([]T, 0, capacity)}
}

// From builds a [Stack] from an existing slice of items, pushing items
// into the stack in the order of the slice.
//
// The stack will be preallocated the size of len(items).
func From[T any](items []T) *Stack[T] {
	stack := WithCapacity[T](len(items))
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

// Capacity returns the capacity of the stack, i.e. the number of items
// it can contain without the need for reallocation.
//
// Use [WithCapacity] to create a stack of a given capacity.
//
//	s := stack.WithCapacity[string](10)
//	s.Capacity() // 10
func (s *Stack[T]) Capacity() int {
	return cap(s.container)
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

// Items returns an iterator over the stack in LIFO order.
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

// String satisfies the [fmt.Stringer] interface and allows a stack to print itself.
func (s *Stack[T]) String() string {
	return fmt.Sprintf("%v", s.container)
}
