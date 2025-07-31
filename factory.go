package gofab

// Define creates a new factory for type T.
func Define[T any](defaults ...Builder[T]) *Factory[T] {
	return &Factory[T]{defaults: defaults}
}

// Build creates an instance using automatic generation and factory defaults.
func (f *Factory[T]) Build(builders ...Builder[T]) T {
	var result T

	autoPopulateFromTags(&result)

	for _, builder := range f.defaults {
		builder(&result)
	}

	for _, builder := range builders {
		builder(&result)
	}

	return result
}

// BuildList creates multiple instances using the factory.
func (f *Factory[T]) BuildList(count int, builders ...Builder[T]) []T {
	result := make([]T, count)
	for i := range result {
		result[i] = f.Build(builders...)
	}

	return result
}
