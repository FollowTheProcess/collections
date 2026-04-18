package orderedmap_test

import (
	"maps"
	"math/rand/v2"
	"slices"
	"testing"
	"testing/quick"

	"go.followtheprocess.codes/collections/orderedmap"
	"go.followtheprocess.codes/test"
)

func TestGetInsert(t *testing.T) {
	m := orderedmap.New[string, string]()

	test.Equal(t, m.Size(), 0, test.Context("Starting size should be 0"))

	missing, ok := m.Get("missing")
	test.False(t, ok, test.Context("Missing item should return ok = false"))
	test.Equal(t, missing, "", test.Context("Value should be zero value"))

	val, existed := m.Insert("new", "item")
	test.False(t, existed, test.Context("Insert of a new item should return false"))
	test.Equal(t, val, "item", test.Context("Insert of new item should return item"))

	test.Equal(t, m.Size(), 1, test.Context("Wrong size, should contain 1 new item"))

	item, ok := m.Get("new")
	test.True(t, ok)
	test.Equal(t, item, "item", test.Context("Retrieved item should be \"item\""))

	old, existed := m.Insert("new", "other item")
	test.True(t, existed, test.Context("Item should have existed"))
	test.Equal(t, old, "item", test.Context("Old item should be item"))

	test.Equal(t, m.Size(), 1, test.Context("Wrong size, should contain 2 new items"))

	val, ok = m.Get("new")
	test.True(t, ok, test.Context("Item should have existed"))
	test.Equal(t, val, "other item", test.Context("The new value should be returned from Get"))
}

func TestInsertRemove(t *testing.T) {
	m := orderedmap.New[int, string]()

	test.Equal(t, m.Size(), 0, test.Context("Wrong initial size"))

	m.Insert(1, "one")

	test.Equal(t, m.Size(), 1, test.Context("Wrong size after inserts"))

	one, existed := m.Remove(1)
	test.True(t, existed, test.Context("1 should have existed in the map"))
	test.Equal(t, one, "one", test.Context("Wrong value returned from Remove"))

	test.Equal(t, m.Size(), 0, test.Context("Wrong size after removal"))
}

func TestRemove(t *testing.T) {
	m := orderedmap.New[int, string]()

	test.Equal(t, m.Size(), 0, test.Context("Wrong initial size"))

	missing, existed := m.Remove(42)
	test.False(t, existed, test.Context("existed should be false"))
	test.Equal(t, missing, "", test.Context("should be zero value"))

	m.Insert(1, "one")
	m.Insert(2, "two")
	m.Insert(3, "three")

	test.Equal(t, m.Size(), 3, test.Context("Wrong size after inserts"))

	two, existed := m.Remove(2)
	test.True(t, existed, test.Context("2 should have existed in the map"))
	test.Equal(t, two, "two", test.Context("Wrong value returned from Remove"))

	test.Equal(t, m.Size(), 2, test.Context("Wrong size after removal"))
}

func TestOldest(t *testing.T) {
	m := orderedmap.New[int, string]()

	oldestKey, oldestValue, ok := m.Oldest()
	test.False(t, ok)
	test.Equal(t, oldestValue, "")
	test.Equal(t, oldestKey, 0)

	// Insert a bunch of stuff
	m.Insert(1, "one")
	m.Insert(2, "two")
	m.Insert(3, "three")
	m.Insert(4, "four")

	oldestKey, oldestValue, ok = m.Oldest()
	test.True(t, ok)
	test.Equal(t, oldestKey, 1, test.Context("Wrong oldest key"))
	test.Equal(t, oldestValue, "one", test.Context("Wrong oldest value"))
}

func TestNewest(t *testing.T) {
	m := orderedmap.New[int, string]()

	newestKey, newestValue, ok := m.Newest()
	test.False(t, ok)
	test.Equal(t, newestValue, "")
	test.Equal(t, newestKey, 0)

	// Insert a bunch of stuff
	m.Insert(1, "one")
	m.Insert(2, "two")
	m.Insert(3, "three")
	m.Insert(4, "four")

	newestKey, newestValue, ok = m.Newest()
	test.True(t, ok)
	test.Equal(t, newestKey, 4, test.Context("Wrong newest key"))
	test.Equal(t, newestValue, "four", test.Context("Wrong newest value"))
}

func TestGetOrInsert(t *testing.T) {
	m := orderedmap.New[string, int]()

	one, existed := m.GetOrInsert("one", 1)
	test.False(t, existed, test.Context("should not have existed"))
	test.Equal(t, one, 1, test.Context("wrong value"))

	// Try again with same value
	one, existed = m.GetOrInsert("one", 1)
	test.True(t, existed, test.Context("should have existed this time"))
	test.Equal(t, one, 1)

	// And again with different value
	one, existed = m.GetOrInsert("one", 100)
	test.True(t, existed, test.Context("should also exist"))
	test.Equal(t, one, 1, test.Context("wrong value"))
}

func TestContains(t *testing.T) {
	m := orderedmap.New[string, int]()
	m.Insert("one", 1)
	m.Insert("two", 2)
	m.Insert("three", 3)

	test.True(t, m.Contains("one"), test.Context("Map should contain \"one\""))
	test.False(t, m.Contains("four"), test.Context("\"four\" is not in the map"))
}

func TestItems(t *testing.T) {
	// Let's use WithCapacity
	m := orderedmap.WithCapacity[string, int](4)
	m.Insert("one", 1)
	m.Insert("two", 2)
	m.Insert("three", 3)
	m.Insert("four", 4)

	items := maps.Collect(m.All())
	want := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
	}

	test.EqualFunc(t, items, want, maps.Equal)
}

func TestKeys(t *testing.T) {
	m := orderedmap.New[string, int]()
	m.Insert("one", 1)
	m.Insert("two", 2)
	m.Insert("three", 3)
	m.Insert("four", 4)

	keys := slices.Collect(m.Keys())
	want := []string{"one", "two", "three", "four"}

	test.EqualFunc(t, keys, want, slices.Equal)
}

func TestValues(t *testing.T) {
	m := orderedmap.New[string, int]()
	m.Insert("one", 1)
	m.Insert("two", 2)
	m.Insert("three", 3)
	m.Insert("four", 4)

	values := slices.Collect(m.Values())
	want := []int{1, 2, 3, 4}

	test.EqualFunc(t, values, want, slices.Equal)
}

func TestInsertGetProperty(t *testing.T) {
	m := orderedmap.New[string, int]()

	// TIL testing/quick exists!

	insert := func(key string, value int) int {
		// If we insert a value against a given key
		m.Insert(key, value)

		return value
	}

	get := func(key string, _ int) int {
		// We should always get it back with Get, regardless
		// of the key or value
		val, _ := m.Get(key)

		return val
	}

	if err := quick.CheckEqual(insert, get, nil); err != nil {
		t.Error(err)
	}
}

func FuzzInsertGet(f *testing.F) {
	// Fuzz is similar but you have to give a hint on the values first
	// and it takes longer
	corpus := [...]string{
		"",
		"a normal sentence",
		"日a本b語ç日ð本Ê語þ日¥本¼語i日©",
		"\xf8\xa1\xa1\xa1\xa1",
		"£$%^&*(((())))",
		"91836347287",
		"日ð本Ê語þ日¥本¼語i",
		"✅🛠️🧠⚡️⚠️😎🪜",
		"\n\n\r\n\t   ",
	}

	for _, item := range corpus {
		f.Add(item, rand.Int())
	}

	m := orderedmap.New[string, int]()

	f.Fuzz(func(t *testing.T, key string, value int) {
		// If we insert a value against a given key
		m.Insert(key, value)

		// We should always get the same value when asking for
		// the same key
		got, ok := m.Get(key)
		if !ok {
			t.Fatalf("key %s not found in map after insertion", key)
		}

		if got != value {
			t.Fatalf("value fetched from map (%d) differs from that inserted (%d)", got, value)
		}
	})
}

func BenchmarkInsert(b *testing.B) {
	b.Run("new", func(b *testing.B) {
		m := orderedmap.New[int, int]()

		i := 0
		for b.Loop() {
			m.Insert(i, i)
			i++
		}
	})

	b.Run("exists", func(b *testing.B) {
		m := orderedmap.New[string, int]()
		m.Insert("hello", 1)

		for b.Loop() {
			m.Insert("hello", 2) // Update the value stored against "hello"
		}
	})
}

func BenchmarkRemove(b *testing.B) {
	b.Run("exists", func(b *testing.B) {
		m := orderedmap.New[string, int]()

		for b.Loop() {
			// I've tried doing various combinations of b.StopTimer() and stuff but
			// it doesn't ever seem to work correctly so we just live with including
			// the insertion too
			m.Insert("hello", 1) // Put the item back again so it always exists on each run
			m.Remove("hello")
		}
	})

	b.Run("missing", func(b *testing.B) {
		m := orderedmap.New[string, int]()

		for b.Loop() {
			m.Remove("hello") // Doesn't exist, remove should be a no-op
		}
	})
}
