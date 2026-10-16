package sample

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReservoirConcurrentUpdatesAndSnapshots(t *testing.T) {
	const workers, updates = 8, 10000
	sample := NewLockFree(32)
	var writers sync.WaitGroup
	start := make(chan struct{})
	for worker := 0; worker < workers; worker++ {
		writers.Add(1)
		go func() {
			defer writers.Done()
			<-start
			for i := 0; i < updates; i++ {
				sample.Update(7)
			}
		}()
	}
	finished := make(chan struct{})
	go func() { writers.Wait(); close(finished) }()
	close(start)
	var observed int64
	for {
		snapshot := sample.Snapshot()
		observed += snapshot.Count()
		for _, value := range snapshot.Values() {
			require.Equal(t, int64(7), value)
		}
		select {
		case <-finished:
			observed += sample.Snapshot().Count()
			require.Equal(t, int64(workers*updates), observed)
			require.Zero(t, sample.Count())
			return
		default:
		}
	}
}

func TestReservoirConcurrentClearAndReads(t *testing.T) {
	sample := NewLockFree(8)
	var group sync.WaitGroup
	for worker := 0; worker < 4; worker++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for i := 0; i < 1000; i++ {
				sample.Update(int64(i))
				sample.Count()
				sample.Values()
				sample.Min()
				sample.Max()
				sample.Mean()
				sample.Variance()
				sample.StdDev()
				sample.Sum()
				sample.Percentiles([]float64{0.5, 0.99})
				sample.Size()
				sample.Clear()
			}
		}()
	}
	group.Wait()
	sample.Clear()
	require.Zero(t, sample.Count())
	require.Empty(t, sample.Snapshot().Values())
}

func TestReservoirSnapshotInterval(t *testing.T) {
	sample := NewLockFree(4)
	sample.Update(2)
	sample.Update(4)
	first := sample.Snapshot()
	require.Equal(t, int64(2), first.Count())
	require.ElementsMatch(t, []int64{2, 4}, first.Values())
	sample.Update(9)
	second := sample.Snapshot()
	require.Equal(t, int64(1), second.Count())
	require.Equal(t, []int64{9}, second.Values())
	require.ElementsMatch(t, []int64{2, 4}, first.Values())
	require.Zero(t, sample.Snapshot().Count())
}

func BenchmarkReservoirConcurrentUpdate(b *testing.B) {
	sample := NewLockFree(1028)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sample.Update(7)
		}
	})
}
