# Collections

[![License](https://img.shields.io/github/license/FollowTheProcess/collections)](https://github.com/FollowTheProcess/collections)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/collections)](https://goreportcard.com/report/github.com/FollowTheProcess/collections)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/collections?logo=github&sort=semver)](https://github.com/FollowTheProcess/collections)
[![CI](https://github.com/FollowTheProcess/collections/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/collections/actions?query=workflow%3ACI)

Collection of generic data structures in Go 📦

* Free software: MIT License

## Project Description

Small, useful, zero dependency implementations of generic collection data structures in Go:

* Hash sets
* Stacks
* Queues

I wrote these primarily for use in some of my other projects but they are useful enough to be applicable in most scenarios (thanks to Go 1.18 and generics!).

## Installation

```shell
go get github.com/FollowTheProcess/collections
```

## Quickstart

### Set

A set is an unordered collection of unique items offering fast lookup and membership checking.

```go
// Initialise a new set with a concrete type
s := set.New[string]()

// Add items to the set
s.Add("hello")
s.Add("sets")
s.Add("in")
s.Add("go")

// All the methods you'd expect
s.Contains("hello") // true
s.Length() // 4

// Remove an item,
s.Remove("go")
s.Length() // 3

// Rich comparison with other sets
other := set.New[string]()
other.Add("hello")
other.Add("more")

// Union: combine both sets into one
set.Union(s, other) // ["hello", "in", "sets", "more"]

// Intersection: all items present in both sets
set.Intersection(s, other) // ["hello"]

// Difference: items in s but not in other
set.Difference(s, other) // ["sets", "in"]
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

s.Length() // 4

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

q.Length() // 4

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

### Credits

This package was created with [cookiecutter] and the [FollowTheProcess/go_cookie] project template.

[cookiecutter]: https://github.com/cookiecutter/cookiecutter
[FollowTheProcess/go_cookie]: https://github.com/FollowTheProcess/go_cookie
