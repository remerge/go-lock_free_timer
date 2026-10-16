package sample

import (
	"sync"

	"math/rand/v2"

	metrics "github.com/rcrowley/go-metrics"
)

type lockFreeSample struct {
	count  int64
	mutex  sync.Mutex
	values []int64
}

func NewLockFree(reservoirSize int) metrics.Sample {
	return &lockFreeSample{
		values: make([]int64, reservoirSize),
	}
}

func (s *lockFreeSample) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count = 0
	s.values = make([]int64, cap(s.values))
}

func (s *lockFreeSample) Count() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.count
}

func (s *lockFreeSample) Max() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleMax(s.values)
}

func (s *lockFreeSample) Mean() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleMean(s.values)
}

func (s *lockFreeSample) Min() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleMin(s.values)
}

func (s *lockFreeSample) Percentile(p float64) float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SamplePercentile(s.values, p)
}

func (s *lockFreeSample) Percentiles(ps []float64) []float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SamplePercentiles(s.values, ps)
}

func (s *lockFreeSample) Size() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.values)
}

func (s *lockFreeSample) Snapshot() metrics.Sample {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	count := s.count
	s.count = 0
	values := make([]int64, min(int(count), len(s.values)))
	copy(values, s.values)
	return metrics.NewSampleSnapshot(count, values)
}

func (s *lockFreeSample) StdDev() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleStdDev(s.values)
}

func (s *lockFreeSample) Sum() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleSum(s.values)
}

func (s *lockFreeSample) Update(v int64) {
	// Count, slot selection and publication must be one operation relative to
	// snapshots and clears. An atomic count alone cannot protect the reservoir.
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count++
	count := s.count
	if int(count) <= len(s.values) {
		s.values[count-1] = v
	} else {
		r := rand.Int64N(count)
		if int(r) < len(s.values) {
			s.values[r] = v
		}
	}
}

func (s *lockFreeSample) Values() []int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	values := make([]int64, len(s.values))
	copy(values, s.values)
	return values
}

func (s *lockFreeSample) Variance() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleVariance(s.values)
}
