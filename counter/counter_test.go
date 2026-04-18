package counter_test

import (
	"maps"
	"slices"
	"testing"

	"go.followtheprocess.codes/collections/counter"
	"go.followtheprocess.codes/test"
)

func TestNewCounter(t *testing.T) {
	c := counter.New[string]()

	test.Equal(t, c.Size(), 0, test.Context("Initial size must be empty"))

	test.Equal(t, c.Add("apple"), 1, test.Context("First add returns 1"))
	test.Equal(t, c.Add("apple"), 2, test.Context("Second add returns 2"))

	test.Equal(t, c.Size(), 1, test.Context("\"apple\" should be the only item"))

	test.Equal(t, c.Add("orange"), 1, test.Context("First add returns 1 (orange)"))

	test.Equal(t, c.Size(), 2, test.Context("should now have 2 items"))

	test.Equal(t, c.Sub("orange"), 0, test.Context("Sub(\"orange\") should remove orange completely"))
	test.Equal(t, c.Sub("apple"), 1, test.Context("Sub(\"apple\") should decrement apple to 1"))
}

func TestCount(t *testing.T) {
	c := counter.New[string]()

	c.Add("human")
	c.Add("human")
	c.Add("dog")

	test.Equal(t, c.Get("human"), 2, test.Context("Wrong number of humans"))
	test.Equal(t, c.Get("dog"), 1, test.Context("Wrong number of dogs"))
	test.Equal(t, c.Get("cats"), 0, test.Context("No cats in the counter"))
}

func TestFrom(t *testing.T) {
	items := []int{1, 5, 2, 4, 8, 5, 4, 4, 6, 2, 3, 12}

	c := counter.From(items)

	test.Equal(t, c.Size(), 8, test.Context("Wrong size"))
	test.Equal(t, c.Get(1), 1, test.Context("Wrong number of 1s"))
	test.Equal(t, c.Get(2), 2, test.Context("Wrong number of 2s"))
	test.Equal(t, c.Get(3), 1, test.Context("Wrong number of 4s"))
	test.Equal(t, c.Get(4), 3, test.Context("Wrong number of 4s"))
	test.Equal(t, c.Get(5), 2, test.Context("Wrong number of 5s"))
	test.Equal(t, c.Get(6), 1, test.Context("Wrong number of 6s"))
	test.Equal(t, c.Get(8), 1, test.Context("Wrong number of 8s"))
	test.Equal(t, c.Get(12), 1, test.Context("Wrong number of 12s"))
}

func TestCollect(t *testing.T) {
	items := []int{1, 5, 2, 4, 8, 5, 4, 4, 6, 2, 3, 12}

	c := counter.Collect(slices.Values(items))

	test.Equal(t, c.Size(), 8, test.Context("Wrong size"))
	test.Equal(t, c.Get(1), 1, test.Context("Wrong number of 1s"))
	test.Equal(t, c.Get(2), 2, test.Context("Wrong number of 2s"))
	test.Equal(t, c.Get(3), 1, test.Context("Wrong number of 4s"))
	test.Equal(t, c.Get(4), 3, test.Context("Wrong number of 4s"))
	test.Equal(t, c.Get(5), 2, test.Context("Wrong number of 5s"))
	test.Equal(t, c.Get(6), 1, test.Context("Wrong number of 6s"))
	test.Equal(t, c.Get(8), 1, test.Context("Wrong number of 8s"))
	test.Equal(t, c.Get(12), 1, test.Context("Wrong number of 12s"))
}

func TestRemove(t *testing.T) {
	type person struct {
		name string
		age  uint
	}

	people := []person{
		{name: "Tom", age: 30},
		{name: "Tom", age: 30},
		{name: "Mark", age: 29},
		{name: "Alice", age: 17},
	}

	c := counter.From(people)

	tom := person{name: "Tom", age: 30}

	test.Equal(t, c.Size(), 3, test.Context("Wrong size"))

	test.Equal(t, c.Remove(tom), 2, test.Context("Wrong number of Toms before remove"))
	test.Equal(t, c.Get(tom), 0, test.Context("Wrong number of Toms after remove"))

	test.Equal(t, c.Remove(person{name: "missing", age: 35}), 0, test.Context("Missing person returns 0"))
}

func TestSum(t *testing.T) {
	fruits := []string{
		"apple",
		"apple",
		"orange",
		"banana",
		"raspberry",
		"raspberry",
		"strawberry",
		"cherry",
		"cherry",
		"pear",
	}

	c := counter.From(fruits)

	test.Equal(t, c.Sum(), len(fruits))
}

func TestReset(t *testing.T) {
	fruits := []string{
		"apple",
		"apple",
		"orange",
		"banana",
		"raspberry",
		"raspberry",
		"strawberry",
		"cherry",
		"cherry",
		"pear",
	}

	c := counter.From(fruits)

	test.Equal(t, c.Size(), 7, test.Context("Wrong size before Reset"))
	test.Equal(t, c.Sum(), len(fruits), test.Context("Wrong sum before Reset"))

	c.Reset()

	test.Equal(t, c.Size(), 0, test.Context("Wrong size after Reset"))
	test.Equal(t, c.Sum(), 0, test.Context("Wrong sum after Reset"))
}

func TestDescending(t *testing.T) {
	names := []string{
		"dave",
		"dave",
		"dave",
		"dave",
		"chris",
		"chris",
		"john",
		"john",
		"john",
		"john",
		"john",
		"mark",
		"alice",
		"alice",
		"alice",
	}

	c := counter.From(names)

	var items []string

	var counts []int

	for item, count := range c.Descending() {
		items = append(items, item)
		counts = append(counts, count)
	}

	wantItems := []string{"john", "dave", "alice", "chris", "mark"}
	wantCounts := []int{5, 4, 3, 2, 1}

	test.EqualFunc(t, items, wantItems, slices.Equal)
	test.EqualFunc(t, counts, wantCounts, slices.Equal)
}

func TestMostCommon(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		c := counter.New[int]()

		test.True(t, c.MostCommon(0) == nil, test.Context("empty counter, n=0"))
		test.True(t, c.MostCommon(1) == nil, test.Context("empty counter, n=1"))
		test.True(t, c.MostCommon(-1) == nil, test.Context("empty counter, n=-1"))
	})

	t.Run("n=0", func(t *testing.T) {
		c := counter.From([]string{"a", "a", "b"})

		test.True(t, c.MostCommon(0) == nil)
	})

	t.Run("n=-1", func(t *testing.T) {
		c := counter.From([]string{"a", "a", "b"})

		test.True(t, c.MostCommon(-1) == nil)
	})

	t.Run("n=1", func(t *testing.T) {
		names := []string{
			"dave", "dave", "dave", "dave",
			"chris", "chris",
			"john", "john", "john", "john", "john",
			"mark",
			"alice", "alice", "alice",
		}

		c := counter.From(names)

		got := c.MostCommon(1)
		want := []counter.Pair[string]{{Item: "john", Count: 5}}

		test.EqualFunc(t, got, want, slices.Equal)
	})

	t.Run("n less than size", func(t *testing.T) {
		names := []string{
			"dave", "dave", "dave", "dave",
			"chris", "chris",
			"john", "john", "john", "john", "john",
			"mark",
			"alice", "alice", "alice",
		}

		c := counter.From(names)

		got := c.MostCommon(3)
		want := []counter.Pair[string]{
			{Item: "john", Count: 5},
			{Item: "dave", Count: 4},
			{Item: "alice", Count: 3},
		}

		test.EqualFunc(t, got, want, slices.Equal)
	})

	t.Run("n equal to size", func(t *testing.T) {
		names := []string{
			"dave", "dave", "dave", "dave",
			"chris", "chris",
			"john", "john", "john", "john", "john",
			"mark",
			"alice", "alice", "alice",
		}

		c := counter.From(names)

		got := c.MostCommon(c.Size())
		want := []counter.Pair[string]{
			{Item: "john", Count: 5},
			{Item: "dave", Count: 4},
			{Item: "alice", Count: 3},
			{Item: "chris", Count: 2},
			{Item: "mark", Count: 1},
		}

		test.EqualFunc(t, got, want, slices.Equal)
	})

	t.Run("n greater than size", func(t *testing.T) {
		names := []string{
			"dave", "dave", "dave", "dave",
			"chris", "chris",
			"john", "john", "john", "john", "john",
			"mark",
			"alice", "alice", "alice",
		}

		c := counter.From(names)

		got := c.MostCommon(c.Size() + 100)
		want := []counter.Pair[string]{
			{Item: "john", Count: 5},
			{Item: "dave", Count: 4},
			{Item: "alice", Count: 3},
			{Item: "chris", Count: 2},
			{Item: "mark", Count: 1},
		}

		test.EqualFunc(t, got, want, slices.Equal)
	})

	t.Run("ties", func(t *testing.T) {
		// "apple" and "orange" both have count 2, "banana" has count 1.
		// The order of the two tied items is explicitly not specified;
		// we only assert that both appear in the top-2 result.
		c := counter.From([]string{"apple", "apple", "orange", "orange", "banana"})

		got := c.MostCommon(2)

		test.Equal(t, len(got), 2)

		items := []string{got[0].Item, got[1].Item}
		slices.Sort(items)

		test.EqualFunc(t, items, []string{"apple", "orange"}, slices.Equal)
		test.Equal(t, got[0].Count, 2)
		test.Equal(t, got[1].Count, 2)
	})
}

func TestAll(t *testing.T) {
	c := counter.New[string]()
	c.Add("one")
	c.Add("two")
	c.Add("two")
	c.Add("three")
	c.Add("three")
	c.Add("three")
	c.Add("four")
	c.Add("four")
	c.Add("four")
	c.Add("four")

	all := maps.Collect(c.All())

	want := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
	}

	test.EqualFunc(t, all, want, maps.Equal)
}

func TestItems(t *testing.T) {
	c := counter.New[string]()
	c.Add("one")
	c.Add("two")
	c.Add("two")
	c.Add("three")
	c.Add("three")
	c.Add("three")
	c.Add("four")
	c.Add("four")
	c.Add("four")
	c.Add("four")

	items := slices.Collect(c.Items())

	want := []string{"one", "two", "three", "four"}

	slices.Sort(items)
	slices.Sort(want)

	test.EqualFunc(t, items, want, slices.Equal)
}

func BenchmarkMostCommon(b *testing.B) {
	names := []string{
		"dave", "dave", "dave", "dave",
		"chris", "chris",
		"john", "john", "john", "john", "john",
		"mark",
		"alice", "alice", "alice",
	}

	c := counter.From(names)

	b.Run("n=1", func(b *testing.B) {
		for b.Loop() {
			c.MostCommon(1)
		}
	})

	b.Run("n=size", func(b *testing.B) {
		n := c.Size()
		for b.Loop() {
			c.MostCommon(n)
		}
	})
}

func BenchmarkDescending(b *testing.B) {
	names := []string{
		"dave",
		"dave",
		"dave",
		"dave",
		"chris",
		"chris",
		"john",
		"john",
		"john",
		"john",
		"john",
		"mark",
		"alice",
		"alice",
		"alice",
	}

	c := counter.From(names)

	for b.Loop() {
		// Just drain the iterator
		for name, count := range c.Descending() {
			_ = name
			_ = count
		}
	}
}
