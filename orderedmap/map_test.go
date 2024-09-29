package orderedmap_test

import (
	"testing"

	"github.com/FollowTheProcess/collections/orderedmap"
	"github.com/FollowTheProcess/test"
)

func TestGetInsert(t *testing.T) {
	m := orderedmap.New[string, string]()

	test.Equal(t, m.Size(), 0) // Starting size should be 0

	missing, ok := m.Get("missing")
	test.False(t, ok)          // Missing item should return ok = false
	test.Equal(t, missing, "") // Value should be zero value

	val, existed := m.Insert("new", "item")
	test.False(t, existed)     // Insert of a new item should return false
	test.Equal(t, val, "item") // Insert of new item should return item

	test.Equal(t, m.Size(), 1) // Wrong size, should contain 1 new item

	item, ok := m.Get("new")
	test.True(t, ok)            // new should exist in the map
	test.Equal(t, item, "item") // Retrieved item should be "item"

	old, existed := m.Insert("new", "other item")
	test.True(t, existed)      // Item should have existed
	test.Equal(t, old, "item") // Old item should be item

	test.Equal(t, m.Size(), 1) // Wrong size, should contain 2 new items

	val, ok = m.Get("new")
	test.True(t, ok)                 // Item should have existed
	test.Equal(t, val, "other item") // The new value should be returned from Get
}

func TestRemove(t *testing.T) {
	m := orderedmap.New[int, string]()

	test.Equal(t, m.Size(), 0) // Wrong initial size

	missing, existed := m.Remove(42)
	test.False(t, existed)     // existed should be false
	test.Equal(t, missing, "") // should be zero value

	m.Insert(1, "one")
	m.Insert(2, "two")
	m.Insert(3, "three")

	test.Equal(t, m.Size(), 3) // Wrong size after inserts

	two, existed := m.Remove(2)
	test.True(t, existed)     // 2 should have existed in the map
	test.Equal(t, two, "two") // Wrong value returned from Remove

	test.Equal(t, m.Size(), 2) // Wrong size after removal
}
