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

		first, err := list.First()
		test.Err(t, err) // Should be list is empty error
		test.Equal(t, first, nil)

		last, err := list.Last()
		test.Err(t, err) // Should be list is empty error
		test.Equal(t, last, nil)

		test.Equal(t, list.Len(), 0)
	})

	t.Run("composite literal", func(t *testing.T) {
		list := &list.List[string]{}

		first, err := list.First()
		test.Err(t, err) // Should be list is empty error
		test.Equal(t, first, nil)

		last, err := list.Last()
		test.Err(t, err) // Should be list is empty error
		test.Equal(t, last, nil)

		test.Equal(t, list.Len(), 0)
	})
}

func TestAppend(t *testing.T) {
	list := list.New[string]()

	list.Append("foo")
	test.Equal(t, list.Len(), 1) // Wrong length after append

	first, err := list.First()
	test.Ok(t, err)
	test.Equal(t, first.Item(), "foo") // First item should be "foo"

	last, err := list.Last()
	test.Ok(t, err)
	test.Equal(t, last.Item(), "foo") // Last item should also be "foo"

	// Append again
	list.Append("bar")
	test.Equal(t, list.Len(), 2)

	first, err = list.First()
	test.Ok(t, err)
	test.Equal(t, first.Item(), "foo") // First should *still* be "foo"

	last, err = list.Last()
	test.Ok(t, err)
	test.Equal(t, last.Item(), "bar") // Last should now be "bar"

	// One more time
	list.Append("baz")
	test.Equal(t, list.Len(), 3)

	first, err = list.First()
	test.Ok(t, err)
	test.Equal(t, first.Item(), "foo") // First should *still* be "foo"

	last, err = list.Last()
	test.Ok(t, err)
	test.Equal(t, last.Item(), "baz") // Last should now be "baz"
}

func TestPrepend(t *testing.T) {
	list := list.New[string]()

	list.Prepend("foo")
	test.Equal(t, list.Len(), 1) // Wrong length after prepend

	first, err := list.First()
	test.Ok(t, err)
	test.Equal(t, first.Item(), "foo") // First element should be "foo"

	last, err := list.Last()
	test.Ok(t, err)
	test.Equal(t, last.Item(), "foo") // Last element should also be "foo"

	// Prepend again
	list.Prepend("bar")
	test.Equal(t, list.Len(), 2)

	first, err = list.First()
	test.Ok(t, err)
	test.Equal(t, first.Item(), "bar") // First should now be "bar"

	last, err = list.Last()
	test.Ok(t, err)
	test.Equal(t, last.Item(), "foo") // Last should still be "foo"

	// One more time
	list.Prepend("baz")
	test.Equal(t, list.Len(), 3)

	first, err = list.First()
	test.Ok(t, err)
	test.Equal(t, first.Item(), "baz") // First should now be "baz"

	last, err = list.Last()
	test.Ok(t, err)
	test.Equal(t, last.Item(), "foo") // Last should still be "foo"
}

func TestPop(t *testing.T) {
	list := list.New[int]()

	list.Append(1)
	list.Append(2)
	list.Append(3)

	first, err := list.First()
	test.Ok(t, err)
	test.Equal(t, list.Len(), 3)
	test.Equal(t, first.Item(), 1)

	last, err := list.Last()
	test.Ok(t, err)
	test.Equal(t, last.Item(), 3)

	three, err := list.Pop()
	test.Ok(t, err)
	test.Equal(t, three.Item(), 3)
	test.Equal(t, list.Len(), 2) // Len should be 2 after Pop

	last, err = list.Last()
	test.Ok(t, err)
	test.Equal(t, last.Item(), 2)

	two, err := list.Pop()
	test.Ok(t, err)
	test.Equal(t, two.Item(), 2)
	test.Equal(t, list.Len(), 1) // Len should be 1 after second Pop

	one, err := list.Pop()
	test.Ok(t, err)
	test.Equal(t, one.Item(), 1)
	test.Equal(t, list.Len(), 0) // Len should be 0 after third Pop

	first, err = list.First()
	test.Err(t, err)
	test.Equal(t, first, nil)

	last, err = list.Last()
	test.Err(t, err)
	test.Equal(t, last, nil)

	// One more Pop, should error
	broke, err := list.Pop()
	test.Err(t, err)
	test.Equal(t, broke, nil)
}

func TestPopFirst(t *testing.T) {
	list := list.New[int]()

	list.Append(1)
	list.Append(2)
	list.Append(3)

	first, err := list.First()
	test.Ok(t, err)
	test.Equal(t, list.Len(), 3)
	test.Equal(t, first.Item(), 1)

	last, err := list.Last()
	test.Ok(t, err)
	test.Equal(t, last.Item(), 3)

	one, err := list.PopFirst()
	test.Ok(t, err)
	test.Equal(t, one.Item(), 1)
	test.Equal(t, list.Len(), 2) // Len should be 2 after PopFirst

	last, err = list.Last()
	test.Ok(t, err)
	test.Equal(t, last.Item(), 3)

	two, err := list.PopFirst()
	test.Ok(t, err)
	test.Equal(t, two.Item(), 2)
	test.Equal(t, list.Len(), 1) // Len should be 1 after second PopFirst

	three, err := list.PopFirst()
	test.Ok(t, err)
	test.Equal(t, three.Item(), 3)
	test.Equal(t, list.Len(), 0) // Len should be 0 after third PopFirst

	first, err = list.First()
	test.Err(t, err)
	test.Equal(t, first, nil)

	last, err = list.Last()
	test.Err(t, err)
	test.Equal(t, last, nil)

	// One more PopFirst, should error
	broke, err := list.PopFirst()
	test.Err(t, err)
	test.Equal(t, broke, nil)
}

func TestRemove(t *testing.T) {
	list := list.New[string]()
	list.Append("one")
	two := list.Append("two")
	three := list.Append("three")
	list.Append("four")

	test.Equal(t, list.Len(), 4)

	list.Remove(two)
	test.Equal(t, list.Len(), 3)

	want := []string{"one", "three", "four"}
	test.EqualFunc(t, slices.Collect(list.All()), want, slices.Equal)

	list.Remove(three)
	test.Equal(t, list.Len(), 2)

	want = []string{"one", "four"}
	test.EqualFunc(t, slices.Collect(list.All()), want, slices.Equal)
}

func TestItems(t *testing.T) {
	list := list.New[int]()

	// Append a bunch of stuff
	for i := range 10 {
		list.Append(i)
	}

	items := slices.Collect(list.All())
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
