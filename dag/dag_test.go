package dag_test

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/collections/dag"
	"github.com/FollowTheProcess/test"
)

func TestAddVertex(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		graph := dag.New[string, string]()
		test.Equal(t, graph.Order(), 0) // Starting order must be 0
		test.Equal(t, graph.Size(), 0)  // Starting size must be 0

		err := graph.AddVertex("v1", "hello")
		test.Ok(t, err)

		test.Equal(t, graph.Order(), 1)          // Must now have 1 vertex
		test.Equal(t, graph.Size(), 0)           // Size must still be 0 -> no edges yet
		test.True(t, graph.ContainsVertex("v1")) // Must contain v1
	})

	t.Run("already exists", func(t *testing.T) {
		graph := dag.New[string, string]()
		test.Equal(t, graph.Order(), 0) // Starting size must be 0

		err := graph.AddVertex("v1", "hello")
		test.Ok(t, err)

		err = graph.AddVertex("v1", "world") // Same id
		test.Err(t, err)
		test.Equal(t, err.Error(), "vertex with id 'v1' already exists")
	})
}

func TestGetVertex(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		graph := dag.New[string, string]()

		err := graph.AddVertex("v1", "hello")
		test.Ok(t, err)

		v1, err := graph.GetVertex("v1")
		test.Ok(t, err)
		test.Equal(t, v1, "hello")
	})

	t.Run("missing", func(t *testing.T) {
		graph := dag.New[string, string]()

		err := graph.AddVertex("v1", "hello")
		test.Ok(t, err)

		_, err = graph.GetVertex("missing")
		test.Err(t, err)
		test.Equal(t, err.Error(), "vertex with id 'missing' not in graph")
	})
}

func TestAddEdge(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		graph := dag.New[string, int]()

		err := graph.AddVertex("one", 1)
		test.Ok(t, err)

		err = graph.AddVertex("two", 2)
		test.Ok(t, err)

		err = graph.AddEdge("one", "two")
		test.Ok(t, err)

		test.Equal(t, graph.Order(), 2) // Should be 2 vertices in the graph
		test.Equal(t, graph.Size(), 1)  // Should be 1 edge: "one" -> "two"
	})

	t.Run("parent missing", func(t *testing.T) {
		graph := dag.New[string, int]()

		err := graph.AddVertex("two", 2)
		test.Ok(t, err)

		err = graph.AddEdge("one", "two")
		test.Err(t, err) // parent "one" not in graph
		test.Equal(t, err.Error(), "parent vertex with id 'one' not in graph")
	})

	t.Run("child missing", func(t *testing.T) {
		graph := dag.New[string, int]()

		err := graph.AddVertex("one", 1)
		test.Ok(t, err)

		err = graph.AddEdge("one", "two")
		test.Err(t, err) // child "two" not in graph
		test.Equal(t, err.Error(), "child vertex with id 'two' not in graph")
	})
}

func TestSort(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		graph := dag.New[string, int]()

		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddVertex("three", 3))
		test.Ok(t, graph.AddVertex("four", 4))
		test.Ok(t, graph.AddVertex("five", 5))

		// two depends on one
		err := graph.AddEdge("one", "two")
		test.Ok(t, err) // AddEdge returned an error ("one", "two")

		// four depends on three
		err = graph.AddEdge("three", "four")
		test.Ok(t, err) // AddEdge returned an error ("three", "four")

		sorted, err := graph.Sort()
		test.Ok(t, err) // Sort returned an error

		// A DAG may have more than one possible topological sort
		possibilities := [][]int{
			{5, 1, 3, 2, 4},
			{1, 3, 5, 2, 4},
			{3, 5, 1, 4, 2},
		}

		test.True(t, isInPossibleSolutions(sorted, possibilities)) // DAG not sorted correctly
	})

	t.Run("not a dag", func(t *testing.T) {
		graph := dag.New[string, int]()

		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddVertex("three", 3))
		test.Ok(t, graph.AddVertex("four", 4))
		test.Ok(t, graph.AddVertex("five", 5))

		// Purposely make it not a DAG (no vertices with an in-degree of 0)
		// easiest way is just complete the cycle when connecting everything

		// two depends on one
		err := graph.AddEdge("one", "two")
		test.Ok(t, err) // AddEdge returned an error ("one", "two")

		// three depends on two
		err = graph.AddEdge("two", "three")
		test.Ok(t, err) // AddEdge returned an error ("two", "three")

		// four depends on three
		err = graph.AddEdge("three", "four")
		test.Ok(t, err) // AddEdge returned an error ("three", "four")

		// four depends on one
		err = graph.AddEdge("one", "four")
		test.Ok(t, err) // AddEdge returned an error ("one", "four")

		// five depends on four
		err = graph.AddEdge("four", "five")
		test.Ok(t, err) // AddEdge returned an error ("four", "five")

		// Complete the cycle: one also depends on five
		err = graph.AddEdge("five", "one")
		test.Ok(t, err) // AddEdge returned an error ("five", "one")

		_, err = graph.Sort()
		test.Err(t, err)
		test.Equal(t, err.Error(), "graph contains a cycle and cannot be sorted")
	})
}

func isInPossibleSolutions[T comparable](result []T, possibles [][]T) bool {
	for _, possible := range possibles {
		if slices.Equal(result, possible) {
			return true
		}
	}

	return false
}

// makeGraph makes a simple DAG with a few connections for things like benchmarks.
func makeGraph(tb testing.TB) *dag.Graph[string, int] {
	tb.Helper()
	graph := dag.WithCapacity[string, int](5)

	test.Ok(tb, graph.AddVertex("one", 1))
	test.Ok(tb, graph.AddVertex("two", 2))
	test.Ok(tb, graph.AddVertex("three", 3))
	test.Ok(tb, graph.AddVertex("four", 4))
	test.Ok(tb, graph.AddVertex("five", 5))

	// two depends on one
	err := graph.AddEdge("one", "two")
	test.Ok(tb, err) // AddEdge returned an error ("one", "two")

	// four depends on three
	err = graph.AddEdge("three", "four")
	test.Ok(tb, err) // AddEdge returned an error ("three", "four")

	return graph
}

func BenchmarkGraphSort(b *testing.B) {
	// Because the graph.Sort method alters the state of the graph (removing edges)
	// a new graph must be constructed for each run meaning this is actually quite slow to run (~1 minute)
	// but we stop and start the timer at the right places to ensure just the sorting code's performance is measured
	for range b.N {
		b.StopTimer()
		graph := makeGraph(b)
		b.StartTimer()
		_, err := graph.Sort()
		if err != nil {
			b.Fatalf("graph.Sort returned an error: %v", err)
		}
	}
}
