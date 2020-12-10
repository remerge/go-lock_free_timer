package sample

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomeBuckets(t *testing.T) {
	assert.EqualValues(t, []float64{20, 60, 80, 90, 100, 110, 120, 140, 180, 260, 420, 740}, DomeBuckets(1, 100, 1000, 10, 2))
}

func TestPartialBucketValues(t *testing.T) {
	updateWithValues := func(sut bucketPartialSample, values ...int64) {
		t.Helper()
		for _, value := range values {
			sut.Update(value)
		}
	}

	t.Run(`without inf`, func(t *testing.T) {
		sut := mustNewBucketPartialSample(10, 20, 30)
		updateWithValues(sut, 5, 15, 16, 23, 28)

		buckets, values := sut.Values()
		require.Len(t, buckets, 3)
		require.Len(t, values, 4)

		assert.EqualValues(t, []int64{1, 3, 5, 5}, values)
	})

	t.Run(`with inf`, func(t *testing.T) {
		sut := mustNewBucketPartialSample(10, 20)
		updateWithValues(sut, 5, 8, 20, 21, 22)

		buckets, values := sut.Values()
		require.Len(t, buckets, 2)
		require.Len(t, values, 3)

		assert.EqualValues(t, []int64{2, 3, 3}, values)
	})

	t.Run(`inf only`, func(t *testing.T) {
		sut := mustNewBucketPartialSample(10)
		updateWithValues(sut, 11, 22, 33)

		buckets, values := sut.Values()
		require.Len(t, buckets, 1)
		require.Len(t, values, 2)
		assert.EqualValues(t, []int64{0, 0}, values)
	})

	t.Run(`panics on no buckets`, func(t *testing.T) {
		assert.Panics(t, func() {
			_ = mustNewBucketPartialSample()
		})
	})
}
