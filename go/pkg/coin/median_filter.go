package coin

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// MedianFilter maintains a moving window of numeric values and provides
// easy access to the current median. It mirrors the C++ template
// implementation used throughout the legacy codebase.
type MedianFilter[T constraints.Integer | constraints.Float] struct {
	size   int
	values []T
	sorted []T
}

// NewMedianFilter creates a filter with the given maximum window size and
// initial value.
func NewMedianFilter[T constraints.Integer | constraints.Float](size int, initial T) *MedianFilter[T] {
	f := &MedianFilter[T]{size: size}
	f.values = append(f.values, initial)
	f.sorted = append(f.sorted, initial)
	return f
}

// Input adds a new sample to the filter, evicting the oldest sample when
// the window is full.
func (f *MedianFilter[T]) Input(v T) {
	if len(f.values) == f.size {
		f.values = f.values[1:]
	}
	f.values = append(f.values, v)

	f.sorted = append([]T(nil), f.values...)
	sort.Slice(f.sorted, func(i, j int) bool { return f.sorted[i] < f.sorted[j] })
}

// Median returns the median value of the current samples.
func (f *MedianFilter[T]) Median() T {
	if len(f.sorted) == 0 {
		var zero T
		return zero
	}
	n := len(f.sorted)
	if n%2 == 1 {
		return f.sorted[n/2]
	}
	return (f.sorted[n/2-1] + f.sorted[n/2]) / 2
}

// Sorted returns a copy of the current sorted sample set.
func (f *MedianFilter[T]) Sorted() []T {
	out := make([]T, len(f.sorted))
	copy(out, f.sorted)
	return out
}

// Size returns the number of samples currently in the filter.
func (f *MedianFilter[T]) Size() int { return len(f.values) }
