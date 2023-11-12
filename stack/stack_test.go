package stack_test

import (
	"reflect"
	"testing"

	"github.com/FollowTheProcess/collections/stack"
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

	if s.IsEmpty() {
		t.Error("IsEmpty should return false")
	}
}

func TestLength(t *testing.T) {
	s := stack.New[string]()
	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	if s.Length() != 4 {
		t.Errorf("wrong length: got %d, wanted %d", s.Length(), 4)
	}
}

func TestCap(t *testing.T) {
	s := stack.New[int](stack.WithCapacity(10))
	if s.Cap() != 10 {
		t.Errorf("wrong capacity: got %d, wanted %d", s.Cap(), 10)
	}

	s2 := stack.New[int]()
	if s2.Cap() != 0 {
		t.Errorf("wrong capacity: got %d, wanted %d", s2.Cap(), 0)
	}

	s3 := stack.New[int](stack.WithCapacity(-1))
	if s3.Cap() != 0 {
		t.Errorf("wrong capacity: got %d, wanted %d", s3.Cap(), 0)
	}
}

func TestPop(t *testing.T) {
	s := stack.New[string]()
	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	item, err := s.Pop()
	if err != nil {
		t.Errorf("Pop() returned an error: %v", err)
	}
	if item != "kenobi" {
		t.Errorf("wrong item popped: got %q, wanted %q", item, "kenobi")
	}

	item, err = s.Pop()
	if err != nil {
		t.Errorf("Pop() returned an error: %v", err)
	}
	if item != "general" {
		t.Errorf("wrong item popped: got %q, wanted %q", item, "general")
	}

	item, err = s.Pop()
	if err != nil {
		t.Errorf("Pop() returned an error: %v", err)
	}
	if item != "there" {
		t.Errorf("wrong item popped: got %q, wanted %q", item, "there")
	}

	item, err = s.Pop()
	if err != nil {
		t.Errorf("Pop() returned an error: %v", err)
	}
	if item != "hello" {
		t.Errorf("wrong item popped: got %q, wanted %q", item, "hello")
	}

	// Try one more pop, should error
	item, err = s.Pop()
	if err == nil {
		t.Error("expected pop from empty stack error, got nil")
	}

	// Err should be a ErrPopFromEmptyStack
	if err != stack.ErrPopFromEmptyStack {
		t.Errorf("wrong error returned: got %v, wanted %v", err, stack.ErrPopFromEmptyStack)
	}

	// Item should be the zero value for the type
	if item != "" {
		t.Errorf("empty pop should be zero value: got %q, wanted %q", item, "")
	}
}

func TestItems(t *testing.T) {
	s := stack.New[string]()
	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	want := []string{"hello", "there", "general", "kenobi"}

	if !reflect.DeepEqual(s.Items(), want) {
		t.Errorf("wrong items: got %v, wanted %v", s.Items(), want)
	}
}

func TestString(t *testing.T) {
	s := stack.New[string]()
	s.Push("hello")
	s.Push("there")
	s.Push("general")
	s.Push("kenobi")

	want := "[hello there general kenobi]"

	if s.String() != want {
		t.Errorf("wrong string: got %s, wanted %s", s.String(), want)
	}
}

func TestNotNew(t *testing.T) {
	s := stack.Stack[int]{}
	s.Push(1)
	s.Push(2)

	first, err := s.Pop()
	if err != nil {
		t.Fatalf("stack.Pop() returned an unexpected error: %v", err)
	}

	if first != 2 {
		t.Errorf("wrong item popped: got %d, wanted %d", first, 2)
	}
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
