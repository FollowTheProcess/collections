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

func TestFrom(t *testing.T) {
	items := []priority.Element[string]{
		{Item: "six", Priority: 6},
		{Item: "one", Priority: 1},
		{Item: "three", Priority: 3},
		{Item: "four", Priority: 4},
		{Item: "two", Priority: 2},
		{Item: "five", Priority: 5},
	}

	q := priority.From(items)

	test.Equal(t, q.Size(), 6)

	item, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, item, "six")
}

func TestFromFunc(t *testing.T) {
	items := []string{"one", "two", "three", "four", "five", "six"}

	priorityFunc := func(item string) int {
		switch item {
		case "one":
			return 1
		case "two":
			return 2
		case "three":
			return 3
		case "four":
			return 4
		case "five":
			return 5
		case "six":
			return 6
		default:
			return 0
		}
	}

	q := priority.FromFunc(items, priorityFunc)

	test.Equal(t, q.Size(), 6)

	item, err := q.Pop()
	test.Ok(t, err)
	test.Equal(t, item, "six")
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

// BenchmarkNew measures the performance of constructing a new empty Queue
// and calling Push to fill it with elements.
func BenchmarkNew(b *testing.B) {
	for range b.N {
		q := priority.New[string]()
		q.Push("one", 1)
		q.Push("two", 2)
		q.Push("three", 3)
		q.Push("four", 4)
		q.Push("five", 5)
		q.Push("six", 6)
	}
}

// BenchmarkFrom measures the performance of constructing a priority Queue
// from a pre-existing slice of Elements.
func BenchmarkFrom(b *testing.B) {
	elements := []priority.Element[string]{
		{Item: "six", Priority: 6},
		{Item: "one", Priority: 1},
		{Item: "three", Priority: 3},
		{Item: "four", Priority: 4},
		{Item: "two", Priority: 2},
		{Item: "five", Priority: 5},
	}

	for range b.N {
		priority.From(elements)
	}
}

// BenchmarkFromFunc measures the performance of constructing a priority Queue
// from a pre-existing slice of items and calculating the priority with a closure.
func BenchmarkFromFunc(b *testing.B) {
	items := []string{"one", "two", "three", "four", "five", "six"}

	priorityFunc := func(item string) int {
		switch item {
		case "one":
			return 1
		case "two":
			return 2
		case "three":
			return 3
		case "four":
			return 4
		case "five":
			return 5
		case "six":
			return 6
		default:
			return 0
		}
	}

	for range b.N {
		priority.FromFunc(items, priorityFunc)
	}
}
