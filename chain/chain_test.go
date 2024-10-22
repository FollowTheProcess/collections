package chain_test

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/collections/chain"
	"github.com/FollowTheProcess/test"
)

func TestNew(t *testing.T) {
	chain := chain.New[int, string]()
	test.Equal(t, chain.Size(), 0) // Wrong chain size before
}

func TestFrom(t *testing.T) {
	maps := []map[string]int{
		{"one": 1, "two": 2},
		{"three": 3, "four": 4},
		{"five": 5, "six": 6},
	}

	chain := chain.From(maps)
	test.Equal(t, chain.Size(), 3) // Wrong size (From)
}

func TestCollect(t *testing.T) {
	maps := []map[string]int{
		{"one": 1, "two": 2},
		{"three": 3, "four": 4},
		{"five": 5, "six": 6},
	}

	chain := chain.Collect(slices.Values(maps))
	test.Equal(t, chain.Size(), 3) // Wrong size (Collect)
}

func TestAppendGet(t *testing.T) {
	chain := chain.New[int, string]()

	// Maps to append
	one := map[int]string{
		1: "one in first map",
		2: "two in first map",
	}
	two := map[int]string{
		1: "one in second map",
		2: "two in second map",
		3: "three in second map",
	}
	three := map[int]string{
		1: "one in third map",
		2: "two in third map",
		3: "three in third map",
		4: "four in third map",
	}

	chain.Append(one)

	test.Equal(t, chain.Size(), 1) // Wrong chain size after append

	got, ok := chain.Get(1)
	test.True(t, ok)
	test.Equal(t, got, "one in first map")

	chain.Append(two)

	got, ok = chain.Get(1)
	test.True(t, ok)
	test.Equal(t, got, "one in first map") // Higher priority

	got, ok = chain.Get(2)
	test.True(t, ok)
	test.Equal(t, got, "two in first map") // Two is in higher priority map

	got, ok = chain.Get(3)
	test.True(t, ok)
	test.Equal(t, got, "three in second map") // Three is in second map

	chain.Append(three)

	got, ok = chain.Get(1)
	test.True(t, ok)
	test.Equal(t, got, "one in first map")

	got, ok = chain.Get(4)
	test.True(t, ok)
	test.Equal(t, got, "four in third map") // Four only exists in 3rd map

	got, ok = chain.Get(5)
	test.False(t, ok) // There is no 5
	test.Equal(t, got, "")
}

func TestPrependGet(t *testing.T) {
	chain := chain.New[int, string]()

	// Maps to prepend
	one := map[int]string{
		1: "one in first map",
		2: "two in first map",
	}
	two := map[int]string{
		1: "one in second map",
		2: "two in second map",
		3: "three in second map",
	}
	three := map[int]string{
		1: "one in third map",
		2: "two in third map",
		3: "three in third map",
		4: "four in third map",
	}

	chain.Prepend(one)

	test.Equal(t, chain.Size(), 1) // Wrong chain size after prepend

	got, ok := chain.Get(1)
	test.True(t, ok)
	test.Equal(t, got, "one in first map")

	chain.Prepend(two)

	got, ok = chain.Get(1)
	test.True(t, ok)
	test.Equal(t, got, "one in second map") // Prepended has higher priority

	got, ok = chain.Get(2)
	test.True(t, ok)
	test.Equal(t, got, "two in second map") // Two is in higher priority map

	got, ok = chain.Get(3)
	test.True(t, ok)
	test.Equal(t, got, "three in second map") // Three is in second map

	chain.Prepend(three)

	got, ok = chain.Get(1)
	test.True(t, ok)
	test.Equal(t, got, "one in third map") // Third map is now highest priority

	got, ok = chain.Get(4)
	test.True(t, ok)
	test.Equal(t, got, "four in third map") // Four only exists in 3rd map

	got, ok = chain.Get(5)
	test.False(t, ok) // There is no 5
	test.Equal(t, got, "")
}

func TestInsert(t *testing.T) {
	maps := []map[int]string{
		{
			1: "one in first map",
			2: "two in first map",
		},
		{
			1: "one in second map",
			2: "two in second map",
			3: "three in second map",
		},
		{
			1: "one in third map",
			2: "two in third map",
			3: "three in third map",
			4: "four in third map",
		},
	}

	chain := chain.From(maps)

	got, existed := chain.Insert(1, "updated one in first map")
	test.True(t, existed)                  // 1 exists in first map
	test.Equal(t, got, "one in first map") // Returned value should be the old value

	// If we get it now we should get the updated one
	got, ok := chain.Get(1)
	test.True(t, ok)
	test.Equal(t, got, "updated one in first map")

	got, existed = chain.Insert(3, "updated three in second map")
	test.True(t, existed)
	test.Equal(t, got, "three in second map")

	got, ok = chain.Get(3)
	test.True(t, ok)
	test.Equal(t, got, "updated three in second map")

	got, existed = chain.Insert(4, "updated four in third map")
	test.True(t, existed)
	test.Equal(t, got, "four in third map")

	got, ok = chain.Get(4)
	test.True(t, ok)
	test.Equal(t, got, "updated four in third map")

	// Brand new insertion goes into the first map
	got, existed = chain.Insert(5, "five brand new in first map")
	test.False(t, existed)
	test.Equal(t, got, "five brand new in first map")

	got, ok = chain.Get(5)
	test.True(t, ok)
	test.Equal(t, got, "five brand new in first map")
}

func TestEmptyMaps(t *testing.T) {
	chain := chain.New[string, int]()

	// Try inserting into an otherwise empty chain
	got, existed := chain.Insert("hello", 1)
	test.False(t, existed)
	test.Equal(t, got, 1)
}
