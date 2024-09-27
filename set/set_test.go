package set_test

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/collections/set"
	"github.com/FollowTheProcess/test"
)

func TestInsert(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		s := set.New[string]()

		test.True(t, s.Insert("foo"))    // Inserting foo for the first time should return true
		test.False(t, s.Insert("foo"))   // Second insert of foo should return false
		test.True(t, s.Contains("foo"))  // Set didn't contain "foo"
		test.False(t, s.Contains("bar")) // Set said it contained "bar" but shouldn't have

		// testing nil safety
		danger := &set.Set[string]{}
		test.True(t, danger.Insert("bar"))
		test.False(t, danger.Insert("bar"))
		test.True(t, danger.Insert("baz"))
	})
	t.Run("ints", func(t *testing.T) {
		s := set.New[int]()

		test.True(t, s.Insert(1))    // Inserting 1 for the first time should return true
		test.False(t, s.Insert(1))   // Second insert of 1 should return false
		test.True(t, s.Contains(1))  // Set didn't contain 1
		test.False(t, s.Contains(2)) // Set said it contained 2 but shouldn't have

		// testing nil safety
		danger := &set.Set[int]{}
		test.True(t, danger.Insert(42))
		test.False(t, danger.Insert(42))
		test.True(t, danger.Insert(69))
	})
}

func TestFrom(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		// testing From
		items := []string{"foo", "bar", "baz"}
		other := set.From(items)

		test.True(t, other.Contains("foo"))      // Set from items didn't contain "foo"
		test.True(t, other.Contains("bar"))      // Set from items didn't contain "bar"
		test.True(t, other.Contains("baz"))      // Set from items didn't contain "baz"
		test.False(t, other.Contains("missing")) // Missing item "missing" reported as present
	})

	t.Run("floats", func(t *testing.T) {
		// testing From
		items := []float64{3.14159, 42.58, 69.73}
		other := set.From(items)

		test.True(t, other.Contains(3.14159)) // Set from items didn't contain pi
		test.True(t, other.Contains(42.58))   // Set from items didn't contain 42.58
		test.True(t, other.Contains(69.73))   // Set from items didn't contain 69.73
		test.False(t, other.Contains(100.1))  // Missing item 100.1 reported as present
	})
}

func TestRemove(t *testing.T) {
	t.Run("structs", func(t *testing.T) {
		type person struct {
			name string
			age  int
		}

		set := set.New[person]()

		missing := person{name: "Missing", age: 42}

		// Remove on an empty set shouldn't panic or do anything bad
		set.Remove(missing)

		tom := person{name: "Tom", age: 30}
		gandalf := person{name: "Gandalf", age: 55000}
		wendy := person{name: "Wendy", age: 12}

		set.Insert(tom)
		set.Insert(gandalf)
		set.Insert(wendy)

		test.Equal(t, set.Size(), 3) // Incorrect size

		// Sorry Wendy
		removed := set.Remove(wendy)
		test.True(t, removed) // Removed should be true

		test.Equal(t, set.Size(), 2) // Incorrect size after killing wendy

		test.False(t, set.Contains(wendy))
	})
}

func TestItems(t *testing.T) {
	items := []string{"cheese", "apples", "oranges", "milk"}
	slices.Sort(items)
	set := set.From(items)

	got := set.Items()
	slices.Sort(got)
	test.EqualFunc(t, got, items, slices.Equal)
}
