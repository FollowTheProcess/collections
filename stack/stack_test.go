package stack_test

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/collections/stack"
	"github.com/FollowTheProcess/test"
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

	test.False(t, s.IsEmpty()) // Empty should not be false
}

func TestPop(t *testing.T) {
	s := stack.New[string]()
	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	item, err := s.Pop()
	test.Ok(t, err)
	test.Equal(t, item, "kenobi")

	item, err = s.Pop()
	test.Ok(t, err)
	test.Equal(t, item, "general")

	item, err = s.Pop()
	test.Ok(t, err)
	test.Equal(t, item, "there")

	item, err = s.Pop()
	test.Ok(t, err)
	test.Equal(t, item, "hello")

	// Try one more pop, should error
	item, err = s.Pop()
	test.Err(t, err)        // Pop from empty stack should error
	test.Equal(t, item, "") // Item should be the zero value
}

func TestNotNew(t *testing.T) {
	s := stack.Stack[int]{}
	s.Push(1)
	s.Push(2)

	first, err := s.Pop()
	test.Ok(t, err)
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
	got := slices.Collect(s.Items())

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

	first, err := s.Pop()
	test.Ok(t, err)
	test.Equal(t, first, "beer")

	second, err := s.Pop()
	test.Ok(t, err)
	test.Equal(t, second, "wine")
}

func TestCollect(t *testing.T) {
	items := []string{"cheese", "apples", "wine", "beer"}

	s := stack.Collect(slices.Values(items))

	test.Equal(t, s.Size(), 4)

	first, err := s.Pop()
	test.Ok(t, err)
	test.Equal(t, first, "beer")

	second, err := s.Pop()
	test.Ok(t, err)
	test.Equal(t, second, "wine")
}

func BenchmarkStack(b *testing.B) {
	s := stack.New[int]()

	for i := 0; i < b.N; i++ {
		s.Push(i)
	}

	for i := 0; i < b.N; i++ {
		_, err := s.Pop()
		if err != nil {
			b.Errorf("Pop() returned an error: %v", err)
		}
	}
}
