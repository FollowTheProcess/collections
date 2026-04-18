package dag

import "errors"

var (
	// ErrVertexExists is returned when adding a vertex whose id is already in use.
	ErrVertexExists = errors.New("vertex already exists")

	// ErrVertexNotFound is returned when referencing a vertex id that is not in the graph.
	ErrVertexNotFound = errors.New("vertex not in graph")

	// ErrEdgeExists is returned when adding an edge that already exists.
	ErrEdgeExists = errors.New("edge already exists")

	// ErrEdgeNotFound is returned when referencing an edge that does not exist.
	ErrEdgeNotFound = errors.New("edge does not exist")

	// ErrSelfLoop is returned when attempting to connect a vertex to itself.
	ErrSelfLoop = errors.New("self-loop not permitted in a DAG")

	// ErrCycle is returned when adding an edge would introduce a cycle.
	ErrCycle = errors.New("would create a cycle")
)
