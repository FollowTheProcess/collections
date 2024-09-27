// Package set implements a simple, generic set data structure.
package set

import (
	"fmt"
)

// Set is a simple, generic implementation of a mathematical set.
type Set[T comparable] struct {
	container map[T]struct{}
}

// New builds and returns a new empty [Set].
//
// The set will grow as needed as items are inserted, an initial small
// size is allocated.
//
// If constructing a set from a pre-existing slice of items, use [From]
// which will preallocate the set with the appropriate size.
func New[T comparable]() *Set[T] {
	return &Set[T]{
		container: make(map[T]struct{}),
	}
}

// From builds a [Set] from an existing slice of items.
//
// The set will be preallocated the size of len(items).
func From[T comparable](items []T) *Set[T] {
	set := &Set[T]{
		container: make(map[T]struct{}, len(items)),
	}
	for _, item := range items {
		set.Insert(item)
	}
	return set
}

// Insert inserts an item into the [Set].
//
// Returns whether the item was newly inserted. Inserting an item that
// is already present is effectively a no-op.
//
//	s := set.New[string]()
//	s.Insert("foo") // true -> set was modified by the insertion
//	s.Insert("foo") // false -> "foo" is already in the set, it was not modified
func (s *Set[T]) Insert(item T) bool {
	// Indexing into a nil map doesn't panic, which is why we can do this
	// first safely
	if _, exists := s.container[item]; exists {
		return false
	}

	// nil safety
	if s.container == nil {
		s.container = make(map[T]struct{})
	}

	s.container[item] = struct{}{}
	return true
}

// Contains reports whether the set contains item.
//
//	s := set.New[int]()
//	s.Contains(1) // false
//	s.Insert(1)
//	s.Contains(1) // true
func (s *Set[T]) Contains(item T) bool {
	_, exists := s.container[item]
	return exists
}

// Remove removes an item from the set.
//
// Returns whether the value was present. Removing an item
// that wasn't in the set is effectively a no-op.
func (s *Set[T]) Remove(item T) bool {
	if _, exists := s.container[item]; !exists {
		return false
	}
	delete(s.container, item)
	return true
}

// Size returns the current size of the set.
func (s *Set[T]) Size() int {
	return len(s.container)
}

// Items returns the set's items as a slice.
//
// The order of the items is non-deterministic, the caller should
// sort the returned slice if order is important.
func (s *Set[T]) Items() []T {
	items := make([]T, 0, s.Size())
	for item := range s.container {
		items = append(items, item)
	}
	return items
}

// String implements [fmt.Stringer] for a [Set] and allows
// it to print itself.
func (s *Set[T]) String() string {
	return fmt.Sprintf("%v", s.Items())
}
