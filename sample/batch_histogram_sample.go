package sample

import (
	"math"
	"sync"

	runtimemetrics "runtime/metrics"

	metrics "github.com/rcrowley/go-metrics"
)

type BatchHistogramSample struct {
	mu      sync.Mutex
	buckets []float64 // Inclusive lower bounds, like runtime/metrics.
	counts  []int64
}

func NewBatchHistogramSample(buckets []float64) metrics.Sample {
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
		counts: make([]int64, len(buckets)-1),
	}
	return h
}

func (h *BatchHistogramSample) BucketsAndValues() (buckets []float64, values []int64) {
	return h.buckets, h.counts
}

func (h *BatchHistogramSample) Clear() {
	for i := range h.counts {
		h.counts[i] = 0
	}
}

func (h *BatchHistogramSample) Count() int64 {
	return int64(len(h.buckets))
}

func (h *BatchHistogramSample) UpdateFromHistogram(his *runtimemetrics.Float64Histogram) {
	counts, buckets := his.Counts, his.Buckets

	h.mu.Lock()
	defer h.mu.Unlock()

	h.Clear()
	// Copy and reduce buckets.
	var j int
	for i, count := range counts {
		h.counts[j] += int64(count)
		if buckets[i+1] >= h.buckets[j+1] {
			j++
		}
	}
}

func (s *BatchHistogramSample) Max() int64 {
	return 0
}

func (s *BatchHistogramSample) Mean() float64 {
	return 0
}

func (s *BatchHistogramSample) Min() int64 {
	return 0
}

func (s *BatchHistogramSample) Percentile(p float64) float64 {
	return 0
}

func (s *BatchHistogramSample) Percentiles(_ []float64) []float64 {
	return nil
}

func (s *BatchHistogramSample) Size() int {
	return 0
}

func (s *BatchHistogramSample) Snapshot() metrics.Sample {
	return nil
}

func (s *BatchHistogramSample) StdDev() float64 {
	return 0
}

func (s *BatchHistogramSample) Sum() int64 {
	return 0
}

func (s *BatchHistogramSample) Update(v int64) {

}

func (s *BatchHistogramSample) Values() []int64 {
	return nil
}

func (s *BatchHistogramSample) Variance() float64 {
	return 0
}
