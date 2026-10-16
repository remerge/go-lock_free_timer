# Lock Free Timer for Go

Provides interval-snapshot timers and samples.

## Timer package

Package `lft` provides a `Timer` compatible with
`github.com/rcrowley/go-metrics`.
Historical public names remain compatible.
The reservoir synchronizes updates with reads, snapshots and clears.
It no longer accepts concurrent unsynchronized writes as a performance trade-off.

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

The figures below describe the earlier racy implementation.
They do not measure the synchronized reservoir.
Benchmark current concurrency and snapshot rates before deployment.

```plain
BenchmarkUpdate-4                3000000               458 ns/op
BenchmarkLockFreeUpdate-4       30000000                41.8 ns/op
```
