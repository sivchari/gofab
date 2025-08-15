package gofab

// Define creates a new factory for type T.
func Define[T any](defaults ...Builder[T]) *Factory[T] {
	return &Factory[T]{
		defaults: defaults,
		traits:   make(map[string][]Builder[T]),
	}
}

// AfterBuild adds a callback that runs after building each instance.
func (f *Factory[T]) AfterBuild(callback Builder[T]) *Factory[T] {
	f.afterBuild = append(f.afterBuild, callback)

	return f
}

// Trait defines a named set of attributes that can be applied when building.
func (f *Factory[T]) Trait(name string, builders ...Builder[T]) *Factory[T] {
	f.traits[name] = builders

	return f
}

// WithTrait returns builders for the specified trait.
func (f *Factory[T]) WithTrait(name string) []Builder[T] {
	return f.traits[name]
}

// WithTraits returns builders for multiple traits combined.
func (f *Factory[T]) WithTraits(names ...string) []Builder[T] {
	var builders []Builder[T]
	for _, name := range names {
		builders = append(builders, f.traits[name]...)
	}

	return builders
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

	for _, callback := range f.afterBuild {
		callback(&result)
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
