# Durable domain revision watermark
Status: selected normative E01 correction, overriding prior domain revision/control/restore prose; not approved architecture, fixture freeze, implementation or runtime evidence. Reviewed candidate ecc7e2e93995bc5147ad3fa240621220ddaf8d7d and [GRO-29 final-candidate-rereview](/GRO/issues/GRO-29#document-final-candidate-rereview), revision 37df43cb-b4cb-4c27-bb66-06fdeb5d24d8. Canonical GRO-21 rev2 remains unchanged. Existing static O1/O2/O3/O5 closures are retained except this new narrow lifecycle binding.

## Decision and invariant
Adopt the integrator's preferred **domain_revision_reserved** control record. The complete current authenticated kit, protected by its current surviving witnesses, is the independent watermark authority. Never derive domain order solely from a restored SQLite database, timestamps, OpID count, or the control sequence.

Let H be the maximum domain_revision across all valid domain_revision_reserved frames in the complete current lineage log; H=0 if none. Let C be the current committed SQLite domain revision. New domain attempt reserves exactly H+1 before its SQLite transaction begins. Domain revisions are safe integers 1..9007199254740991; permanent unused gaps are valid. Control sequence separately increments for every frame of any kind; its number is never used as the domain revision.

Invariant: every committed CRM domain transition has exactly one prior durably protected reservation R, and each reserved R can authorize at most one successful domain transaction. Therefore every prior committed revision <= H. After B0 restore the next attempted domain transaction uses H+1, strictly greater than every prior commit, including commits absent from B0. A reservation alone proves neither transaction execution nor success.

Allocate only when all current kit/witness states can be validated, with existing installation/kit locks and exclusive control serialization. No second writer, remote anchor or new state service is introduced.

## Closed persistent byte amendment
Add exactly one Kind to control-protocol v1:
- kind: domain_revision_reserved
- payload, exact ordered keys: {domain_revision:UInt,recorded_at:Instant}.
- domain_revision is positive, <=9007199254740991. recorded_at follows the existing real UTC second grammar.
- No OpID, source/person reference, request/plan digest, result, content or commit-status field is present in this permanent payload.

Use the existing exact ordered envelope version,lineage,key_id,sequence,previous,kind,payload; same canonical ASCII restrictions, 8-byte big-endian length, 65536-byte bound, domain-separated HMAC-SHA256, frame hash, genesis and authenticated head/witness format. No head/witness field is added: the authenticated complete chain determines H. The reader validates that each revision reservation is exactly the previous reservation value plus 1, even with intervening operation/denial/revocation frames. Reject repeats, decreases, jumps, zero, unsafe integers, unknown/duplicate keys and invalid MAC/chain. A larger syntactically valid value does not override chain validation.

Illustrative payload only (not a calculated cryptographic vector):
{"domain_revision":43,"recorded_at":"2026-10-20T12:00:00Z"}

H is recomputed during startup from the complete current stream; a cached SQLite high-watermark is never an authority. Missing prefixes, truncation required by a witness, incomplete lineage history, corruption or inability to establish current custody keeps serving CLOSED. Existing valid-frame forward reconciliation and coordinated-all-copy-rollback limits remain unchanged.

## Which transactions consume a revision
Every attempted transaction that could change authoritative CRM identities, sources, assertions/history, commitments, domain privacy enforcement or domain migration/import state reserves a fresh revision first, including owner setup and trusted recovery/maintenance writes. All changes in one atomic transaction use the same R. Explicit no-op retry returning a known receipt uses no new reservation and causes no domain effect.

Pure control-checkpoint/reservation mirroring, preview persistence, terminal-error bookkeeping, derived-index rebuilding, read/query/export and unchanged validation are not domain transactions and consume no domain revision. They may advance applied_control_sequence/head but never fabricate C. A code path cannot classify a domain mutation as bookkeeping to avoid reservation. Restoring/importing historical rows preserves their original known times/order; a recovery enforcement/migration transaction gets its new operational revision without rewriting history as newly learned.

Each domain transaction is associated with one unique authenticated server-issued OpID, including trusted internal maintenance attempts (allocate/reserve such identity under the existing contract). SQLite stores the operation, reserved revision, control reservation sequence/head, committed effects and receipt atomically. Enforce unique successful effect per OpID and unique consumed domain revision. The permanent control frame deliberately omits that association; after deleted-dependent state is scrubbed, it retains no source/person/operation-result map. Reusing an OpID with missing outcome after B0 stays unavailable; a revision reservation must never make it fresh.

## Ordered write and acknowledgement protocol
1. Validate current authority, eligible prepared plan, base_revision/targets and trusted destructive grant if needed; reconcile any unknown prior attempt first. Acquire the single domain-attempt serialization barrier. No model runs inside commit. Existing current OpID reservation must already be durably protected; internal maintenance obtains one before new effects.
2. Compute R=H+1 after capacity/deadline checks. Append the exact domain_revision_reserved frame; fsync log, atomically update kit head and fsync its parent, then advance every required witness with existing file/fsync/rename/directory-fsync ordering. Hold serialization throughout. No SQLite domain BEGIN or domain writes until all required protection is durable.
3. Mirror the reservation and current applied-control head into SQLite via control-only bookkeeping if necessary. Mark the attempt ready only within current process state; a reservation recovered without an eligible durable outcome is unavailable. Advance H in memory from the protected chain, regardless of later SQL outcome.
4. Begin the domain transaction; revalidate base_revision, identity/source policy and plan inside the writer boundary. Atomically bind R and unique OpID, perform all effects/history updates, set C=R, consume that reservation and store known committed receipt. WAL/FULL and foreign keys remain required. A failed validation, rollback, timeout before BEGIN, crash or failed commit never permits reuse of R.
5. Acknowledge committed only after SQLite durable commit and existing privacy enforcement/publication requirements. Lost reply reconciles the receipt. Commit uncertainty stays applying/unknown until actual SQLite recovery determines it; do not automatically allocate another R or repeat effects. If rollback is positively established on intact current state and exact retry is allowed by the API, a later domain attempt reserves a new revision. No old reservation is recycled.

Domain order is serialized even if up to two other supported requests overlap. New control frames cannot be interleaved into a fixed contiguous family scope if pending.first_sequence relies on that scope; the family control barrier remains held. Unrelated new domain attempts do not begin while family completion is pending.

## Crash and recovery states
| Interruption | Required behavior |
| --- | --- |
| Before any complete durable revision frame | No SQL domain attempt; reconcile log/tail under existing rules before deciding H. |
| Complete valid frame before kit head/witness advancement | Recover forward and protect all required heads; R is permanently burned, no domain effect inferred. |
| Any required witness update fails | No SQL domain attempt; remain CLOSED/pending until reconciled. R is never reused. |
| Protected frame before SQL BEGIN or before committed receipt | After loss of active attempt, preserve unavailable reservation. SQL inspection may establish a known committed receipt; otherwise reservation alone supplies no committed status. |
| SQL rollback/verified noncommit | Gap persists. Explicit permissible retry uses fresh R; no automatic restore retry. |
| SQL durable commit before reply | Recover same receipt/effect and C=R, zero new revision/effect for matching retry. |
| SQLite C exceeds H or references absent/mismatched revision reservation | Corruption/stale authority: CLOSED, never raise H from restored DB. |
| H exceeds C after crash/B0 | Expected gaps/lost domain history; preserve C as restored committed snapshot, preserve H independently, next R=H+1. Never fabricate missing commits or histories. |

Startup scans and reconciles authenticated control, derives H, imports unknown OpID reservations and applies outstanding denials before serving. Missing database domain outcomes remain unavailable even if H proves some revision was reserved. H cannot determine how many updates were lost; report verified recoverable commit point and uncertainty separately.

## Deletion, revocation and pending family interaction
Keep existing CLOSED serving barrier and durable pending scope BEFORE any family source_deleted frame. Append all required denial frames and protect kit/witness heads under the existing contiguous scope protocol. Then reserve a fresh R for the domain enforcement transaction, append/protect it, enforce all required privacy removals/sanitized outcomes and consume R atomically. Revision allocation is subsequent control work, not a source_deleted member: pending.first_sequence still denotes the first source denial, and exactly the sorted scope refs determine its denial frames.

Do not remove pending.json until the whole scope has durable denial frames, domain enforcement is committed with valid R, applied head agrees and all privacy validation succeeds. Interrupted enforcement burns R; startup completes the same pending denials without duplicate frames where already present, then reserves another fresh R for a necessary new enforcement transaction. If valid completed enforcement and receipt survive intact, do not perform a second domain effect merely to close bookkeeping. If all source denials survived but pending scope/enforcement cannot be verified, existing fail-closed rule applies.

Revocation follows the same ordering: durable denial, then protected revision, then domain enforcement and validation. If revision/frame/media/deadline budget cannot complete enforcement, do not acknowledge success; serving stays CLOSED, and denial is never rolled back. Owner-grant expiry/restart does not undo already accepted durable privacy intent; trusted recovery finishes it from validated pending/control state without asking a model to approve it.

## B0, portable import, migration and version semantics
B0 remains quarantined. Obtain the complete current kit plus independently current surviving witness, derive H, validate restored C and every retained known-order boundary <= its committed snapshot, reapply current denials/reservations, and rebuild indexes. Missing post-B0 commits remain lost/unavailable; their reserved revisions still contribute to H.

Portable export snapshot_revision is the actual committed snapshot C, not H. Historical KnownOrder values must remain bounded by C; gaps are allowed. Existing control_head binds the coherent current authority snapshot used for publication. No manifest watermark can override the authenticated kit. The watermark need not be added to closed export fields; it is derived from current control frames. An old manifest cannot supply current H or authorize activation.

Same-lineage import preserves IDs, original known times, known-order values and original snapshot C in quarantine. Every referenced historical domain revision must have a corresponding reservation in current authenticated lineage history, and no imported value may exceed H. Missing required historical reservations or foreign lineage rejects activation; do not seed H from an export. If import/migration/restore privacy reconciliation requires a new domain-changing transaction, reserve R>H through the normal protocol and commit its effects there; preserve old historical event values, and label any new actual domain changes with their real recovery commit time/order. Pure verified promotion of unchanged restored state does not consume a revision.

Invalidate old prepared plans, session grants and cursors as existing restore rules require. Next ordinary commit from restored C must allocate above H, not C+1. Counter gaps never mean complete no-history/zero lost updates. New empty lineage starts H=C=0 only when explicitly abandoning old content under fail-closed reset, never as automatic counter rollover.

## Exhaustion, support budgets and upgrade limits
Every new revision reservation consumes a frame and bytes within existing 1,000,000-frame/512 MiB lineage limits; ordinary attempts obey the 448 MiB ordinary ceiling, volume utilization and independent denial reserve. Privacy enforcement and recovery revision frames must be budgeted with denial/pending/SQL completion, including frame-count as well as byte headroom. Before accepting a family deletion reserve capacity for its bounded source frames plus at least its first enforcement revision; subsequent failures may exhaust support and leave CLOSED. No unlimited crash-retry guarantee.

All append/head/witness/SQLite waits count inside unchanged deadlines. Reserve immediately before the domain attempt to avoid model-time wastage, but do not exclude allocation I/O from commit latency. No revision allocation during automatic model repair (none is allowed). Existing OpID rate/plan admission bounds remain; every actual reattempt consumes another revision/frame even with same OpID. Measure both reservation kinds in kit-growth forecasts.

If H reaches 9007199254740991, or the relevant control/media budget cannot protect another frame, refuse new domain attempts before SQL writes. Do not wrap, reuse or infer a lower value. If deletion cannot durably finish, close serving. Eligible read/export may continue only when there is no pending privacy enforcement and current authority validates; no silent new lineage or compaction. Current preimplementation documents have no deployed consumers; this new closed kind must be integrated into the reviewed v1 grammar before implementation, with updated vectors/reader tests and frozen document identities. Existing old control streams from fixtures cannot be asserted to protect nonzero historical commits without matching revision frames.

## Exact golden and recovery fixture requirements
Integrator/QA must add independently checked fourth-kind vectors, not implementation-generated expectations. Reuse the public synthetic key/lineage:
1. Continue original three-kind sequence at control sequence 4, previous equal to the published third head; first domain reservation payload domain_revision=1 and a fixed UTC timestamp. Freeze exact P, L, frame hex, HMAC, new head and checkpoint MAC. Append a non-domain control frame then next domain_revision=2 to prove sequences are separate.
2. Freeze pending-family fixture including source frames then domain reservation, SQLite enforcement and pending removal; inject every durable boundary and a burned revision followed by a larger successful enforcement.
3. B0 at C=1; reserve/commit R=2; reserve R=3 and fail; restore B0 with complete current kit. Require H=3, C=1, next successful attempt R=4, absent post-B0 OpID unavailable, no false “committed,” old historical order unchanged. Then lost reply at R=4 yields one effect/receipt and no R=5 allocation for its replay.
4. Reject valid-MAC revision records with repeated/skipped/decreasing/zero/overflow/noninteger values, wrong kind/field/order, wrong lineage/key/previous; keep CLOSED for C>H or missing historical reservations. Distinguish expected H>C gaps.
5. Same-lineage export/import preserves old KnownOrder and C while allocating a necessary new recovery change above current H. Old manifest or kit must not lower H; missing kit/witness blocks activation.
6. Exercise safe exhaustion at integer, ordinary-byte/frame and denial-completion budgets, witness failure before SQL, uncertain SQL commit and corrupt/truncated tails under existing rules. Record no reused revisions, no duplicate effects and no privacy disclosure.

No cryptographic values were calculated and no code/models/downloads or repository changes were performed for this amendment. E01 is addressed as a concrete proposed lifecycle contract; QE may close only after integration and bounded independent re-review of the immutable changed candidate. Existing maintainer agreement/approval/default-branch merge gate remains.