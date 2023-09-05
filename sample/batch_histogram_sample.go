package sample

import (
	"fmt"
	"math"
	"sync"

	runtimemetrics "runtime/metrics"

	metrics "github.com/rcrowley/go-metrics"
)

// If we want to update sample from Histogram, we cannot use bucketPartialSample
// But this implementation still implements BucketsAndValues interface
// Inspired by batch histogram https://github.com/prometheus/client_golang/blob/5e78d5f66b851fef874b783814b2e884df2798d0/prometheus/go_collector_latest.go/#L454-L455
type BatchHistogramSample struct {
	mu      sync.Mutex
	buckets []float64 // Inclusive lower bounds, like runtime/metrics.
	counts  []int64
	sum     int64
}

func NewBatchHistogramSample(buckets []float64) metrics.Sample {
	// We need to remove -Inf values. runtime/metrics keeps them around.
	// But -Inf bucket should not be allowed for prometheus histograms.
	if buckets[0] == math.Inf(-1) {
		buckets = buckets[1:]
	}
	if buckets[len(buckets)-1] == math.Inf(+1) {
		fmt.Println(len(buckets))
		buckets = buckets[:len(buckets)-1]
		fmt.Println(len(buckets))
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

// Returns bucket upper bounds & theirs values + 1 extra bucket for Inf+
func (h *BatchHistogramSample) BucketsAndValues() (buckets []float64, values []int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	retBuckets := h.buckets[1:]
	retCount := make([]int64, len(h.counts)+1)
	for i := range h.counts {
		retCount[i] = h.counts[i]
		if i > 0 {
			retCount[i] += retCount[i-1]
		}
	}
	retCount[len(h.counts)] = retCount[len(h.counts)-1]

	return retBuckets, retCount
}

func (h *BatchHistogramSample) Clear() {
	for i := range h.counts {
		h.counts[i] = 0
	}
	h.sum = 0
}

// Count of observations
func (h *BatchHistogramSample) Count() int64 {
	return h.sum
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
		h.sum += int64(count)
		if buckets[i+1] >= h.buckets[j+1] {
			j++
		}
	}
}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) Max() int64 {
	return 0
}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) Mean() float64 {
	return 0
}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) Min() int64 {
	return 0
}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) Percentile(p float64) float64 {
	return 0
}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) Percentiles(_ []float64) []float64 {
	return nil
}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) Size() int {
	return 0
}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) Snapshot() metrics.Sample {
	return metrics.NewSampleSnapshot(0, []int64{})
}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) StdDev() float64 {
	return 0
}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) Sum() int64 {
	return 0
}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) Update(v int64) {

}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) Values() []int64 {
	return nil
}

// Returning NOOP because we rely on prometheus for aggregations
func (s *BatchHistogramSample) Variance() float64 {
	return 0
}
