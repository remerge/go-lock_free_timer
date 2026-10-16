# Test Design

Reproduce the unmodified pinned dependency race in concurrent update/snapshot tests.
Verify exact total observation count across concurrent snapshots and nonzero published values.
Exercise concurrent updates, clears and all readers under the race detector.
Verify sequential snapshot reset and ownership explicitly.
Run the entire dependency suite with race instrumentation repeatedly.
Measure the same concurrent-update microbenchmark before and after synchronization.
Then use this isolated dependency in the bidder workspace and repeat the unchanged native request-path test and benchmark.
Keep functional correctness, sampling overhead and full-bidder performance as separate claims.
