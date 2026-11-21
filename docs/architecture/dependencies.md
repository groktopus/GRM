# Dependency identities and acquisition contract

Status: proposed architecture dependency selection. Metadata hashes are publisher expectations or hashes of fetched metadata, not verification of installed model/binary bytes. The governing baseline and architecture gate remain unchanged.

## Integrator disposition

Select Go go1.27.4, official MCP SDK v1.7.0, modernc.org/sqlite v1.62.0 (SQLite 3.54.1) with its required modernc.org/libc v1.80.1. Select golang.org/x/crypto v0.59.0, source commit 6c966a7517bfc4fa0f268afb1fe663013ed932a5; module h1:xS2jhqYbzNrpvkWP2YTDf4TJwBy3ykR03i8vFDxEc7c= and go.mod h1:veDyUwATGs3c+vN3VfTbzViMGXeTBk3WyhBOkn+F/4U=. Direct primary metadata: https://proxy.golang.org/golang.org/x/crypto/@v/v0.59.0.info, https://proxy.golang.org/golang.org/x/crypto/@v/v0.59.0.mod and https://sum.golang.org/lookup/golang.org/x/crypto@v0.59.0. Read-only metadata inspection confirms Go >=1.26.0 and dependent module declarations; no Go module built or checksum-database inclusion proof verified. Runtime vulnerability and compatibility checks remain required before release.

For native llama.cpp source acquisition, choose verified Git commit/tree objects as the source identity, NOT the generated tarball byte stream. Pin the full commit and tree in Data metadata below, verify Git object reachability/object hashes and the listed SHA256 documentation/license files after later acquisition, reject mismatches. No tarball is part of the supported source acquisition path, so no archive SHA256 is asserted or required. This explicitly supersedes the earlier archive-hash requirement; it does not relabel an ETag or SHA1 as SHA256. A future archive distribution needs its own reviewed content digest. Compiler family GCC 13 and CMake 3.28 target Ubuntu24.04; exact security-patched compiler/CMake package build IDs, complete transitive Go module closure, build flags and produced binary hashes are pinned with the first authorized build under Go contribution standards, before that build is accepted. They do not permit changing model/runtime/driver/protocol/privacy choices. No build manifest or workflow is introduced here.

All selected host/model/source identities below are normative expected identities. Later installation must verify acquired bytes/objects and test embedded tokenizer/template equivalence, selected protocol, actual runtime source ID and no-egress/resource behavior. Failed verification cannot fall back to moving refs. The initial prompt/typed plan/support rule and schema identity must be frozen in documentation before implementation; model acquisition/execution remains after the maintainer gate.

## Source data metadata handoff

# Dependency identity manifest and data review — cb91c20
Status: documentation-only evidence and changes-needed review, 2026-11-21. For [GRO-24 PR #3](https://github.com/groktopus/GRM/pull/3), immutable candidate `cb91c20e1fa4f69f15f95a1bbac4a2ad75de1c37`. Read docs/architecture/mvp.md FIRST at that commit; it takes precedence over retained specialist-era sections. No repository changes, model weights, archives, builds or inference were downloaded/executed. Small metadata/config/tokenizer/license/documentation bodies were read in memory; model LFS values below came from publisher metadata, not downloaded weights.

## Runtime source identity and explicit archive gap
Select native ggml-org/llama.cpp at commit **b1b770d90d94cad9e9ea5b993bdbc383a522fb28**, Git tree **03d5ee6748f415d4069dc570764c24ff3576a0f1**. Metadata resolution of v1.0.0: annotated tag object **a6d1a74dbebce3a16f4f2ceddbf38c21a08d8674**; pin the peeled commit, not the tag. Release API's target_commitish is a different value (2d8a32502909ca065c41c4cf5f2fb1c427c047fc); it is not used as source identity. [Commit metadata](https://api.github.com/repos/ggml-org/llama.cpp/commits/b1b770d90d94cad9e9ea5b993bdbc383a522fb28), [tag metadata](https://api.github.com/repos/ggml-org/llama.cpp/git/ref/tags/v1.0.0), [release metadata](https://api.github.com/repos/ggml-org/llama.cpp/releases/tags/v1.0.0).

| Runtime metadata file at selected commit | Bytes | SHA-256 of fetched body |
| --- | ---: | --- |
| LICENSE | 1078 | 94f29bbed6a22c35b992c5c6ebf0e7c92f13b836b90f36f461c9cf2f0f1d010d |
| tools/server/README.md | 139914 | 0d16e91b78b8fe1306e6f0dc69b0bd8571f0fa2dbc1ef10252299a67db9dcdbf |
| docs/build.md | 46342 | eee33d5549db002540704813766f75f2f9e86205a4bbfd9fd9a72cf9446fbc5e |
| grammars/README.md | 18051 | e1874c0f08abbcfe45eff0b7b0290faa58f630981bf94274e2d99751ef3b0208 |
| CMakeLists.txt | 11159 | ede5bbf33de939489087579c2b212a097d3bbb208146b12077d375ae5ee7dc58 |

[Immutable server documentation](https://github.com/ggml-org/llama.cpp/blob/b1b770d90d94cad9e9ea5b993bdbc383a522fb28/tools/server/README.md) contains offline, api-key-file, pooling, no-cache-prompt and json_schema controls. This is documentation inspection, not feature execution or assurance that every selected JSON Schema keyword is supported.

Source archive URL selected: https://api.github.com/repos/ggml-org/llama.cpp/tarball/b1b770d90d94cad9e9ea5b993bdbc383a522fb28 .
**source_archive_sha256: unavailable from inspected metadata; not populated.** HEAD returned ETag `45bc3a999ba6a581936b9c86b224f074f65209acd8db314beec975269e0d93b7`, but no declared digest semantics. Do not relabel ETag, Git SHA-1/tree identity or nightly-tag.txt's asset digest as the archive SHA-256. Release metadata exposes only a nightly-tag.txt asset, no source-archive digest. Under this assignment's read-only-metadata-only restriction I did not retrieve/hash archive bytes. Therefore the manifest remains incomplete against mvp.md's archive-hash approval gate. Sole author must either obtain a trustworthy published checksum, separately authorize source-archive byte verification, or explicitly propose a reviewed Git-object/tree-based source acquisition contract. No gate is silently relaxed here.

## Publisher model artifact identities
Pinned Qwen GGUF repository commit: **bb5d59e06d9551d752d08b292a50eb208b07ab1f**.
[Immutable publisher metadata](https://huggingface.co/api/models/Qwen/Qwen2.5-7B-Instruct-GGUF/revision/bb5d59e06d9551d752d08b292a50eb208b07ab1f?blobs=true).

| Exact file | Bytes | LFS SHA-256 / expected weight SHA-256 | Git pointer blob ID |
| --- | ---: | --- | --- |
| qwen2.5-7b-instruct-q4_k_m-00001-of-00002.gguf | 3993201344 | dfce12e3862a5283ccfb88221b48480e58745165de856439950d0f22590580db | 72de309249b41bf41019dd2efc5cb8e4a19c8c8e |
| qwen2.5-7b-instruct-q4_k_m-00002-of-00002.gguf | 689872288 | 539cf93f78e887edea1c04e2d7d8cdaca9d01dae9c9025bcb8accbe29df3d72a | c550b63c2e50a9984b5535c5c21ccd168c0cb08f |

Pinned Nomic GGUF repository commit: **0188c9bf409793f810680a5a431e7b899c46104c**.
[Immutable publisher metadata](https://huggingface.co/api/models/nomic-ai/nomic-embed-text-v1.5-GGUF/revision/0188c9bf409793f810680a5a431e7b899c46104c?blobs=true).

| Exact file | Bytes | LFS SHA-256 / expected weight SHA-256 | Git pointer blob ID |
| --- | ---: | --- | --- |
| nomic-embed-text-v1.5.f16.gguf | 274290560 | f7af6f66802f4df86eda10fe9bbcfc75c39562bed48ef6ace719a251cf1c2fdb | d7f9650c6271af9d27f6951932f1e78598c09dff |

LFS content digest and Git pointer blob are different identities. These are expected hashes supplied by metadata, not locally verified artifact bytes. Later operator acquisition must verify full length/hash and both Qwen shards before readiness, with no moving-ref fallback.

## Config, tokenizer, template and license identities
Use Qwen/Qwen2.5-7B-Instruct metadata reference commit **a09a35458c702b33eeacc393d103063234e8bc28**. Each hash below was computed over the small fetched file bytes at that full commit, not guessed from a filename. [Pinned metadata](https://huggingface.co/api/models/Qwen/Qwen2.5-7B-Instruct/revision/a09a35458c702b33eeacc393d103063234e8bc28?blobs=true).

| File | Bytes | SHA-256 |
| --- | ---: | --- |
| config.json | 663 | `7463bb0ea78315365e6c6b74de4e73bbcc8359dfb0c5a737584e077d42c0b03c` |
| tokenizer.json | 7031645 | `c0382117ea329cdf097041132f6d735924b697924d6f6fc3945713e96ce87539` |
| tokenizer_config.json | 7305 | `5b5d4f65d0acd3b2d56a35b56d374a36cbc1c8fa5cf3b3febbbfabf22f359583` |
| vocab.json | 2776833 | `ca10d7e9fb3ed18575dd1e277a2579c16d108e32f27439684afa0e10b1440910` |
| merges.txt | 1671839 | `599bab54075088774b1733fde865d5bd747cbcc7a547c5bc12610e874e26f5e3` |
| LICENSE | 11343 | `832dd9e00a68dd83b3c3fb9f5588dad7dcf337a0db50f7d9483f310cd292e92e` |

Qwen chat-template identity: SHA-256 **cd8e9439f0570856fd70470bf8889ebd8b5d1107207f67a5efb46e342330527f**, over UTF-8 bytes of the decoded `chat_template` string from that tokenizer_config.json, without added newline. License file is Apache-2.0; the GGUF repository LICENSE has the same Git blob identity 6634c8cc3133b3848ec74b9f275acaaa1ea618ab.

Use nomic-ai/nomic-embed-text-v1.5 metadata reference commit **e9b6763023c676ca8431644204f50c2b100d9aab**. [Pinned metadata](https://huggingface.co/api/models/nomic-ai/nomic-embed-text-v1.5/revision/e9b6763023c676ca8431644204f50c2b100d9aab?blobs=true).

| File | Bytes | SHA-256 |
| --- | ---: | --- |
| config.json | 2538 | `9ab00bd92cee80a569f708140b7b6c1661a65891ff3765b1519e181ba2f2c92b` |
| tokenizer.json | 711396 | `d241a60d5e8f04cc1b2b3e9ef7a4921b27bf526d9f6050ab90f9267a1f9e5c66` |
| tokenizer_config.json | 1191 | `d7e0000bcc80134debd2222220427e6bf5fa20a669f40a0d0d1409cc18e0a9bc` |
| vocab.txt | 231508 | `07eced375cec144d27c900241f3e339478dec958f92fddbc551f295c992038a3` |
| special_tokens_map.json | 695 | `5d5b662e421ea9fac075174bb0688ee0d9431699900b90662acd44b2a350503a` |
| 1_Pooling/config.json | 286 | `3d6180f7233df2f529d65ea60c8db4af614fb78cbe99fd81103dcfeb9f8166a5` |
| modules.json | 255 | `f93c503b32fa7bc2e186843f0add961caa14eb68f7be516e327df3cadbcd541a` |
| README.md | 71843 | `4430178908fbd574fdb6457e302d23e643d733ed036b5777180b0806946767b2` |

Nomic has no chat template in this contract: uncased BERT tokenizer plus literal search_document/search_query prefixes and mean pooling. Its Apache-2.0 declaration is publisher model-card metadata, not an invented standalone LICENSE file (none listed in the inspected base/GGUF metadata). Retain the pinned model card as license declaration and have redistribution review preserve required Apache notices. Runtime MIT license identity is above.

**Provenance limitation:** these base-model reference snapshots are independently pinned reference metadata, not evidence of the exact conversion inputs used by the older publisher GGUF commits. Hashing a base tokenizer does not prove that the selected GGUF embeds it unchanged. No GGUF header/body was downloaded. Require later embedded-tokenizer/template/pooling/config equivalence probes; any mismatch returns to documentation review, not silent runtime substitution.

## Concrete Go SQLite recommendation
Select **modernc.org/sqlite v1.62.0**, pure-Go driver/transpiled SQLite, for Ubuntu linux/amd64. Source commit **d04d243c4dfdb9f9fb19bd1c0c8c37d3923f0589**. Its own changelog and platform table bind **SQLite 3.54.1**, not the machine's sqlite CLI version. Pin **modernc.org/libc v1.80.1** exactly as its go.mod requires. Its module declares **Go 1.26.0** minimum; the governing toolchain choice must meet that, not auto-upgrade during a build. [Pinned changelog](https://gitlab.com/cznic/sqlite/-/blob/d04d243c4dfdb9f9fb19bd1c0c8c37d3923f0589/CHANGELOG.md), [module metadata](https://proxy.golang.org/modernc.org/sqlite/@v/v1.62.0.mod), [versioned package documentation](https://pkg.go.dev/modernc.org/sqlite@v1.62.0).

Module checksum identities from [Go checksum service](https://sum.golang.org/lookup/modernc.org/sqlite@v1.62.0):
- module content: `h1:gMaCjSM1qpl+iOjiWhR1LYP3zqrYDI+5GYLGt6veP2g=`
- go.mod: `h1:Jzjh1mt4s0HuQ0KvtB1LIuHiGRmU2D5pgj8dUhW7/SE=`

These are Go module h1 checksums, not raw ZIP SHA-256 or a completed local module verification. Driver license BSD-3-Clause; pinned LICENSE fetched as 1,489 bytes, SHA-256 c6fe05491a60ae13bcd223088d2705e36dede24e5587226231d2459ada5c4822. SQLite is public-domain upstream; this driver release also carries sqlite-vec v0.1.9 (MIT) according to its changelog/license documentation. Do not enable or depend on that extension for the proposed exact-vector path; nevertheless include its notice/dependency review.

Actual upstream SQLite 3.54.1 identity: SQLITE_SOURCE_ID **2026-10-16 14:45:20 555e07ef33f69a0957d24d8e555b702464a131d6f1607f4c3c1e4551e8e6ca76**. Upstream publishes sqlite3.c SHA3-256 **e2b869d3003b33c942b7dd7063604fbc53553e1b66609570df39b9b7b558ee81**. Do not call it SHA-256. [SQLite release record](https://sqlite.org/releaselog/3_54_1.html). Driver documentation lists 3.54.1; matching the eventual runtime sqlite_version()/sqlite_source_id() and compile options remains required qualification.

Compared cgo alternative: **github.com/mattn/go-sqlite3 v1.14.53**, commit **93c7e0a2ec6762a38ab371755d49a2b3f4a12e67**; inspected bundled header declares 3.54.1 and the same source ID. [Pinned bundled header](https://github.com/mattn/go-sqlite3/blob/93c7e0a2ec6762a38ab371755d49a2b3f4a12e67/sqlite3-binding.h), [module metadata](https://proxy.golang.org/github.com/mattn/go-sqlite3/@v/v1.14.53.info).
- module h1: `mnTM0owzUMrKuLo+0w51NUiB0ir2Fo7kFXlPt/l3BGs=`
- go.mod h1: `6JTjA44L93a0QCyJef5YvlPoKXntQPjzWv5gtm9sB6w=`

Modernc follows integrator preference and avoids adding a C compiler/linker to GRM's own Go build. That does not prove faster or more reliable behavior. Mattn is a credible MIT cgo alternative with a bundled upstream C engine; it adds C toolchain/ABI/build-flag review and could improve some workloads, which was not measured. Do not choose system-library linking opportunistically. Modernc's libc coupling and Go minimum are real maintenance constraints. Both remain subject to approved dependency review, vulnerability checks and actual WAL/FULL/foreign-key/locking/cancellation/online-backup/crash tests; no driver was installed, built or run.

## Immutable candidate review findings
Reviewed mvp.md and relevant registry definitions at exactly cb91c20. Findings are proposals for sole-author disposition, not unilateral architecture edits.

| ID / severity | Evidence and reproducible risk | Required disposition / observable check |
| --- | --- | --- |
| D01 HIGH — permanent OpID privacy | mvp.md reserves every accepted random OpID permanently, without linkage; registry PrepareInput nevertheless requires caller-provided OpID matching only op_ plus 1–64 base64url characters. Examples include op_delete_S7. A model/client can submit op_cobalt_origami (or encode private text into a well-formed 128-bit-looking value); permanent authority then retains private content despite S7 deletion. Shape/length does not prove independent randomness. | Specify trusted server/local-launch issuance of content-independent random operation IDs, and validate issuance before permanent reservation. Ordinary client-provided labels must not enter the forever control log. Preserve lost-response/restart semantics with a concrete issuance/retry contract and scrub sensitive invalid-input diagnostics. Test canary topic in OpID and ensure it never enters kit/witness/reservations export. |
| D02 HIGH — ambiguous correction target in export | mvp.md defines claim version key (claim.id,claim.version), but Claim.supersedes is only FactID, not version tuple. For one logical fact with versions v1/v2/v3, supersedes alone cannot identify which prior version was corrected; a same-ID value is self-referential unless version selection rules are added. | Bind exact prior version tuple or immutable assertion-version ID, or specify a unique validated predecessor ordering that survives export. Roundtrip three corrections and verify every original target, known interval and supported account; reject ambiguous/cyclic predecessor links. Do not invent target at import. |
| D03 MEDIUM — temporal ordering omitted from export | mvp.md says server commit sequence breaks UTC timestamp ties; Instant registry is seconds-only. Export rows preserve timestamps and opaque Version but no per-row commit sequence or defined sortable revision relation. Two versions within one second cannot recover that tie-break ordering from the closed fields alone. | Include explicit monotonic sequence/order semantics for exported versions (with coherent snapshot bound), or constrain updates to a documented representable ordering. Test two same-second changes and correction reconstruction after import; no lexical comparison of arbitrary Version strings. |
| D04 HIGH — typed conflict/attribution portability underspecified | Closed Claim includes predicate/participants/statement/valid/known/evidence/supersedes, but no explicit exclusivity/conflict status or span-to-speaker mapping. StoredEvidence has source and offsets only. The frozen S4 exclusive-employer conflict must remain unresolved, and S1 two spans have different attribution. A free-form statement can describe this, but closed native export does not yet bind a deterministic lossless representation. | Define supported typed semantics: e.g. separate claims per attributed span, explicit speaker role versus unknown represented by absent speaker, and a stable way to preserve exclusive-employer conflict and resolution state. If existing fields suffice, document exact synthetic records and validator rules, not implementation discretion. Roundtrip S1 unknown speaker/S4 and require same supported speakers/conflict without re-running a model. |
| D05 MEDIUM — unbounded reservation retention | Never-expiring operation reservations are outside the declared entity/assertion/fragment caps. A single model/session can submit many otherwise valid prepares, forcing permanent kit growth and fsync/witness work even with no domain effect. Disk free threshold alone does not bound recovery cost. | Declare reservation/control-log admission and rate/storage limits before accepting/reserving IDs, with deletion/recovery reserve kept usable. Never evict prior reservations to satisfy the bound. Test exhaustion yields BUSY/resource error before new permanent effect and old reserved IDs remain unavailable. |

SourceID binding is a substantive improvement: immutable body per wire SourceID plus Unicode code-point half-open spans avoids old citations silently resolving new text. Preserve it. Verify body replacement allocates new SourceID; old eligible citation still resolves old body; family deletion enumerates exact immutable IDs and dependent copies in the signed preview/digest, while independently supplied source families remain intact. Do not use body equality as identity or cross-family deletion scope. No evidence here proves that implemented mapping exists.

Export's closed fields, six-stream file manifest, current-kit gate and non-executable reservation import are sensible boundaries. D01–D04 prevent declaring lossless/private roundtrip qualified. Also require source-family referential rules, unknown grant handling, permitted unavailable-marker shape, and removal of identifying alias/locator/claim strings derived solely from deleted evidence in deterministic fixtures. These are additional validation cases, not permission to retain a deleted family map.

Independent current OpID reservation is the right protection against replaying effects absent from B0, provided issuance is truly content-independent. Crash after reservation/before preview intentionally yields unavailable; do not label it committed. Every kit-reserved ID must participate in restore checks even if never present in the database or portable stream. Authoritative current reservations cannot be replaced by B0's/export's incomplete list.

## Result and remaining evidence
Manifest metadata/model file/config/tokenizer/license entries and concrete driver/SQLite versions are populated above. **Archive SHA-256 remains genuinely unavailable under the allowed metadata-only procedure; the approval blocker is not closed.** Embedded GGUF metadata equivalence, actual dependency vulnerability/feature/host/latency/quality tests remain not run by scope. Compiler/CMake, whole module closure, separate Goose/Ollama artifact identities and prompt/schema revisions remain integration-owner/other-specialist manifest work; this contribution does not certify their pins.

Recommendation to integrator: changes needed on D01–D04 before governing approval, with D05 admission bounds documented. Record dispositions and request bounded independent re-review. Read-only source and metadata checks passed; repository status remained clean. No implementation, approval or merge is claimed.

## Host metadata handoff

# Host, Go and Ollama identity metadata for PR #3

Status: **read-only metadata expectations**, fetched 2026-11-21 for integrated candidate **cb91c20e1fa4f69f15f95a1bbac4a2ad75de1c37** in [PR #3](https://github.com/groktopus/GRM/pull/3). This populates the requested API-owned portion of the approval-blocking dependency manifest. It does not change the repository, approve the architecture, acquire a release binary/model, verify its bytes, install/run a toolchain, start inference or prove qualification. The [contract review](/GRO/issues/GRO-26#document-integrated-contract-review) remains changes-needed.

Primary sources were public GitHub release/tag metadata, the official Go download/module mirror/checksum database, and Ollama registry's small manifest/config/template metadata. No CRM/user content was sent. Source text/model-template text was treated only as data.

## Exact host/runtime release candidates

| Component | Selected release source commit | Linux x86-64 asset and length | Expected SHA-256 |
| --- | --- | --- | --- |
| Goose CLI v1.60.2 | fb55b5a6d7f5b5df96c8f817f1cde6a81880d9f1 | goose-x86_64-unknown-linux-gnu.tar.bz2; **56872060 bytes**; asset ID 771077525 | 83fef08e07c63c8cd41f9702e5e840487463f3b4e1290711a94ade5beab2d082 |
| Ollama v0.40.17 | 23b94a2daf2cd1fbf6b4dcfbfd3efbaca424b7bf | ollama-linux-amd64.tar.zst; **1545470208 bytes**; asset ID 774304976 | 1df9d6228bac5a22e420b242708c0d10e5b886db97bab99d446a41760f4e79af |

Sources: [Goose release API](https://api.github.com/repos/aaif-goose/goose/releases/tags/v1.60.2), [Goose tag commit](https://github.com/aaif-goose/goose/commit/fb55b5a6d7f5b5df96c8f817f1cde6a81880d9f1), [Goose selected asset metadata](https://api.github.com/repos/aaif-goose/goose/releases/assets/771077525), [Ollama release API](https://api.github.com/repos/ollama/ollama/releases/tags/v0.40.17), [Ollama tag commit](https://github.com/ollama/ollama/commit/23b94a2daf2cd1fbf6b4dcfbfd3efbaca424b7bf), [Ollama selected asset metadata](https://api.github.com/repos/ollama/ollama/releases/assets/774304976). Human-readable release pages: [Goose](https://github.com/aaif-goose/goose/releases/tag/v1.60.2), [Ollama](https://github.com/ollama/ollama/releases/tag/v0.40.17).

These are the non-Vulkan GNU/Linux Goose CLI archive and generic amd64 Ollama archive, not Goose ACP/Desktop/musl or an Ollama ROCm/MLX add-on. The operating profile still requires CPU-only behavior and later actual qualification; selecting an asset does not enable GPU support or execute it.

GitHub reported immutable:false for both releases. Therefore neither mutable tag nor filename alone is an immutable pin. Preserve commit + release/asset IDs + expected size/hash; after the approval gate, acquisition must reject any mismatch. The API-provided digest is a published expectation, **not an independently recomputed checksum of the binary**. Source commit and release binary correspondence is publisher metadata, not a reproducible-build/attestation verification.

Exact release download identities, supplied as references only:
- [Goose CLI asset](https://github.com/aaif-goose/goose/releases/download/v1.60.2/goose-x86_64-unknown-linux-gnu.tar.bz2)
- [Ollama runtime asset](https://github.com/ollama/ollama/releases/download/v0.40.17/ollama-linux-amd64.tar.zst)

## Go SDK v1.7.0

Module: **github.com/modelcontextprotocol/go-sdk**, version **v1.7.0**, source commit **bc72835f62eb94d0fb484439f886b6885b075f36**. [Git tag commit](https://github.com/modelcontextprotocol/go-sdk/commit/bc72835f62eb94d0fb484439f886b6885b075f36), [module mirror .info](https://proxy.golang.org/github.com/modelcontextprotocol/go-sdk/@v/v1.7.0.info), [module .mod](https://proxy.golang.org/github.com/modelcontextprotocol/go-sdk/@v/v1.7.0.mod), [checksum database lookup](https://sum.golang.org/lookup/github.com/modelcontextprotocol/go-sdk@v1.7.0), [release notes](https://github.com/modelcontextprotocol/go-sdk/releases/tag/v1.7.0).

Official mirror .info records source origin/time at 2026-07-27T15:20:53Z; GitHub tag lookup independently returned the same source commit. The .mod declares **go 1.25.0**. This is its declared module language/toolchain floor, not the selected GRM build toolchain or a build result.

Published Go module checksums:
```text
github.com/modelcontextprotocol/go-sdk v1.7.0 h1:yqjY2dsbKAC0LSuWZVBMrHgiG8ukXv6NRo0JiALay44=
github.com/modelcontextprotocol/go-sdk v1.7.0/go.mod h1:dL7u98E/zjJTGzEq+j30jQ8K2k1mb6LeAH4inEcSGts=
```

These h1 values are Go module-directory/go.mod checksums, not raw source ZIP SHA-256 values. No module ZIP was acquired; no Go sumdb inclusion/signature verification, module build or dependency/vulnerability check was executed. The .mod metadata text itself was read and its observed SHA-256 was 8679c6a4a8447d6d2d7e3ac99e922161ba1a0eaa881d7666d503cc473cdbc8a3; do not substitute this text hash for either h1 value.

SDK v1.7.0 documents backward compatibility with MCP 2025-11-25. Actual Goose-to-GRM dated negotiation and cancellation/error presentation still need authorized tests; the new SDK's latest/default protocol must not silently replace the selected contract.

## Concrete Go toolchain candidate for integrator disposition

The reviewed core has **not yet selected a Go toolchain version**. Recommend pinning **go1.27.4**, rather than leaving a “current Go” implementation choice. This is a proposed manifest completion for CTO/integrator acceptance, not a retroactive claim that cb91c20 already contains the pin.

Official metadata at [Go downloads JSON](https://go.dev/dl/?mode=json) and [release table](https://go.dev/dl/) lists:

| Field | Exact value |
| --- | --- |
| version | go1.27.4 |
| platform | linux/amd64 |
| filename | go1.27.4.linux-amd64.tar.gz |
| byte length | 70680823 |
| expected SHA-256 | 7226736c2339719cda28fbf3a54d71aa63bcefe376e98cef00034268dab24840 |
| official artifact reference | [go1.27.4.linux-amd64.tar.gz](https://go.dev/dl/go1.27.4.linux-amd64.tar.gz) |

No toolchain archive was downloaded or executed. This version exceeds the SDK's declared go1.25.0 floor, but compatibility, native driver/linker behavior and the whole dependency graph remain untested. Pin direct/indirect modules, including terminal key-parser and database driver choices, and review checksums/vulnerabilities separately before claiming a complete build manifest. Do not infer crypto/SQLite/CMake/compiler versions from this Go pin.

## Full Ollama host-model manifest identity

Requested host model: **qwen2.5:7b-instruct-q4_K_M**. This is the host's Ollama-distributed model, **not** the separate publisher Qwen GGUF shards used by GRM's internal llama.cpp instances. They must not be substituted for one another.

Observed manifest endpoint: [Ollama registry metadata](https://registry.ollama.ai/v2/library/qwen2.5/manifests/7b-instruct-q4_K_M). Two read-only observations in this review returned the same bytes/hash.

- Manifest media type: application/vnd.docker.distribution.manifest.v2+json; schemaVersion 2.
- Manifest byte length: **858**.
- SHA-256 recomputed over the exact fetched manifest metadata bytes: **845dbda0ea48ed749caafd9e6037047aa19acfcfd82e704d7ca97d631a0b697e**.
- The model page's earlier abbreviated 845dbda0ea48 is now resolved to this full manifest identity.
- Response did not provide a Docker-Content-Digest header in the read path; the value above is a hash of the actual received small manifest.
- A direct request to manifests/sha256:845dbda0ea48ed749caafd9e6037047aa19acfcfd82e704d7ca97d631a0b697e returned **404**. Do not promise digest-addressed acquisition works. A later allowed acquisition may fetch the named tag only if it verifies the exact expected manifest bytes/hash and every expected blob; a tag change must fail rather than repin automatically.

| Descriptor | Media type | Byte length | Full expected content digest |
| --- | --- | --- | --- |
| config | application/vnd.docker.container.image.v1+json | 487 | sha256:2f15b3218f0552c60647ce60ada83632d2c09755b16259b13e3e4458e9ae419d |
| layer 1 | application/vnd.ollama.image.model | 4683073952 | sha256:2bada8a7450677000f678be90653b85d364de7db25eb5ea54136ada5f3933730 |
| layer 2 | application/vnd.ollama.image.system | 68 | sha256:66b9ea09bd5b7099cbb4fc820f31b575c0366fa439b08245566692c6784e281e |
| layer 3 | application/vnd.ollama.image.template | 1482 | sha256:eb4402837c7829a690fa845de4d7f3fd842c2adee476d5341da8a46ea9255175 |
| layer 4 | application/vnd.ollama.image.license | 11343 | sha256:832dd9e00a68dd83b3c3fb9f5588dad7dcf337a0db50f7d9483f310cd292e92e |

The model-weight layer's **4,683,073,952 bytes** were **not fetched or hashed locally**. Layer lengths/digests in this table are the manifest's publisher expectations. Small config metadata (487 bytes) and template metadata (1482 bytes) were read by their content-addressed blob URLs and each recomputed hash matched its descriptor. The model metadata describes gguf/qwen2/7.6B/Q4_K_M on linux/amd64. This is not inference or proof of parser/runtime compatibility.

The fetched config's rootfs.diff_ids list is not substituted for the manifest's explicit layer descriptors; the observed system/template entries differ. Preserve the exact config+manifest tuple and require upstream/runtime format interpretation at qualification instead of inventing corrected metadata. No conclusion about invalidity is asserted from that metadata difference alone.

Exact raw manifest text (858 UTF-8 bytes, **no trailing newline** in the fetched body; the fenced display may add presentation newlines):
```json
{"schemaVersion":2,"mediaType":"application/vnd.docker.distribution.manifest.v2+json","config":{"mediaType":"application/vnd.docker.container.image.v1+json","digest":"sha256:2f15b3218f0552c60647ce60ada83632d2c09755b16259b13e3e4458e9ae419d","size":487},"layers":[{"mediaType":"application/vnd.ollama.image.model","digest":"sha256:2bada8a7450677000f678be90653b85d364de7db25eb5ea54136ada5f3933730","size":4683073952},{"mediaType":"application/vnd.ollama.image.system","digest":"sha256:66b9ea09bd5b7099cbb4fc820f31b575c0366fa439b08245566692c6784e281e","size":68},{"mediaType":"application/vnd.ollama.image.template","digest":"sha256:eb4402837c7829a690fa845de4d7f3fd842c2adee476d5341da8a46ea9255175","size":1482},{"mediaType":"application/vnd.ollama.image.license","digest":"sha256:832dd9e00a68dd83b3c3fb9f5588dad7dcf337a0db50f7d9483f310cd292e92e","size":11343}]}
```

The template/license/system digests are separately pinned, not just the large weight layer. No claim is made that the system text or template can grant tool authority. The standard provider must still expose only the five approved GRM tools, keep source text untrusted, and meet protocol/no-egress/context/quality constraints.

## What remains adverse or unverified

1. This report resolves full expected host/runtime release asset identities, SDK origin/checksums, a concrete Go toolchain candidate, and the full Ollama manifest/config/layer tuple. It does not edit the PR manifest; the integrator must incorporate these exact expectations and explicitly accept the proposed Go toolchain.
2. The wider core manifest remains incomplete until Data/Reliability/Security resolve their assigned llama.cpp source/build/compiler/CMake/SQLite/driver/crypto dependency and both internal publisher Qwen shards/Nomic/tokenizer/template/license/prompt-schema identities. This report does not fabricate or choose those other owners' values.
3. No binary/model-weight/toolchain bytes were acquired, signature/attestation checks run or software executed. GitHub release metadata is mutable; Ollama digest-addressed manifest fetch was unavailable. Preserve exact metadata expectations and verify actual later bytes before readiness.
4. No supported-host, local-provider, restricted tool inventory, no-egress offline launch, 32k host context, 8 GiB residency, 10-second query deadline, semantic quality or cancellation pass is claimed.
5. Metadata resolution cannot clear API-F1–F3 or the governing documentation-first maintainer agreement/approval/default-branch merge gate. If exact approved artifacts cannot later be obtained/verified or qualify, return an adverse result; no moving-tag/model/provider fallback.

Read-only identity-resolution work is complete with the limitations above.
