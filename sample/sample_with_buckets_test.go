package sample

import (
	"testing"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithBuckets_Values(t *testing.T) {
	buckets := []float64{10, 20}
	sut := WithBuckets(NewLockFree(10), buckets...)

	sut.Update(5)
	sut.Update(10)
	sut.Update(15)
	sut.Update(21)

	t.Run(`Check buckets and values`, func(t *testing.T) {
		buckets, values := sut.BucketsAndValues()
		require.Equal(t, []float64{10, 20}, buckets)
		require.Equal(t, []int64{2, 3, 3}, values)
	})

	t.Run(`Legacy Values`, func(t *testing.T) {
		assert.EqualValues(t, []int64{5, 10, 15, 21, 0, 0, 0, 0, 0, 0}, sut.Values())
	})

	t.Run(`Bucket Values`, func(t *testing.T) {
		sutWithBuckets := sut.(SampleWithBuckets)

		require.NotNil(t, sutWithBuckets)
		buckets, values := sutWithBuckets.BucketsAndValues()

		assert.EqualValues(t, []float64{10, 20}, buckets)
		assert.EqualValues(t, []int64{2, 3, 3}, values)
	})
}

func TestWithBuckets_InterfaceCompatibility(t *testing.T) {
	buckets := []float64{10, 20}
	sut := metrics.NewHistogram(
		WithBuckets(NewLockFree(10), buckets...),
	)

	registry := metrics.NewRegistry()
	registry.Register(`histogram with custom sample`, sut)

	histogram := registry.Get(`histogram with custom sample`).(metrics.Histogram)
	assert.Implements(t, (*SampleWithBuckets)(nil), histogram.Sample())
}
