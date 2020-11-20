package lft

import (
	metrics "github.com/rcrowley/go-metrics"
	sample "github.com/remerge/go-lock_free_timer/sample"
)

func NewLockFreeSample(reservoirSize int) metrics.Sample {
	return sample.NewLockFree(reservoirSize)
}

func NewLockFreeSampleWithBuckets(buckets []float64) metrics.Sample {
	const reservoirSize = 2048
	return sample.WithBuckets(
		sample.NewLockFree(reservoirSize), buckets...,
	)
}
