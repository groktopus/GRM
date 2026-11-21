# Internal interpretation grammar contract v1

Proposed documentation-only schema specification, paired with the literal EXTRACT-1/PLAN-1/SUPPORT-1 prompts in [retrieval](retrieval.md). No runtime has executed them. These output schemas are internal untrusted model boundaries, not MCP tools or direct database plans. The server assigns authority, IDs, revisions, eligible candidates, evidence, preconditions, digest and receipts after validation.

Compose the definitions below with the consolidated MCP registry's exact transitive IdentitySlot/Action definitions at the frozen architecture commit. They do not introduce network schema resolution or let a model select its schema. Materialize the same self-contained local-reference closure rule as MCP. No extra fields, SQL, endpoint, tool or instruction payload is accepted as executable data. Null fields below are explicit internal route grammar values, separate from MCP's documented no-null/readiness exception.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$defs": {
    "ExtractCandidate": {
      "type": "object",
      "properties": {
        "state": {"const": "candidate"},
        "identities": {"type": "array", "items": {"$ref": "#/$defs/IdentitySlot"}, "minItems": 0, "maxItems": 20},
        "actions": {"type": "array", "items": {"$ref": "#/$defs/Action"}, "minItems": 1, "maxItems": 100}
      },
      "required": ["state", "identities", "actions"],
      "additionalProperties": false
    },
    "ExtractClarification": {
      "type": "object",
      "properties": {
        "state": {"const": "needs_clarification"},
        "reason": {"enum": ["ambiguous_identity", "missing_evidence", "missing_time", "missing_copy_authorization", "unsupported_intent"]}
      },
      "required": ["state", "reason"],
      "additionalProperties": false
    },
    "ExtractOutput": {"oneOf": [{"$ref": "#/$defs/ExtractCandidate"}, {"$ref": "#/$defs/ExtractClarification"}]},
    "PlanOutput": {
      "type": "object",
      "properties": {
        "route": {"enum": ["exact", "topic", "path", "history", "commitments", "unsupported"]},
        "topic": {"oneOf": [{"type": "string", "minLength": 1, "maxLength": 4096}, {"type": "null"}]},
        "requested_attribute": {"enum": ["discussion", "expertise", "other", null]},
        "subject_slots": {"type": "array", "items": {"type": "integer", "minimum": 0, "maximum": 19}, "minItems": 0, "maxItems": 20, "uniqueItems": true},
        "need_clarification": {"type": "boolean"}
      },
      "required": ["route", "topic", "requested_attribute", "subject_slots", "need_clarification"],
      "additionalProperties": false
    },
    "SupportLabel": {"enum": ["supported", "irrelevant", "uncertain"]}
  }
}
```

PLAN-1 additionally checks topic token cap <=512, exact substring membership in original question, topic-route nonnull/attribute rules, non-topic nulls, authorized subject-slot membership, ambiguity and deterministic cutoff validation. Character cap does not substitute for tokenizer cap. No model may set new temporal/identity scope by inventing a slot.

SUPPORT-1 grammar is deterministically materialized from **N admitted whole candidates**, N=1..20, in query order. Root type object; properties are exactly decimal strings `0` through `N-1`, each SupportLabel; required is that entire ordered list; additionalProperties false. N=0 skips inference and uses authoritative coverage logic, never treats absence of candidates as a universal negative. Keys with leading zeros, negative indices, skipped/extra keys or non-enum values are rejected; JSON member order is not a support decision. The same dynamic schema and input candidate list are used for a paired arm; model cannot add or choose candidates.

Two-candidate complete shape example:

```json
{"0":"supported","1":"irrelevant"}
```

This is not a claim that arbitrary first candidate is supported. Gold support/irrelevance labels derive from the frozen synthetic evidence and independent rubric. Missing/uncertain/truncated/invalid classifier output yields explicit failure/partial per retrieval contract, not an accepted prefix. All returned sources/spans/claims/speakers are selected by deterministic joins to existing eligible records; the model emits none of them.

EXTRACT-1's IdentitySlot/Action output is only a candidate. Go regenerates final candidate lists, checks exact mention substrings/creation evidence/target revisions and model-unsupported operations, then constructs the authoritative preview. It does not trust model candidates, grant requirements or permissions. Destructive deletion bypasses extraction and uses fixed typed target scope with trusted terminal grant. Export KnownOrder and exact predecessor tuples are server-owned, never chosen by the model.

Changing literal prompts, normalization version/parameters, schema/grammar, model tuple or classification template changes the evaluated revision; preserve previous evidence and re-freeze before comparison. These are bounds and observable contracts, not measured model quality or deadline guarantees.
