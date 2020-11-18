package sample

import (
	"sync"
	"sync/atomic"

	metrics "github.com/rcrowley/go-metrics"
	rand "github.com/remerge/go-xorshift"
)

type LegacyLockFreeSample struct {
	count  int64
	mutex  sync.Mutex
	values []int64
}

func NewLegacy(reservoirSize int) metrics.Sample {
	return &LegacyLockFreeSample{
		values: make([]int64, reservoirSize),
	}
}

func (s *LegacyLockFreeSample) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count = 0
	s.values = make([]int64, cap(s.values))
}

func (s *LegacyLockFreeSample) Count() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.count
}

func (s *LegacyLockFreeSample) Max() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleMax(s.values)
}

func (s *LegacyLockFreeSample) Mean() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleMean(s.values)
}

func (s *LegacyLockFreeSample) Min() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleMin(s.values)
}

func (s *LegacyLockFreeSample) Percentile(p float64) float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SamplePercentile(s.values, p)
}

func (s *LegacyLockFreeSample) Percentiles(ps []float64) []float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SamplePercentiles(s.values, ps)
}

func (s *LegacyLockFreeSample) Size() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.values)
}

func (s *LegacyLockFreeSample) Snapshot() metrics.Sample {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	count := atomic.SwapInt64(&s.count, 0)
	values := make([]int64, min(int(count), len(s.values)))
	copy(values, s.values)
	return metrics.NewSampleSnapshot(count, values)
}

func (s *LegacyLockFreeSample) StdDev() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleStdDev(s.values)
}

func (s *LegacyLockFreeSample) Sum() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleSum(s.values)
}

func (s *LegacyLockFreeSample) Update(v int64) {
	// we accept a data race here to reduce lock
	// contention and to increase performance
	count := atomic.AddInt64(&s.count, 1)
	if int(count) <= len(s.values) {
		s.values[count-1] = v
	} else {
		r := rand.Int63n(count)
		if int(r) < len(s.values) {
			s.values[r] = v
		}
	}
}

func (s *LegacyLockFreeSample) Values() []int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	values := make([]int64, len(s.values))
	copy(values, s.values)
	return values
}

func (s *LegacyLockFreeSample) Variance() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleVariance(s.values)
}

func min(s, v int) int {
	if s <= v {
		return s
	}
	return v
}
