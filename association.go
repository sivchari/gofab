package gofab

// AssociationFactory creates instances with associated objects.
// R is the type that holds all associations.
type AssociationFactory[T any, R any] struct {
	defaults     []Builder[T]
	traits       map[string][]Builder[T]
	associations func() R
	linker       func(*T, R)
}

// DefineWithAssociations creates a new factory with association support.
// R is a struct type that holds all related objects.
func DefineWithAssociations[T any, R any](
	associations func() R,
	linker func(*T, R),
	defaults ...Builder[T],
) *AssociationFactory[T, R] {
	return &AssociationFactory[T, R]{
		defaults:     defaults,
		traits:       make(map[string][]Builder[T]),
		associations: associations,
		linker:       linker,
	}
}

// Trait defines a named set of attributes that can be applied when building.
func (f *AssociationFactory[T, R]) Trait(name string, builders ...Builder[T]) *AssociationFactory[T, R] {
	f.traits[name] = builders

	return f
}

// WithTrait returns builders for the specified trait.
func (f *AssociationFactory[T, R]) WithTrait(name string) []Builder[T] {
	return f.traits[name]
}

// WithTraits returns builders for multiple traits combined.
func (f *AssociationFactory[T, R]) WithTraits(names ...string) []Builder[T] {
	var builders []Builder[T]
	for _, name := range names {
		builders = append(builders, f.traits[name]...)
	}

	return builders
}

// Build creates an instance with automatically created associations.
func (f *AssociationFactory[T, R]) Build(builders ...Builder[T]) T {
	var result T

	autoPopulateFromTags(&result)

	for _, builder := range f.defaults {
		builder(&result)
	}

	// Create and link associations
	if f.associations != nil && f.linker != nil {
		assoc := f.associations()
		f.linker(&result, assoc)
	}

	for _, builder := range builders {
		builder(&result)
	}

	return result
}

// BuildWithAssociations creates an instance and returns both the instance and its associations.
func (f *AssociationFactory[T, R]) BuildWithAssociations(builders ...Builder[T]) (T, R) {
	var (
		result T
		assoc  R
	)

	autoPopulateFromTags(&result)

	for _, builder := range f.defaults {
		builder(&result)
	}

	// Create and link associations
	if f.associations != nil {
		assoc = f.associations()

		if f.linker != nil {
			f.linker(&result, assoc)
		}
	}

	for _, builder := range builders {
		builder(&result)
	}

	return result, assoc
}

// BuildList creates multiple instances with automatically created associations.
func (f *AssociationFactory[T, R]) BuildList(count int, builders ...Builder[T]) []T {
	result := make([]T, count)
	for i := range result {
		result[i] = f.Build(builders...)
	}

	return result
}

// BuildListWithAssociations creates multiple instances and returns both instances and associations.
func (f *AssociationFactory[T, R]) BuildListWithAssociations(count int, builders ...Builder[T]) ([]T, []R) {
	results := make([]T, count)
	assocs := make([]R, count)

	for i := range results {
		results[i], assocs[i] = f.BuildWithAssociations(builders...)
	}

	return results, assocs
}

// WithCustomAssociations allows overriding the default associations for a single build.
func (f *AssociationFactory[T, R]) WithCustomAssociations(customAssoc R, builders ...Builder[T]) T {
	var result T

	autoPopulateFromTags(&result)

	for _, builder := range f.defaults {
		builder(&result)
	}

	// Link custom associations
	if f.linker != nil {
		f.linker(&result, customAssoc)
	}

	for _, builder := range builders {
		builder(&result)
	}

	return result
}
