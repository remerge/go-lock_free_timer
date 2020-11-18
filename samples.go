package lft

import (
	metrics "github.com/rcrowley/go-metrics"
	sample "github.com/remerge/go-lock_free_timer/sample"
)

func NewLockFreeSample(reservoirSize int) metrics.Sample {
	// These buckets might be pulled out of a hat but are ok for now as a first approximation.
	// 9435590846921020 is a max possible metrics's value in Prometheus on the day this code is written and
	// was split to 8 equal ranges. Subject for a further tuning.
	weirdDefaultBuckets := []float64{
		0,
		1347941549560145,
		2695883099120290,
		4043824648680435,
		5391766198240580,
		6739707747800725,
		8087649297360870,
		9435590846921020,
	}
	return sample.NewLegacyWithBuckets(reservoirSize, weirdDefaultBuckets...)
}

func NewLockFreeSampleWithBuckets(reservoirSize int, buckets []float64) metrics.Sample {
	return sample.NewLegacyWithBuckets(reservoirSize, buckets...)
}
