package gofab

import (
	"sync/atomic"
)

// sequenceCounter manages sequential values.
type sequenceCounter struct {
	value int64
}

// next returns the next value in the sequence.
func (s *sequenceCounter) next() int64 {
	return atomic.AddInt64(&s.value, 1)
}

// Sequence creates a sequential value builder.
func Sequence[T any, V any](setter func(*T, V), generator func(int64) V) Builder[T] {
	counter := &sequenceCounter{value: -1}

	return func(obj *T) {
		setter(obj, generator(counter.next()))
	}
}
