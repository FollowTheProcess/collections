package set

const defaultMapSize = 10

// options is a struct that holds the options for a set.
type options struct {
	size int // The starting size of the set
}

// Option is a function that configures a set.
type Option func(*options)

// WithSize sets the starting size of the set.
//
// If the size is negative, the default of 10 is used.
//
//	s := set.New[string](set.WithSize(10))
//	s.Cap() // 10
func WithSize(size int) Option {
	return func(o *options) {
		// Size must be a positive number
		if size > 0 {
			o.size = size
		} else {
			o.size = defaultMapSize
		}
	}
}
