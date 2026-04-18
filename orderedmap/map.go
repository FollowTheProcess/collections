// Package orderedmap implements an ordered map, that is; a map that remembers the order in which
// key, value pairs were inserted.
//
// The map is not safe for concurrent access across goroutines, the caller is responsible for
// synchronising concurrent access.
package orderedmap // import "go.followtheprocess.codes/collections/orderedmap"

import "iter"

// none is the sentinel used for "no neighbour" and "no free slot". Real
// slot indices are always non-negative.
const none = -1

// entry is one slot in the arena. When the slot is live, prev and next
// are indices of the neighbouring live entries in insertion order (or
// none at the ends). When the slot is on the freelist, next points to
// the next free slot and prev is unused.
type entry[K comparable, V any] struct {
	key   K
	value V
	prev  int
	next  int
}

// Map is an ordered map.
type Map[K comparable, V any] struct {
	inner   map[K]int     // key -> index into entries
	entries []entry[K, V] // arena of entries
	head    int           // index of oldest live entry, or none
	tail    int           // index of newest live entry, or none
	free    int           // head of the freelist, or none
	size    int           // count of live entries
}

// New creates and returns a new ordered map.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		inner: make(map[K]int),
		head:  none,
		tail:  none,
		free:  none,
	}
}

// WithCapacity creates and returns a new ordered [Map] with the given capacity.
//
// This can be a useful performance improvement when the expected maximum size of the map
// is known ahead of time as it eliminates the need for reallocation.
func WithCapacity[K comparable, V any](capacity int) *Map[K, V] {
	return &Map[K, V]{
		inner:   make(map[K]int, capacity),
		entries: make([]entry[K, V], 0, capacity),
		head:    none,
		tail:    none,
		free:    none,
	}
}

// allocSlot returns the index of a slot to use for a new entry. It pops
// from the freelist if possible, otherwise grows the arena.
func (m *Map[K, V]) allocSlot() int {
	if m.free != none {
		idx := m.free
		m.free = m.entries[idx].next

		return idx
	}

	m.entries = append(m.entries, entry[K, V]{})

	return len(m.entries) - 1
}

// linkAtTail appends the slot at idx to the tail of the insertion-order
// list. The slot's key and value must already be set by the caller.
func (m *Map[K, V]) linkAtTail(idx int) {
	if m.tail == none {
		m.entries[idx].prev = none
		m.entries[idx].next = none
		m.head = idx
		m.tail = idx

		return
	}

	m.entries[idx].prev = m.tail
	m.entries[idx].next = none
	m.entries[m.tail].next = idx
	m.tail = idx
}

// unlink removes the slot at idx from the insertion-order list, clears
// its key/value so the backing array cannot retain references, and
// pushes it onto the freelist.
func (m *Map[K, V]) unlink(idx int) {
	prev := m.entries[idx].prev
	next := m.entries[idx].next

	if prev != none {
		m.entries[prev].next = next
	} else {
		m.head = next
	}

	if next != none {
		m.entries[next].prev = prev
	} else {
		m.tail = prev
	}

	var zeroK K

	var zeroV V

	m.entries[idx].key = zeroK
	m.entries[idx].value = zeroV
	m.entries[idx].prev = none
	m.entries[idx].next = m.free
	m.free = idx
}

// Get returns the value stored against the given key in the map and a boolean
// to indicate presence, like the standard Go map.
//
// If the requested key wasn't in the map, the zero value for the item and false are returned.
// If the key was present, the item and true are returned.
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	idx, exists := m.inner[key]
	if !exists {
		var zero V

		return zero, false
	}

	return m.entries[idx].value, true
}

// Contains reports whether the map contains the given key.
func (m *Map[K, V]) Contains(key K) bool {
	_, exists := m.inner[key]

	return exists
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
	if idx, exists := m.inner[key]; exists {
		oldValue := m.entries[idx].value
		m.entries[idx].value = value

		return oldValue, true
	}

	idx := m.allocSlot()
	m.entries[idx].key = key
	m.entries[idx].value = value
	m.linkAtTail(idx)
	m.inner[key] = idx
	m.size++

	return value, false
}

// Remove removes a key from the map, returning the stored value and
// a boolean to indicate whether it was in the map to begin with.
//
// If the value was in the map, the removed value and true are returned, if not
// the zero value for the value type and false are returned.
func (m *Map[K, V]) Remove(key K) (value V, existed bool) {
	idx, exists := m.inner[key]
	if !exists {
		var zero V

		return zero, false
	}

	val := m.entries[idx].value
	delete(m.inner, key)
	m.unlink(idx)
	m.size--

	return val, true
}

// Size returns the number of items currently stored in the map. This operation
// is O(1).
func (m *Map[K, V]) Size() int {
	return m.size
}

// Oldest returns the oldest key, value pair in the map, i.e. the pair
// that was inserted first. Note that in place modifications do not update the order.
func (m *Map[K, V]) Oldest() (key K, value V, ok bool) {
	if m.head == none {
		return key, value, ok
	}

	e := m.entries[m.head]

	return e.key, e.value, true
}

// Newest returns the newest key, value pair in the map, i.e. the pair that
// was inserted last. Note that in place modifications do not update the order.
func (m *Map[K, V]) Newest() (key K, value V, ok bool) {
	if m.tail == none {
		return key, value, ok
	}

	e := m.entries[m.tail]

	return e.key, e.value, true
}

// GetOrInsert fetches a value by it's key if it is present in the map, and if not
// inserts the passed in value against that key instead.
//
// The returned boolean reports whether the key already existed.
func (m *Map[K, V]) GetOrInsert(key K, value V) (val V, existed bool) {
	if idx, exists := m.inner[key]; exists {
		return m.entries[idx].value, true
	}

	idx := m.allocSlot()
	m.entries[idx].key = key
	m.entries[idx].value = value
	m.linkAtTail(idx)
	m.inner[key] = idx
	m.size++

	return value, false
}

// All returns an iterator over the entries in the map
// in the order in which they were inserted.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := m.head; i != none; i = m.entries[i].next {
			e := m.entries[i]
			if !yield(e.key, e.value) {
				return
			}
		}
	}
}

// Keys returns an iterator over the keys in the map
// in the order in which they were inserted.
func (m *Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for i := m.head; i != none; i = m.entries[i].next {
			if !yield(m.entries[i].key) {
				return
			}
		}
	}
}

// Values returns an iterator over the values in the map
// in the order in which they were inserted.
func (m *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := m.head; i != none; i = m.entries[i].next {
			if !yield(m.entries[i].value) {
				return
			}
		}
	}
}
