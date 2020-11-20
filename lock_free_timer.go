package lft

import (
	"time"

	metrics "github.com/rcrowley/go-metrics"
)

func GetOrRegisterLockFreeTimer(name string, r metrics.Registry) metrics.Timer {
	if nil == r {
		r = metrics.DefaultRegistry
	}
	return r.GetOrRegister(name, NewLockFreeTimer).(metrics.Timer)
}

func NewLockFreeTimer() metrics.Timer {
	return &LockFreeTimer{
		counter:   metrics.NewCounter(),
		histogram: metrics.NewHistogram(NewLockFreeSample(1028)),
	}
}

type LockFreeTimer struct {
	counter   metrics.Counter
	histogram metrics.Histogram
}

func (t *LockFreeTimer) Count() int64 {
	return t.counter.Count()
}

func (t *LockFreeTimer) Max() int64 {
	return t.histogram.Max()
}

func (t *LockFreeTimer) Mean() float64 {
	return t.histogram.Mean()
}

func (t *LockFreeTimer) Min() int64 {
	return t.histogram.Min()
}

func (t *LockFreeTimer) Percentile(p float64) float64 {
	return t.histogram.Percentile(p)
}

func (t *LockFreeTimer) Percentiles(ps []float64) []float64 {
	return t.histogram.Percentiles(ps)
}

func (t *LockFreeTimer) Rate1() float64 {
	return 0.0
}

func (t *LockFreeTimer) Rate5() float64 {
	return 0.0
}

func (t *LockFreeTimer) Rate15() float64 {
	return 0.0
}

func (t *LockFreeTimer) RateMean() float64 {
	return 0.0
}

func (t *LockFreeTimer) Snapshot() metrics.Timer {
	return &LockFreeTimerSnapshot{
		counter:   t.counter.Snapshot().(metrics.CounterSnapshot),
		histogram: t.histogram.Snapshot().(*metrics.HistogramSnapshot),
	}
}

func (t *LockFreeTimer) StdDev() float64 {
	return t.histogram.StdDev()
}

func (t *LockFreeTimer) Stop() {}

func (t *LockFreeTimer) Sum() int64 {
	return t.histogram.Sum()
}

func (t *LockFreeTimer) Time(f func()) {
	ts := time.Now()
	f()
	t.Update(time.Since(ts))
}

func (t *LockFreeTimer) Update(d time.Duration) {
	t.counter.Inc(1)
	t.histogram.Update(int64(d))
}

func (t *LockFreeTimer) UpdateSince(ts time.Time) {
	t.Update(time.Since(ts))
}

func (t *LockFreeTimer) Variance() float64 {
	return t.histogram.Variance()
}

type LockFreeTimerSnapshot struct {
	counter   metrics.CounterSnapshot
	histogram *metrics.HistogramSnapshot
}

func (t *LockFreeTimerSnapshot) Count() int64 { return t.counter.Count() }

func (t *LockFreeTimerSnapshot) Max() int64 { return t.histogram.Max() }

func (t *LockFreeTimerSnapshot) Mean() float64 { return t.histogram.Mean() }

func (t *LockFreeTimerSnapshot) Min() int64 { return t.histogram.Min() }

func (t *LockFreeTimerSnapshot) Percentile(p float64) float64 {
	return t.histogram.Percentile(p)
}

func (t *LockFreeTimerSnapshot) Percentiles(ps []float64) []float64 {
	return t.histogram.Percentiles(ps)
}

func (t *LockFreeTimerSnapshot) Rate1() float64 { return 0.0 }

func (t *LockFreeTimerSnapshot) Rate5() float64 { return 0.0 }

func (t *LockFreeTimerSnapshot) Rate15() float64 { return 0.0 }

func (t *LockFreeTimerSnapshot) RateMean() float64 { return 0.0 }

func (t *LockFreeTimerSnapshot) Snapshot() metrics.Timer { return t }

func (t *LockFreeTimerSnapshot) StdDev() float64 { return t.histogram.StdDev() }

func (t *LockFreeTimerSnapshot) Stop() {}

func (t *LockFreeTimerSnapshot) Sum() int64 { return t.histogram.Sum() }

func (*LockFreeTimerSnapshot) Time(func()) {
	panic("Time called on a LockFreeTimerSnapshot")
}

func (*LockFreeTimerSnapshot) Update(time.Duration) {
	panic("Update called on a LockFreeTimerSnapshot")
}

func (*LockFreeTimerSnapshot) UpdateSince(time.Time) {
	panic("UpdateSince called on a LockFreeTimerSnapshot")
}

func (t *LockFreeTimerSnapshot) Variance() float64 {
	return t.histogram.Variance()
}
