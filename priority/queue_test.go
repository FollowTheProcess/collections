package priority_test

import (
	"testing"

	"github.com/FollowTheProcess/collections/priority"
	"github.com/FollowTheProcess/test"
)

func TestNew(t *testing.T) {
	q := priority.New[string]()
	test.Equal(t, q.Size(), 0)     // Initial size should be empty
	test.Equal(t, q.Empty(), true) // Should be empty

	item, err := q.Pop()
	test.Err(t, err) // Pop from empty queue
	test.Equal(t, item, "")
}

func TestSize(t *testing.T) {
	q := priority.New[string]()
	test.Equal(t, q.Size(), 0)     // Initial size should be empty
	test.Equal(t, q.Empty(), true) // Should be empty

	q.Push("one", 1)
	q.Push("two", 2)
	q.Push("three", 3)
	q.Push("four", 4)

	test.Equal(t, q.Size(), 4)      // Wrong size after push
	test.Equal(t, q.Empty(), false) // Should not be empty
}

func TestPushPop(t *testing.T) {
	q := priority.New[string]()

	q.Push("two", 2)
	q.Push("one", 1)
	q.Push("three", 3)
	q.Push("four", 4)

	test.Equal(t, q.Size(), 4) // Incorrect size after push

	first, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, first, "four") // four has highest priority

	test.Equal(t, q.Size(), 3) // Incorrect size after pop

	second, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, second, "three") // Next priority is in three

	third, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, third, "two") // Next highest is two

	fourth, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, fourth, "one") // One is least

	fifth, err := q.Pop()
	test.Err(t, err) // Pop from empty queue
	test.Equal(t, fifth, "")
}
