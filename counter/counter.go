// Package counter implements a convenient data structure for counting comparable values.
//
// Inspired by python's collections.Counter.
//
// The counter is not safe for concurrent access across goroutines, the caller is responsible for
// synchronising concurrent access.
package counter

import (
	"iter"

	"github.com/FollowTheProcess/collections/priority"
)

// Pair holds a countable item along with its count.
type Pair[T comparable] struct {
	Item  T   // The actual item
	Count int // The number of times it was counted
}

// Counter is a convenient construct for counting comparable values.
type Counter[T comparable] struct {
	counts map[T]int
}

// New constructs and returns a new [Counter].
func New[T comparable]() *Counter[T] {
	return &Counter[T]{counts: make(map[T]int)}
}

// From builds a [Counter] from an existing slice of items, counting
// them in the order they are given.
//
//	items := []string{"apple", "apple", "orange", "banana"}
//	counts := counter.From(items)
func From[T comparable](items []T) *Counter[T] {
	counter := &Counter[T]{
		counts: make(map[T]int, len(items)), // Preallocate the known size
	}

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

// Count returns the count of item, or 0 if it's not yet been seen.
func (c *Counter[T]) Count(item T) int {
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

// MostCommon returns the n most common items in descending order.
func (c *Counter[T]) MostCommon(n int) []Pair[T] {
	queue := priority.New[T]()
	for item, count := range c.counts {
		queue.Push(item, count)
	}

	results := make([]Pair[T], 0, n)
	// Pop off the queue in priority (count) order
	for range n {
		item, _ := queue.Pop() //nolint: errcheck // Only error is pop from empty queue which we know we won't hit
		results = append(results, Pair[T]{Item: item, Count: c.counts[item]})
	}

	return results
}
