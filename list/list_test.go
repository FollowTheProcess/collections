package list_test

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/collections/list"
	"github.com/FollowTheProcess/test"
)

func TestEmptyList(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		list := list.New[string]()

		first, ok := list.First()
		test.False(t, ok)
		test.Equal(t, first, "")

		last, ok := list.Last()
		test.False(t, ok)
		test.Equal(t, last, "")

		test.Equal(t, list.Len(), 0)
	})

	t.Run("composite literal", func(t *testing.T) {
		list := &list.List[string]{}

		first, ok := list.First()
		test.False(t, ok)
		test.Equal(t, first, "")

		last, ok := list.Last()
		test.False(t, ok)
		test.Equal(t, last, "")

		test.Equal(t, list.Len(), 0)
	})
}

func TestAppend(t *testing.T) {
	list := list.New[string]()

	list.Append("foo")
	test.Equal(t, list.Len(), 1) // Wrong length after append

	first, ok := list.First()
	test.True(t, ok)
	test.Equal(t, first, "foo") // First element should be "foo"

	last, ok := list.Last()
	test.True(t, ok)
	test.Equal(t, last, "foo") // Last element should also be "foo"

	// Append again
	list.Append("bar")
	test.Equal(t, list.Len(), 2)

	first, ok = list.First()
	test.True(t, ok)
	test.Equal(t, first, "foo") // First should *still* be "foo"

	last, ok = list.Last()
	test.True(t, ok)
	test.Equal(t, last, "bar") // Last should now be "bar"

	// One more time
	list.Append("baz")
	test.Equal(t, list.Len(), 3)

	first, ok = list.First()
	test.True(t, ok)
	test.Equal(t, first, "foo") // First should *still* be "foo"

	last, ok = list.Last()
	test.True(t, ok)
	test.Equal(t, last, "baz") // Last should now be "baz"
}

func TestPrepend(t *testing.T) {
	list := list.New[string]()

	list.Prepend("foo")
	test.Equal(t, list.Len(), 1) // Wrong length after prepend

	first, ok := list.First()
	test.True(t, ok)
	test.Equal(t, first, "foo") // First element should be "foo"

	last, ok := list.Last()
	test.True(t, ok)
	test.Equal(t, last, "foo") // Last element should also be "foo"

	// Prepend again
	list.Prepend("bar")
	test.Equal(t, list.Len(), 2)

	first, ok = list.First()
	test.True(t, ok)
	test.Equal(t, first, "bar") // First should now be "bar"

	last, ok = list.Last()
	test.True(t, ok)
	test.Equal(t, last, "foo") // Last should still be "foo"

	// One more time
	list.Prepend("baz")
	test.Equal(t, list.Len(), 3)

	first, ok = list.First()
	test.True(t, ok)
	test.Equal(t, first, "baz") // First should now be "baz"

	last, ok = list.Last()
	test.True(t, ok)
	test.Equal(t, last, "foo") // Last should still be "foo"
}

func TestPop(t *testing.T) {
	list := list.New[int]()

	list.Append(1)
	list.Append(2)
	list.Append(3)

	first, ok := list.First()
	test.True(t, ok)
	test.Equal(t, list.Len(), 3)
	test.Equal(t, first, 1)

	last, ok := list.Last()
	test.True(t, ok)
	test.Equal(t, last, 3)

	three, err := list.Pop()
	test.Ok(t, err)
	test.Equal(t, three, 3)
	test.Equal(t, list.Len(), 2) // Len should be 2 after Pop

	last, ok = list.Last()
	test.True(t, ok)
	test.Equal(t, last, 2)

	two, err := list.Pop()
	test.Ok(t, err)
	test.Equal(t, two, 2)
	test.Equal(t, list.Len(), 1) // Len should be 1 after second Pop

	one, err := list.Pop()
	test.Ok(t, err)
	test.Equal(t, one, 1)
	test.Equal(t, list.Len(), 0) // Len should be 0 after third Pop

	first, ok = list.First()
	test.False(t, ok)
	test.Equal(t, first, 0)

	last, ok = list.Last()
	test.False(t, ok)
	test.Equal(t, last, 0)

	// One more Pop, should error
	broke, err := list.Pop()
	test.Err(t, err)
	test.Equal(t, broke, 0)
}

func TestPopFirst(t *testing.T) {
	list := list.New[int]()

	list.Append(1)
	list.Append(2)
	list.Append(3)

	first, ok := list.First()
	test.True(t, ok)
	test.Equal(t, list.Len(), 3)
	test.Equal(t, first, 1)

	last, ok := list.Last()
	test.True(t, ok)
	test.Equal(t, last, 3)

	one, err := list.PopFirst()
	test.Ok(t, err)
	test.Equal(t, one, 1)
	test.Equal(t, list.Len(), 2) // Len should be 2 after PopFirst

	last, ok = list.Last()
	test.True(t, ok)
	test.Equal(t, last, 3)

	two, err := list.PopFirst()
	test.Ok(t, err)
	test.Equal(t, two, 2)
	test.Equal(t, list.Len(), 1) // Len should be 1 after second PopFirst

	three, err := list.PopFirst()
	test.Ok(t, err)
	test.Equal(t, three, 3)
	test.Equal(t, list.Len(), 0) // Len should be 0 after third PopFirst

	first, ok = list.First()
	test.False(t, ok)
	test.Equal(t, first, 0)

	last, ok = list.Last()
	test.False(t, ok)
	test.Equal(t, last, 0)

	// One more PopFirst, should error
	broke, err := list.PopFirst()
	test.Err(t, err)
	test.Equal(t, broke, 0)
}

func TestItems(t *testing.T) {
	list := list.New[int]()

	// Append a bunch of stuff
	for i := range 10 {
		list.Append(i)
	}

	items := slices.Collect(list.Items())
	want := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	test.EqualFunc(t, items, want, slices.Equal)
}

func TestReverse(t *testing.T) {
	list := list.New[int]()

	// Append a bunch of stuff
	for i := range 10 {
		list.Append(i)
	}

	items := slices.Collect(list.Backwards())
	want := []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}

	test.EqualFunc(t, items, want, slices.Equal)
}
