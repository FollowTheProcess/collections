package stack

// options is a struct that holds the options for a stack.
type options struct {
	capacity int // The initial capacity of the stack.
}

// Option is a function that configures a stack.
type Option func(*options)

// WithCapacity sets the initial capacity of the stack.
//
// If the capacity is negative, the default of 0 is used.
//
//	s := stack.New[string](stack.WithCapacity(10))
//	s.Cap() // 10
func WithCapacity(capacity int) Option {
	return func(o *options) {
		// Capacity must be a positive number
		if capacity > 0 {
			o.capacity = capacity
		}
	}
}
