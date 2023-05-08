// Package stack implements a LIFO stack generic over any type.
//
// The stack is implemented using an internal slice pointer, and is not thread safe.
package stack

import (
	"errors"
	"fmt"
)

// ErrPopFromEmptyStack is returned when Pop() is called on an empty stack.
var ErrPopFromEmptyStack = errors.New("pop from empty stack")

// Stack is a LIFO stack generic over any type.
//
// A Stack should be instantiated by the New function and not directly,
// doing so will result in a nil pointer dereference.
type Stack[T any] struct {
	container *[]T // Underlying slice, reference to allow mutation.
}

// New constructs and returns a new stack.
func New[T any]() Stack[T] {
	container := make([]T, 0)
	return Stack[T]{container: &container}
}

// Push adds an item to the top of stack.
//
//	s := stack.New[string]()
//	s.Push("hello")
func (s Stack[T]) Push(item T) {
	*s.container = append(*s.container, item)
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
func (s Stack[T]) Pop() (T, error) {
	l := len(*s.container)
	if l == 0 {
		var none T
		return none, ErrPopFromEmptyStack
	}
	item := (*s.container)[l-1]
	*s.container = (*s.container)[:l-1]

	return item, nil
}

// Length returns the number of items in the stack.
//
//	s := stack.New[string]()
//	s.Push("hello")
//	s.Push("there")
//	s.Length() // 2
func (s Stack[T]) Length() int {
	return len(*s.container)
}

// IsEmpty returns whether or not the stack is empty.
//
//	s := stack.New[string]()
//	s.IsEmpty() // true
//	s.Push("hello")
//	s.IsEmpty() // false
func (s Stack[T]) IsEmpty() bool {
	return len(*s.container) == 0
}

// Items returns the items in the stack as a new slice (copy).
//
//	s := stack.New[string]()
//	s.Push("hello")
//	s.Push("there")
//	s.Items() // [hello there]
func (s Stack[T]) Items() []T {
	return append([]T{}, *s.container...)
}

// String satisfies the stringer interface and allows a stack to be printed.
func (s Stack[T]) String() string {
	return fmt.Sprintf("%v", *s.container)
}
