// Package gofab provides a simple factory pattern implementation in Go.
package gofab

// Builder is a function that configures an instance of type T.
type Builder[T any] func(*T)

// Factory holds factory definition for a specific type.
type Factory[T any] struct {
	defaults []Builder[T]
}

// Define creates a new factory for type T.
func Define[T any]() *Factory[T] {
	return &Factory[T]{}
}

// Default adds a default builder to the factory.
func (f *Factory[T]) Default(builder Builder[T]) *Factory[T] {
	f.defaults = append(f.defaults, builder)

	return f
}

// Build creates an instance using the factory defaults.
func (f *Factory[T]) Build(builders ...Builder[T]) T {
	var result T

	// Apply factory defaults
	for _, builder := range f.defaults {
		builder(&result)
	}

	// Apply provided builders (override defaults)
	for _, builder := range builders {
		builder(&result)
	}

	return result
}

// Create creates an instance of type T with optional builders
// This is a standalone function for creating instances without factory defaults.
func Create[T any](builders ...Builder[T]) T {
	var result T

	// Apply provided builders
	for _, builder := range builders {
		builder(&result)
	}

	return result
}
