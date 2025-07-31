package gofab

// Build creates a single instance with automatic field population based on struct tags
func Build[T any](customizers ...func(*T)) T {
	var result T
	
	// Auto-populate using struct tags
	autoPopulateFromTags(&result)
	
	// Apply any custom overrides
	for _, customizer := range customizers {
		customizer(&result)
	}
	
	return result
}

// BuildList creates multiple instances with automatic field population
func BuildList[T any](count int, customizers ...func(*T)) []T {
	result := make([]T, count)
	for i := range result {
		result[i] = Build[T](customizers...)
	}
	return result
}