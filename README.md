# Collections

[![License](https://img.shields.io/github/license/FollowTheProcess/collections)](https://github.com/FollowTheProcess/collections)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/collections)](https://goreportcard.com/report/github.com/FollowTheProcess/collections)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/collections?logo=github&sort=semver)](https://github.com/FollowTheProcess/collections)
[![CI](https://github.com/FollowTheProcess/collections/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/collections/actions?query=workflow%3ACI)
[![Go Reference](https://pkg.go.dev/badge/github.com/FollowTheProcess/collections.svg)](https://pkg.go.dev/github.com/FollowTheProcess/collections)
[![codecov](https://codecov.io/gh/FollowTheProcess/collections/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/collections)

Collection of generic data structures in Go 📦

> [!TIP]
> Most collections support the Go 1.23 functional iterator pattern

## Project Description

Small, useful, zero dependency implementations of generic collection data structures in Go:

* **Set:** Offers fast membership checking as well as difference, intersection etc.
* **Stack:** Simple LIFO stack
* **Queue:** Simple FIFO queue
* **List:** A doubly-linked list
* **OrderedMap:** A map that remembers the order in which keys were inserted
* **DAG:** A generic directed acyclic graph
* **Counter:** A convenient construct for counting occurrences of things (similar to Python's [collections.Counter])
* **Chain:** A chain of maps, lookups first look in one map, then the next, then the next, returning the first result found (similar to Python's [collections.ChainMap])

## Installation

```shell
go get github.com/FollowTheProcess/collections@latest
```

## Quickstart

### Set

A set is an unordered collection of unique items offering fast lookup and membership checking.

```go
// Initialise a new set with a concrete type
s := set.New[string]()

// Insert items to the set
s.Insert("hello")
s.Insert("sets")
s.Insert("in")
s.Insert("go")

// All the methods you'd expect
s.Contains("hello") // true
s.Size() // 4

// Remove an item,
s.Remove("go")
s.Size() // 3

// Rich comparison with other sets
other := set.New[string]()
other.Insert("hello")
other.Insert("more")

// Union: combine both sets into one
set.Union(s, other) // ["hello", "in", "sets", "more", "go"]

// Intersection: all items present in both sets
set.Intersection(s, other) // ["hello"]

// Difference: items in s but not in other
set.Difference(s, other) // ["sets", "in", "go"]
```

### Stack

A stack is a LIFO data structure useful in a variety of situations.

```go
// Initialise a new stack with a concrete type
s := stack.New[string]()

// Push items onto the stack
s.Push("hello")
s.Push("stacks")
s.Push("in")
s.Push("go")

s.Size() // 4

// Pop items off the stack in LIFO order
item, _ := s.Pop()
fmt.Println(item) // "go"

item, _ = s.Pop()
fmt.Println(item) // "in"

item, _ = s.Pop()
fmt.Println(item) // "stacks"

item, _ = s.Pop()
fmt.Println(item) // "hello"

// Popping from an empty stack returns an error
_, err := s.Pop()
fmt.Println(err) // "pop from empty stack"
```

### Queue

A queue is a FIFO data structure useful in a variety of situations.

```go
// Initialise a new queue with a concrete type
q := queue.New[string]()

// Push items into the back of the queue
q.Push("hello")
q.Push("queues")
q.Push("in")
q.Push("go")

q.Size() // 4

// Pop items off the front of the queue
item, _ := q.Pop()
fmt.Println(item) // "hello"

item, _ = q.Pop()
fmt.Println(item) // "queues"

item, _ = q.Pop()
fmt.Println(item) // "in"

item, _ = q.Pop()
fmt.Println(item) // "go"

// Popping from an empty queue returns an error
_, err := q.Pop()
fmt.Println(err) // "pop from empty queue"
```

### List

A doubly linked list is a data structure where nodes wrap the data and point to their next and previous nodes. It offers cheap insertion and removal.

```go
// Initialise a new list holding a string as the data
l := list.New[string]()

// Bolt things on the end
l.Append("one")
l.Append("two")

// Push things at the start
l.Prepend("before")

last, err := l.Last()
// Handle err.. means empty list
fmt.Println(last.Item()) // <- Last is a Node, so you must call .Item() to get underlying data
```

### Ordered Map

An ordered map is like the Go standard map, except it remembers the order in which items were inserted.

```go
m := orderedmap.New[int, string]()

// Insert key value pairs
m.Insert(1, "one")
m.Insert(2, "two")
m.Insert(3, "three")

one, ok := m.Get(1) // Fetch them back out, same API as go map
if !ok {
    fmt.Println("1 was missing!")
}

two, existed := m.Remove(2) // Removal returns what was in the map

oldestKey, oldestVal, ok := m.Oldest() // Get the first inserted thing (there's also a Newest())
```

### DAG

A DAG ([Directed Acyclic Graph]) is an ordered graph ideal for task orchestration and dependency management.

```go
// Create a new DAG storing integers as the vertex data type, and a unique ID
// for each vertex of a string (this must uniquely identify a single vertex in the graph)
graph := dag.New[string, int]()

_ = graph.AddVertex("one", 1) // Add a vertex named "one" storing the integer 1
_ = graph.AddVertex("two", 2) // Add a vertex named "two" storing the integer 2

// Connect the two vertices, "two" depends on "one"
_ = graph.AddEdge("one", "two")

// Topologically sort the graph
order, err := graph.Sort()
```

### Counter

A convenient construct to count occurrences of comparable items.

```go
counts := counter.New[string]()

// Count fruits
counts.Add("apple")
counts.Add("apple")
counts.Add("apple")
counts.Add("orange")
counts.Add("orange")
counts.Add("raspberry")

// How many apples?
counts.Count("apple") // 3

// How many fruits in total?
counts.Sum() // 6

// What's the most common fruit
counts.MostCommon(1) // [{Item: "apple", Count: 3}]
```

### Chain

A chain of maps who's values are looked up in order. If the value isn't in the first map, it falls through to the second etc. Fresh inserts always go to the first map, updates update the value in whichever map it's first found in.

> [!TIP]
> A `Chain` is very useful for structured lookups of different priorities e.g. taking configuration from command line args which have precedence over env vars, and then falling back to default values

```go
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

    // 1 is in the first map in the chain
    chain.Get(1) // -> "one in first map", true 

    // To get 4, we look through every map until it's found
    chain.Get(4) // -> "four in third map", true

    // 5 isn't in any map
    chain.Get(5) // -> "", false
```

[Directed Acyclic Graph]: https://en.wikipedia.org/wiki/Directed_acyclic_graph
[collections.Counter]: https://docs.python.org/3/library/collections.html#collections.Counter
[collections.ChainMap]: https://docs.python.org/3/library/collections.html#collections.ChainMap
