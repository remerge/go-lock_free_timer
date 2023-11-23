# Lock Free Timer for Go

Provides lock free timer and sample

## Timer package

Package `lft` provides a `Timer` compatible with
`github.com/rcrowley/go-metrics` without a mutex on the hot `Update` code path.

The implementation accepts a data race in exchange for much lower mutex
contention and latency impact on high volume code paths.

It's important to note that as the samples used for the timer are essentially
cleared during each snapshot, the lock free timer is only suitable for cases
where `Update` is called frequently between snapshots.

## Sample package

Package `sample` provides two different samples compatible with
`github.com/rcrowley/go-metrics`

- `BatchHistogramSample` - can be updated from histogram (`Float64Histogram`).
  Updates all bucket counts directly in a batch
- `SampleWithBuckets` - can be updated using single values. Update single bucket
  counts by passing observed values to Update

Both allow setting fixed amount of buckets.

## Install

```bash
go get github.com/remerge/go-lock_free_timer
```

## Usage

```go
start := time.Now()
...
t := lft.NewLockFreeTimer()
t.UpdateSince(start)
```

## Benchmark

```plain
BenchmarkUpdate-4                3000000               458 ns/op
BenchmarkLockFreeUpdate-4       30000000                41.8 ns/op
```
