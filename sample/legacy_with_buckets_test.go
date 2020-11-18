package sample

import (
	"testing"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyWithBuckets_Values(t *testing.T) {
	buckets := []float64{10, 20}
	sut := NewLegacyWithBuckets(10, buckets...)

	sut.Update(5)
	sut.Update(10)
	sut.Update(15)
	sut.Update(21)

	t.Run(`Legacy Values`, func(t *testing.T) {
		assert.EqualValues(t, []int64{5, 10, 15, 21, 0, 0, 0, 0, 0, 0}, sut.Values())
	})

	t.Run(`Bucket Values`, func(t *testing.T) {
		sutWithBuckets := sut.(LegacyWithBuckets)

		require.NotNil(t, sutWithBuckets)
		buckets, values := sutWithBuckets.BucketsAndValues()

		assert.EqualValues(t, []float64{10, 20}, buckets)
		assert.EqualValues(t, []int64{2, 3, 3}, values)
	})
}

func TestLegacyWithBuckets_InterfaceCompatibility(t *testing.T) {
	buckets := []float64{10, 20}
	sut := metrics.NewHistogram(
		NewLegacyWithBuckets(10, buckets...),
	)

	registry := metrics.NewRegistry()
	registry.Register(`histogram with custom sample`, sut)

	histogram := registry.Get(`histogram with custom sample`).(metrics.Histogram)
	assert.Implements(t, (*LegacyWithBuckets)(nil), histogram.Sample())
}
