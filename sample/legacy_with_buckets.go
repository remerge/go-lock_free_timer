package sample

import (
	metrics "github.com/rcrowley/go-metrics"
)

type LegacyWithBuckets interface {
	metrics.Sample
	BucketsAndValues() (buckets []float64, values []int64)
}

type legacyWithBuckets struct {
	metrics.Sample
	bucketPartialSample
}

func NewLegacyWithBuckets(reservoirSize int, buckets ...float64) metrics.Sample {
	return &legacyWithBuckets{
		Sample:              NewLegacy(reservoirSize),
		bucketPartialSample: mustNewBucketPartialSample(buckets...),
	}
}

func (s *legacyWithBuckets) Update(val int64) {
	s.Sample.Update(val)
	s.bucketPartialSample.Update(val)
}

func (s *legacyWithBuckets) Values() []int64 {
	return s.Sample.Values()
}

func (s *legacyWithBuckets) BucketsAndValues() ([]float64, []int64) {
	return s.bucketPartialSample.Values()
}
