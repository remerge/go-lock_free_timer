package sample

import (
	metrics "github.com/rcrowley/go-metrics"
)

type SampleWithBuckets interface {
	metrics.Sample
	BucketsAndValues() (buckets []float64, values []int64)
}

func WithBuckets(decoratee metrics.Sample, buckets ...float64) SampleWithBuckets {
	return sampleWithBuckets{
		Sample:              decoratee,
		bucketPartialSample: mustNewBucketPartialSample(buckets...),
	}
}

type sampleWithBuckets struct {
	metrics.Sample
	bucketPartialSample
}

func (s sampleWithBuckets) Update(val int64) {
	s.Sample.Update(val)
	s.bucketPartialSample.Update(val)
}

func (s sampleWithBuckets) Values() []int64 {
	return s.Sample.Values()
}

func (s sampleWithBuckets) BucketsAndValues() ([]float64, []int64) {
	return s.bucketPartialSample.Values()
}
