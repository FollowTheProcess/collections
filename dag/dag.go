// Package dag implements a [Directed Acyclic Graph]
//
// A [Graph] is not safe for concurrent use. The caller is responsible for
// synchronising concurrent access.
//
// [Directed Acyclic Graph]: https://en.wikipedia.org/wiki/Directed_acyclic_graph
package dag // import "go.followtheprocess.codes/collections/dag"

import (
	"fmt"
	"iter"
	"maps"
)

// Graph is a generic directed acyclic graph, generic over 'K' which is a comparable
// type to be used as the unique ID for each vertex and 'T' which is the data
// you wish to store in each vertex of the graph.
//
// The ID must be unique within a [Graph].
type Graph[K comparable, T any] struct {
	vertices map[K]T              // The map of id -> item
	children map[K]map[K]struct{} // id -> set of child ids
	parents  map[K]map[K]struct{} // id -> set of parent ids
	edges    int                  // Edge count
}

// New creates and returns a new [Graph].
//
// It is generic over 'K' which is a comparable type to be used as the unique ID
// for each vertex, and 'T' which is the data you wish to store in each vertex of the graph.
//
// So for a graph storing integers with a unique ID that is a string, the signature would be:
//
//	graph := dag.New[string, int]()
//
// The ID must be unique within a [Graph].
func New[K comparable, T any]() *Graph[K, T] {
	return &Graph[K, T]{
		vertices: make(map[K]T),
		children: make(map[K]map[K]struct{}),
		parents:  make(map[K]map[K]struct{}),
	}
}

// WithCapacity creates and returns a new [Graph] with the given capacity.
//
// This can be a useful performance improvement if the expected maximum number of elements
// the graph will hold is known ahead of time as it eliminates the need for reallocation.
func WithCapacity[K comparable, T any](capacity int) *Graph[K, T] {
	return &Graph[K, T]{
		vertices: make(map[K]T, capacity),
		children: make(map[K]map[K]struct{}, capacity),
		parents:  make(map[K]map[K]struct{}, capacity),
	}
}

// Order returns the number of vertices in the graph.
func (g *Graph[K, T]) Order() int {
	return len(g.vertices)
}

// Size returns the number of edges in the graph.
func (g *Graph[K, T]) Size() int {
	return g.edges
}

// AddVertex adds an item to the graph as a vertex (or node) in the graph.
//
// If the vertex already exists, [ErrVertexExists] will be returned.
//
//	graph := dag.New[string, int]()
//	graph.AddVertex("one", 1)
//
// The ID must uniquely identify a single vertex in the [Graph].
func (g *Graph[K, T]) AddVertex(id K, item T) error {
	if _, exists := g.vertices[id]; exists {
		return fmt.Errorf("vertex with id '%v': %w", id, ErrVertexExists)
	}

	g.vertices[id] = item
	g.children[id] = make(map[K]struct{})
	g.parents[id] = make(map[K]struct{})

	return nil
}

// RemoveVertex removes the vertex with the given id and all edges connected to it.
//
// If the vertex does not exist, [ErrVertexNotFound] will be returned.
func (g *Graph[K, T]) RemoveVertex(id K) error {
	if !g.ContainsVertex(id) {
		return fmt.Errorf("vertex with id '%v': %w", id, ErrVertexNotFound)
	}

	for child := range g.children[id] {
		delete(g.parents[child], id)
		g.edges--
	}

	for parent := range g.parents[id] {
		delete(g.children[parent], id)
		g.edges--
	}

	delete(g.vertices, id)
	delete(g.children, id)
	delete(g.parents, id)

	return nil
}

// GetVertex returns the item stored in a vertex.
//
// If the vertex does not exist, [ErrVertexNotFound] will be returned.
func (g *Graph[K, T]) GetVertex(id K) (T, error) {
	var zero T

	item, exists := g.vertices[id]
	if !exists {
		return zero, fmt.Errorf("vertex with id '%v': %w", id, ErrVertexNotFound)
	}

	return item, nil
}

// ContainsVertex reports whether a vertex with the given id is present in the graph.
func (g *Graph[K, T]) ContainsVertex(id K) bool {
	_, exists := g.vertices[id]
	return exists
}

// Vertices returns an iterator over all vertices in the graph as (id, item) pairs.
//
// The iteration order is non-deterministic.
func (g *Graph[K, T]) Vertices() iter.Seq2[K, T] {
	return maps.All(g.vertices)
}

// AddEdge creates a directed edge from the vertex with id 'from' to the vertex with id 'to'.
//
// For the canonical use of a DAG as a dependency graph, where task "two" depends on task "one":
//
//	AddEdge("one", "two")
//
// Self-loops return [ErrSelfLoop], duplicate edges return [ErrEdgeExists], and edges
// that would introduce a cycle return [ErrCycle].
func (g *Graph[K, T]) AddEdge(from, to K) error {
	if from == to {
		return fmt.Errorf("'%v' -> '%v': %w", from, to, ErrSelfLoop)
	}

	if !g.ContainsVertex(from) {
		return fmt.Errorf("parent vertex with id '%v': %w", from, ErrVertexNotFound)
	}

	if !g.ContainsVertex(to) {
		return fmt.Errorf("child vertex with id '%v': %w", to, ErrVertexNotFound)
	}

	if _, exists := g.children[from][to]; exists {
		return fmt.Errorf("'%v' -> '%v': %w", from, to, ErrEdgeExists)
	}

	if g.canReach(to, from) {
		return fmt.Errorf("adding '%v' -> '%v': %w", from, to, ErrCycle)
	}

	g.children[from][to] = struct{}{}
	g.parents[to][from] = struct{}{}
	g.edges++

	return nil
}

// RemoveEdge removes the directed edge from 'from' to 'to'.
//
// If either vertex does not exist, [ErrVertexNotFound] will be returned.
// If the edge does not exist, [ErrEdgeNotFound] will be returned.
func (g *Graph[K, T]) RemoveEdge(from, to K) error {
	if !g.ContainsVertex(from) {
		return fmt.Errorf("vertex with id '%v': %w", from, ErrVertexNotFound)
	}

	if !g.ContainsVertex(to) {
		return fmt.Errorf("vertex with id '%v': %w", to, ErrVertexNotFound)
	}

	if !g.HasEdge(from, to) {
		return fmt.Errorf("'%v' -> '%v': %w", from, to, ErrEdgeNotFound)
	}

	delete(g.children[from], to)
	delete(g.parents[to], from)
	g.edges--

	return nil
}

// HasEdge reports whether a directed edge exists from 'from' to 'to'.
func (g *Graph[K, T]) HasEdge(from, to K) bool {
	_, exists := g.children[from][to]
	return exists
}

// Children returns an iterator over the direct children (immediate dependents) of the vertex with the given id.
//
// If the vertex does not exist, [ErrVertexNotFound] will be returned.
func (g *Graph[K, T]) Children(id K) (iter.Seq[K], error) {
	if !g.ContainsVertex(id) {
		return nil, fmt.Errorf("vertex with id '%v': %w", id, ErrVertexNotFound)
	}

	return maps.Keys(g.children[id]), nil
}

// Parents returns an iterator over the direct parents (immediate dependencies) of the vertex with the given id.
//
// If the vertex does not exist, [ErrVertexNotFound] will be returned.
func (g *Graph[K, T]) Parents(id K) (iter.Seq[K], error) {
	if !g.ContainsVertex(id) {
		return nil, fmt.Errorf("vertex with id '%v': %w", id, ErrVertexNotFound)
	}

	return maps.Keys(g.parents[id]), nil
}

// Descendants returns an iterator over all transitive descendants of the vertex with the given id,
// i.e. all vertices reachable by following directed edges forward from id.
//
// Each vertex is yielded at most once. If the vertex does not exist, [ErrVertexNotFound] will be returned.
func (g *Graph[K, T]) Descendants(id K) (iter.Seq[K], error) {
	if !g.ContainsVertex(id) {
		return nil, fmt.Errorf("vertex with id '%v': %w", id, ErrVertexNotFound)
	}

	return func(yield func(K) bool) {
		visited := make(map[K]struct{}, len(g.vertices))
		stack := []K{id}

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if _, seen := visited[current]; seen {
				continue
			}
			visited[current] = struct{}{}

			if current != id {
				if !yield(current) {
					return
				}
			}

			for child := range g.children[current] {
				stack = append(stack, child)
			}
		}
	}, nil
}

// Ancestors returns an iterator over all transitive ancestors of the vertex with the given id,
// i.e. all vertices from which id is reachable by following directed edges forward.
//
// Each vertex is yielded at most once. If the vertex does not exist, [ErrVertexNotFound] will be returned.
func (g *Graph[K, T]) Ancestors(id K) (iter.Seq[K], error) {
	if !g.ContainsVertex(id) {
		return nil, fmt.Errorf("vertex with id '%v': %w", id, ErrVertexNotFound)
	}

	return func(yield func(K) bool) {
		visited := make(map[K]struct{}, len(g.vertices))
		stack := []K{id}

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if _, seen := visited[current]; seen {
				continue
			}
			visited[current] = struct{}{}

			if current != id {
				if !yield(current) {
					return
				}
			}

			for parent := range g.parents[current] {
				stack = append(stack, parent)
			}
		}
	}, nil
}

// Sort returns the topological sort of the graph, returning the underlying items
// in a valid dependency order.
//
// A DAG may have multiple valid topological sorts; the one returned from this function
// is guaranteed to be valid but is not deterministic. Sort never returns an error because
// the graph is guaranteed to be acyclic by construction — any edge that would create a
// cycle is rejected by [Graph.AddEdge].
func (g *Graph[K, T]) Sort() []T {
	// Kahn's algorithm: https://en.wikipedia.org/wiki/Topological_sorting
	n := len(g.vertices)

	// Assign each vertex a contiguous integer index so inDegree lives in a plain
	// []int rather than a map. The inner loop then pays one map read (child → index)
	// instead of two map ops (read-modify-write on a map[K]int) per edge processed,
	// which matters on dense graphs where the inner loop runs O(E) = O(V²) times.
	index := make(map[K]int, n)
	keys := make([]K, 0, n)
	for id := range g.vertices {
		index[id] = len(keys)
		keys = append(keys, id)
	}

	inDegree := make([]int, n)
	for i, id := range keys {
		inDegree[i] = len(g.parents[id])
	}

	queue := make([]int, 0, n)
	for i, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, i)
		}
	}

	result := make([]T, 0, n)

	head := 0
	for head < len(queue) {
		i := queue[head]
		head++

		result = append(result, g.vertices[keys[i]])

		for child := range g.children[keys[i]] {
			ci := index[child]
			inDegree[ci]--
			if inDegree[ci] == 0 {
				queue = append(queue, ci)
			}
		}
	}

	return result
}

// canReach reports whether 'to' can reach 'from' via existing edges.
// If true, adding the edge from -> to would close a cycle.
func (g *Graph[K, T]) canReach(to, from K) bool {
	visited := make(map[K]struct{}, len(g.vertices))
	stack := []K{to}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if current == from {
			return true
		}

		if _, seen := visited[current]; seen {
			continue
		}
		visited[current] = struct{}{}

		for child := range g.children[current] {
			stack = append(stack, child)
		}
	}

	return false
}
