package sample

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// Provides a partial implementation of metrics.Sample interface to be used as a mixin.
// The whole implementation is based on an underlying `prometheus.Histogram`.
type bucketPartialSample struct {
	// We completely rely on prometheus's original implementation here because it is completely lock-free
	impl prometheus.Histogram
}

func mustNewBucketPartialSample(buckets ...float64) bucketPartialSample {
	if 0 == len(buckets) {
		panic(`No buckets specified`)
	}

	return bucketPartialSample{
		impl: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    `dummy`,
			Buckets: buckets,
		}),
	}
}

func (b bucketPartialSample) Update(val int64) {
	b.impl.Observe(float64(val))
}

// Returns bucket upper bounds & theirs values + 1 extra bucket for Inf+
func (b bucketPartialSample) Values() (buckets []float64, values []int64) {
	snapshot := b.mustSnapshot()

	buckets = make([]float64, len(snapshot))
	values = make([]int64, len(snapshot)+1)

	for idx, bucket := range snapshot {
		buckets[idx] = bucket.GetUpperBound()
		values[idx] = int64(bucket.GetCumulativeCount())
	}

	lastBucket := values[len(buckets)-1]
	values[len(buckets)] = lastBucket

	return buckets, values
}

func (b bucketPartialSample) mustSnapshot() []*dto.Bucket {
	snapshot := dto.Metric{}
	if err := b.impl.Write(&snapshot); nil != err {
		panic(err)
	}

	return snapshot.Histogram.GetBucket()
}
