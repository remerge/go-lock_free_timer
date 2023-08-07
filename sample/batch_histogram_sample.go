package sample

import (
	"math"
	"runtime/metrics"
	"sync"
)

type BatchHistogramSample struct {
	mu      sync.Mutex
	buckets []float64 // Inclusive lower bounds, like runtime/metrics.
	counts  []uint64
}

func NewBatchHistogramSample(buckets []float64) *BatchHistogramSample {
	// We need to remove -Inf values. runtime/metrics keeps them around.
	// But -Inf bucket should not be allowed for prometheus histograms.
	if buckets[0] == math.Inf(-1) {
		buckets = buckets[1:]
	}
	h := &BatchHistogramSample{
		buckets: buckets,
		// Because buckets follows runtime/metrics conventions, there's
		// 1 more value in the buckets list than there are buckets represented,
		// because in runtime/metrics, the bucket values represent *boundaries*,
		// and non-Inf boundaries are inclusive lower bounds for that bucket.
		counts: make([]uint64, len(buckets)-1),
	}
	return h
}

func (h *BatchHistogramSample) update(his *metrics.Float64Histogram, sum float64) {
	counts, buckets := his.Counts, his.Buckets

	h.mu.Lock()
	defer h.mu.Unlock()

	// Clear buckets.
	for i := range h.counts {
		h.counts[i] = 0
	}
	// Copy and reduce buckets.
	var j int
	for i, count := range counts {
		h.counts[j] += count
		if buckets[i+1] >= h.buckets[j+1] {
			j++
		}
	}
	if h.hasSum {
		h.sum = sum
	}
}
