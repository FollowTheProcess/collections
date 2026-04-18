package queue_test

import (
	"slices"
	"testing"

	"go.followtheprocess.codes/collections/queue"
	"go.followtheprocess.codes/test"
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

// TestIsEmptyAfterDrain exercises IsEmpty after pushing and popping all
// items back out.
func TestIsEmptyAfterDrain(t *testing.T) {
	q := queue.New[int]()
	for i := range 10 {
		q.Push(i)
	}

	test.False(t, q.IsEmpty())

	for range 10 {
		_, ok := q.Pop()
		test.True(t, ok)
	}

	test.True(t, q.IsEmpty(), test.Context("Queue should be empty after full drain"))
}

func TestSize(t *testing.T) {
	q := queue.New[string]()
	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	test.Equal(t, q.Size(), 4)
}

func TestPop(t *testing.T) {
	q := queue.New[string]()
	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	item, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, item, "hello")

	item, ok = q.Pop()
	test.True(t, ok)
	test.Equal(t, item, "there")

	item, ok = q.Pop()
	test.True(t, ok)
	test.Equal(t, item, "general")

	item, ok = q.Pop()
	test.True(t, ok)
	test.Equal(t, item, "kenobi")

	// Try one more pop, should report empty
	_, ok = q.Pop()
	test.False(t, ok)
}

// TestMixedPushPop exercises interleaved Push/Pop patterns to make sure
// the ring buffer indices track correctly across wraparound.
func TestMixedPushPop(t *testing.T) {
	q := queue.New[int]()

	// Prime with 3 items, pop 2, push 3 more — drives the head/tail
	// indices through a wraparound.
	q.Push(1)
	q.Push(2)
	q.Push(3)

	first, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, first, 1)

	second, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, second, 2)

	q.Push(4)
	q.Push(5)
	q.Push(6)

	got := slices.Collect(q.All())
	test.EqualFunc(t, got, []int{3, 4, 5, 6}, slices.Equal)
}

// TestLongCycle pushes and pops many more items than the queue ever
// holds at once to catch any growth-on-pop regression.
func TestLongCycle(t *testing.T) {
	q := queue.New[int]()

	for i := range 10_000 {
		q.Push(i)

		got, ok := q.Pop()
		test.True(t, ok)
		test.Equal(t, got, i)
	}

	test.True(t, q.IsEmpty(), test.Context("Queue should be empty after equal Push/Pop count"))
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

	first, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, first, "cheese")

	second, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, second, "apples")
}

func TestCollect(t *testing.T) {
	items := []string{"cheese", "apples", "wine", "beer"}

	q := queue.Collect(slices.Values(items))

	test.Equal(t, q.Size(), 4)

	first, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, first, "cheese")

	second, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, second, "apples")
}

func TestWrapAroundFIFO(t *testing.T) {
	q := queue.WithCapacity[int](4)
	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Push(4)
	q.Pop()
	q.Pop()
	q.Push(5)
	q.Push(6)

	got := slices.Collect(q.All())
	want := []int{3, 4, 5, 6}
	test.EqualFunc(t, got, want, slices.Equal)
}

func TestNotNew(t *testing.T) {
	q := queue.Queue[int]{}
	q.Push(1)
	q.Push(2)

	first, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, first, 1)
}

func BenchmarkQueue(b *testing.B) {
	s := queue.New[int]()

	for b.Loop() {
		s.Push(1)

		_, ok := s.Pop()
		if !ok {
			b.Error("Pop() returned ok=false")
		}
	}
}
