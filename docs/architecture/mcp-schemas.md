# GRM MCP v1 consolidated schema registry

Proposed governing architecture; not implementation approval. The [corrected MCP contract](mcp-corrections.md) and [governing integration](mvp.md) define semantics. Five tools remain. This registry composes the reviewed source schemas with GRO-26 corrective-contract-supplement revision530b03fa-ddc6-4458-aa4c-52cd17f8b83a; it replaces the old registry, not an implementer-selected alternative.

JSON Schema2020-12. Each input/output uses a self-contained transitive $defs closure rooted at its named definition; never use the registry alone as permissive tool input. All objects closed unless explicitly stated. Enable actual date/date-time format assertion; check checker availability. Shape validation does not establish source support, signatures, authenticated issuance, valid identity/scope/time, control freshness, global budgets or domain effects.

| Tool | Input | Output |
| --- | --- | --- |
| grm_prepare_change_v1 | PrepareInput | PrepareResult |
| grm_commit_change_v1 | CommitInput | CommitResult |
| grm_query_v1 | QueryInput | QueryResult |
| grm_read_v1 | ReadInput | ReadResult |
| grm_operation_status_v1 | StatusInput | StatusResult |

TerminalGrant remains a local operator artifact outside every MCP input closure. Signed session/base fields remain required. OpID cryptographic issuance is a domain check, not inferred from its pattern. current_index_revision null is the sole documented readiness exception to no-null in the prior contract. Export known_order and exact FactVersionRef types compose with the closed native bundle below.

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
      "minLength": 4,
      "maxLength": 83,
      "pattern": "^op_[A-Za-z0-9_-]{1,80}$"
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
          "$ref": "#/$defs/FactVersionRef"
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
            },
            "evidence": {
              "type": "array",
              "items": {
                "$ref": "#/$defs/Evidence"
              },
              "minItems": 1,
              "maxItems": 20
            }
          },
          "required": [
            "type",
            "entity",
            "new_name",
            "evidence"
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
            },
            "evidence": {
              "type": "array",
              "items": {
                "$ref": "#/$defs/Evidence"
              },
              "minItems": 1,
              "maxItems": 20
            }
          },
          "required": [
            "type",
            "commitment",
            "evidence"
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
      "oneOf": [
        {
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
                "complete"
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
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "operation_id": {
              "$ref": "#/$defs/OpID"
            },
            "intent": {
              "const": "delete"
            },
            "targets": {
              "type": "array",
              "items": {
                "$ref": "#/$defs/SourceID"
              },
              "minItems": 1,
              "maxItems": 1
            }
          },
          "required": [
            "contract",
            "operation_id",
            "intent",
            "targets"
          ],
          "additionalProperties": false
        }
      ]
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
        },
        "version": {
          "$ref": "#/$defs/Version"
        }
      },
      "required": [
        "contract",
        "ref"
      ],
      "additionalProperties": false
    },
    "StatusInput": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "contract": {
              "const": "1"
            },
            "allocate": {
              "const": true
            }
          },
          "required": [
            "contract",
            "allocate"
          ],
          "additionalProperties": false
        },
        {
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
        }
      ]
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
        "INSTALLATION_BUSY",
        "OPERATION_ID_INVALID",
        "OPERATION_ID_EXPIRED"
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
        },
        "creation_evidence": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/Evidence"
          },
          "minItems": 1,
          "maxItems": 20
        }
      },
      "required": [
        "slot",
        "mention",
        "kind",
        "candidates",
        "allow_new"
      ],
      "additionalProperties": false,
      "allOf": [
        {
          "if": {
            "properties": {
              "allow_new": {
                "const": true
              }
            },
            "required": [
              "allow_new"
            ]
          },
          "then": {
            "required": [
              "creation_evidence"
            ]
          }
        }
      ]
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
          "$ref": "#/$defs/GrantVersionToken"
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
        },
        "instruction_source": {
          "oneOf": [
            {
              "$ref": "#/$defs/SourceID"
            },
            {
              "$ref": "#/$defs/Temp"
            }
          ]
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
        },
        "semantic_readiness": {
          "type": "string",
          "enum": [
            "pending",
            "ready",
            "unavailable"
          ]
        },
        "current_index_revision": {
          "oneOf": [
            {
              "$ref": "#/$defs/Version"
            },
            {
              "type": "null"
            }
          ]
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
        "replayed",
        "semantic_readiness",
        "current_index_revision"
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
          "$ref": "#/$defs/FactVersionRef"
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
            "$ref": "#/$defs/FactVersionRef"
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
        },
        {
          "$ref": "#/$defs/Allocated"
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
    },
    "Allocated": {
      "type": "object",
      "properties": {
        "state": {
          "const": "allocated"
        },
        "operation_id": {
          "$ref": "#/$defs/OpID"
        },
        "expires_at": {
          "$ref": "#/$defs/Instant"
        }
      },
      "required": [
        "state",
        "operation_id",
        "expires_at"
      ],
      "additionalProperties": false
    },
    "FactVersionRef": {
      "type": "object",
      "properties": {
        "id": {
          "$ref": "#/$defs/FactID"
        },
        "version": {
          "$ref": "#/$defs/Version"
        }
      },
      "required": [
        "id",
        "version"
      ],
      "additionalProperties": false
    },
    "Sequence": {
      "type": "integer",
      "minimum": 0,
      "maximum": 9007199254740991
    },
    "KnownOrder": {
      "type": "object",
      "properties": {
        "from": {
          "$ref": "#/$defs/Sequence"
        },
        "until": {
          "$ref": "#/$defs/Sequence"
        }
      },
      "required": [
        "from"
      ],
      "additionalProperties": false
    },
    "ExportEntityVersion": {
      "type": "object",
      "properties": {
        "record": {
          "const": "entity_version"
        },
        "entity": {
          "$ref": "#/$defs/Entity"
        },
        "known_from": {
          "$ref": "#/$defs/Instant"
        },
        "known_order": {
          "$ref": "#/$defs/KnownOrder"
        },
        "evidence": {
          "type": "array",
          "items": {
            "$ref": "#/$defs/StoredEvidence"
          },
          "minItems": 0,
          "maxItems": 20
        },
        "known_until": {
          "$ref": "#/$defs/Instant"
        }
      },
      "required": [
        "record",
        "entity",
        "known_from",
        "known_order",
        "evidence"
      ],
      "additionalProperties": false
    },
    "ExportClaimVersion": {
      "type": "object",
      "properties": {
        "record": {
          "const": "claim_version"
        },
        "claim": {
          "$ref": "#/$defs/Claim"
        },
        "known_order": {
          "$ref": "#/$defs/KnownOrder"
        }
      },
      "required": [
        "record",
        "claim",
        "known_order"
      ],
      "additionalProperties": false
    },
    "ExportCommitmentVersion": {
      "type": "object",
      "properties": {
        "record": {
          "const": "commitment_version"
        },
        "commitment": {
          "$ref": "#/$defs/Commitment"
        },
        "known_order": {
          "$ref": "#/$defs/KnownOrder"
        }
      },
      "required": [
        "record",
        "commitment",
        "known_order"
      ],
      "additionalProperties": false
    },
    "ExportUnavailableVersion": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "record": {
              "const": "unavailable_version"
            },
            "ref": {
              "oneOf": [
                {
                  "$ref": "#/$defs/EntityID"
                },
                {
                  "$ref": "#/$defs/FactID"
                },
                {
                  "$ref": "#/$defs/CommitmentID"
                }
              ]
            },
            "version": {
              "$ref": "#/$defs/Version"
            },
            "state": {
              "const": "deleted"
            },
            "recorded_at": {
              "$ref": "#/$defs/Instant"
            },
            "deleted_at": {
              "$ref": "#/$defs/Instant"
            }
          },
          "required": [
            "record",
            "ref",
            "version",
            "state",
            "recorded_at",
            "deleted_at"
          ],
          "additionalProperties": false
        },
        {
          "type": "object",
          "properties": {
            "record": {
              "const": "unavailable_version"
            },
            "ref": {
              "oneOf": [
                {
                  "$ref": "#/$defs/EntityID"
                },
                {
                  "$ref": "#/$defs/FactID"
                },
                {
                  "$ref": "#/$defs/CommitmentID"
                }
              ]
            },
            "version": {
              "$ref": "#/$defs/Version"
            },
            "state": {
              "const": "ineligible"
            },
            "recorded_at": {
              "$ref": "#/$defs/Instant"
            }
          },
          "required": [
            "record",
            "ref",
            "version",
            "state",
            "recorded_at"
          ],
          "additionalProperties": false
        }
      ]
    }
  }
}
```

## Complete synthetic shape examples

These are full schema instances, not executed product results. All-A operation tokens/signatures are deliberately unauthentic shape fixtures; real domain fixtures obtain authenticated server-issued IDs with fixed synthetic clock/key/RNG. They are never valid authority by shape alone.

Definition **StatusInput**:
```json
{
  "contract": "1",
  "allocate": true
}
```

Definition **StatusResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "state": "allocated",
    "operation_id": "op_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
    "expires_at": "2026-10-02T12:15:00Z"
  }
}
```

Definition **StatusInput**:
```json
{
  "contract": "1",
  "operation_id": "op_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
}
```

Definition **PrepareInput**:
```json
{
  "contract": "1",
  "operation_id": "op_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
  "intent": "rename",
  "text": "Rename Morgan Lee to Morgan Vale.",
  "targets": [
    "ent_M1"
  ]
}
```

Definition **PrepareResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "operation_id": "op_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
    "plan_id": "plan_rename",
    "plan_digest": "fixture_digest",
    "expires_at": "2026-10-03T12:00:00Z",
    "base_revision": "revision_10",
    "identities": [],
    "actions": [
      {
        "type": "add_source",
        "source": {
          "id": "tmp_owner_instruction",
          "mode": "owner_note",
          "locator": "owner instruction",
          "body": "Rename Morgan Lee to Morgan Vale."
        }
      },
      {
        "type": "rename",
        "entity": "ent_M1",
        "new_name": "Morgan Vale",
        "evidence": [
          {
            "source": "tmp_owner_instruction",
            "span_start": 0,
            "span_end": 33
          }
        ]
      }
    ],
    "preconditions": [
      {
        "ref": "ent_M1",
        "version": "v1"
      }
    ],
    "summary": "Rename selected M1 and retain its former alias.",
    "authorization": {
      "mode": "delegated",
      "state": "not_required"
    },
    "instruction_source": "tmp_owner_instruction"
  }
}
```

Definition **CommitInput**:
```json
{
  "contract": "1",
  "operation_id": "op_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
  "plan_id": "plan_rename",
  "plan_digest": "fixture_digest",
  "choices": [],
  "confirm": true
}
```

Definition **CommitResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "operation_id": "op_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
    "plan_id": "plan_rename",
    "state": "committed",
    "committed_at": "2026-10-02T12:00:00Z",
    "revision": "revision_11",
    "changes": [
      {
        "ref": "src_owner_instruction",
        "effect": "created"
      },
      {
        "ref": "ent_M1",
        "effect": "renamed"
      }
    ],
    "id_map": {
      "tmp_owner_instruction": "src_owner_instruction"
    },
    "replayed": false,
    "semantic_readiness": "pending",
    "current_index_revision": "revision_10"
  }
}
```

Definition **StatusResult**:
```json
{
  "contract": "1",
  "kind": "ok",
  "data": {
    "state": "committed",
    "receipt": {
      "operation_id": "op_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
      "plan_id": "plan_rename",
      "state": "committed",
      "committed_at": "2026-10-02T12:00:00Z",
      "revision": "revision_11",
      "changes": [
        {
          "ref": "src_owner_instruction",
          "effect": "created"
        },
        {
          "ref": "ent_M1",
          "effect": "renamed"
        }
      ],
      "id_map": {
        "tmp_owner_instruction": "src_owner_instruction"
      },
      "replayed": true,
      "semantic_readiness": "ready",
      "current_index_revision": "revision_12"
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
    "operation_id": "op_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
    "plan_id": "plan_rename",
    "state": "committed",
    "committed_at": "2026-10-02T12:00:00Z",
    "revision": "revision_11",
    "changes": [
      {
        "ref": "src_owner_instruction",
        "effect": "created"
      },
      {
        "ref": "ent_M1",
        "effect": "renamed"
      }
    ],
    "id_map": {
      "tmp_owner_instruction": "src_owner_instruction"
    },
    "replayed": true,
    "semantic_readiness": "unavailable",
    "current_index_revision": null
  }
}
```

Definition **PrepareInput**:
```json
{
  "contract": "1",
  "operation_id": "op_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
  "intent": "delete",
  "targets": [
    "src_S7"
  ]
}
```

Definition **ReadInput**:
```json
{
  "contract": "1",
  "ref": "fact_S2",
  "version": "v1"
}
```

Definition **FactVersionRef**:
```json
{
  "id": "fact_S2",
  "version": "v1"
}
```

Definition **ExportClaimVersion**:
```json
{
  "record": "claim_version",
  "claim": {
    "id": "fact_S2",
    "version": "v1",
    "predicate": "role_at",
    "participants": [
      {
        "role": "subject",
        "entity": "ent_M1"
      },
      {
        "role": "organization",
        "entity": "ent_Cedar"
      }
    ],
    "statement": "M1 researcher at Cedar Lab from September 1.",
    "valid": {
      "from": {
        "precision": "day",
        "value": "2026-09-01"
      }
    },
    "known_from": "2026-10-02T12:00:00Z",
    "known_until": "2026-10-02T12:00:00Z",
    "evidence": [
      {
        "source": "src_S2",
        "span_start": 0,
        "span_end": 1
      }
    ]
  },
  "known_order": {
    "from": 10,
    "until": 11
  }
}
```

Definition **ExportUnavailableVersion**:
```json
{
  "record": "unavailable_version",
  "ref": "fact_S2",
  "version": "v1",
  "state": "deleted",
  "recorded_at": "2026-10-02T12:00:00Z",
  "deleted_at": "2026-10-21T12:00:00Z"
}
```

Definition **StatusResult**:
```json
{
  "contract": "1",
  "kind": "error",
  "error": {
    "code": "OPERATION_ID_INVALID",
    "message": "Operation identity is invalid.",
    "retry": "new_request"
  }
}
```

Definition **StatusResult**:
```json
{
  "contract": "1",
  "kind": "error",
  "error": {
    "code": "REFERENCE_UNAVAILABLE",
    "message": "Reserved operation outcome is unavailable; no replay was executed.",
    "retry": "never",
    "operation_id": "op_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
  }
}
```

## Verification

See [schema verification](schema-verification.md). Recompute checks on this consolidated registry before final freeze. No implementation/runtime evidence is claimed.
