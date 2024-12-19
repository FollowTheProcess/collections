// Package set implements a simple, generic set data structure.
//
// The set is not safe for concurrent access across goroutines, the caller is responsible for
// synchronising concurrent access.
package set

import (
	"fmt"
	"iter"
	"maps"
	"math"
	"slices"
)

// TODO(@FollowTheProcess): IsDisjoint, IsSubset, IsSuperSet, SymmetricDifference

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
// which will preallocate the set with the appropriate size. Or to collect
// an iterator into a [Set], use [Collect].
func New[T comparable]() *Set[T] {
	return &Set[T]{
		container: make(map[T]struct{}),
	}
}

// WithCapacity builds and returns a new [Set] with the given capacity.
//
// This can be a useful performance improvement when the expected maximum size of the set
// is known ahead of time as it eliminates the need for reallocation.
func WithCapacity[T comparable](capacity int) *Set[T] {
	return &Set[T]{
		container: make(map[T]struct{}, capacity),
	}
}

// From builds a [Set] from an existing slice of items.
//
// The set will be preallocated the size of len(items).
func From[T comparable](items []T) *Set[T] {
	set := WithCapacity[T](len(items))
	for _, item := range items {
		// Note: intentionally not using Insert here as we don't need
		// the checks it provides
		set.container[item] = struct{}{}
	}
	return set
}

// Collect builds a [Set] from an iterator of items.
func Collect[T comparable](items iter.Seq[T]) *Set[T] {
	set := New[T]()
	for item := range items {
		// Note: intentionally not using Insert here as we don't need
		// the checks it provides
		set.container[item] = struct{}{}
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

// Size returns the number of items currently in the set.
//
//	s := set.New[int]()
//	s.Insert(1)
//	s.Insert(2)
//	s.Size() // 2
func (s *Set[T]) Size() int {
	return len(s.container)
}

// Items returns the an iterator over the sets items.
//
// The order of the items is non-deterministic, the caller should collect
// and sort the returned items if order is important.
func (s *Set[T]) Items() iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range s.container {
			if !yield(item) {
				return
			}
		}
	}
}

// IsEmpty reports whether the set is empty.
//
//	s := set.New[int]()
//	s.IsEmpty() // true
//	s.Insert(42)
//	s.IsEmpty() // false
func (s *Set[T]) IsEmpty() bool {
	return len(s.container) == 0
}

// String implements [fmt.Stringer] for a [Set] and allows
// it to print itself.
func (s *Set[T]) String() string {
	return fmt.Sprintf("%v", slices.Collect(maps.Keys(s.container)))
}

// Equal returns whether two sets are equal to one another, i.e. they are exactly
// the same size and contain exactly the same elements.
func Equal[T comparable](a, b *Set[T]) bool {
	if a == nil || b == nil {
		return false
	}

	if len(a.container) != len(b.container) {
		return false
	}

	for item := range a.container {
		if _, ok := b.container[item]; !ok {
			return false
		}
	}

	return true
}

// Union returns a set that is the combination of all the input sets, i.e. all
// the items from all the sets in one new set, without duplicates.
func Union[T comparable](sets ...*Set[T]) *Set[T] {
	union := New[T]()
	for _, set := range sets {
		for item := range set.container {
			// Don't need the additional checks of Insert
			union.container[item] = struct{}{}
		}
	}

	return union
}

// Intersection returns a set containing all the items present in all the input sets, without duplicates.
func Intersection[T comparable](sets ...*Set[T]) *Set[T] {
	if sets == nil {
		return New[T]()
	}

	if len(sets) == 1 {
		// Just return the input set as there is nothing to intersect with
		return sets[0]
	}

	intersection := New[T]()

	n := len(sets) // Number of sets we've been passed

	var smallest *Set[T]
	minSize := math.MaxInt
	for _, set := range sets {
		// If any set is empty, we can immediately return the empty set because
		// no matter what the other sets contain, anything intersection empty set
		// should return the empty set
		if set.IsEmpty() {
			return New[T]()
		}

		if len(set.container) < minSize {
			smallest = set
			minSize = len(set.container)
		}
	}

	for item := range smallest.container {
		numContains := 0 // The number of sets that contain this item
		for _, other := range sets {
			if other.Contains(item) {
				numContains++
			}
		}

		// If the number of other sets that contain the item is equal
		// to the total number of sets we have been passed, then it's
		// present in all of them and should be included in the intersection
		if numContains == n {
			// Don't need the additional checks of Insert
			intersection.container[item] = struct{}{}
		}
	}

	return intersection
}

// Difference returns a set containing the items present in set that are not contained in any of the others.
func Difference[T comparable](set *Set[T], others ...*Set[T]) *Set[T] {
	if set == nil || others == nil {
		return New[T]()
	}

	if set.IsEmpty() {
		return New[T]()
	}

	difference := New[T]()

	n := len(others)

	for item := range set.container {
		numNotContains := 0 // The number of sets that do not contain this item
		for _, other := range others {
			if !other.Contains(item) {
				numNotContains++
			}
		}

		// If the number of other sets that don't contain the item is equal
		// to the total number of other sets we've been passed, then the item is
		// unique to set and should be added to the difference
		if numNotContains == n {
			difference.container[item] = struct{}{}
		}
	}

	return difference
}

// IsDisjoint returns whether the sets have no items in common with one another.
//
// It is equivalent to checking for the empty intersection but is significantly faster
// than calling [Intersection] because it does not construct the result set and does no allocation.
func IsDisjoint[T comparable](sets ...*Set[T]) bool {
	// Easy to handle early and guards against indexing later
	if len(sets) == 0 {
		return false
	}

	if len(sets) == 1 {
		return false // It has every item in common with itself
	}

	var smallestIndex int // The index in sets where the smallest set is located
	minSize := math.MaxInt
	for index, set := range sets {
		if len(set.container) < minSize {
			minSize = len(set.container)
			smallestIndex = index
		}
	}

	smallest := sets[smallestIndex]

	// Remove the smallest one so we're left with a list of "others"
	// otherwise disjoint will always be false as the smallest set
	// will, by definition, include the value and will be included in sets
	sets = slices.Delete(sets, smallestIndex, smallestIndex+1)

	for item := range smallest.container {
		for _, other := range sets {
			if other.Contains(item) {
				return false
			}
		}
	}

	return true
}
