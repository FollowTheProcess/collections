// Package orderedmap implements an ordered map, that is; a map that remembers the order in which
// key, value pairs were inserted.
package orderedmap

import "github.com/FollowTheProcess/collections/list"

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

// Size returns the number of items currently stored in the map. This operation
// is O(1).
func (m *Map[K, V]) Size() int {
	return m.list.Len()
}
