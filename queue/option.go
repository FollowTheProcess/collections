package queue

// options is a struct that holds the options for a queue.
type options struct {
	capacity int // The initial capacity of the queue.
}

// Option is a function that configures a queue.
type Option func(*options)

// WithCapacity sets the initial capacity of the queue.
//
// If the capacity is negative, the default of 0 is used which is equivalent to
// the capacity go allocates new slices.
//
//	q := queue.New[string](queue.WithCapacity(10))
//	q.Cap() // 10
func WithCapacity(capacity int) Option {
	return func(o *options) {
		// Capacity must be a positive number
		if capacity > 0 {
			o.capacity = capacity
		}
	}
}
