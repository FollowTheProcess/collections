package set_test

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/FollowTheProcess/collections/set"
)

func TestBasics(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		s := set.New[string]()

		if !s.IsEmpty() {
			t.Error("IsEmpty should have returned true")
		}

		s.Add("hello")
		s.Add("tom")
		s.Add("hello")

		if s.IsEmpty() {
			t.Error("IsEmpty should have returned false")
		}

		if s.Length() != 2 {
			t.Errorf("wrong length: got %d, wanted %d", s.Length(), 2)
		}

		if !s.Contains("tom") {
			t.Error("expected item 'tom' not found in set")
		}

		if !s.Contains("hello") {
			t.Error("expected item 'hello' not found in set")
		}

		items := s.Items()
		want := []string{"hello", "tom"}
		sort.Strings(items)

		if !reflect.DeepEqual(want, items) {
			t.Errorf("slice mismatch: got %#v, wanted %#v", items, want)
		}

		s.Remove("tom")
		if s.Contains("tom") {
			t.Error("set contained deleted item after remove")
		}
	})

	t.Run("ints", func(t *testing.T) {
		s := set.New[int]()

		if !s.IsEmpty() {
			t.Error("IsEmpty should have returned true")
		}

		s.Add(100)
		s.Add(27)
		s.Add(100)

		if s.IsEmpty() {
			t.Error("IsEmpty should have returned false")
		}

		if s.Length() != 2 {
			t.Errorf("wrong length: got %d, wanted %d", s.Length(), 2)
		}

		if !s.Contains(100) {
			t.Error("expected item '27' not found in set")
		}

		if !s.Contains(27) {
			t.Error("expected item '100' not found in set")
		}

		items := s.Items()
		want := []int{27, 100}
		sort.Ints(items)

		if !reflect.DeepEqual(want, items) {
			t.Errorf("slice mismatch: got %#v, wanted %#v", items, want)
		}

		s.Remove(100)
		if s.Contains(100) {
			t.Error("set contained deleted item after remove")
		}
	})
}

func TestWithSize(t *testing.T) {
	s := set.New[string](set.WithSize(10))
	s.Add("hello")
	s.Add("there")
	s.Add("general")
	s.Add("kenobi")

	if s.Length() != 4 {
		t.Errorf("wrong length: got %d, wanted %d", s.Length(), 4)
	}

	s2 := set.New[string](set.WithSize(-10)) // Shouldn't panic
	s2.Add("hello")
	s2.Add("there")
	s2.Add("general")
	s2.Add("kenobi")

	if s2.Length() != 4 {
		t.Errorf("wrong length: got %d, wanted %d", s2.Length(), 4)
	}
}

func TestUnion(t *testing.T) {
	this := set.New[string]()
	that := set.New[string]()

	this.Add("hello")
	this.Add("there")
	this.Add("general")
	this.Add("kenobi")

	that.Add("hello")
	that.Add("to")
	that.Add("you")
	that.Add("too")

	union := set.Union(this, that).Items()
	sort.Strings(union)

	want := []string{"general", "hello", "kenobi", "there", "to", "too", "you"}

	if !reflect.DeepEqual(union, want) {
		t.Errorf("wrong union: got %#v, wanted %#v", union, want)
	}
}

func TestIntersection(t *testing.T) {
	this := set.New[string]()
	that := set.New[string]()

	this.Add("hello")
	this.Add("there")
	this.Add("general")
	this.Add("kenobi")

	that.Add("hello")
	that.Add("to")
	that.Add("you")
	that.Add("too")

	intersection := set.Intersection(this, that).Items()
	sort.Strings(intersection)

	want := []string{"hello"}

	if !reflect.DeepEqual(intersection, want) {
		t.Errorf("wrong intersection: got %#v, wanted %#v", intersection, want)
	}
}

func TestDifference(t *testing.T) {
	this := set.New[string]()
	that := set.New[string]()

	this.Add("hello")
	this.Add("there")
	this.Add("general")
	this.Add("kenobi")

	that.Add("hello")
	that.Add("to")
	that.Add("you")
	that.Add("too")

	difference := set.Difference(this, that).Items()
	sort.Strings(difference)

	want := []string{"general", "kenobi", "there"}

	if !reflect.DeepEqual(difference, want) {
		t.Errorf("wrong difference: got %#v, wanted %#v", difference, want)
	}
}

func TestString(t *testing.T) {
	s := set.New[string]()

	s.Add("hello")
	s.Add("there")
	s.Add("general")
	s.Add("kenobi")

	got := s.String()

	// A set is an unordered collection and it's pointless to sort just for string representation
	// so we just check the existence of the items
	targets := []string{"hello", "there", "general", "kenobi"}
	for _, target := range targets {
		if !strings.Contains(got, target) {
			t.Errorf("string representation does not contain %q", target)
		}
	}
}

func ExampleUnion() {
	this := set.New[string]()
	that := set.New[string]()

	this.Add("hello")
	this.Add("there")

	that.Add("general")
	that.Add("kenobi")
	that.Add("says")
	that.Add("hello")

	// Get the items in a slice of strings
	union := set.Union(this, that).Items()
	sort.Strings(union) // A set is an unordered collection

	fmt.Println(union)
	// Output: [general hello kenobi says there]
}

func ExampleIntersection() {
	this := set.New[string]()
	that := set.New[string]()

	this.Add("hello")
	this.Add("there")

	that.Add("general")
	that.Add("kenobi")
	that.Add("says")
	that.Add("hello")

	// Get the items in a slice of strings
	intersection := set.Intersection(this, that).Items()
	sort.Strings(intersection) // A set is an unordered collection

	fmt.Println(intersection)
	// Output: [hello]
}

func ExampleDifference() {
	this := set.New[string]()
	that := set.New[string]()

	this.Add("hello")
	this.Add("there")

	that.Add("general")
	that.Add("kenobi")
	that.Add("says")
	that.Add("hello")

	// Get the items in a slice of strings
	difference := set.Difference(this, that).Items()
	sort.Strings(difference) // A set is an unordered collection

	fmt.Println(difference)
	// Output: [there]
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := set.New[int]()
		for j := 0; j < 1000; j++ {
			s.Add(j)
		}
	}
}

func BenchmarkIntersection(b *testing.B) {
	s1 := set.New[int]()
	s2 := set.New[int]()

	for i := 0; i < 1000; i++ {
		s1.Add(i)
		s2.Add(i + 500)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		set.Intersection(s1, s2)
	}
}

func BenchmarkUnion(b *testing.B) {
	s1 := set.New[int]()
	s2 := set.New[int]()

	for i := 0; i < 1000; i++ {
		s1.Add(i)
		s2.Add(i + 500)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		set.Union(s1, s2)
	}
}

func BenchmarkDifference(b *testing.B) {
	s1 := set.New[int]()
	s2 := set.New[int]()

	for i := 0; i < 1000; i++ {
		s1.Add(i)
		s2.Add(i + 500)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		set.Difference(s1, s2)
	}
}
