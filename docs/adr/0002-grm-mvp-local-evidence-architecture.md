# ADR 0002: Local transactional relationship evidence for MVP

Status: Proposed. Date: 2026-11-21. Integrator: Systems Architect; technical ratification: CTO; product authority: CPO; approval/merge: maintainer.

## Context and decision

The canonical MVP requires conversational capture/query, stable identity, attributed paraphrase recall, supported paths, valid/known time, correction/deletion and safe self-hosted recovery. [ADR 0001](0001-choose-go-for-backend-services.md) decides language only. Adopt the [governing MVP proposal](../architecture/mvp.md) and its linked contracts: one Go stdio server/session, SQLite WAL/FULL authority on local ext4, selected immutable source excerpts, typed temporal assertions, exact/hybrid retrieval and same-host inference/assistant dependencies.

Use five task-oriented MCP tools; GRM translates capture and query internally. No automatic merge, remote fallback, outreach, connectors or update queue (A10 omitted). Use explicit source-copy eligibility, current-policy publication barriers and independent current control authority plus witness before any restored content becomes usable. Reserve operation identities outside rollbackable content snapshots; deletion scrubs private comparison material and forbids replay effects. Native portability is manifest plus closed JSONL streams.

## Alternatives and consequences

Dedicated graph/vector services improve some scaling options but create a second consistency, deletion and restore boundary without measured need. PostgreSQL remains a credible future choice if the declared local envelope fails; changing it needs a new approved ADR. Lexical-only is a comparison control and may fail paraphrase recall; local dense retrieval is a candidate, not a proven improvement. Hosted inference reduces hardware demands but adds a private-data destination and is excluded absent CPO disposition. Multiple serving sessions/remote HTTP need a new authentication/topology decision. Warning-only old-backup restoration violates privacy, so loss of current control state blocks full restoration rather than republishing deleted material.

Costs include local CPU/RAM and dependency maintenance, restricted host tools, a separate terminal grant for deletion, required current recovery kit/witness availability, conservative limits and possible unavailable replay detail after restore. Simultaneous rollback/loss of every trusted head cannot be detected or repaired automatically. No quality, latency, compatibility, durability or recovery target is asserted as measured.

## Confirmation, approval and evolution

Before implementation: resolve the dependency identity manifest and independent Q/AE findings, freeze technical fixtures against canonical baseline revision, obtain CTO technical disposition and maintainer agreement/approval/default-branch merge of this documentation-only PR. Later implementation cites merged documents and uses Go standards. Qualify all applicable A01–A13 and R1–R6 individually on immutable released artifacts/selected host and real boundaries; compare retrieval support/citations independently of fluency, measure resource/deadline fit and exercise deletion-safe B0 restore. See [coverage and handoff](../architecture/mvp.md#fitness-evidence-review-and-implementation-handoff).

Reopen for changed privacy boundary, host/session/transport, storage/retrieval, resource limits or lifecycle semantics. Failure to meet the frozen oracle is a defect or design revision, never permission to weaken the expected facts. This ADR is not accepted until the maintainer gate is satisfied.
