package set_test

import (
	"fmt"
	"reflect"
	"sort"
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
