package coin

import (
	"sync"
	"time"
)

// Time provides adjusted time logic similar to the C++ singleton.
// It uses a MedianFilter to smooth differences between peer clocks.
type Time struct {
	mu     sync.Mutex
	filter *MedianFilter[int64]
	offset int64
}

var globalTime = &Time{filter: NewMedianFilter[int64](200, 0)}

// Instance returns the global time singleton.
func InstanceTime() *Time { return globalTime }

// GetAdjusted returns the current time plus any offset calculated from peer samples.
func (t *Time) GetAdjusted() uint64 {
	t.mu.Lock()
	off := t.offset
	t.mu.Unlock()
	return uint64(time.Now().Unix()) + uint64(off)
}

// AddSample adds a peer timestamp to the median filter. The addr parameter is
// ignored for now; it's kept to mimic the original interface.
func (t *Time) AddSample(addr string, timestamp uint64) {
	offSample := int64(timestamp) - time.Now().Unix()

	t.mu.Lock()
	t.filter.Input(offSample)
	if t.filter.Size() >= 5 && t.filter.Size()%2 == 1 {
		median := t.filter.Median()
		sorted := t.filter.Sorted()
		if Abs64(median) < 70*60 {
			t.offset = median
		} else {
			t.offset = 0
			found := false
			for _, v := range sorted {
				if v != 0 && Abs64(v) < 5*60 {
					found = true
					break
				}
			}
			if !found {
				// In the C++ code this would log an incorrect clock warning.
			}
		}
	}
	t.mu.Unlock()
}
