package queue_test

import (
	"reflect"
	"testing"

	"github.com/FollowTheProcess/collections/queue"
)

func TestIsEmpty(t *testing.T) {
	q := queue.New[string]()

	if !q.IsEmpty() {
		t.Error("IsEmpty should return true")
	}

	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	if q.IsEmpty() {
		t.Error("IsEmpty should return false")
	}
}

func TestLength(t *testing.T) {
	q := queue.New[string]()
	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	if q.Length() != 4 {
		t.Errorf("wrong length: got %d, wanted %d", q.Length(), 4)
	}
}

func TestPop(t *testing.T) {
	q := queue.New[string]()
	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	item, err := q.Pop()
	if err != nil {
		t.Errorf("Pop() returned an error: %v", err)
	}
	if item != "hello" {
		t.Errorf("wrong item popped: got %q, wanted %q", item, "hello")
	}

	item, err = q.Pop()
	if err != nil {
		t.Errorf("Pop() returned an error: %v", err)
	}
	if item != "there" {
		t.Errorf("wrong item popped: got %q, wanted %q", item, "there")
	}

	item, err = q.Pop()
	if err != nil {
		t.Errorf("Pop() returned an error: %v", err)
	}
	if item != "general" {
		t.Errorf("wrong item popped: got %q, wanted %q", item, "general")
	}

	item, err = q.Pop()
	if err != nil {
		t.Errorf("Pop() returned an error: %v", err)
	}
	if item != "kenobi" {
		t.Errorf("wrong item popped: got %q, wanted %q", item, "kenobi")
	}

	// Try one more pop, should error
	_, err = q.Pop()
	if err == nil {
		t.Error("expected pop from empty queue, got nil")
	}
}

func TestItems(t *testing.T) {
	q := queue.New[string]()
	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	want := []string{"hello", "there", "general", "kenobi"}

	if !reflect.DeepEqual(q.Items(), want) {
		t.Errorf("wrong items: got %v, wanted %v", q.Items(), want)
	}
}

func TestString(t *testing.T) {
	q := queue.New[string]()
	q.Push("hello")
	q.Push("there")
	q.Push("general")
	q.Push("kenobi")

	want := "[hello there general kenobi]"

	if q.String() != want {
		t.Errorf("wrong string: got %s, wanted %s", q.String(), want)
	}
}
