package set_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/FollowTheProcess/collections/set"
	"github.com/FollowTheProcess/test"
)

func TestInsert(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		s := set.New[string]()

		test.True(t, s.IsEmpty()) // Initial set was not empty

		test.True(t, s.Insert("foo"))    // Inserting foo for the first time should return true
		test.False(t, s.Insert("foo"))   // Second insert of foo should return false
		test.True(t, s.Contains("foo"))  // Set didn't contain "foo"
		test.False(t, s.Contains("bar")) // Set said it contained "bar" but shouldn't have

		// testing nil safety
		danger := &set.Set[string]{}
		test.True(t, danger.Insert("bar"))
		test.False(t, danger.Insert("bar"))
		test.True(t, danger.Insert("baz"))
	})
	t.Run("ints", func(t *testing.T) {
		s := set.New[int]()

		test.True(t, s.Insert(1))    // Inserting 1 for the first time should return true
		test.False(t, s.Insert(1))   // Second insert of 1 should return false
		test.True(t, s.Contains(1))  // Set didn't contain 1
		test.False(t, s.Contains(2)) // Set said it contained 2 but shouldn't have

		// testing nil safety
		danger := &set.Set[int]{}
		test.True(t, danger.Insert(42))
		test.False(t, danger.Insert(42))
		test.True(t, danger.Insert(69))
	})
}

func TestFrom(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		// testing From
		items := []string{"foo", "bar", "baz"}
		other := set.From(items)

		test.True(t, other.Contains("foo"))      // Set from items didn't contain "foo"
		test.True(t, other.Contains("bar"))      // Set from items didn't contain "bar"
		test.True(t, other.Contains("baz"))      // Set from items didn't contain "baz"
		test.False(t, other.Contains("missing")) // Missing item "missing" reported as present
	})

	t.Run("floats", func(t *testing.T) {
		// testing From
		items := []float64{3.14159, 42.58, 69.73}
		other := set.From(items)

		test.True(t, other.Contains(3.14159)) // Set from items didn't contain pi
		test.True(t, other.Contains(42.58))   // Set from items didn't contain 42.58
		test.True(t, other.Contains(69.73))   // Set from items didn't contain 69.73
		test.False(t, other.Contains(100.1))  // Missing item 100.1 reported as present
	})
}

func TestRemove(t *testing.T) {
	t.Run("structs", func(t *testing.T) {
		type person struct {
			name string
			age  int
		}

		set := set.New[person]()

		missing := person{name: "Missing", age: 42}

		// Remove on an empty set shouldn't panic or do anything bad
		set.Remove(missing)

		tom := person{name: "Tom", age: 30}
		gandalf := person{name: "Gandalf", age: 55000}
		wendy := person{name: "Wendy", age: 12}

		set.Insert(tom)
		set.Insert(gandalf)
		set.Insert(wendy)

		test.Equal(t, set.Size(), 3) // Incorrect size

		// Sorry Wendy
		removed := set.Remove(wendy)
		test.True(t, removed) // Removed should be true

		test.Equal(t, set.Size(), 2) // Incorrect size after killing wendy

		test.False(t, set.Contains(wendy))
	})
}

func TestItems(t *testing.T) {
	items := []string{"cheese", "apples", "oranges", "milk"}
	slices.Sort(items)
	set := set.From(items)

	got := slices.Sorted(set.All())
	test.EqualFunc(t, got, items, slices.Equal)
}

func TestCollect(t *testing.T) {
	items := []string{"cheese", "apples", "oranges", "milk"}
	slices.Sort(items)

	set := set.Collect(slices.Values(items))

	got := slices.Sorted(set.All())
	test.EqualFunc(t, got, items, slices.Equal)
}

func TestEqual(t *testing.T) {
	tests := []struct {
		a, b *set.Set[string] // Sets the compare
		name string           // Name of the test case
		want bool             // Whether they should be considered equal
	}{
		{
			name: "nil",
			a:    nil,
			b:    nil,
			want: false,
		},
		{
			name: "empty",
			a:    set.New[string](),
			b:    set.New[string](),
			want: true,
		},
		{
			name: "equal",
			a:    set.From([]string{"hello", "there"}),
			b:    set.From([]string{"there", "hello"}),
			want: true,
		},
		{
			name: "not equal same length",
			a:    set.From([]string{"hello", "there"}),
			b:    set.From([]string{"goodbye", "yes"}),
			want: false,
		},
		{
			name: "not equal different length",
			a:    set.From([]string{"hello", "there"}),
			b:    set.From([]string{"goodbye", "yes", "more", "things"}),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, set.Equal(tt.a, tt.b), tt.want)
		})
	}
}

func TestUnion(t *testing.T) {
	tests := []struct {
		want *set.Set[string]   // The expected union set
		name string             // The name of the test case
		sets []*set.Set[string] // The sets to pass to Union
	}{
		{
			name: "nil",
			sets: nil,
			want: set.New[string](),
		},
		{
			name: "empty",
			sets: []*set.Set[string]{set.New[string]()},
			want: set.New[string](),
		},
		{
			name: "one empty",
			sets: []*set.Set[string]{
				set.New[string](),
				set.From([]string{"hello", "there"}),
			},
			want: set.From([]string{"hello", "there"}),
		},
		{
			name: "two empty",
			sets: []*set.Set[string]{
				set.New[string](),
				set.New[string](),
			},
			want: set.New[string](),
		},
		{
			name: "three full",
			sets: []*set.Set[string]{
				set.From([]string{"hello", "there", "general", "kenobi"}),
				set.From([]string{"hello", "to", "you", "too"}),
				set.From([]string{"hello", "again", "from", "another"}),
			},
			want: set.From([]string{
				"again",
				"another",
				"from",
				"general",
				"hello",
				"kenobi",
				"there",
				"to",
				"too",
				"you",
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.EqualFunc(t, set.Union(tt.sets...), tt.want, set.Equal)
		})
	}
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		want *set.Set[string]   // The expected intersection set
		name string             // The name of the test case
		sets []*set.Set[string] // The sets to pass to Intersection
	}{
		{
			name: "nil",
			sets: nil,
			want: set.New[string](),
		},
		{
			name: "empty",
			sets: []*set.Set[string]{set.New[string]()},
			want: set.New[string](),
		},
		{
			name: "one empty",
			sets: []*set.Set[string]{
				set.New[string](),
				set.From([]string{"hello", "there"}),
			},
			want: set.New[string](),
		},
		{
			name: "two empty",
			sets: []*set.Set[string]{
				set.New[string](),
				set.New[string](),
			},
			want: set.New[string](),
		},
		{
			name: "three full",
			sets: []*set.Set[string]{
				set.From([]string{"hello", "there", "general", "kenobi"}),
				set.From([]string{"hello", "to", "you", "too"}),
				set.From([]string{"hello", "from", "another", "set"}),
			},
			want: set.From([]string{"hello"}), // hello is the only item common to all
		},
		{
			name: "no common items",
			sets: []*set.Set[string]{
				set.From([]string{"hello", "there", "general", "kenobi"}),
				set.From([]string{"random", "other", "words"}),
				set.From([]string{"oh", "no", "nothing", "matches"}),
				set.From([]string{"these", "sets", "don't"}),
				set.From([]string{"share", "anything"}),
			},
			want: set.New[string](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.EqualFunc(t, set.Intersection(tt.sets...), tt.want, set.Equal)
		})
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		target *set.Set[string]   // The target set
		want   *set.Set[string]   // The expected difference set
		name   string             // The name of the test case
		others []*set.Set[string] // The sets to pass to Difference
	}{
		{
			name:   "nil",
			target: nil,
			others: nil,
			want:   set.New[string](),
		},
		{
			name:   "empty",
			target: set.New[string](),
			others: []*set.Set[string]{set.New[string]()},
			want:   set.New[string](),
		},
		{
			name:   "empty target",
			target: set.New[string](),
			others: []*set.Set[string]{
				set.From([]string{"some", "stuff", "here"}),
				set.From([]string{"more", "here", "too"}),
			},
			want: set.New[string](),
		},
		{
			name:   "full target empty others",
			target: set.From([]string{"target", "has", "items"}),
			others: []*set.Set[string]{
				set.New[string](),
			},
			want: set.From([]string{"target", "has", "items"}),
		},
		{
			name:   "three full",
			target: set.From([]string{"hello", "there", "general", "kenobi"}),
			others: []*set.Set[string]{
				set.From([]string{"hello", "to", "you", "to"}),
				set.From([]string{"hello", "from", "another", "set", "kenobi"}),
			},
			want: set.From([]string{"general", "there"}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.EqualFunc(t, set.Difference(tt.target, tt.others...), tt.want, set.Equal)
		})
	}
}

func TestIsDisjoint(t *testing.T) {
	tests := []struct {
		name string          // Name of the test case
		sets []*set.Set[int] // The sets to pass to IsDisjoint
		want bool            // Expected answer
	}{
		{
			name: "nil",
			sets: nil,
			want: false,
		},
		{
			name: "empty",
			sets: []*set.Set[int]{set.New[int]()},
			want: false,
		},
		{
			name: "one empty",
			sets: []*set.Set[int]{
				set.New[int](),
				set.From([]int{1, 2}),
			},
			want: true,
		},
		{
			name: "two empty",
			sets: []*set.Set[int]{
				set.New[int](),
				set.New[int](),
			},
			want: true,
		},
		{
			name: "three full",
			sets: []*set.Set[int]{
				set.From([]int{1, 2, 3, 4}),
				set.From([]int{1, 5, 6, 7}),
				set.From([]int{5, 2, 6, 3}),
			},
			want: false,
		},
		{
			name: "no common items",
			sets: []*set.Set[int]{
				set.From([]int{1, 2, 3, 4}),
				set.From([]int{5, 6, 7, 8, 9}),
				set.From([]int{10, 11, 12, 13}),
				set.From([]int{14}),
				set.From([]int{15, 16}),
			},
			want: true, // Nothing in common between any set
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, set.IsDisjoint(tt.sets...), tt.want)
		})
	}
}

func TestIsSubset(t *testing.T) {
	tests := []struct {
		a, b *set.Set[string] // The sets to compare
		name string           // Name of the test case
		want bool             // Expected answer
	}{
		{
			name: "nil",
			a:    nil,
			b:    nil,
			want: false,
		},
		{
			name: "both empty",
			a:    set.New[string](),
			b:    set.New[string](),
			want: false,
		},
		{
			name: "a empty",
			a:    set.New[string](),
			b:    set.From([]string{"one", "two", "three"}),
			want: false,
		},
		{
			name: "b empty",
			a:    set.From([]string{"one", "two", "three"}),
			b:    set.New[string](),
			want: false,
		},
		{
			name: "valid subset",
			a:    set.From([]string{"one", "two", "three"}),
			b:    set.From([]string{"four", "two", "five", "one", "three"}),
			want: true,
		},
		{
			name: "not a subset",
			a:    set.From([]string{"one", "two", "three"}),
			b:    set.From([]string{"four", "two", "five", "one"}), // Missing "three"
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, set.IsSubset(tt.a, tt.b), tt.want)
		})
	}
}

func TestIsSuperset(t *testing.T) {
	tests := []struct {
		a, b *set.Set[string] // The sets to compare
		name string           // Name of the test case
		want bool             // Expected answer
	}{
		{
			name: "nil",
			a:    nil,
			b:    nil,
			want: false,
		},
		{
			name: "both empty",
			a:    set.New[string](),
			b:    set.New[string](),
			want: false,
		},
		{
			name: "a empty",
			a:    set.New[string](),
			b:    set.From([]string{"one", "two", "three"}),
			want: false,
		},
		{
			name: "b empty",
			a:    set.From([]string{"one", "two", "three"}),
			b:    set.New[string](),
			want: false,
		},
		{
			name: "not a superset",
			a:    set.From([]string{"one", "two"}),
			b:    set.From([]string{"zero", "one"}),
			want: false,
		},
		{
			name: "valid superset",
			a:    set.From([]string{"zero", "one", "two"}),
			b:    set.From([]string{"one", "two"}),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, set.IsSuperset(tt.a, tt.b), tt.want)
		})
	}
}

func TestSymmetricDifference(t *testing.T) {
	tests := []struct {
		a, b *set.Set[string] // The sets to compare
		want *set.Set[string] // Expected answer
		name string           // Name of the test case
	}{
		{
			name: "nil",
			a:    nil,
			b:    nil,
			want: set.New[string](),
		},
		{
			name: "a empty",
			a:    set.New[string](),
			b:    set.From([]string{"some", "stuff", "here"}),
			want: set.From([]string{"some", "stuff", "here"}),
		},
		{
			name: "b empty",
			a:    set.From([]string{"some", "stuff", "here"}),
			b:    set.New[string](),
			want: set.From([]string{"some", "stuff", "here"}),
		},
		{
			name: "both empty",
			a:    set.New[string](),
			b:    set.New[string](),
			want: set.New[string](),
		},
		{
			name: "actual difference",
			a:    set.From([]string{"one", "two", "three", "four"}),
			b:    set.From([]string{"two", "three", "four", "five"}),
			want: set.From([]string{"one", "five"}), // "one" is only in a and "five" is only in b
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.EqualFunc(t, set.SymmetricDifference(tt.a, tt.b), tt.want, set.Equal)
		})
	}
}

func TestString(t *testing.T) {
	s := set.New[string]()

	s.Insert("cheese")
	s.Insert("apples")
	s.Insert("oranges")
	s.Insert("wine")

	got := s.String()

	// A set is an unordered collection and it's pointless to sort just for string representation
	// so we just check the existence of the items
	targets := []string{"cheese", "apples", "oranges", "wine"}
	for _, target := range targets {
		test.True(t, strings.Contains(got, target))
	}
}

func ExampleUnion() {
	this := set.New[string]()
	that := set.New[string]()

	this.Insert("hello")
	this.Insert("there")

	that.Insert("general")
	that.Insert("kenobi")
	that.Insert("says")
	that.Insert("hello")

	// Get the union in a slice of strings
	union := slices.Sorted(set.Union(this, that).All()) // A set is unordered

	fmt.Println(union)
	// Output: [general hello kenobi says there]
}

func ExampleIntersection() {
	this := set.New[string]()
	that := set.New[string]()

	this.Insert("hello")
	this.Insert("there")

	that.Insert("general")
	that.Insert("kenobi")
	that.Insert("says")
	that.Insert("hello")

	// Get the items in a slice of strings
	intersection := slices.Sorted(set.Intersection(this, that).All())

	fmt.Println(intersection)
	// Output: [hello]
}

func ExampleDifference() {
	this := set.New[string]()
	that := set.New[string]()

	this.Insert("hello")
	this.Insert("there")

	that.Insert("general")
	that.Insert("kenobi")
	that.Insert("says")
	that.Insert("hello")

	// Get the items in a slice of strings
	difference := slices.Sorted(set.Difference(this, that).All())

	fmt.Println(difference)
	// Output: [there]
}

func ExampleSymmetricDifference() {
	this := set.From([]int{1, 2, 3, 4})
	that := set.From([]int{2, 3, 4, 5})

	// Symmetric difference is the items that are in "this" or "that"
	// but not both
	difference := slices.Sorted(set.SymmetricDifference(this, that).All())

	fmt.Println(difference)
	// Output: [1 5]
}

func BenchmarkIntersection(b *testing.B) {
	s1 := set.New[int]()
	s2 := set.New[int]()

	for i := 0; i < 1000; i++ {
		s1.Insert(i)
		s2.Insert(i + 500)
	}

	b.ResetTimer()
	for range b.N {
		set.Intersection(s1, s2)
	}
}

func BenchmarkInsert(b *testing.B) {
	s := set.New[int]()
	for range b.N {
		s.Insert(b.N)
	}
}

func BenchmarkContains(b *testing.B) {
	s := set.New[int]()

	// So some will be present and others won't
	for i := range 10000 {
		s.Insert(i)
	}

	b.ResetTimer()
	for range b.N {
		s.Contains(b.N)
	}
}

func BenchmarkUnion(b *testing.B) {
	s1 := set.New[int]()
	s2 := set.New[int]()

	for i := 0; i < 1000; i++ {
		s1.Insert(i)
		s2.Insert(i + 500)
	}

	b.ResetTimer()
	for range b.N {
		set.Union(s1, s2)
	}
}

func BenchmarkDifference(b *testing.B) {
	s1 := set.New[int]()
	s2 := set.New[int]()

	for i := 0; i < 1000; i++ {
		s1.Insert(i)
		s2.Insert(i + 500)
	}

	b.ResetTimer()
	for range b.N {
		set.Difference(s1, s2)
	}
}

func BenchmarkIsDisjoint(b *testing.B) {
	s1 := set.New[int]()
	s2 := set.New[int]()

	for i := 0; i < 1000; i++ {
		s1.Insert(i)
		s2.Insert(i + 500)
	}

	b.ResetTimer()
	for range b.N {
		set.IsDisjoint(s1, s2)
	}
}

func BenchmarkIsSubset(b *testing.B) {
	s1 := set.New[int]()
	s2 := set.New[int]()

	for i := 0; i < 1000; i++ {
		s1.Insert(i)
		s2.Insert(i + 500)
	}

	b.ResetTimer()
	for range b.N {
		set.IsSubset(s1, s2)
	}
}

func BenchmarkSymmetricDifference(b *testing.B) {
	s1 := set.New[int]()
	s2 := set.New[int]()

	for i := 0; i < 1000; i++ {
		s1.Insert(i)
		s2.Insert(i + 500)
	}

	b.ResetTimer()
	for range b.N {
		set.SymmetricDifference(s1, s2)
	}
}
