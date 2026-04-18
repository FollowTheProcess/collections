package dag_test

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"go.followtheprocess.codes/collections/dag"
	"go.followtheprocess.codes/test"
)

func TestAddVertex(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		graph := dag.New[string, string]()
		test.Equal(t, graph.Order(), 0)
		test.Equal(t, graph.Size(), 0)

		err := graph.AddVertex("v1", "hello")
		test.Ok(t, err)

		test.Equal(t, graph.Order(), 1)
		test.Equal(t, graph.Size(), 0)
		test.True(t, graph.ContainsVertex("v1"))
	})

	t.Run("already exists", func(t *testing.T) {
		graph := dag.New[string, string]()

		err := graph.AddVertex("v1", "hello")
		test.Ok(t, err)

		err = graph.AddVertex("v1", "world")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrVertexExists))
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
		test.True(t, errors.Is(err, dag.ErrVertexNotFound))
	})
}

func TestContainsVertex(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.True(t, graph.ContainsVertex("a"))
	})

	t.Run("absent", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Equal(t, graph.ContainsVertex("missing"), false)
	})
}

func TestAddEdge(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		graph := dag.New[string, int]()

		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddEdge("one", "two"))

		test.Equal(t, graph.Order(), 2)
		test.Equal(t, graph.Size(), 1)
	})

	t.Run("parent missing", func(t *testing.T) {
		graph := dag.New[string, int]()

		test.Ok(t, graph.AddVertex("two", 2))

		err := graph.AddEdge("one", "two")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrVertexNotFound))
	})

	t.Run("child missing", func(t *testing.T) {
		graph := dag.New[string, int]()

		test.Ok(t, graph.AddVertex("one", 1))

		err := graph.AddEdge("one", "two")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrVertexNotFound))
	})

	t.Run("self loop", func(t *testing.T) {
		graph := dag.New[string, int]()

		test.Ok(t, graph.AddVertex("one", 1))

		err := graph.AddEdge("one", "one")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrSelfLoop))
	})

	t.Run("duplicate edge", func(t *testing.T) {
		graph := dag.New[string, int]()

		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddEdge("one", "two"))

		err := graph.AddEdge("one", "two")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrEdgeExists))
	})

	t.Run("would create cycle", func(t *testing.T) {
		graph := dag.New[string, int]()

		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddVertex("three", 3))
		test.Ok(t, graph.AddEdge("one", "two"))
		test.Ok(t, graph.AddEdge("two", "three"))

		err := graph.AddEdge("three", "one")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrCycle))
	})
}

func TestHasEdge(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddEdge("one", "two"))
		test.True(t, graph.HasEdge("one", "two"))
	})

	t.Run("does not exist", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Equal(t, graph.HasEdge("one", "two"), false)
	})

	t.Run("unknown vertex returns false", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Equal(t, graph.HasEdge("unknown", "also-unknown"), false)
	})
}

func TestRemoveEdge(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddEdge("one", "two"))
		test.Equal(t, graph.Size(), 1)

		test.Ok(t, graph.RemoveEdge("one", "two"))

		test.Equal(t, graph.Size(), 0)
		test.Equal(t, graph.HasEdge("one", "two"), false)
	})

	t.Run("missing from", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("two", 2))

		err := graph.RemoveEdge("one", "two")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrVertexNotFound))
	})

	t.Run("missing to", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("one", 1))

		err := graph.RemoveEdge("one", "two")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrVertexNotFound))
	})

	t.Run("edge does not exist", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))

		err := graph.RemoveEdge("one", "two")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrEdgeNotFound))
	})

	t.Run("allows re-adding after remove", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.RemoveEdge("a", "b"))
		test.Ok(t, graph.AddEdge("a", "b"))
	})
}

func TestRemoveVertex(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddEdge("one", "two"))

		test.Ok(t, graph.RemoveVertex("one"))

		test.Equal(t, graph.Order(), 1)
		test.Equal(t, graph.Size(), 0)
		test.Equal(t, graph.ContainsVertex("one"), false)
	})

	t.Run("removes all connected edges", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.AddEdge("a", "c"))
		test.Equal(t, graph.Size(), 2)

		test.Ok(t, graph.RemoveVertex("a"))

		test.Equal(t, graph.Size(), 0)
		test.Equal(t, graph.HasEdge("a", "b"), false)
		test.Equal(t, graph.HasEdge("a", "c"), false)
	})

	t.Run("missing", func(t *testing.T) {
		graph := dag.New[string, int]()
		err := graph.RemoveVertex("missing")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrVertexNotFound))
	})

	t.Run("re-add after remove", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.RemoveVertex("one"))
		test.Ok(t, graph.AddVertex("one", 99))

		v, err := graph.GetVertex("one")
		test.Ok(t, err)
		test.Equal(t, v, 99)
	})
}

func TestVertices(t *testing.T) {
	t.Run("iterates all vertices", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddVertex("three", 3))

		var ids []string
		var vals []int
		for id, val := range graph.Vertices() {
			ids = append(ids, id)
			vals = append(vals, val)
		}
		slices.Sort(ids)
		slices.Sort(vals)

		test.True(t, slices.Equal(ids, []string{"one", "three", "two"}))
		test.True(t, slices.Equal(vals, []int{1, 2, 3}))
	})

	t.Run("empty graph", func(t *testing.T) {
		graph := dag.New[string, int]()
		count := 0
		for range graph.Vertices() {
			count++
		}
		test.Equal(t, count, 0)
	})

	t.Run("early termination", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))

		count := 0
		for range graph.Vertices() {
			count++
			break
		}
		test.Equal(t, count, 1)
	})
}

func TestChildren(t *testing.T) {
	t.Run("has children", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.AddEdge("a", "c"))

		children, err := graph.Children("a")
		test.Ok(t, err)

		got := slices.Sorted(children)
		test.True(t, slices.Equal(got, []string{"b", "c"}))
	})

	t.Run("leaf has no children", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddEdge("a", "b"))

		children, err := graph.Children("b")
		test.Ok(t, err)
		test.Equal(t, len(slices.Collect(children)), 0)
	})

	t.Run("missing vertex", func(t *testing.T) {
		graph := dag.New[string, int]()
		_, err := graph.Children("missing")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrVertexNotFound))
	})

	t.Run("early termination", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddVertex("d", 4))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.AddEdge("a", "c"))
		test.Ok(t, graph.AddEdge("a", "d"))

		children, err := graph.Children("a")
		test.Ok(t, err)

		count := 0
		for range children {
			count++
			break
		}
		test.Equal(t, count, 1)
	})
}

func TestParents(t *testing.T) {
	t.Run("has parents", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddEdge("a", "c"))
		test.Ok(t, graph.AddEdge("b", "c"))

		parents, err := graph.Parents("c")
		test.Ok(t, err)

		got := slices.Sorted(parents)
		test.True(t, slices.Equal(got, []string{"a", "b"}))
	})

	t.Run("root has no parents", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddEdge("a", "b"))

		parents, err := graph.Parents("a")
		test.Ok(t, err)
		test.Equal(t, len(slices.Collect(parents)), 0)
	})

	t.Run("missing vertex", func(t *testing.T) {
		graph := dag.New[string, int]()
		_, err := graph.Parents("missing")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrVertexNotFound))
	})
}

func TestDescendants(t *testing.T) {
	t.Run("chain", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.AddEdge("b", "c"))

		desc, err := graph.Descendants("a")
		test.Ok(t, err)

		got := slices.Sorted(desc)
		test.True(t, slices.Equal(got, []string{"b", "c"}))
	})

	t.Run("diamond — each node yielded once", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddVertex("d", 4))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.AddEdge("a", "c"))
		test.Ok(t, graph.AddEdge("b", "d"))
		test.Ok(t, graph.AddEdge("c", "d"))

		desc, err := graph.Descendants("a")
		test.Ok(t, err)

		got := slices.Sorted(desc)
		test.True(t, slices.Equal(got, []string{"b", "c", "d"}))
	})

	t.Run("leaf has no descendants", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddEdge("a", "b"))

		desc, err := graph.Descendants("b")
		test.Ok(t, err)
		test.Equal(t, len(slices.Collect(desc)), 0)
	})

	t.Run("missing vertex", func(t *testing.T) {
		graph := dag.New[string, int]()
		_, err := graph.Descendants("missing")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrVertexNotFound))
	})

	t.Run("early termination", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddVertex("d", 4))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.AddEdge("b", "c"))
		test.Ok(t, graph.AddEdge("c", "d"))

		desc, err := graph.Descendants("a")
		test.Ok(t, err)

		count := 0
		for range desc {
			count++
			break
		}
		test.Equal(t, count, 1)
	})
}

func TestAncestors(t *testing.T) {
	t.Run("chain", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.AddEdge("b", "c"))

		anc, err := graph.Ancestors("c")
		test.Ok(t, err)

		got := slices.Sorted(anc)
		test.True(t, slices.Equal(got, []string{"a", "b"}))
	})

	t.Run("diamond — each node yielded once", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddVertex("d", 4))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.AddEdge("a", "c"))
		test.Ok(t, graph.AddEdge("b", "d"))
		test.Ok(t, graph.AddEdge("c", "d"))

		anc, err := graph.Ancestors("d")
		test.Ok(t, err)

		got := slices.Sorted(anc)
		test.True(t, slices.Equal(got, []string{"a", "b", "c"}))
	})

	t.Run("root has no ancestors", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddEdge("a", "b"))

		anc, err := graph.Ancestors("a")
		test.Ok(t, err)
		test.Equal(t, len(slices.Collect(anc)), 0)
	})

	t.Run("missing vertex", func(t *testing.T) {
		graph := dag.New[string, int]()
		_, err := graph.Ancestors("missing")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrVertexNotFound))
	})

	t.Run("early termination", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddVertex("d", 4))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.AddEdge("b", "c"))
		test.Ok(t, graph.AddEdge("c", "d"))

		anc, err := graph.Ancestors("d")
		test.Ok(t, err)

		count := 0
		for range anc {
			count++
			break
		}
		test.Equal(t, count, 1)
	})
}

func TestSort(t *testing.T) {
	t.Run("empty graph", func(t *testing.T) {
		graph := dag.New[string, int]()
		sorted := graph.Sort()
		test.Equal(t, len(sorted), 0)
	})

	t.Run("single vertex", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 42))
		sorted := graph.Sort()
		test.Equal(t, len(sorted), 1)
		test.Equal(t, sorted[0], 42)
	})

	t.Run("linear chain", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.AddEdge("b", "c"))

		sorted := graph.Sort()
		test.Equal(t, len(sorted), 3)
		test.True(t, isValidSort(sorted, [][2]int{{1, 2}, {2, 3}}))
	})

	t.Run("diamond", func(t *testing.T) {
		graph := dag.New[string, int]()
		test.Ok(t, graph.AddVertex("a", 1))
		test.Ok(t, graph.AddVertex("b", 2))
		test.Ok(t, graph.AddVertex("c", 3))
		test.Ok(t, graph.AddVertex("d", 4))
		test.Ok(t, graph.AddEdge("a", "b"))
		test.Ok(t, graph.AddEdge("a", "c"))
		test.Ok(t, graph.AddEdge("b", "d"))
		test.Ok(t, graph.AddEdge("c", "d"))

		sorted := graph.Sort()
		test.Equal(t, len(sorted), 4)
		test.True(t, isValidSort(sorted, [][2]int{{1, 2}, {1, 3}, {2, 4}, {3, 4}}))
	})

	t.Run("happy path", func(t *testing.T) {
		graph := dag.WithCapacity[string, int](5)

		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddVertex("three", 3))
		test.Ok(t, graph.AddVertex("four", 4))
		test.Ok(t, graph.AddVertex("five", 5))
		test.Ok(t, graph.AddEdge("one", "two"))
		test.Ok(t, graph.AddEdge("three", "four"))

		sorted := graph.Sort()
		test.Equal(t, len(sorted), 5)
		test.True(t, isValidSort(sorted, [][2]int{{1, 2}, {3, 4}}))
	})

	t.Run("cycle rejected at AddEdge", func(t *testing.T) {
		graph := dag.New[string, int]()

		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddVertex("three", 3))
		test.Ok(t, graph.AddVertex("four", 4))
		test.Ok(t, graph.AddVertex("five", 5))
		test.Ok(t, graph.AddEdge("one", "two"))
		test.Ok(t, graph.AddEdge("two", "three"))
		test.Ok(t, graph.AddEdge("three", "four"))
		test.Ok(t, graph.AddEdge("one", "four"))
		test.Ok(t, graph.AddEdge("four", "five"))

		err := graph.AddEdge("five", "one")
		test.Err(t, err)
		test.True(t, errors.Is(err, dag.ErrCycle))
	})

	t.Run("sort is idempotent", func(t *testing.T) {
		graph := dag.WithCapacity[string, int](5)

		test.Ok(t, graph.AddVertex("one", 1))
		test.Ok(t, graph.AddVertex("two", 2))
		test.Ok(t, graph.AddVertex("three", 3))
		test.Ok(t, graph.AddVertex("four", 4))
		test.Ok(t, graph.AddVertex("five", 5))
		test.Ok(t, graph.AddEdge("one", "two"))
		test.Ok(t, graph.AddEdge("three", "four"))

		edges := [][2]int{{1, 2}, {3, 4}}

		test.True(t, isValidSort(graph.Sort(), edges))
		test.True(t, isValidSort(graph.Sort(), edges))
	})
}

// isValidSort verifies that result is a valid topological sort by checking that
// for every edge [from, to], from appears before to in result. It does not
// enumerate possible orderings, so it works correctly for any graph shape.
func isValidSort[T comparable](result []T, edges [][2]T) bool {
	pos := make(map[T]int, len(result))
	for i, v := range result {
		pos[v] = i
	}
	for _, edge := range edges {
		fromPos, fromOk := pos[edge[0]]
		toPos, toOk := pos[edge[1]]
		if !fromOk || !toOk {
			return false
		}
		if fromPos >= toPos {
			return false
		}
	}
	return true
}

// makeGraph makes a simple DAG with a few connections for benchmarks.
func makeGraph(tb testing.TB) *dag.Graph[string, int] {
	tb.Helper()

	graph := dag.WithCapacity[string, int](5)

	test.Ok(tb, graph.AddVertex("one", 1))
	test.Ok(tb, graph.AddVertex("two", 2))
	test.Ok(tb, graph.AddVertex("three", 3))
	test.Ok(tb, graph.AddVertex("four", 4))
	test.Ok(tb, graph.AddVertex("five", 5))
	test.Ok(tb, graph.AddEdge("one", "two"))
	test.Ok(tb, graph.AddEdge("three", "four"))

	return graph
}

// makeLargeGraph makes a 50-node linear chain for benchmarks that need a more
// representative graph size.
func makeLargeGraph(tb testing.TB) *dag.Graph[string, int] {
	tb.Helper()

	const n = 50

	graph := dag.WithCapacity[string, int](n)
	for i := range n {
		test.Ok(tb, graph.AddVertex(fmt.Sprintf("node%d", i), i))
	}
	for i := range n - 1 {
		test.Ok(tb, graph.AddEdge(fmt.Sprintf("node%d", i), fmt.Sprintf("node%d", i+1)))
	}

	return graph
}

// makeDenseDAG builds a complete DAG where every vertex points to all later vertices,
// producing N*(N-1)/2 edges. This is the pathological case for Sort: Kahn's algorithm
// processes every edge once in the inner loop, so edge count directly drives runtime.
// A linear chain of the same N has only N-1 edges — O(N²) vs O(N).
func makeDenseDAG(tb testing.TB, n int) *dag.Graph[int, int] {
	tb.Helper()

	graph := dag.WithCapacity[int, int](n)
	for i := range n {
		test.Ok(tb, graph.AddVertex(i, i))
	}
	for i := range n {
		for j := i + 1; j < n; j++ {
			test.Ok(tb, graph.AddEdge(i, j))
		}
	}

	return graph
}

func BenchmarkGraphSortDense(b *testing.B) {
	graph := makeDenseDAG(b, 50)
	for b.Loop() {
		graph.Sort()
	}
}

func BenchmarkGraphSort(b *testing.B) {
	graph := makeGraph(b)
	for b.Loop() {
		graph.Sort()
	}
}

func BenchmarkGraphSortLarge(b *testing.B) {
	graph := makeLargeGraph(b)
	for b.Loop() {
		graph.Sort()
	}
}

func BenchmarkDescendants(b *testing.B) {
	graph := makeLargeGraph(b)
	for b.Loop() {
		desc, err := graph.Descendants("node0")
		if err != nil {
			b.Fatalf("Descendants returned an error: %v", err)
		}
		// sink outside the loop prevents the compiler from eliminating it as dead code,
		// ensuring we actually benchmark the full traversal without allocating a result slice.
		var sink string
		for id := range desc {
			sink = id
		}
		_ = sink
	}
}
