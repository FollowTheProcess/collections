package queue_test

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/collections/queue"
	"github.com/FollowTheProcess/test"
)

func TestIsEmpty(t *testing.T) {
	q := queue.New[string]()
	test.True(t, q.IsEmpty())

	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	test.False(t, q.IsEmpty())
}

func TestSize(t *testing.T) {
	q := queue.New[string]()
	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	test.Equal(t, q.Size(), 4)
}

func TestCapacity(t *testing.T) {
	q := queue.WithCapacity[int](10)
	test.Equal(t, q.Capacity(), 10)
}

func TestPop(t *testing.T) {
	q := queue.New[string]()
	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	item, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, item, "hello")

	item, err = q.Pop()
	test.Ok(t, err)
	test.Equal(t, item, "there")

	item, err = q.Pop()
	test.Ok(t, err)
	test.Equal(t, item, "general")

	item, err = q.Pop()
	test.Ok(t, err)
	test.Equal(t, item, "kenobi")

	// Try one more pop, should error
	_, err = q.Pop()
	test.Err(t, err)
}

func TestItems(t *testing.T) {
	q := queue.New[string]()
	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	want := []string{"hello", "there", "general", "kenobi"}
	got := slices.Collect(q.All())

	test.EqualFunc(t, got, want, slices.Equal)
}

func TestString(t *testing.T) {
	q := queue.New[string]()
	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	want := "[hello there general kenobi]"

	test.Equal(t, q.String(), want)
}

func TestFrom(t *testing.T) {
	items := []string{"cheese", "apples", "wine", "beer"}

	q := queue.From(items)

	test.Equal(t, q.Size(), 4)

	first, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, first, "cheese")

	second, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, second, "apples")
}

func TestCollect(t *testing.T) {
	items := []string{"cheese", "apples", "wine", "beer"}

	q := queue.Collect(slices.Values(items))

	test.Equal(t, q.Size(), 4)

	first, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, first, "cheese")

	second, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, second, "apples")
}

func TestNotNew(t *testing.T) {
	q := queue.Queue[int]{}
	q.Push(1)
	q.Push(2)

	first, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, first, 1)
}

func BenchmarkQueue(b *testing.B) {
	s := queue.New[int]()

	for range b.N {
		s.Push(1)

		_, err := s.Pop()
		if err != nil {
			b.Errorf("Pop() returned an error: %v", err)
		}
	}
}
