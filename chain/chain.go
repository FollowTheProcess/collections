// Package chain provides an implementation of chained maps in a single view, where lookups traverse
// the layers of maps until a match is found.
package chain

import (
	"iter"
	"slices"
)

// Chain is a single view over chained maps.
type Chain[K comparable, V any] struct {
	maps []map[K]V
}

// New constructs a new [Chain].
func New[K comparable, V any]() *Chain[K, V] {
	return &Chain[K, V]{
		maps: []map[K]V{},
	}
}

// From constructs a new [Chain] from an existing slice of maps.
//
// The order of priority is in the order of the slice so a slice of maps
// [a, b, c] will result in a [Chain] who's lookup order is a, b, c.
func From[K comparable, V any](maps []map[K]V) *Chain[K, V] {
	return &Chain[K, V]{
		maps: maps,
	}
}

// Collect constructs a new [Chain] from an iterator of maps.
//
// The order of priority is in the order of the slice so a iterator of maps
// yielding [a, b, c] will result in a [Chain] who's lookup order is a, b, c.
func Collect[K comparable, V any](maps iter.Seq[map[K]V]) *Chain[K, V] {
	return &Chain[K, V]{
		maps: slices.Collect(maps),
	}
}

// Append adds a map to the end of the [Chain] (lowest lookup priority).
func (c *Chain[K, V]) Append(m map[K]V) {
	c.maps = append(c.maps, m)
}

// Prepend adds a map to the start of the [Chain] (highest lookup priority).
func (c *Chain[K, V]) Prepend(m map[K]V) {
	c.maps = append(c.maps, m)
	copy(c.maps[1:], c.maps)
	c.maps[0] = m
}

// Size returns the number of maps in the chain.
func (c *Chain[K, V]) Size() int {
	return len(c.maps)
}

// Get returns the value stored against the given key in the chain of maps and
// a boolean to indicate presence, like the standard Go map.
//
// The key is looked up in each map in the chain in order, and the value returned
// is from the first one with the key present.
//
// If the requested key wasn't in any of the maps in the chain the zero value for the
// value type and false are returned.
func (c Chain[K, V]) Get(key K) (value V, ok bool) {
	for _, m := range c.maps {
		val, exists := m[key]
		if exists {
			// Return the first one to have it
			return val, true
		}
	}

	// Wasn't in any of the maps
	var zero V

	return zero, false
}

// Insert inserts a new value into the [Chain] against the given key, returning the previous
// value and a boolean to indicate presence.
//
// If the key did not exist in any of the maps before the call to Insert, the value will
// be inserted into the first map in the chain, Insert will return the value just inserted and false.
//
// If any map in the chain did have this key, it will be updated in place in that same map and
// Insert will return the previous value and true.
func (c *Chain[K, V]) Insert(key K, value V) (val V, existed bool) {
	for _, m := range c.maps {
		if old, exists := m[key]; exists {
			// The item exists in one of the maps, this is therefore an update
			m[key] = value

			return old, true
		}
	}

	// The item didn't exist, so insert it into the first map
	// If we haven't got a list of maps yet, create one
	if len(c.maps) == 0 {
		c.maps = []map[K]V{make(map[K]V)}
	}

	c.maps[0][key] = value

	return value, false
}

// Remove removes a key from the [Chain], returning the stored value and
// a boolean to indicate whether it was in the map to begin with.
//
// If the value was in the chain, the removed value and true are returned, if not
// the zero value for the value type and false are returned.
//
// The value removed will be the first one encountered.
func (c *Chain[K, V]) Remove(key K) (value V, existed bool) {
	for _, m := range c.maps {
		if val, exists := m[key]; exists {
			delete(m, key)

			return val, true
		}
	}

	// Didn't exist
	var zero V

	return zero, false
}
