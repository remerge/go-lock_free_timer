package sample

import (
	"testing"

	runtimeMetrics "runtime/metrics"

	"github.com/stretchr/testify/require"
)

func TestBatchHistogramSample(t *testing.T) {

	h := NewBatchHistogramSample([]float64{10, 20, 30})
	batchH, ok := h.(*BatchHistogramSample)
	require.True(t, ok)

	float64Hist := &runtimeMetrics.Float64Histogram{}
	float64Hist.Buckets = []float64{10, 20, 30}
	float64Hist.Counts = []uint64{1, 2}

	batchH.UpdateFromHistogram(float64Hist)

	buckets, values := batchH.BucketsAndValues()
	require.EqualValues(t, buckets, []float64{20, 30})
	require.EqualValues(t, values, []int64{1, 3, 3})

	require.EqualValues(t, batchH.Count(), 3)

	float64Hist.Buckets = []float64{5, 10, 20, 30}
	float64Hist.Counts = []uint64{1, 2, 4}

	batchH.UpdateFromHistogram(float64Hist)

	buckets, values = batchH.BucketsAndValues()
	require.EqualValues(t, buckets, []float64{20, 30})
	require.EqualValues(t, values, []int64{3, 7, 7})

	require.EqualValues(t, batchH.Count(), 7)
}
