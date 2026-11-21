> Proposed governing documentation. Read [MVP integration](mvp.md) first: its explicit dispositions supersede conflicting specialist-era wording and close the listed integration gaps. No implementation is authorized until maintainer agreement, approval and default-branch merge.

# GRM MCP v1 JSON Schema companion — reconciled refinement

Status: documentation proposal, companion to the [contract proposal](/GRO/issues/GRO-26#document-mcp-contract-proposal). Amended for GRO-24 reconciliation revision 4ebf319a-8eb8-4ce9-9b9b-248483bf5a60. No implementation, approved architecture or supported runtime claimed.

JSON Schema **2020-12** definitions below govern wire shapes. The narrative additionally governs byte/global aggregate limits, valid calendar/time/span relationships, intent/action/identity consistency, grant signatures/expiry/authority, enrollment, single-server lock, replay, current privacy eligibility and frozen semantic oracle. Enable date/date-time format assertion. Schema shape validation cannot establish those domain or cryptographic properties.

## Tool schema bindings and bounded materialization

For each tools/list inputSchema/outputSchema, materialize a self-contained object with type:"object", this bundle's $schema, $ref:"#/$defs/<definition>", and **only the transitive $defs reference closure reachable from that definition**, excluding the definition itself from the root properties but retaining it in $defs. Omit bundle $id to avoid resolver collisions. No network references. Do not embed every registry definition into every tool: the five tools' ten minimal schema objects total **41725 compact UTF-8 bytes** for this revision; whole-registry repetition unnecessarily consumes model context. Actual tokenizer/context fit remains untested.

| Tool | input definition | output definition |
| --- | --- | --- |
| grm_prepare_change_v1 | PrepareInput | PrepareResult |
| grm_commit_change_v1 | CommitInput | CommitResult |
| grm_query_v1 | QueryInput | QueryResult |
| grm_read_v1 | ReadInput | ReadResult |
| grm_operation_status_v1 | StatusInput | StatusResult |

The bundle root is a registry, not a permissive input schema. Never use it alone to validate requests. TerminalGrant/GrantPayload/Signature describe the **trusted local operator artifact only**, not MCP tool input or model approval; they are unreachable from every tool input. Additional properties are rejected, including model-supplied owner_grant/token/signature fields.

## Refinement equivalence

- Required Plan.authorization is delegated/not_required for ordinary actions or owner_terminal/required for delete_source. Missing/invalid local grant leaves a deletion plan prepared and returns OWNER_GRANT_REQUIRED/OWNER_GRANT_INVALID with unchanged retry identity. Client capabilities are checked separately. confirm:true still cannot grant deletion.
- Single active stdio server: second startup rejects INSTALLATION_BUSY before discovery/database access; the startup failure is a content-free stderr reason/nonzero exit, not a private MCP response. Overlapping requests within the sole server and subsequent restarts retain the documented replay contract.
- Limits is one closed constant object reflecting conservative corpus/request/response/hop/traversal/path/hydration/evidence/deadline ceilings; required in self and Answer. Paths allow at most four entities (three hops). Global <=10 paths, <=20 hydrated fragments and <=8000 evidence tokens require semantic aggregate checks over the query, not each page.
- CommittedUnavailable exposes only opaque OpID/committed timestamp after private receipt payload/fingerprints are scrubbed. It never reapplies a mutation and carries no plan/digest/source/person association.
- Token values are ASCII base64url 1..512. Grant's authority epoch/key/nonce are protected operator values, not source-derived hashes. The 86-character example signature below is **deliberately nonauthentic**: shape acceptance is not signature verification.
- Concrete selected candidate: Goose CLI v1.60.2 + Ollama v0.40.17/local Qwen 7B path on Ubuntu 24.04 x86-64, official Go SDK v1.7.0, dated MCP 2025-11-25. Exact negotiated interoperability/full model identity/offline behavior remain adverse qualification gaps in the narrative. No code/config/download/service action occurred.

## Registry

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "urn:grm:mcp:contract:1",
  "$defs": {
    "Text": {
      "type": "string",
      "minLength": 1,
      "maxLength": 4096
    },
    "Label": {
      "type": "string",
      "minLength": 1,
      "maxLength": 200
    },
    "NL": {
      "type": "string",
      "minLength": 1,
      "maxLength": 16000
    },
    "Body": {
      "type": "string",
      "minLength": 1,
      "maxLength": 16000
    },
    "Token": {
      "type": "string",
      "minLength": 1,
      "maxLength": 512,
      "pattern": "^[A-Za-z0-9_-]{1,512}$"
    },
    "Version": {
      "type": "string",
      "minLength": 1,
      "maxLength": 128
    },
    "EntityID": {
      "type": "string",
      "minLength": 1,
      "maxLength": 70,
      "pattern": "^ent_[A-Za-z0-9_-]{1,64}$"
    },
    "SourceID": {
      "type": "string",
      "minLength": 1,
      "maxLength": 70,
      "pattern": "^src_[A-Za-z0-9_-]{1,64}$"
    },
    "FactID": {
      "type": "string",
      "minLength": 1,
      "maxLength": 70,
      "pattern": "^fact_[A-Za-z0-9_-]{1,64}$"
    },
    "CommitmentID": {
      "type": "string",
      "minLength": 1,
      "maxLength": 70,
      "pattern": "^com_[A-Za-z0-9_-]{1,64}$"
    },
    "OpID": {
      "type": "string",
      "minLength": 1,
      "maxLength": 70,
      "pattern": "^op_[A-Za-z0-9_-]{1,64}$"
    },
    "PlanID": {
      "type": "string",
      "minLength": 1,
      "maxLength": 70,
      "pattern": "^plan_[A-Za-z0-9_-]{1,64}$"
    },
    "QueryID": {
      "type": "string",
      "minLength": 1,
      "maxLength": 70,
      "pattern": "^qry_[A-Za-z0-9_-]{1,64}$"
    },
    "Temp": {
      "type": "string",
      "minLength": 1,
      "maxLength": 68,
      "pattern": "^tmp_[A-Za-z0-9_-]{1,64}$"
    },
    "SlotID": {
      "type": "string",
      "minLength": 1,
      "maxLength": 69,
      "pattern": "^slot_[A-Za-z0-9_-]{1,64}$"
    },
    "Ref": {
      "oneOf": [
        {
          "$ref": "#/$defs/EntityID"
        },
        {
          "$ref": "#/$defs/SourceID"
        },
        {
          "$ref": "#/$defs/FactID"
        },
        {
          "$ref": "#/$defs/CommitmentID"
        }
      ]
    },
    "EntityLink": {
      "oneOf": [
        {
          "$ref": "#/$defs/EntityID"
        },
        {
          "$ref": "#/$defs/SlotID"
        }
      ]
    },
    "Instant": {
      "type": "string",
      "minLength": 20,
      "maxLength": 20,
      "pattern": "^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$",
      "format": "date-time"
    },
    "Day": {
      "type": "string",
      "minLength": 10,
      "maxLength": 10,
      "format": "date",
      "pattern": "^[0-9]{4}-[0-9]{2}-[0-9]{2}$"
    },
    "Month": {
      "type": "string",
      "minLength": 7,
      "maxLength": 7,
      "pattern": "^[0-9]{4}-(0[1-9]|1[0-2])$"
    },
    "Time": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "precision": {
              "const": "instant"
            },
            "value": {
              "$ref": "#/$defs/Instant"
            }
          },
          "required": [
            "precision",
            "value"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "precision": {
              "const": "day"
            },
            "value": {
              "$ref": "#/$defs/Day"
            }
          },
          "required": [
            "precision",
            "value"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "precision": {
              "const": "month"
            },
            "value": {
              "$ref": "#/$defs/Month"
            }
          },
          "required": [
            "precision",
            "value"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "precision": {
              "const": "unknown"
            }
          },
          "required": [
            "precision"
          ],
          "additionalProperties": false
        }
      ]
    },
    "Valid": {
      "type": "object",
      "properties": {
        "from": {
          "$ref": "#/$defs/Time"
        },
        "until": {
          "$ref": "#/$defs/Time"
        }
      },
      "required": [
        "from"
      ],
      "additionalProperties": false
    },
    "Kind": {
      "type": "string",
      "enum": [
        "person",
        "organization",
        "place",
        "thing",
        "event"
      ]
    },
    "Role": {
      "type": "string",
      "enum": [
        "subject",
        "object",
        "attendee",
        "speaker",
        "lender",
        "borrower",
        "item",
        "event",
        "place",
        "organization"
      ]
    },
    "Predicate": {
      "type": "string",
      "enum": [
        "knows",
        "attended",
        "conversed",
        "located_at",
        "role_at",
        "lent",
        "returned",
        "discussed",
        "moved_to",
        "statement"
      ]
    },
    "Evidence": {
      "type": "object",
      "properties": {
        "source": {
          "oneOf": [
            {
              "$ref": "#/$defs/SourceID"
            },
            {
              "$ref": "#/$defs/Temp"
            }
          ]
        },
        "span_start": {
          "type": "integer",
          "minimum": 0
        },
        "span_end": {
          "type": "integer",
          "minimum": 1
        }
      },
      "required": [
        "source",
        "span_start",
        "span_end"
      ],
      "additionalProperties": false
    },
    "StoredEvidence": {
      "type": "object",
      "properties": {
        "source": {
          "$ref": "#/$defs/SourceID"
        },
        "span_start": {
          "type": "integer",
          "minimum": 0
        },
        "span_end": {
          "type": "integer",
          "minimum": 1
        }
      },
      "required": [
        "source",
        "span_start",
        "span_end"
      ],
      "additionalProperties": false
    },
    "Participant": {
      "type": "object",
      "properties": {
        "role": {
          "$ref": "#/$defs/Role"
        },
        "entity": {
          "$ref": "#/$defs/EntityLink"
        }
      },
      "required": [
        "role",
        "entity"
      ],
      "additionalProperties": false
    },
    "StoredParticipant": {
      "type": "object",
      "properties": {
        "role": {
          "$ref": "#/$defs/Role"
        },
        "entity": {
          "$ref": "#/$defs/EntityID"
        }
      },
      "required": [
        "role",
        "entity"
      ],
      "additionalProperties": false
    },
    "SourceDraft": {
      "type": "object",
      "properties": {
        "id": {
          "$ref": "#/$defs/Temp"
        },
        "mode": {
          "type": "string",
          "enum": [
            "owner_note",
            "authorized_copy",
            "linked"
          ]
        },
        "locator": {
          "$ref": "#/$defs/Label"
        },
        "body": {
          "$ref": "#/$defs/Body"
        }
      },
      "required": [
        "id",
        "mode",
        "locator",
        "body"
      ],
      "additionalProperties": false
    },
    "ClaimDraft": {
      "type": "object",
      "properties": {
        "id": {
          "$ref": "#/$defs/Temp"
        },
        "predicate": {
          "$ref": "#/$defs/Predicate"
        },
        "participants": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/Participant"
          },
          "minItems": 0,
          "maxItems": 20
        },
        "statement": {
          "$ref": "#/$defs/Text"
        },
        "valid": {
          "$ref": "#/$defs/Valid"
        },
        "evidence": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/Evidence"
          },
          "minItems": 1,
          "maxItems": 20
        },
        "supersedes": {
          "$ref": "#/$defs/FactID"
        }
      },
      "required": [
        "id",
        "predicate",
        "participants",
        "statement",
        "valid",
        "evidence"
      ],
      "additionalProperties": false
    },
    "CommitmentDraft": {
      "type": "object",
      "properties": {
        "id": {
          "$ref": "#/$defs/Temp"
        },
        "promisor": {
          "$ref": "#/$defs/EntityLink"
        },
        "beneficiary": {
          "$ref": "#/$defs/EntityLink"
        },
        "description": {
          "$ref": "#/$defs/Text"
        },
        "evidence": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/Evidence"
          },
          "minItems": 1,
          "maxItems": 20
        },
        "due_on": {
          "$ref": "#/$defs/Day"
        }
      },
      "required": [
        "id",
        "promisor",
        "beneficiary",
        "description",
        "evidence"
      ],
      "additionalProperties": false
    },
    "Action": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "add_source"
            },
            "source": {
              "$ref": "#/$defs/SourceDraft"
            }
          },
          "required": [
            "type",
            "source"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "add_claim"
            },
            "claim": {
              "$ref": "#/$defs/ClaimDraft"
            }
          },
          "required": [
            "type",
            "claim"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "add_commitment"
            },
            "commitment": {
              "$ref": "#/$defs/CommitmentDraft"
            }
          },
          "required": [
            "type",
            "commitment"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "rename"
            },
            "entity": {
              "$ref": "#/$defs/EntityID"
            },
            "new_name": {
              "$ref": "#/$defs/Label"
            }
          },
          "required": [
            "type",
            "entity",
            "new_name"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "complete"
            },
            "commitment": {
              "$ref": "#/$defs/CommitmentID"
            }
          },
          "required": [
            "type",
            "commitment"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "delete_source"
            },
            "source": {
              "$ref": "#/$defs/SourceID"
            }
          },
          "required": [
            "type",
            "source"
          ],
          "additionalProperties": false
        }
      ]
    },
    "Binding": {
      "type": "object",
      "properties": {
        "mention": {
          "$ref": "#/$defs/Label"
        },
        "entity": {
          "$ref": "#/$defs/EntityID"
        }
      },
      "required": [
        "mention",
        "entity"
      ],
      "additionalProperties": false
    },
    "PrepareInput": {
      "type": "object",
      "properties": {
        "contract": {
          "const": "1"
        },
        "operation_id": {
          "$ref": "#/$defs/OpID"
        },
        "intent": {
          "type": "string",
          "enum": [
            "capture",
            "correct",
            "rename",
            "complete",
            "delete"
          ]
        },
        "text": {
          "$ref": "#/$defs/NL"
        },
        "bindings": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/Binding"
          },
          "minItems": 0,
          "maxItems": 20
        },
        "targets": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/Ref"
          },
          "minItems": 0,
          "maxItems": 20
        },
        "sources": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/SourceDraft"
          },
          "minItems": 0,
          "maxItems": 10
        }
      },
      "required": [
        "contract",
        "operation_id",
        "intent",
        "text"
      ],
      "additionalProperties": false
    },
    "CommitInput": {
      "type": "object",
      "properties": {
        "contract": {
          "const": "1"
        },
        "operation_id": {
          "$ref": "#/$defs/OpID"
        },
        "plan_id": {
          "$ref": "#/$defs/PlanID"
        },
        "plan_digest": {
          "$ref": "#/$defs/Token"
        },
        "choices": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "slot": {
                "$ref": "#/$defs/SlotID"
              },
              "selection": {
                "oneOf": [
                  {
                    "$ref": "#/$defs/EntityID"
                  },
                  {
                    "const": "new"
                  }
                ]
              }
            },
            "required": [
              "slot",
              "selection"
            ],
            "additionalProperties": false
          },
          "minItems": 0,
          "maxItems": 20
        },
        "confirm": {
          "const": true
        }
      },
      "required": [
        "contract",
        "operation_id",
        "plan_id",
        "plan_digest",
        "choices",
        "confirm"
      ],
      "additionalProperties": false
    },
    "QueryInput": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "question": {
              "$ref": "#/$defs/NL"
            },
            "bindings": {
              "type": "array",
              "items": {
                "$ref": "#/$defs/Binding"
              },
              "minItems": 0,
              "maxItems": 20
            },
            "valid_at": {
              "$ref": "#/$defs/Instant"
            },
            "known_at": {
              "$ref": "#/$defs/Instant"
            },
            "page_size": {
              "type": "integer",
              "minimum": 1,
              "maximum": 20
            }
          },
          "required": [
            "contract",
            "question"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "cursor": {
              "$ref": "#/$defs/Token"
            }
          },
          "required": [
            "contract",
            "cursor"
          ],
          "additionalProperties": false
        }
      ]
    },
    "ReadInput": {
      "type": "object",
      "properties": {
        "contract": {
          "const": "1"
        },
        "ref": {
          "oneOf": [
            {
              "$ref": "#/$defs/Ref"
            },
            {
              "const": "self"
            }
          ]
        },
        "known_at": {
          "$ref": "#/$defs/Instant"
        },
        "offset": {
          "type": "integer",
          "minimum": 0
        },
        "length": {
          "type": "integer",
          "minimum": 1,
          "maximum": 2000
        }
      },
      "required": [
        "contract",
        "ref"
      ],
      "additionalProperties": false
    },
    "StatusInput": {
      "type": "object",
      "properties": {
        "contract": {
          "const": "1"
        },
        "operation_id": {
          "$ref": "#/$defs/OpID"
        }
      },
      "required": [
        "contract",
        "operation_id"
      ],
      "additionalProperties": false
    },
    "ErrorCode": {
      "type": "string",
      "enum": [
        "INVALID_INPUT",
        "UNAUTHORIZED",
        "REFERENCE_NOT_FOUND",
        "REFERENCE_DELETED",
        "REFERENCE_UNAVAILABLE",
        "REPLAY_CONFLICT",
        "STALE_PLAN",
        "PLAN_EXPIRED",
        "PLAN_ALREADY_COMMITTED",
        "IDENTITY_REQUIRED",
        "LIMIT_EXCEEDED",
        "DEPENDENCY_UNAVAILABLE",
        "CURSOR_EXPIRED",
        "CURSOR_INVALIDATED",
        "INTERNAL_ERROR",
        "OWNER_GRANT_REQUIRED",
        "OWNER_GRANT_INVALID",
        "CAPABILITY_DENIED",
        "INSTALLATION_BUSY"
      ]
    },
    "Error": {
      "type": "object",
      "properties": {
        "code": {
          "$ref": "#/$defs/ErrorCode"
        },
        "message": {
          "$ref": "#/$defs/Text"
        },
        "retry": {
          "type": "string",
          "enum": [
            "same_request",
            "after_status",
            "new_request",
            "never"
          ]
        },
        "field": {
          "$ref": "#/$defs/Label"
        },
        "operation_id": {
          "$ref": "#/$defs/OpID"
        }
      },
      "required": [
        "code",
        "message",
        "retry"
      ],
      "additionalProperties": false
    },
    "Entity": {
      "type": "object",
      "properties": {
        "id": {
          "$ref": "#/$defs/EntityID"
        },
        "version": {
          "$ref": "#/$defs/Version"
        },
        "kind": {
          "$ref": "#/$defs/Kind"
        },
        "label": {
          "$ref": "#/$defs/Label"
        },
        "aliases": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/Label"
          },
          "minItems": 0,
          "maxItems": 20
        }
      },
      "required": [
        "id",
        "version",
        "kind",
        "label",
        "aliases"
      ],
      "additionalProperties": false
    },
    "IdentitySlot": {
      "type": "object",
      "properties": {
        "slot": {
          "$ref": "#/$defs/SlotID"
        },
        "mention": {
          "$ref": "#/$defs/Label"
        },
        "kind": {
          "$ref": "#/$defs/Kind"
        },
        "candidates": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "entity": {
                "$ref": "#/$defs/EntityID"
              },
              "label": {
                "$ref": "#/$defs/Label"
              },
              "context": {
                "$ref": "#/$defs/Text"
              }
            },
            "required": [
              "entity",
              "label",
              "context"
            ],
            "additionalProperties": false
          },
          "minItems": 0,
          "maxItems": 20
        },
        "allow_new": {
          "type": "boolean"
        }
      },
      "required": [
        "slot",
        "mention",
        "kind",
        "candidates",
        "allow_new"
      ],
      "additionalProperties": false
    },
    "Plan": {
      "type": "object",
      "properties": {
        "operation_id": {
          "$ref": "#/$defs/OpID"
        },
        "plan_id": {
          "$ref": "#/$defs/PlanID"
        },
        "plan_digest": {
          "$ref": "#/$defs/Token"
        },
        "expires_at": {
          "$ref": "#/$defs/Instant"
        },
        "base_revision": {
          "$ref": "#/$defs/Version"
        },
        "identities": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/IdentitySlot"
          },
          "minItems": 0,
          "maxItems": 20
        },
        "actions": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/Action"
          },
          "minItems": 1,
          "maxItems": 100
        },
        "preconditions": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "ref": {
                "$ref": "#/$defs/Ref"
              },
              "version": {
                "$ref": "#/$defs/Version"
              }
            },
            "required": [
              "ref",
              "version"
            ],
            "additionalProperties": false
          },
          "minItems": 0,
          "maxItems": 100
        },
        "summary": {
          "$ref": "#/$defs/Text"
        },
        "authorization": {
          "$ref": "#/$defs/AuthorizationRequirement"
        }
      },
      "required": [
        "operation_id",
        "plan_id",
        "plan_digest",
        "expires_at",
        "base_revision",
        "identities",
        "actions",
        "preconditions",
        "summary",
        "authorization"
      ],
      "additionalProperties": false
    },
    "Receipt": {
      "type": "object",
      "properties": {
        "operation_id": {
          "$ref": "#/$defs/OpID"
        },
        "plan_id": {
          "$ref": "#/$defs/PlanID"
        },
        "state": {
          "const": "committed"
        },
        "committed_at": {
          "$ref": "#/$defs/Instant"
        },
        "revision": {
          "$ref": "#/$defs/Version"
        },
        "changes": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "ref": {
                "$ref": "#/$defs/Ref"
              },
              "effect": {
                "type": "string",
                "enum": [
                  "created",
                  "renamed",
                  "superseded",
                  "completed",
                  "deleted"
                ]
              }
            },
            "required": [
              "ref",
              "effect"
            ],
            "additionalProperties": false
          },
          "minItems": 1,
          "maxItems": 120
        },
        "id_map": {
          "type": "object",
          "propertyNames": {
            "pattern": "^(tmp|slot)_[A-Za-z0-9_-]{1,64}$"
          },
          "additionalProperties": {
            "$ref": "#/$defs/Ref"
          },
          "maxProperties": 120
        },
        "replayed": {
          "type": "boolean"
        }
      },
      "required": [
        "operation_id",
        "plan_id",
        "state",
        "committed_at",
        "revision",
        "changes",
        "id_map",
        "replayed"
      ],
      "additionalProperties": false
    },
    "Claim": {
      "type": "object",
      "properties": {
        "id": {
          "$ref": "#/$defs/FactID"
        },
        "version": {
          "$ref": "#/$defs/Version"
        },
        "predicate": {
          "$ref": "#/$defs/Predicate"
        },
        "participants": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/StoredParticipant"
          },
          "minItems": 0,
          "maxItems": 20
        },
        "statement": {
          "$ref": "#/$defs/Text"
        },
        "valid": {
          "$ref": "#/$defs/Valid"
        },
        "known_from": {
          "$ref": "#/$defs/Instant"
        },
        "evidence": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/StoredEvidence"
          },
          "minItems": 1,
          "maxItems": 20
        },
        "known_until": {
          "$ref": "#/$defs/Instant"
        },
        "supersedes": {
          "$ref": "#/$defs/FactID"
        }
      },
      "required": [
        "id",
        "version",
        "predicate",
        "participants",
        "statement",
        "valid",
        "known_from",
        "evidence"
      ],
      "additionalProperties": false
    },
    "Commitment": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "id": {
              "$ref": "#/$defs/CommitmentID"
            },
            "version": {
              "$ref": "#/$defs/Version"
            },
            "promisor": {
              "$ref": "#/$defs/EntityID"
            },
            "beneficiary": {
              "$ref": "#/$defs/EntityID"
            },
            "description": {
              "$ref": "#/$defs/Text"
            },
            "state": {
              "const": "open"
            },
            "evidence": {
              "type": "array",
              "items": {
                "$ref": "#/$defs/StoredEvidence"
              },
              "minItems": 1,
              "maxItems": 20
            },
            "known_from": {
              "$ref": "#/$defs/Instant"
            },
            "due_on": {
              "$ref": "#/$defs/Day"
            },
            "known_until": {
              "$ref": "#/$defs/Instant"
            }
          },
          "required": [
            "id",
            "version",
            "promisor",
            "beneficiary",
            "description",
            "state",
            "evidence",
            "known_from"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "id": {
              "$ref": "#/$defs/CommitmentID"
            },
            "version": {
              "$ref": "#/$defs/Version"
            },
            "promisor": {
              "$ref": "#/$defs/EntityID"
            },
            "beneficiary": {
              "$ref": "#/$defs/EntityID"
            },
            "description": {
              "$ref": "#/$defs/Text"
            },
            "state": {
              "const": "completed"
            },
            "evidence": {
              "type": "array",
              "items": {
                "$ref": "#/$defs/StoredEvidence"
              },
              "minItems": 1,
              "maxItems": 20
            },
            "known_from": {
              "$ref": "#/$defs/Instant"
            },
            "completed_at": {
              "$ref": "#/$defs/Instant"
            },
            "due_on": {
              "$ref": "#/$defs/Day"
            },
            "known_until": {
              "$ref": "#/$defs/Instant"
            }
          },
          "required": [
            "id",
            "version",
            "promisor",
            "beneficiary",
            "description",
            "state",
            "evidence",
            "known_from",
            "completed_at"
          ],
          "additionalProperties": false
        }
      ]
    },
    "Citation": {
      "type": "object",
      "properties": {
        "source": {
          "$ref": "#/$defs/SourceID"
        },
        "locator": {
          "$ref": "#/$defs/Label"
        },
        "span_start": {
          "type": "integer",
          "minimum": 0
        },
        "span_end": {
          "type": "integer",
          "minimum": 1
        },
        "passage": {
          "$ref": "#/$defs/Text"
        }
      },
      "required": [
        "source",
        "locator",
        "span_start",
        "span_end",
        "passage"
      ],
      "additionalProperties": false
    },
    "AnswerItem": {
      "type": "object",
      "properties": {
        "claim": {
          "$ref": "#/$defs/Text"
        },
        "entities": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/EntityID"
          },
          "minItems": 0,
          "maxItems": 20
        },
        "facts": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/FactID"
          },
          "minItems": 0,
          "maxItems": 20
        },
        "commitments": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/CommitmentID"
          },
          "minItems": 0,
          "maxItems": 20
        },
        "path": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/EntityID"
          },
          "minItems": 0,
          "maxItems": 4
        },
        "citations": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/Citation"
          },
          "minItems": 1,
          "maxItems": 20
        },
        "certainty": {
          "type": "string",
          "enum": [
            "supported",
            "conflicting",
            "unknown_speaker",
            "imprecise"
          ]
        }
      },
      "required": [
        "claim",
        "entities",
        "facts",
        "commitments",
        "path",
        "citations",
        "certainty"
      ],
      "additionalProperties": false
    },
    "Answer": {
      "type": "object",
      "properties": {
        "query_id": {
          "$ref": "#/$defs/QueryID"
        },
        "as_of": {
          "$ref": "#/$defs/Instant"
        },
        "valid_at": {
          "$ref": "#/$defs/Instant"
        },
        "known_at": {
          "$ref": "#/$defs/Instant"
        },
        "owner_timezone": {
          "$ref": "#/$defs/Label"
        },
        "completeness": {
          "type": "string",
          "enum": [
            "complete",
            "partial",
            "needs_clarification"
          ]
        },
        "reasons": {
          "type": "array",
          "items": {
            "type": "string",
            "enum": [
              "limit_reached",
              "dependency_unavailable",
              "source_unavailable",
              "ambiguous_identity",
              "unsupported_question"
            ]
          },
          "minItems": 0,
          "maxItems": 5
        },
        "summary": {
          "$ref": "#/$defs/Text"
        },
        "items": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/AnswerItem"
          },
          "minItems": 0,
          "maxItems": 20
        },
        "next_cursor": {
          "$ref": "#/$defs/Token"
        },
        "clarification": {
          "$ref": "#/$defs/Text"
        },
        "limits": {
          "$ref": "#/$defs/Limits"
        }
      },
      "required": [
        "query_id",
        "as_of",
        "valid_at",
        "known_at",
        "owner_timezone",
        "completeness",
        "reasons",
        "summary",
        "items",
        "limits"
      ],
      "additionalProperties": false
    },
    "View": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "self"
            },
            "owner": {
              "$ref": "#/$defs/EntityID"
            },
            "owner_timezone": {
              "$ref": "#/$defs/Label"
            },
            "contract": {
              "const": "1"
            },
            "revision": {
              "$ref": "#/$defs/Version"
            },
            "limits": {
              "$ref": "#/$defs/Limits"
            }
          },
          "required": [
            "type",
            "owner",
            "owner_timezone",
            "contract",
            "revision",
            "limits"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "entity"
            },
            "entity": {
              "$ref": "#/$defs/Entity"
            }
          },
          "required": [
            "type",
            "entity"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "claim"
            },
            "claim": {
              "$ref": "#/$defs/Claim"
            }
          },
          "required": [
            "type",
            "claim"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "commitment"
            },
            "commitment": {
              "$ref": "#/$defs/Commitment"
            }
          },
          "required": [
            "type",
            "commitment"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "source"
            },
            "id": {
              "$ref": "#/$defs/SourceID"
            },
            "version": {
              "$ref": "#/$defs/Version"
            },
            "mode": {
              "type": "string",
              "enum": [
                "owner_note",
                "authorized_copy",
                "linked"
              ]
            },
            "locator": {
              "$ref": "#/$defs/Label"
            },
            "recorded_at": {
              "$ref": "#/$defs/Instant"
            },
            "offset": {
              "type": "integer",
              "minimum": 0
            },
            "body": {
              "$ref": "#/$defs/Text"
            },
            "total_codepoints": {
              "type": "integer",
              "minimum": 1
            },
            "next_offset": {
              "type": "integer",
              "minimum": 0
            }
          },
          "required": [
            "type",
            "id",
            "version",
            "mode",
            "locator",
            "recorded_at",
            "offset",
            "body",
            "total_codepoints"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "type": {
              "const": "unavailable"
            },
            "ref": {
              "$ref": "#/$defs/Ref"
            },
            "state": {
              "type": "string",
              "enum": [
                "deleted",
                "ineligible"
              ]
            },
            "recorded_at": {
              "$ref": "#/$defs/Instant"
            },
            "deleted_at": {
              "$ref": "#/$defs/Instant"
            }
          },
          "required": [
            "type",
            "ref",
            "state",
            "recorded_at"
          ],
          "additionalProperties": false
        }
      ]
    },
    "Applying": {
      "type": "object",
      "properties": {
        "state": {
          "const": "applying"
        },
        "operation_id": {
          "$ref": "#/$defs/OpID"
        }
      },
      "required": [
        "state",
        "operation_id"
      ],
      "additionalProperties": false
    },
    "CommitSuccess": {
      "oneOf": [
        {
          "$ref": "#/$defs/Receipt"
        },
        {
          "$ref": "#/$defs/Applying"
        },
        {
          "$ref": "#/$defs/CommittedUnavailable"
        }
      ]
    },
    "Operation": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "state": {
              "const": "prepared"
            },
            "plan": {
              "$ref": "#/$defs/Plan"
            }
          },
          "required": [
            "state",
            "plan"
          ],
          "additionalProperties": false
        },
        {
          "$ref": "#/$defs/Applying"
        },
        {
          "type": "object",
          "properties": {
            "state": {
              "const": "committed"
            },
            "receipt": {
              "$ref": "#/$defs/Receipt"
            }
          },
          "required": [
            "state",
            "receipt"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "state": {
              "const": "aborted"
            },
            "operation_id": {
              "$ref": "#/$defs/OpID"
            }
          },
          "required": [
            "state",
            "operation_id"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "state": {
              "const": "unknown"
            },
            "operation_id": {
              "$ref": "#/$defs/OpID"
            }
          },
          "required": [
            "state",
            "operation_id"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "state": {
              "const": "failed"
            },
            "operation_id": {
              "$ref": "#/$defs/OpID"
            },
            "error": {
              "$ref": "#/$defs/Error"
            }
          },
          "required": [
            "state",
            "operation_id",
            "error"
          ],
          "additionalProperties": false
        },
        {
          "$ref": "#/$defs/CommittedUnavailable"
        }
      ]
    },
    "PrepareResult": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "kind": {
              "const": "ok"
            },
            "data": {
              "$ref": "#/$defs/Plan"
            }
          },
          "required": [
            "contract",
            "kind",
            "data"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "kind": {
              "const": "error"
            },
            "error": {
              "$ref": "#/$defs/Error"
            }
          },
          "required": [
            "contract",
            "kind",
            "error"
          ],
          "additionalProperties": false
        }
      ]
    },
    "CommitResult": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "kind": {
              "const": "ok"
            },
            "data": {
              "$ref": "#/$defs/CommitSuccess"
            }
          },
          "required": [
            "contract",
            "kind",
            "data"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "kind": {
              "const": "error"
            },
            "error": {
              "$ref": "#/$defs/Error"
            }
          },
          "required": [
            "contract",
            "kind",
            "error"
          ],
          "additionalProperties": false
        }
      ]
    },
    "QueryResult": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "kind": {
              "const": "ok"
            },
            "data": {
              "$ref": "#/$defs/Answer"
            }
          },
          "required": [
            "contract",
            "kind",
            "data"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "kind": {
              "const": "error"
            },
            "error": {
              "$ref": "#/$defs/Error"
            }
          },
          "required": [
            "contract",
            "kind",
            "error"
          ],
          "additionalProperties": false
        }
      ]
    },
    "ReadResult": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "kind": {
              "const": "ok"
            },
            "data": {
              "$ref": "#/$defs/View"
            }
          },
          "required": [
            "contract",
            "kind",
            "data"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "kind": {
              "const": "error"
            },
            "error": {
              "$ref": "#/$defs/Error"
            }
          },
          "required": [
            "contract",
            "kind",
            "error"
          ],
          "additionalProperties": false
        }
      ]
    },
    "StatusResult": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "kind": {
              "const": "ok"
            },
            "data": {
              "$ref": "#/$defs/Operation"
            }
          },
          "required": [
            "contract",
            "kind",
            "data"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "kind": {
              "const": "error"
            },
            "error": {
              "$ref": "#/$defs/Error"
            }
          },
          "required": [
            "contract",
            "kind",
            "error"
          ],
          "additionalProperties": false
        }
      ]
    },
    "AuthorizationRequirement": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "mode": {
              "const": "delegated"
            },
            "state": {
              "const": "not_required"
            }
          },
          "required": [
            "mode",
            "state"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "mode": {
              "const": "owner_terminal"
            },
            "state": {
              "const": "required"
            }
          },
          "required": [
            "mode",
            "state"
          ],
          "additionalProperties": false
        }
      ]
    },
    "Limits": {
      "type": "object",
      "const": {
        "max_entities": 10000,
        "max_assertion_versions": 50000,
        "max_eligible_fragments": 10000,
        "max_request_bytes": 65536,
        "max_response_bytes": 65536,
        "max_path_hops": 3,
        "max_relationship_versions": 1000,
        "max_returned_paths": 10,
        "max_hydrated_fragments": 20,
        "max_evidence_tokens": 8000,
        "query_timeout_ms": 10000,
        "prepare_timeout_ms": 30000,
        "commit_timeout_ms": 10000
      }
    },
    "CommittedUnavailable": {
      "type": "object",
      "properties": {
        "state": {
          "const": "committed_unavailable"
        },
        "operation_id": {
          "$ref": "#/$defs/OpID"
        },
        "committed_at": {
          "$ref": "#/$defs/Instant"
        }
      },
      "required": [
        "state",
        "operation_id",
        "committed_at"
      ],
      "additionalProperties": false
    },
    "GrantVersionToken": {
      "type": "string",
      "minLength": 1,
      "maxLength": 128,
      "pattern": "^[A-Za-z0-9_-]{1,128}$"
    },
    "Signature": {
      "type": "string",
      "minLength": 86,
      "maxLength": 86,
      "pattern": "^[A-Za-z0-9_-]{86}$"
    },
    "GrantPayload": {
      "type": "object",
      "properties": {
        "version": {
          "const": "1"
        },
        "installation": {
          "$ref": "#/$defs/Token"
        },
        "client": {
          "$ref": "#/$defs/Token"
        },
        "operation_id": {
          "$ref": "#/$defs/OpID"
        },
        "plan_id": {
          "$ref": "#/$defs/PlanID"
        },
        "plan_digest": {
          "$ref": "#/$defs/Token"
        },
        "action": {
          "const": "delete_source"
        },
        "target": {
          "$ref": "#/$defs/SourceID"
        },
        "authority_epoch": {
          "$ref": "#/$defs/GrantVersionToken"
        },
        "key_id": {
          "$ref": "#/$defs/Token"
        },
        "nonce": {
          "$ref": "#/$defs/Token"
        },
        "issued_at": {
          "$ref": "#/$defs/Instant"
        },
        "expires_at": {
          "$ref": "#/$defs/Instant"
        },
        "session": {
          "$ref": "#/$defs/Token"
        },
        "base_revision": {
          "$ref": "#/$defs/GrantVersionToken"
        }
      },
      "required": [
        "version",
        "installation",
        "client",
        "operation_id",
        "plan_id",
        "plan_digest",
        "action",
        "target",
        "authority_epoch",
        "key_id",
        "nonce",
        "issued_at",
        "expires_at",
        "session",
        "base_revision"
      ],
      "additionalProperties": false
    },
    "TerminalGrant": {
      "type": "object",
      "properties": {
        "payload": {
          "$ref": "#/$defs/GrantPayload"
        },
        "signature": {
          "$ref": "#/$defs/Signature"
        }
      },
      "required": [
        "payload",
        "signature"
      ],
      "additionalProperties": false
    }
  }
}
```

## Complete synthetic shape examples

These are schema instances, not executed GRM outcomes. Baseline revision 20d061ee-ed64-4f9b-8ecd-ae2854a69219 and the proposal's A01–A13/R1–R6 table still govern meaning. Operator-grant example does not authorize any action and is never submitted to MCP.

Definition **PrepareInput**:
```json
{
  "contract": "1",
  "operation_id": "op_rename",
  "intent": "rename",
  "text": "Rename M1 to Morgan Vale",
  "targets": [
    "ent_M1"
  ]
}
```

Definition **CommitInput**:
```json
{
  "contract": "1",
  "operation_id": "op_rename",
  "plan_id": "plan_rename",
  "plan_digest": "fixture-rename",
  "choices": [],
  "confirm": true
}
```

Definition **QueryInput**:
```json
{
  "contract": "1",
  "question": "What role was supported on September 10?",
  "valid_at": "2026-09-10T12:00:00Z",
  "known_at": "2026-10-21T12:00:00Z"
}
```

Definition **ReadInput**:
```json
{
  "contract": "1",
  "ref": "src_S1",
  "offset": 0,
  "length": 1000
}
```

Definition **StatusInput**:
```json
{
  "contract": "1",
  "operation_id": "op_rename"
}
```

Definition **PrepareResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "operation_id": "op_rename",
    "plan_id": "plan_rename",
    "plan_digest": "fixture-rename",
    "expires_at": "2026-10-03T12:00:00Z",
    "base_revision": "r2",
    "identities": [],
    "actions": [
      {
        "type": "rename",
        "entity": "ent_M1",
        "new_name": "Morgan Vale"
      }
    ],
    "preconditions": [
      {
        "ref": "ent_M1",
        "version": "v1"
      }
    ],
    "summary": "Rename M1 to Morgan Vale; retain Morgan Lee as prior name.",
    "authorization": {
      "mode": "delegated",
      "state": "not_required"
    }
  }
}
```

Definition **CommitResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "operation_id": "op_rename",
    "plan_id": "plan_rename",
    "state": "committed",
    "committed_at": "2026-10-02T12:00:00Z",
    "revision": "r3",
    "changes": [
      {
        "ref": "ent_M1",
        "effect": "renamed"
      }
    ],
    "id_map": {},
    "replayed": false
  }
}
```

Definition **QueryResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "query_id": "qry_role",
    "as_of": "2026-10-21T12:00:00Z",
    "valid_at": "2026-09-10T12:00:00Z",
    "known_at": "2026-10-21T12:00:00Z",
    "owner_timezone": "America/New_York",
    "completeness": "complete",
    "reasons": [],
    "summary": "No supported Cedar Lab role on September 10 after S3 shifted the start to September 15. No earlier employer is supplied.",
    "items": [],
    "limits": {
      "max_entities": 10000,
      "max_assertion_versions": 50000,
      "max_eligible_fragments": 10000,
      "max_request_bytes": 65536,
      "max_response_bytes": 65536,
      "max_path_hops": 3,
      "max_relationship_versions": 1000,
      "max_returned_paths": 10,
      "max_hydrated_fragments": 20,
      "max_evidence_tokens": 8000,
      "query_timeout_ms": 10000,
      "prepare_timeout_ms": 30000,
      "commit_timeout_ms": 10000
    }
  }
}
```

Definition **ReadResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "type": "unavailable",
    "ref": "src_S7",
    "state": "deleted",
    "recorded_at": "2026-10-02T12:00:00Z",
    "deleted_at": "2026-10-21T12:00:00Z"
  }
}
```

Definition **StatusResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "state": "applying",
    "operation_id": "op_rename"
  }
}
```

Definition **CommitResult**:
```json
{
  "contract": "1",
  "kind": "error",
  "error": {
    "code": "REPLAY_CONFLICT",
    "message": "Operation ID already binds a different reviewed request.",
    "retry": "new_request",
    "operation_id": "op_rename"
  }
}
```

Definition **CommitResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "state": "applying",
    "operation_id": "op_rename"
  }
}
```

Definition **PrepareResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "operation_id": "op_delete_S7",
    "plan_id": "plan_delete_S7",
    "plan_digest": "fixture_delete_S7",
    "expires_at": "2026-10-22T12:00:00Z",
    "base_revision": "r20",
    "identities": [],
    "actions": [
      {
        "type": "delete_source",
        "source": "src_S7"
      }
    ],
    "preconditions": [
      {
        "ref": "src_S7",
        "version": "v3"
      }
    ],
    "summary": "Delete the selected source and all covered private derivatives, including historical routes.",
    "authorization": {
      "mode": "owner_terminal",
      "state": "required"
    }
  }
}
```

Definition **CommitInput**:
```json
{
  "contract": "1",
  "operation_id": "op_delete_S7",
  "plan_id": "plan_delete_S7",
  "plan_digest": "fixture_delete_S7",
  "choices": [],
  "confirm": true
}
```

Definition **CommitResult**:
```json
{
  "contract": "1",
  "kind": "error",
  "error": {
    "code": "OWNER_GRANT_REQUIRED",
    "message": "Obtain an exact local owner-terminal grant, then retry the unchanged commit.",
    "retry": "same_request",
    "operation_id": "op_delete_S7"
  }
}
```

Definition **CommitResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "state": "committed_unavailable",
    "operation_id": "op_old_capture",
    "committed_at": "2026-10-02T12:00:00Z"
  }
}
```

Definition **StatusResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "state": "committed_unavailable",
    "operation_id": "op_old_capture",
    "committed_at": "2026-10-02T12:00:00Z"
  }
}
```

Definition **ReadResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "type": "self",
    "owner": "ent_O",
    "owner_timezone": "America/New_York",
    "contract": "1",
    "revision": "r20",
    "limits": {
      "max_entities": 10000,
      "max_assertion_versions": 50000,
      "max_eligible_fragments": 10000,
      "max_request_bytes": 65536,
      "max_response_bytes": 65536,
      "max_path_hops": 3,
      "max_relationship_versions": 1000,
      "max_returned_paths": 10,
      "max_hydrated_fragments": 20,
      "max_evidence_tokens": 8000,
      "query_timeout_ms": 10000,
      "prepare_timeout_ms": 30000,
      "commit_timeout_ms": 10000
    }
  }
}
```

Definition **TerminalGrant**:
```json
{
  "payload": {
    "version": "1",
    "installation": "installation_fixture",
    "client": "client_fixture",
    "operation_id": "op_delete_S7",
    "plan_id": "plan_delete_S7",
    "plan_digest": "fixture_delete_S7",
    "action": "delete_source",
    "target": "src_S7",
    "authority_epoch": "epoch_2",
    "key_id": "owner_key_1",
    "nonce": "nonce_fixture",
    "issued_at": "2026-10-21T12:00:00Z",
    "expires_at": "2026-10-21T12:05:00Z",
    "session": "fixture_session_1",
    "base_revision": "fixture_revision_1"
  },
  "signature": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
}
```

## Recorded refinement checks and limits

Later integration verification is recorded in [schema verification](schema-verification.md): full 2020-12 checks with active date-time format assertion now ran for the current registry. The earlier dependency-missing statement below is historical, not the current tooling status. Any next schema revision requires revalidation.

- Registry JSON serialization/parse round-trip: passed; 66 definitions, 195 local references, zero unresolved.
- Minimal input/output materialization: 41725 compact UTF-8 schema bytes across ten objects; actual token fit in the selected host is not measured.
- Bounded in-memory structural probes over the used keyword subset: **19 positive instances accepted, 13 invalid instances rejected**. Negatives include forged MCP owner_grant, confirm:false, mixed cursor/question, malformed reference, null OpID, extra input field, missing Plan authorization, model-approved authorization state, four-hop limit substitution/five-entity path, private committed_unavailable/deleted-view payload and malformed grant signature length.
- These custom probes are not JSON Schema meta-schema validation, signature verification or provider/host conformance. The prior Python jsonschema dependency was missing; no validator package was installed or repository/environment change made. Full standards validation and domain/crypto/host/runtime tests remain not run.
- Required later tests include absent/expired/mismatched/revoked owner grant, changed digest/target/client/epoch, replay and deletion-source fingerprint purge, second-process lock rejection, restart and old-backup grant rejection. Preserve exact canonical semantic assertions.