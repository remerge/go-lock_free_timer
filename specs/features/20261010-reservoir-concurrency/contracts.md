# Contracts

Use the existing reservoir mutex for the whole count/selection/publication operation and for snapshot/reset.
Each update belongs to exactly one adjacent snapshot interval when no Clear intervenes.
Snapshot values must have been published before that interval was captured; no uninitialized slot or concurrent slice mutation is allowed.
Clear intentionally discards pending interval observations.
Keep constructor names, reservoir capacity, public sample methods and timer API compatible.
No package is published and no production consumer changes without separate authorization.
