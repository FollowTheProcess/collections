// Package collections implements a variety of generic data structures e.g. hashsets, stacks and queues.
//
// There are 2 general rules for the data structures in this package:
//   - They should be considered not thread safe unless explicitly stated otherwise.
//   - They should each be instantiated using the `New“ constructor function, and not directly. Doing so will likely result in a nil pointer dereference.
package collections
