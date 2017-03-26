# go-lock_free_timer

Package `lft` provides a `Timer` compatible with
`github.com/rcrowley/go-metrics` without a mutex on the hot `Update` code path.

The implementation accepts a data race in exchange for much lower mutex
contention and latency impact on high volume code paths.

## Install

```bash
go get github.com/remerge/go-lock_free_timer
```

## Usage

TBD

## Benchmark

```
BenchmarkUpdate-4                3000000               458 ns/op
BenchmarkLockFreeUpdate-4       30000000                41.8 ns/op
```
