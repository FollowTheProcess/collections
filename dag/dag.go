package dag

import (
	"fmt"

	"github.com/FollowTheProcess/collections/queue"
	"github.com/FollowTheProcess/collections/set"
)

// vertex is a single node in the graph, and holds the underlying data
// we want to represent in the graph.
type vertex[T any] struct {
	parents  *set.Set[*vertex[T]] // The direct parents of this vertex
	children *set.Set[*vertex[T]] // The direct children of this vertex
	item     T                    // The actual data
}

// newVertex creates and returns a new vertex containing item.
func newVertex[T any](item T) *vertex[T] {
	return &vertex[T]{
		parents:  set.New[*vertex[T]](),
		children: set.New[*vertex[T]](),
		item:     item,
	}
}

// inDegree returns the number of inbound edges to the vertex.
func (v vertex[T]) inDegree() int {
	return v.parents.Size()
}

// Graph is a generic directed acyclic graph, generic over 'K' which is a comparable
// type to be used as the unique ID for each vertex and 'T' which is the data
// you wish to store in each vertex of the graph.
type Graph[K comparable, T any] struct {
	vertices map[K]*vertex[T] // The map of id -> vertex
	edges    int              // The current number of edges in the graph
}

// New creates and returns a new [Graph].
func New[K comparable, T any]() *Graph[K, T] {
	return &Graph[K, T]{
		vertices: make(map[K]*vertex[T]),
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
// If the vertex already exists, an error will be returned.
func (g *Graph[K, T]) AddVertex(id K, item T) error {
	if _, exists := g.vertices[id]; exists {
		return fmt.Errorf("vertex with id '%v' already exists", id)
	}

	g.vertices[id] = newVertex(item)
	return nil
}

// AddEdge creates a connection from the vertex with id 'from' and one
// with id 'to'.
func (g *Graph[K, T]) AddEdge(from, to K) error {
	parent, exists := g.vertices[from]
	if !exists {
		return fmt.Errorf("parent vertex with id '%v' not in graph", from)
	}

	child, exists := g.vertices[to]
	if !exists {
		return fmt.Errorf("child vertex with id '%v' not in graph", to)
	}

	// Create the connection
	parent.children.Insert(child)
	child.parents.Insert(parent)
	g.edges++

	return nil
}

// Sort returns the topological sort of the graph, returning the underlying items
// in the correct order.
func (g *Graph[K, T]) Sort() ([]T, error) {
	// Note: this is kahns algorithm
	// https://en.wikipedia.org/wiki/Topological_sorting
	zeroInDegreeQueue := queue.New[*vertex[T]]()
	result := make([]T, 0, len(g.vertices))

	for _, vertex := range g.vertices {
		// Put all vertices with a 0 in-degree into the queue
		if vertex.inDegree() == 0 {
			zeroInDegreeQueue.Push(vertex)
		}
	}

	// If there is not at least 1 vertex with 0 in-degree, then it's not
	// a DAG and cannot be sorted
	if zeroInDegreeQueue.Empty() {
		return nil, fmt.Errorf("graph contains a cycle and cannot be sorted")
	}

	// While queue is not empty
	for !zeroInDegreeQueue.Empty() {
		vert, _ := zeroInDegreeQueue.Pop() //nolint: errcheck // Only error is pop from empty queue

		// Add its item to the result slice
		result = append(result, vert.item)

		// For each child, remove 'vert' as a parent and check if it
		// now has an in-degree of 0
		for child := range vert.children.Items() {
			child.parents.Remove(vert)

			// If it now has an in-degree of 0, add it to the queue
			if child.inDegree() == 0 {
				zeroInDegreeQueue.Push(child)
			}
		}
	}

	return result, nil
}
