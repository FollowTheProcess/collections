package stack_test

import (
	"slices"
	"testing"

	"go.followtheprocess.codes/collections/stack"
	"go.followtheprocess.codes/test"
)

func TestIsEmpty(t *testing.T) {
	s := stack.New[string]()

	if !s.IsEmpty() {
		t.Error("IsEmpty should return true")
	}

	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	test.False(t, s.IsEmpty(), test.Context("Empty should not be false"))
}

func TestPop(t *testing.T) {
	s := stack.New[string]()
	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	item, ok := s.Pop()
	test.True(t, ok)
	test.Equal(t, item, "kenobi")

	item, ok = s.Pop()
	test.True(t, ok)
	test.Equal(t, item, "general")

	item, ok = s.Pop()
	test.True(t, ok)
	test.Equal(t, item, "there")

	item, ok = s.Pop()
	test.True(t, ok)
	test.Equal(t, item, "hello")

	// Try one more pop, should report empty
	item, ok = s.Pop()
	test.False(t, ok, test.Context("Pop from empty stack should return ok=false"))
	test.Equal(t, item, "", test.Context("Item should be the zero value"))
}

// TestPopPushReuse covers the slot-zeroing behaviour in Pop: draining the
// stack and pushing again should put a fresh value into the reused slot,
// so the stack must not silently surface any retained prior value.
func TestPopPushReuse(t *testing.T) {
	s := stack.New[string]()
	s.Push("a")
	s.Push("b")
	s.Push("c")

	// Drain
	for range 3 {
		_, ok := s.Pop()
		test.True(t, ok)
	}

	test.True(t, s.IsEmpty())

	// Push new values into the reused backing array
	s.Push("x")
	s.Push("y")

	got, ok := s.Pop()
	test.True(t, ok)
	test.Equal(t, got, "y")

	got, ok = s.Pop()
	test.True(t, ok)
	test.Equal(t, got, "x")
}

// TestAllAfterPartialPops covers iteration order after a mid-drain pop.
func TestAllAfterPartialPops(t *testing.T) {
	s := stack.New[int]()
	for i := 1; i <= 5; i++ {
		s.Push(i)
	}

	// Pop two, leaving 1, 2, 3 in the stack
	for range 2 {
		_, ok := s.Pop()
		test.True(t, ok)
	}

	got := slices.Collect(s.All())
	test.EqualFunc(t, got, []int{3, 2, 1}, slices.Equal)
}

func TestNotNew(t *testing.T) {
	s := stack.Stack[int]{}
	s.Push(1)
	s.Push(2)

	first, ok := s.Pop()
	test.True(t, ok)
	test.Equal(t, first, 2)
}

func TestSize(t *testing.T) {
	s := stack.New[string]()
	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	test.Equal(t, s.Size(), 4)
}

func TestCapacity(t *testing.T) {
	s := stack.WithCapacity[int](10)
	test.Equal(t, s.Capacity(), 10)
}

func TestItems(t *testing.T) {
	s := stack.New[string]()
	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	want := []string{"kenobi", "general", "there", "hello"}
	got := slices.Collect(s.All())

	test.EqualFunc(t, got, want, slices.Equal)
}

func TestString(t *testing.T) {
	s := stack.New[string]()
	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	want := "[hello there general kenobi]"

	test.Equal(t, s.String(), want)
}

func TestFrom(t *testing.T) {
	items := []string{"cheese", "apples", "wine", "beer"}

	s := stack.From(items)

	test.Equal(t, s.Size(), 4)

	first, ok := s.Pop()
	test.True(t, ok)
	test.Equal(t, first, "beer")

	second, ok := s.Pop()
	test.True(t, ok)
	test.Equal(t, second, "wine")
}

func TestCollect(t *testing.T) {
	items := []string{"cheese", "apples", "wine", "beer"}

	s := stack.Collect(slices.Values(items))

	test.Equal(t, s.Size(), 4)

	first, ok := s.Pop()
	test.True(t, ok)
	test.Equal(t, first, "beer")

	second, ok := s.Pop()
	test.True(t, ok)
	test.Equal(t, second, "wine")
}

func BenchmarkStack(b *testing.B) {
	s := stack.New[int]()

	for b.Loop() {
		s.Push(1)

		_, ok := s.Pop()
		if !ok {
			b.Error("Pop() returned ok=false")
		}
	}
}
