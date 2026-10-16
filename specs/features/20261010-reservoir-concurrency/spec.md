# Concurrent Reservoir Access

Keep the existing public timer/sample interfaces and bounded uniform reservoir algorithm.
Make Update, Clear, Snapshot and all readers safe to call concurrently.
Snapshot captures the count and current sample, then resets the interval.
The returned values are independently owned.
Do not silently replace interval snapshots with lifetime sampling.
Timer-level counter/snapshot atomicity is not guaranteed.
The safety of unrelated sample implementations is outside this change.
