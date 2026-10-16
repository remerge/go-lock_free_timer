# Concurrent Reservoir Access

Keep the existing public timer/sample interfaces and bounded uniform reservoir algorithm.
Make Update, Clear, Snapshot and all readers safe to call concurrently.
Snapshot preserves interval semantics: capture the count and current sample, reset the interval count, and return independently owned values.
Do not silently replace interval snapshots with lifetime sampling.
Do not claim timer-level counter/snapshot atomicity or safety of unrelated sample implementations.
