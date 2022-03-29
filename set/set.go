// Package set implements a hash set, generic over any comparable type and associated functionality
// such as intersection, union and difference.
//
// The hash set implemented here is an unordered collection and as such the ordering of the items
// when printing or converting to a slice is non-deterministic, the user should therefore sort the
// results if a deterministic order is required.
package set

import (
	"fmt"
)

// Set is a hash set generic over any comparable type.
type Set[T comparable] struct {
	container map[T]struct{} // Underlying map with empty struct for minimal memory use
}

// New constructs and returns a new set for a specific comparable type.
func New[T comparable]() *Set[T] {
	return &Set[T]{container: make(map[T]struct{})}
}

// Add adds an item to the set, if the item is
// already present, Add becomes a no-op.
//
//  s := set.New[string]()
//  s.Add("hello")
func (s *Set[T]) Add(item T) {
	if _, ok := s.container[item]; !ok {
		s.container[item] = struct{}{}
	}
}

// Remove removes an item from the set, if the item is
// not present, Remove becomes a no-op.
//
//  s := set.New[string]()
//  s.Add("hello")
//  s.Add("there")
//  s.Remove("there")
//  s.Items() // [hello]
func (s *Set[T]) Remove(item T) {
	delete(s.container, item)
}

// Contains returns whether or not the set contains the given item.
//
//  s := set.New[int]()
//  s.Contains(27) // false
//  s.Add(27)
//  s.Contains(27) // true
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.container[item]
	return ok
}

// Items returns all the items in the set as a slice.
//
//  s := set.New[string]()
//  s.Add("hello")
//  s.Add("there")
//  s.Items() // [hello there]
func (s *Set[T]) Items() []T {
	items := make([]T, 0, len(s.container))
	for k := range s.container {
		items = append(items, k)
	}
	return items
}

// Length returns the number of elements in the set.
//
//  s := set.New[int]()
//  s.Add(42)
//  s.Add(27)
//  s.Length() // 2
func (s *Set[T]) Length() int {
	return len(s.container)
}

// IsEmpty returns whether or not the set is empty.
//
//  s := set.New[string]()
//  s.IsEmpty() // true
//  s.Add("a thing")
//  s.IsEmpty() // false
func (s *Set[T]) IsEmpty() bool {
	return len(s.container) == 0
}

// String satisfies the stringer interface and allows a set to be printed.
func (s Set[T]) String() string {
	return fmt.Sprintf("%v", s.Items())
}

// Union returns a set that is the combination of a and b.
func Union[S *Set[T], T comparable](a, b *Set[T]) *Set[T] {
	result := New[T]()
	for item := range a.container {
		result.Add(item)
	}

	for item := range b.container {
		if !result.Contains(item) {
			result.Add(item)
		}
	}

	return result
}

// Intersection returns a set containing all the items present in both a and b.
func Intersection[S *Set[T], T comparable](a, b *Set[T]) *Set[T] {
	result := New[T]()
	for item := range a.container {
		if b.Contains(item) {
			result.Add(item)
		}
	}

	return result
}

// Difference returns a set containing the items present in a but not b.
func Difference[S *Set[T], T comparable](a, b *Set[T]) *Set[T] {
	result := New[T]()
	for item := range a.container {
		if !b.Contains(item) {
			result.Add(item)
		}
	}

	return result
}
