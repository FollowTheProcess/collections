package stack_test

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/collections/stack"
	"github.com/FollowTheProcess/test"
)

func TestIsEmpty(t *testing.T) {
	s := stack.New[string]()

	if !s.Empty() {
		t.Error("IsEmpty should return true")
	}

	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	test.False(t, s.Empty()) // Empty should not be false
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

func TestLength(t *testing.T) {
	s := stack.New[string]()
	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	test.Equal(t, s.Length(), 4)
}

func TestItems(t *testing.T) {
	s := stack.New[string]()
	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	want := []string{"hello", "there", "general", "kenobi"}

	test.EqualFunc(t, s.Items(), want, slices.Equal)
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
