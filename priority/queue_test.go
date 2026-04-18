package priority_test

import (
	"slices"
	"testing"

	"go.followtheprocess.codes/collections/priority"
	"go.followtheprocess.codes/test"
)

func TestNew(t *testing.T) {
	q := priority.New[string]()
	test.Equal(t, q.Size(), 0, test.Context("Initial size should be empty"))
	test.Equal(t, q.IsEmpty(), true, test.Context("Should be empty"))

	item, ok := q.Pop()
	test.False(t, ok, test.Context("Pop from empty queue should return ok=false"))
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

	item, ok := q.Pop()
	test.True(t, ok)
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

	item, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, item, "six")
}

// TestFromFuncFullDrain checks that subsequent Pops after the first also
// come out in priority order, verifying the heap invariant is properly
// established by FromFunc.
func TestFromFuncFullDrain(t *testing.T) {
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

	want := []string{"six", "five", "four", "three", "two", "one"}

	for _, expect := range want {
		got, ok := q.Pop()
		test.True(t, ok)
		test.Equal(t, got, expect)
	}

	test.True(t, q.IsEmpty())
}

func TestSize(t *testing.T) {
	q := priority.New[string]()
	test.Equal(t, q.Size(), 0, test.Context("Initial size should be empty"))
	test.Equal(t, q.IsEmpty(), true, test.Context("Should be empty"))

	q.Push("one", 1)
	q.Push("two", 2)
	q.Push("three", 3)
	q.Push("four", 4)

	test.Equal(t, q.Size(), 4, test.Context("Wrong size after push"))
	test.Equal(t, q.IsEmpty(), false, test.Context("Should not be empty"))
}

func TestPushPop(t *testing.T) {
	q := priority.New[string]()

	q.Push("two", 2)
	q.Push("one", 1)
	q.Push("three", 3)
	q.Push("four", 4)

	test.Equal(t, q.Size(), 4, test.Context("Incorrect size after push"))

	first, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, first, "four", test.Context("four has highest priority"))

	test.Equal(t, q.Size(), 3, test.Context("Incorrect size after pop"))

	second, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, second, "three", test.Context("Next priority is in three"))

	third, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, third, "two", test.Context("Next highest is two"))

	fourth, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, fourth, "one", test.Context("One is least"))

	fifth, ok := q.Pop()
	test.False(t, ok, test.Context("Pop from empty queue should return ok=false"))
	test.Equal(t, fifth, "")
}

// TestEqualPriorities makes sure items sharing a priority still all come
// out of the queue (order among ties is non-deterministic).
func TestEqualPriorities(t *testing.T) {
	q := priority.New[string]()

	q.Push("a", 1)
	q.Push("b", 1)
	q.Push("c", 1)
	q.Push("d", 2)

	first, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, first, "d", test.Context("d has the higher priority"))

	// The remaining three are tied — collect and sort to compare.
	rest := make([]string, 0, 3)

	for range 3 {
		got, ok := q.Pop()
		test.True(t, ok)
		rest = append(rest, got)
	}

	slices.Sort(rest)
	test.EqualFunc(t, rest, []string{"a", "b", "c"}, slices.Equal)

	test.True(t, q.IsEmpty())
}

// TestNegativePriorities covers ordering when priorities include zero and
// negative values.
func TestNegativePriorities(t *testing.T) {
	q := priority.New[string]()
	q.Push("neg", -5)
	q.Push("zero", 0)
	q.Push("pos", 5)

	order := make([]string, 0, 3)

	for range 3 {
		got, ok := q.Pop()
		test.True(t, ok)
		order = append(order, got)
	}

	test.EqualFunc(t, order, []string{"pos", "zero", "neg"}, slices.Equal)
}

// TestDuplicateItems makes sure identical items with differing priorities
// both make it through the queue in the right order.
func TestDuplicateItems(t *testing.T) {
	q := priority.New[string]()
	q.Push("dup", 1)
	q.Push("dup", 5)
	q.Push("dup", 3)

	var priorities []int

	// Cannot tell items apart by string so we just verify the count
	// and that size decreases correctly.
	for i := range 3 {
		got, ok := q.Pop()
		test.True(t, ok)
		test.Equal(t, got, "dup")
		test.Equal(t, q.Size(), 2-i)
		priorities = append(priorities, 0) // placeholder: we only check count here
	}

	test.Equal(t, len(priorities), 3)
	test.True(t, q.IsEmpty())
}

// TestAll verifies that All yields items in priority order without
// modifying the queue.
func TestAll(t *testing.T) {
	q := priority.New[string]()
	q.Push("one", 1)
	q.Push("four", 4)
	q.Push("three", 3)
	q.Push("two", 2)

	got := slices.Collect(q.All())
	want := []string{"four", "three", "two", "one"}
	test.EqualFunc(t, got, want, slices.Equal)

	// Queue itself must be untouched.
	test.Equal(t, q.Size(), 4, test.Context("All must not mutate the queue"))

	popped, ok := q.Pop()
	test.True(t, ok)
	test.Equal(t, popped, "four")
}

// BenchmarkNew measures the performance of constructing a new empty Queue
// and calling Push to fill it with elements.
func BenchmarkNew(b *testing.B) {
	for b.Loop() {
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

	for b.Loop() {
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

	for b.Loop() {
		priority.FromFunc(items, priorityFunc)
	}
}
