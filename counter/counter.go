// Package counter implements a convenient data structure for counting comparable values.
//
// Inspired by python's collections.Counter.
//
// The counter is not safe for concurrent access across goroutines, the caller is responsible for
// synchronising concurrent access.
package counter

import (
	"cmp"
	"iter"
	"slices"
)

// Counter is a convenient construct for counting comparable values.
type Counter[T comparable] struct {
	counts map[T]int
}

// New constructs and returns a new [Counter].
func New[T comparable]() *Counter[T] {
	return &Counter[T]{counts: make(map[T]int)}
}

// WithCapacity constructs and returns a new [Counter] with the given capacity.
//
// This can be a useful performance improvement when the number of unique items to count
// is known ahead of time as it eliminates the need for reallocation.
func WithCapacity[T comparable](capacity int) *Counter[T] {
	return &Counter[T]{counts: make(map[T]int, capacity)}
}

// From builds a [Counter] from an existing slice of items, counting
// them in the order they are given.
//
//	items := []string{"apple", "apple", "orange", "banana"}
//	counts := counter.From(items)
func From[T comparable](items []T) *Counter[T] {
	counter := WithCapacity[T](len(items))

	for _, item := range items {
		counter.Add(item)
	}

	return counter
}

// Collect builds a [Counter] from an iterator of items, counting
// them in the order they are iterated through.
//
//	items := []string{"apple", "apple", "orange", "banana"}
//	counts := counter.Collect(slices.Values(items))
func Collect[T comparable](items iter.Seq[T]) *Counter[T] {
	counter := New[T]()

	for item := range items {
		counter.Add(item)
	}

	return counter
}

// Size returns the current number of items in the [Counter].
func (c *Counter[T]) Size() int {
	return len(c.counts)
}

// Add adds an item to the counter, incrementing it's count and returning the new count.
//
// If the item doesn't exist, it is added to the counter with the count of 1, and 1 will be returned.
func (c *Counter[T]) Add(item T) int {
	v, exists := c.counts[item]
	if !exists {
		// Not previously seen -> item: 1
		c.counts[item] = 1
		return 1
	}

	// Already existed, increment it's count
	v += 1
	c.counts[item] = v
	return v
}

// Sub subtracts an item from the counter, decrementing it's count and returning the new count.
//
// If the decrement would set the item's count to 0, it is then removed
// entirely and 0 is returned.
//
// If the item doesn't exist, this is a no-op returning 0.
func (c *Counter[T]) Sub(item T) int {
	if v, exists := c.counts[item]; exists {
		v -= 1
		if v == 0 {
			// If it's now 0, remove it entirely
			delete(c.counts, item)
			return 0
		}

		// Otherwise, store the new value back
		c.counts[item] = v
		return v
	}

	return 0
}

// Remove completely removes an item from the counter, returning it's count if
// it was present, or 0 if not.
//
// If the item didn't exist, Remove is a no-op.
func (c *Counter[T]) Remove(item T) int {
	v, exists := c.counts[item]
	if !exists {
		return 0
	}

	delete(c.counts, item)
	return v
}

// Get returns the count of item, or 0 if it's not yet been seen.
func (c *Counter[T]) Get(item T) int {
	v, exists := c.counts[item]
	if !exists {
		return 0
	}

	return v
}

// Sum returns the sum of all the item counts in the [Counter], effectively
// the overall number of items including duplicates.
func (c *Counter[T]) Sum() int {
	sum := 0
	for _, count := range c.counts {
		sum += count
	}

	return sum
}

// Reset resets the [Counter], removing all items and freeing the memory.
func (c *Counter[T]) Reset() {
	clear(c.counts)
}

// MostCommon returns the item with the highest count, along with the count itself.
//
// If the Counter is empty it returns the zero value for the item type and 0 for the count.
func (c *Counter[T]) MostCommon() (item T, count int) {
	if len(c.counts) == 0 {
		var zero T
		return zero, 0
	}

	var mostCommon T
	highestCount := 0
	for item, count := range c.counts {
		if count > highestCount {
			mostCommon = item
			highestCount = count
		}
	}

	return mostCommon, highestCount
}

// Descending returns an iterator of the item, count pairs in the Counter, yielding them
// in descending order (i.e. highest count first).
func (c *Counter[T]) Descending() iter.Seq2[T, int] {
	type count struct {
		item  T
		value int
	}
	counts := make([]count, 0, len(c.counts))
	for item, value := range c.counts {
		counts = append(counts, count{item: item, value: value})
	}

	// Sort by value in descending order
	slices.SortStableFunc(counts, func(a, b count) int {
		return cmp.Compare(b.value, a.value)
	})

	return func(yield func(T, int) bool) {
		for _, count := range counts {
			if !yield(count.item, count.value) {
				return
			}
		}
	}
}

// Counts returns an iterator over the item, count pairs in the Counter, yielding them
// in a non-deterministic order.
func (c *Counter[T]) Counts() iter.Seq2[T, int] {
	return func(yield func(T, int) bool) {
		for item, count := range c.counts {
			if !yield(item, count) {
				return
			}
		}
	}
}

// Items returns an iterator over the items in the Counter, yielding them
// in a non-deterministic order.
func (c *Counter[T]) Items() iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range c.counts {
			if !yield(item) {
				return
			}
		}
	}
}
