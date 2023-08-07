package sample

import (
	"fmt"
	"testing"

	runtimeMetrics "runtime/metrics"

	"github.com/stretchr/testify/require"
)

func TestBatchHistogramSample(t *testing.T) {

	h := NewBatchHistogramSample([]float64{10, 20, 30})
	batchH, ok := h.(*BatchHistogramSample)
	require.True(t, ok)

	float64Hist := &runtimeMetrics.Float64Histogram{}

	batchH.UpdateFromHistogram(float64Hist)
	batchH.UpdateFromHistogram(float64Hist)
	batchH.UpdateFromHistogram(float64Hist)
	batchH.UpdateFromHistogram(float64Hist)

	fmt.Println(batchH.BucketsAndValues())
}
