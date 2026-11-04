# Test Design

Reproduce the unmodified pinned dependency race in concurrent update/snapshot tests.
Verify exact observation totals across concurrent snapshots.
Verify that published values are nonzero.
Exercise concurrent updates, clears and all readers under the race detector.
Verify sequential snapshot reset and ownership explicitly.
Run the entire dependency suite with race instrumentation repeatedly.
Measure the same concurrent-update microbenchmark before and after synchronization.
Use this isolated dependency in the bidder workspace.
Repeat the unchanged native request-path test and benchmark.
Report correctness, sampling overhead and bidder performance separately.
