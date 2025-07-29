// Package ptr provides utility functions for working with pointers in Go.
package ptr

// To converts a value of type T to a pointer of type *T.
func To[T any](v T) *T {
	return &v
}

// Deref dereferences a pointer of type *T and returns the value of type T.
func Deref[T any](ptr *T) T {
	var zero T
	if ptr == nil {
		return zero
	}

	return *ptr
}
