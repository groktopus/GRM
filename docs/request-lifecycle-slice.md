# Internal request admission and lifecycle slice

This extends PR4's cooperative worker with a shared process-local owner. It is an
internal library and synthetic dependency tests, not a serving CRM or full MCP
integration. No storage, model runtime, OS lock, privacy gate, cryptography,
transport, deployment, or administrative operation is implemented.

## Governing architecture

Maintainer-merged [PR3](https://github.com/groktopus/GRM/pull/3), commit
`ae89cf5a2e7be55e79e1f0f272a548e9c5b20059`, governs this implementation despite
historical proposed-status headings. See
[ADR0002](https://github.com/groktopus/GRM/blob/ae89cf5a2e7be55e79e1f0f272a548e9c5b20059/docs/adr/0002-grm-mvp-local-evidence-architecture.md),
[MVP limits](https://github.com/groktopus/GRM/blob/ae89cf5a2e7be55e79e1f0f272a548e9c5b20059/docs/architecture/mvp.md),
[MCP concurrency/cancellation](https://github.com/groktopus/GRM/blob/ae89cf5a2e7be55e79e1f0f272a548e9c5b20059/docs/architecture/mcp.md#replay-interruptions-concurrency-and-cancellation),
[operations resource limits](https://github.com/groktopus/GRM/blob/ae89cf5a2e7be55e79e1f0f272a548e9c5b20059/docs/architecture/operations.md),
and [installation/session ownership](https://github.com/groktopus/GRM/blob/ae89cf5a2e7be55e79e1f0f272a548e9c5b20059/docs/architecture/trust.md).
Architecture files and dependencies are unchanged. GRO-21/GRO-29 acceptance and
independent review gates remain unchanged.

## Supported behavior and ownership

Create exactly one `admission.Owner` after obtaining the installation lock and
route all supported interactive calls through it. Multiple owners do not share
permits: process wiring and second-process exclusion remain integration duties.
The zero value rejects work; New is the only constructor. Never copy an owner.

`Run` represents one query, prepare or commit invocation, with one homogeneous
boundary (ordinary, writer or model) and a finite batch of up to nine internal
jobs. The existing worker primitive runs one job at a time for that request;
this intentionally uses fewer than its two available worker slots. Two admitted
requests can execute concurrently. The eight internal pending batch entries are
not eight additional admitted interactive requests. Callers retain immutable job
storage until return. Invalid operation/boundary, nil dependency or oversized
batch rejects before dependency calls. Empty batches have no effects.

| Approved bound | Implementation and evidence |
| --- | --- |
| 2 executing / 8 waiting requests | Mutex-protected shared admission, FIFO waiting list; simultaneous 20 callers with two held executing show eight waiting and ten busy rejections, exactly two dependency entries. |
| Busy before effects | Rejection precedes worker creation and dependency dispatch. ErrBusy identifies retryable capacity exhaustion; transport mapping remains unimplemented. |
| Query 10s / prepare 30s / commit 10s | Context timeout created before admission; earlier parent deadline wins. Admission and dependency permit waits consume that same context. Tests verify deadline bounds and queued/permit-wait expiry. |
| 1 writer / 1 internal model job | Separate shared permits wrap dependency calls across concurrent Run invocations. Waiting for a permit occupies an executing request slot and consumes its deadline. Model and writer can run concurrently with one another. |
| Cancellation/error identity | errors.Is/As retain dependency errors. On caller timeout/cancellation or shutdown, context error takes return precedence. Shutdown errors also match ErrClosed. No automatic retries. |
| Shutdown stop/cancel/drain/join | Closed flag and cancellation set under admission mutex; queued callers remove themselves, executing workers join before release, Shutdown waits for the last admitted request to release. Concurrent/repeated Shutdown is safe and cannot reopen admission. |

FIFO promotion prevents new arrivals overtaking queued requests. Cancellation of
queued work removes its entry; promotion of already canceled work cannot invoke
its dependency after the context check. Cancellation can race a dispatch or
successful completion, so dependencies must check/use their context. There is no
product latency or general OS scheduler fairness guarantee. Capacity can remain
occupied until cooperative work returns. Tests cover cancellation, failure and
FIFO progress under this supported behavior; they do not establish throughput.

The caller owns Run; Run owns its worker; worker joins each synchronous dependency
call. Dependencies must honor context across all I/O, must not detach work, and
must accurately declare their boundary. The API does not support reentrant calls
to this owner or Shutdown from a dependency (self-wait/deadlock). No arbitrary
callback sandbox or forced I/O termination exists. Panic is not converted into an
ordinary request error. Noncooperative I/O can delay Run/Shutdown indefinitely;
the test holds such I/O, observes Shutdown waiting, then releases it and verifies
zero active calls after shutdown. No goroutine exists solely to schedule, promote,
or wait for shutdown.

## Commit outcome and integration limits

A cancellation/deadline/closed error describes request lifecycle, never durable
rollback. A real writer may commit just before cancellation; that outcome remains
unknown to the caller until the approved operation-status/receipt reconciliation.
This slice implements neither transactions nor reconciliation and supplies no
fake aborted/committed receipt. Prepare/query callbacks must not commit facts;
real adapters must enforce domain permissions, transaction boundaries and the
approved persistent serving/privacy gate. A Writer permit is serialization, not
transactional correctness. A Model permit is exclusion, not model readiness,
resource qualification or the model protocol implementation.

Homogeneous batches and the three interactive deadlines are the supported API.
Administrative restore/export, multi-stage model-plus-writer orchestration,
reopening after a privacy barrier, startup readiness, storage lock waits, SDK
cancellation notifications and shutdown signal wiring remain outside this slice.
Unsupported lifecycle operations cannot be selected through Run. A closed owner
always rejects valid calls; invalid configurations return ErrInput without effects.
No control-state persistence or recoverability claim is made.

## Verification

Use the existing verified Go 1.27.4 toolchain with GOTOOLCHAIN=local and pinned
`golang.org/x/vuln/cmd/govulncheck` v1.10.0. No module change or dependency upgrade.
Run gofmt, go vet ./..., go test -count=1 -v ./..., go test -race -count=1 ./...,
go test -race -count=20 ./internal/admission and govulncheck ./.... Synthetic tests
use entered/completed handshakes; polling only observes locked scheduler state and
has a failure timeout. Short deadlines exercise expiration, not latency targets.
Race results cover exercised schedules, not exhaustive proof. Logs and source/tool
hashes are delivered with GRO-40 rather than committed as product files.
