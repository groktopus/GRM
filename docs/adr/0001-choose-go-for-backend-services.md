# ADR 0001: Choose Go for GRM backend services

**Status:** Accepted (language decision only)

**Date:** 2026-11-21

**Decision owner:** Magnus Hedemark

## Context

GRM is a self-hosted, single-user relationship CRM. Its backend will accept MCP requests, turn model suggestions into validated records, preserve when facts changed, and retrieve evidence for answers. Most of the expected work is coordinating calls to storage and model endpoints, not running inference or constructing a storage engine in the service itself. We have not benchmarked the workload. Coding agents will write much of the code, so simple diffs, static checks, and behavior a maintainer can verify matter more than the speed of generating a first draft.

## Considered options

| Option | What it offers | Why we chose it or set it aside |
| --- | --- | --- |
| Go | Compiled service, official MCP SDK, explicit errors, standard formatting and testing tools | Chosen. It gives us useful compile-time checks and a small, readable service without introducing a separate application interpreter. It suits the expected I/O-heavy workload. |
| Rust | A compiled service, official MCP SDK, ownership checks, and strong types for modeling states | A credible alternative. Its guarantees are attractive for embedded storage or intensive local processing. The initial GRM workload doesn’t justify the extra ownership and async complexity we’d have to review and maintain. |
| Python and Ruby (interpreted alternatives) | Rapid iteration; Python in particular offers convenient model and retrieval libraries | Ruled out for GRM backend services. We want to ship a compiled service without interpreter and package-environment upgrades becoming part of operating GRM. Model calls can go over an API without bringing an interpreter into this service. |

## Decision

Use **Go for the proof of concept**. Go is the language of choice for **all GRM backend services**, including any GRM-owned background worker that ships separately. Do not introduce an interpreted backend or sidecar to avoid this choice. A substantive exception requires a new approved architecture decision before implementation.

This decision does not choose a database, MCP transport, inference provider, worker topology, deployment layout, or SDK version. It does not govern third-party databases or model runtimes used as external dependencies. It does not authorize coding before the rest of the governing architecture is approved and merged.

## Consequences and confirmation

A compiled service removes an application interpreter from the deployment, not supply-chain risk. Follow the [Go contribution standards](../go-standards.md) for keeping build tools and modules reviewable.

Go’s compiler will not catch a mistaken identity match, a bad historical answer, or a duplicate write. Once implementation is authorized, use synthetic acceptance cases to verify at least these behaviors: replay doesn’t duplicate a record, a later correction preserves what was known at the earlier date, and retrying an update acknowledgment affects only the intended item. Check that the service builds and runs in a supported self-hosted environment without requiring a Python or Ruby runtime. These are planned checks, not results.

[Official Go MCP SDK](https://go.sdk.modelcontextprotocol.io/) · [Official Rust MCP SDK](https://rust.sdk.modelcontextprotocol.io/) · [Rust ownership](https://doc.rust-lang.org/book/ch04-00-understanding-ownership.html)
