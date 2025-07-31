// Package gofab provides a simple factory pattern implementation in Go.
package gofab

// Builder is a function that configures an instance of type T.
type Builder[T any] func(*T)

// Factory holds factory definition for a specific type.
type Factory[T any] struct {
	defaults []Builder[T]
}
