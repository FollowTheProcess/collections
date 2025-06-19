// Package orderedmap implements an ordered map, that is; a map that remembers the order in which
// key, value pairs were inserted.
package orderedmap // import "go.followtheprocess.codes/collections/orderedmap"

import (
	"iter"

	"go.followtheprocess.codes/collections/list"
)

// entry is a single key, value pair entry in the map.
type entry[K comparable, V any] struct {
	key   K                        // The key, used to look the entry up in the inner map
	value V                        // The value
	node  *list.Node[*entry[K, V]] // The node in the linked list storing this entry
}

// Map is an ordered map.
type Map[K comparable, V any] struct {
	inner map[K]*entry[K, V]       // The backing hashmap
	list  *list.List[*entry[K, V]] // The linked list keeping track of insertion order
}

// New creates and returns a new ordered map.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		inner: make(map[K]*entry[K, V]),
		list:  list.New[*entry[K, V]](),
	}
}

// WithCapacity creates and returns a new ordered [Map] with the given capacity.
//
// This can be a useful performance improvement when the expected maximum size of the map
// is known ahead of time as it eliminates the need for reallocation.
func WithCapacity[K comparable, V any](capacity int) *Map[K, V] {
	return &Map[K, V]{
		inner: make(map[K]*entry[K, V], capacity),
		list:  list.New[*entry[K, V]](),
	}
}

// Get returns the value stored against the given key in the map and a boolean
// to indicate presence, like the standard Go map.
//
// If the requested key wasn't in the map, the zero value for the item and false are returned.
// If the key was present, the item and true are returned.
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	var zero V

	val, exists := m.inner[key]
	if !exists {
		return zero, false
	}

	return val.value, true
}

// Contains reports whether the map contains the given key.
func (m *Map[K, V]) Contains(key K) bool {
	if _, exists := m.inner[key]; exists {
		return true
	}

	return false
}

// Insert inserts a new value into the map against the given key, returning the previous
// value and a boolean to indicate presence.
//
// If the map did not have this key present before the call to Insert, it will return the
// value just inserted and false.
//
// If the map did have this key, and this call to Insert is therefore an update of an existing value,
// then the old value and true are returned.
func (m *Map[K, V]) Insert(key K, value V) (val V, existed bool) {
	if old, exists := m.inner[key]; exists {
		// The item exists, this is therefore an update
		oldValue := old.value // Take a copy so we can return it
		old.value = value     // Set the new value back

		return oldValue, true
	}

	// The item didn't exist, this is a brand new insertion
	e := &entry[K, V]{
		key:   key,
		value: value,
	}

	e.node = m.list.Append(e)
	m.inner[key] = e

	return value, false
}

// Remove removes a key from the map, returning the stored value and
// a boolean to indicate whether it was in the map to begin with.
//
// If the value was in the map, the removed value and true are returned, if not
// the zero value for the value type and false are returned.
func (m *Map[K, V]) Remove(key K) (value V, existed bool) {
	if entry, existed := m.inner[key]; existed {
		m.list.Remove(entry.node) // Drop it from our list
		delete(m.inner, key)      // And the map

		return entry.value, true
	}

	// Didn't exist, just return
	var zero V

	return zero, false
}

// Size returns the number of items currently stored in the map. This operation
// is O(1).
func (m *Map[K, V]) Size() int {
	return m.list.Len()
}

// Oldest returns the oldest key, value pair in the map, i.e. the pair
// that was inserted first. Note that in place modifications do not update the order.
func (m *Map[K, V]) Oldest() (key K, value V, ok bool) {
	var zeroKey K

	var zeroVal V

	node, err := m.list.First()
	if err != nil {
		// Empty list
		return zeroKey, zeroVal, false
	}

	return node.Item().key, node.Item().value, true
}

// Newest returns the newest key, value pair in the map, i.e. the pair that
// was inserted last. Note that in place modifications do not update the order.
func (m *Map[K, V]) Newest() (key K, value V, ok bool) {
	var zeroKey K

	var zeroVal V

	node, err := m.list.Last()
	if err != nil {
		// Empty list
		return zeroKey, zeroVal, false
	}

	return node.Item().key, node.Item().value, true
}

// GetOrInsert fetches a value by it's key if it is present in the map, and if not
// inserts the passed in value against that key instead.
//
// The returned boolean reports whether the key already existed.
func (m *Map[K, V]) GetOrInsert(key K, value V) (val V, existed bool) {
	if entry, exists := m.inner[key]; exists {
		// Already in the map, return the value
		return entry.value, true
	}

	// The item didn't exist, this is a brand new insertion
	e := &entry[K, V]{
		key:   key,
		value: value,
	}

	e.node = m.list.Append(e)
	m.inner[key] = e

	return value, false
}

// All returns an iterator over the entries in the map
// in the order in which they were inserted.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for item := range m.list.All() {
			if !yield(item.key, item.value) {
				return
			}
		}
	}
}

// Keys returns an iterator over the keys in the map
// in the order in which they were inserted.
func (m *Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for item := range m.list.All() {
			if !yield(item.key) {
				return
			}
		}
	}
}

// Values returns an iterator over the values in the map
// in the order in which they were inserted.
func (m *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for item := range m.list.All() {
			if !yield(item.value) {
				return
			}
		}
	}
}
