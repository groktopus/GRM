# Go contribution standards

These standards apply once implementation is authorized. They do not start implementation or select additional architecture.

## Make code easy to review

- Run `gofmt`. Follow the [Go review conventions](https://go.dev/wiki/CodeReviewComments). Use clear domain names and keep packages and interfaces small. Prefer the standard library and concrete types over a framework or abstraction with no present use.
- Keep transport, model-proposed input, domain validation, and storage distinct. Decode external and model data into typed values; validate at the boundary before changing state. Do not let model text choose how or where writes execute.
- Handle errors explicitly, wrap them with useful context without exposing private data, and propagate request cancellation and deadlines across I/O. Don’t introduce a goroutine without a clear owner, bounds, and shutdown path.

## Keep builds auditable

- Pin the supported Go toolchain and versioned build dependencies when the build is introduced. Commit and review `go.mod` and `go.sum` changes. Explain new direct dependencies, including the problem they solve; keep unrelated upgrades out of feature PRs.
- Keep generated files reproducible from reviewable inputs and pinned tools. Don’t hand-edit output or introduce new dependency, runtime, or toolchain requirements quietly. Run a pinned [`govulncheck`](https://go.dev/doc/security/vuln/) once modules exist; resolve applicable findings or document an approved exception.

## Prove behavior, not just compilation

- Baseline checks when code exists: formatting, `go vet ./...`, and `go test ./...`. Use `go test -race ./...` where supported for concurrency-sensitive changes. Use table-driven tests for clear edge cases and fuzz parsers where malformed input is a meaningful risk.
- Use fixed clocks, synthetic records, and test doubles at external model boundaries. Run integration tests for the actual storage and MCP boundaries once selected. Test error paths and replay, not just success. Avoid personal information or transcripts in fixtures, logs, or test output.
- Keep the change and the evidence small enough to review. Do not weaken a check or rewrite an acceptance test merely to pass CI. If the intended behavior changes, explain it and get the expectation reviewed. State which checks ran and which did not.
