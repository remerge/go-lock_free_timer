# Contracts

Use the existing mutex for count, selection and publication together.
Use the same mutex for snapshot and reset.
Each update belongs to one snapshot interval unless Clear intervenes.
Snapshot values must have been published before the interval was captured.
Uninitialized slots and concurrent slice mutation are forbidden.
Clear intentionally discards pending interval observations.
Keep constructor names, reservoir capacity and public interfaces compatible.
Publishing and production consumer changes require separate authorization.
