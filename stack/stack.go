// Package stack implements a LIFO stack generic over any type.
package stack

import (
	"errors"
	"fmt"
)

// Stack is a LIFO stack generic over any type.
type Stack[T any] struct {
	container []T // Underlying slice
}

// New constructs and returns a new stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{container: make([]T, 0)}
}

// Push adds an item to the top of stack.
//
//	s := stack.New[string]()
//	s.Push("hello")
func (s *Stack[T]) Push(item T) {
	s.container = append(s.container, item)
}

// Pop removes an item from the top of the stack, if the stack
// is empty an error will be returned.
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

// Length returns the number of items in the stack.
//
//	s := stack.New[string]()
//	s.Push("hello")
//	s.Push("there")
//	s.Length() // 2
func (s *Stack[T]) Length() int {
	return len(s.container)
}

// IsEmpty returns whether or not the stack is empty.
//
//	s := stack.New[string]()
//	s.IsEmpty() // true
//	s.Push("hello")
//	s.IsEmpty() // false
func (s *Stack[T]) IsEmpty() bool {
	return len(s.container) == 0
}

// Items returns the items in the stack as a slice.
//
//	s := stack.New[string]()
//	s.Push("hello")
//	s.Push("there")
//	s.Items() // [hello there]
func (s *Stack[T]) Items() []T {
	return s.container
}

// String satisfies the stringer interface and allows a stack to be printed.
func (s Stack[T]) String() string {
	return fmt.Sprintf("%v", s.container)
}
