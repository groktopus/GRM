# Contributing to GRM

GRM is in discovery. Bring a concrete relationship problem, a useful source, or an example that exposes a weakness in the design. “I remembered the promise but forgot who offered the introduction” is more useful than “add more integrations.”

## Architecture first

**The governing architecture must be documented, agreed upon, approved by a maintainer, and merged into the default branch before any implementation starts.** The same rule applies to later changes that affect the architecture.

1. Open a documentation-only PR explaining the problem, options, tradeoffs, and how we would tell whether the design works.
2. Resolve review comments, get maintainer approval, and merge the architecture documentation.
3. Only then start the implementation in a separate PR. Link the merged architecture document in that PR. If the design changes, return to step 1.

The research documents are inputs to that process, not approval to implement. Do not submit speculative scaffolding or implementation prototypes in place of architecture review.

## Go backend contributions

Go is the language of choice for all GRM backend services, including the proof of concept. See [ADR 0001](docs/adr/0001-choose-go-for-backend-services.md) for the rationale and the [Go contribution standards](docs/go-standards.md) for style, dependencies, and checks. Significant decisions go in the [ADR log](docs/adr/README.md). Use a new ADR to change an accepted decision rather than rewriting why it was made.

For research contributions, link the primary source and explain which design decision it informs. For examples, use fictional people and synthetic records. Do not submit real contacts, private meeting transcripts, sensitive events, or credentials.
