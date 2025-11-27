package gofab

import (
	"sync"
	"sync/atomic"
)

type sequenceCounter struct {
	value int64
}

func (s *sequenceCounter) next() int64 {
	return atomic.AddInt64(&s.value, 1)
}

func (s *sequenceCounter) reset(value int64) {
	atomic.StoreInt64(&s.value, value)
}

// sequenceRegistry manages all sequence counters for thread-safe access.
type sequenceRegistry struct {
	mu       sync.RWMutex
	counters map[string]*sequenceCounter
}

var globalRegistry = &sequenceRegistry{
	counters: make(map[string]*sequenceCounter),
}

// getOrCreate returns an existing counter or creates a new one for the given key.
func (r *sequenceRegistry) getOrCreate(key string, startValue int64) *sequenceCounter {
	r.mu.RLock()
	counter, exists := r.counters[key]
	r.mu.RUnlock()

	if exists {
		return counter
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock
	if counter, exists = r.counters[key]; exists {
		return counter
	}

	counter = &sequenceCounter{value: startValue - 1}
	r.counters[key] = counter

	return counter
}

// reset resets a specific counter to the given value.
func (r *sequenceRegistry) reset(key string, value int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if counter, exists := r.counters[key]; exists {
		counter.reset(value - 1)
	}
}

// resetAll resets all counters to their initial state.
func (r *sequenceRegistry) resetAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counters = make(map[string]*sequenceCounter)
}

// ResetSequence resets a specific sequence counter by key.
// The key format is typically "TypeName.FieldName" for tag-based sequences.
func ResetSequence(key string) {
	globalRegistry.reset(key, 1)
}

// ResetSequenceWithValue resets a specific sequence counter to start from the given value.
func ResetSequenceWithValue(key string, startValue int64) {
	globalRegistry.reset(key, startValue)
}

// ResetAllSequences resets all sequence counters.
// Useful for test cleanup between test cases.
func ResetAllSequences() {
	globalRegistry.resetAll()
}

// Sequence creates a sequential value builder.
func Sequence[T any, V any](setter func(*T, V), generator func(int64) V) Builder[T] {
	counter := &sequenceCounter{value: -1}

	return func(obj *T) {
		setter(obj, generator(counter.next()))
	}
}

// SequenceWithStart creates a sequential value builder starting from a custom value.
func SequenceWithStart[T any, V any](startValue int64, setter func(*T, V), generator func(int64) V) Builder[T] {
	counter := &sequenceCounter{value: startValue - 1}

	return func(obj *T) {
		setter(obj, generator(counter.next()))
	}
}
