# SQL Server driver — known limitations

This driver implements both roles: `AsWarehouse` (export to Parquet, ingested by
DuckDB) and `AsOLAP` (metrics views query SQL Server directly, selected with
`olap_connector: mssql`). The warehouse path is the one exercised in production
today; the OLAP path is implemented but less travelled, and the first entry below
matters only there.

The list was produced by an adversarial code review (OpenAI Codex, GPT-5.x,
2026-08-30) prompted to read the diff as a Rill maintainer would. Every entry was
then checked by hand against the sibling drivers. Items are kept, not fixed, when
the sibling drivers behave the same way: this driver deliberately matches the
existing conventions rather than diverging from them. Fixing any of these is a
change to the shared pattern and belongs in a separate PR that changes every
driver at once.

## Shared with the sibling drivers — left as-is on purpose

### 1. Parameter placeholders are `?`, which `sqlserver` does not accept

`dialect.go` embeds `drivers.BaseDialect`, which emits `?` placeholders
(`runtime/drivers/dialect.go:267` and `:446`), reached from the metrics view query
builder at `runtime/metricsview/ast.go:1427`. The `sqlserver` `database/sql` driver
requires `@p1…@pN` and nothing calls `sqlx.Rebind` anywhere in the runtime.

A metrics view with a time dimension queried through `olap_connector: mssql` will
therefore fail with a syntax error. The warehouse path does not go through this
code and is unaffected.

**`postgres` has the identical defect**: `postgres/dialect.go` is a bare
`BaseDialect` with only an identifier escaper, and it opens `pgx`, which requires
`$1`. `pinot` is the only driver that overrides the placeholder
(`GetTimeDimensionParameter`). Verified by reading the code, not by running a
Postgres OLAP query.

Fixing this properly means teaching `BaseDialect` to emit ordinal placeholders,
which changes Postgres too.

### 2. A malformed DSN can leak the password

`ResolveDSN` returns the secret unvalidated and the `sqlx.Open` error is propagated
verbatim. `go-mssqldb` parses with `net/url`, whose errors quote the whole URL
including userinfo. In `QueryAsFiles` the message also reaches the span status.

`mysql` (`sqlx.Open("mysql", dsn)`) and `postgres` (`sqlx.Connect("pgx", …)`) do
exactly the same. No driver in the tree redacts DSNs.

### 3. Encryption has no enforced default

`encrypt` carries a `Placeholder: "mandatory"` but no `Default`, and the property
is omitted from the DSN when empty, so `go-mssqldb` falls back to its own default.
The spec's advertised `trust_server_certificate=false` is never serialised either.

`mysql`'s `ssl-mode` is defined the same way — `Placeholder: "require"`, no
`Default`. Note that a spec `Default` is a form hint; it is not applied by the
runtime, so a real fix has to live in config resolution, not in the property list.

### 4. `BINARY`/`VARBINARY`/`IMAGE` are reported as `CODE_STRING`

`olap.go:187` maps them to string and `prepareScanDest` scans into
`sql.NullString`, so non-UTF-8 bytes can be mangled in OLAP results. The warehouse
path is correct and writes Arrow Binary.

`mysql` groups `BINARY`, `VARBINARY` and the `BLOB` family with `CHAR`/`VARCHAR`
onto `CODE_STRING` in the same way. `postgres` gets it right (`BYTEA` →
`CODE_BYTES`), so the tree is split; this driver follows `mysql`.

### 5. The 64 MB threshold does not bound batch memory

`rowGroupBufferSize` is checked only after a full 8192-row Arrow batch is built, so
a table of `VARCHAR(MAX)`/`VARBINARY(MAX)` values can allocate far more than the
limit suggests.

Copied from `snowflake`, which uses the same 64 MB constant and the same
row-count-driven batching.

### 6. Parquet close errors are discarded

`writer.Close()` and `fw.Close()` run under `defer` and their errors are dropped,
so a footer that fails to write can still be reported as a successful file.

`snowflake` does the same (`defer writer.Close()`, `defer fw.Close()`).

### 7. `Lookup` reports every column as nullable

The schema query does not select `c.is_nullable` and calls
`databaseTypeToPB(typ, true)` unconditionally. (The `is_nullable` reference at
`information_schema.go:240` is inside the DDL-generating query and serves a
different purpose.)

Neither `mysql` nor `postgres` retrieves nullability at all.

## Specific to this driver

### 8. `TIME` and `UNIQUEIDENTIFIER` reach DuckDB as `VARCHAR`

Both are written as Arrow String, so DuckDB types them `VARCHAR`, while the OLAP
schema declares `CODE_TIME` and `CODE_UUID`. The UUID value itself is normalised
correctly through `NullUniqueIdentifier`.

This is a deliberate divergence from `snowflake`, which **rejects** `TIME` columns
outright (`warehouse.go:219`, "TIME data type … is not currently supported"). A
readable string was preferred over a failed ingestion, at the cost of the type. Any
change here is a behaviour change for existing projects.

### 9. Timestamps without a zone shift if the DSN sets `timezone`

`go-mssqldb` honours a `timezone` DSN parameter and returns `DATE`, `DATETIME`,
`SMALLDATETIME` and `DATETIME2` in that location. The conversion normalises to UTC,
which is right for `DATETIMEOFFSET` but wrong for the zoneless types: a midnight
`DATE` in `Europe/Madrid` can land on the previous day.

Not reachable through the structured config, which never sets `timezone`; it needs
a hand-written DSN.

### 10. Dry-run ignores query arguments

`sp_describe_first_result_set` is given the SQL only, so `stmt.Args` are dropped and
a valid parameterised query fails validation while succeeding on execution.

### 11. The connector form allows a combination the backend rejects

The "Use DSN" branch does not forbid `trust_server_certificate`, but `ResolveDSN`
rejects a DSN accompanied by that flag set to true, so the user passes form
validation and then gets `invalid config`.

## Fixed rather than documented

Three test assertions could not fail and were corrected in this PR: two scans
inside a bare `for res.Next()` in `model_executor_warehouse_self_mssql_test.go`,
which are skipped entirely on an empty result set, and a cancellation assertion
whose `errors.Is(ctx.Err(), context.Canceled)` operand is true by construction
after `cancel()`.

The warehouse round-trip test still omits `VARBINARY`, `UNIQUEIDENTIFIER`, `TIME`
and `DATETIMEOFFSET` — the four types entries 4, 8 and 9 above are about. Extending
it is worthwhile whenever those entries are revisited.

The assertion style used here is the majority one in the tree (76 occurrences of
`require.True(t, …Next())` against 27 bare loops), so this is a move towards the
convention, not away from it. The vulnerable shape does survive in 24 other driver
tests, including the two this file was modelled on
(`model_executor_sqlstore_self_mysql_test.go` and `…_postgres_test.go`, where an
identical `for res.Next()` wraps a row-count assertion that an empty `sink` would
skip). Cleaning those up is a separate, tree-wide change.
