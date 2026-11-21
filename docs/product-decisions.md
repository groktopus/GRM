# GRM product and domain decisions

These are the choices that shape GRM. They describe what the product must represent and how the main assistant uses it, not how to implement it. [ADR 0001](adr/0001-choose-go-for-backend-services.md) separately records the Go decision.

## Product boundary

GRM is a self-hosted personal CRM for one owner. It is not a social network or a separate assistant the owner has to talk to. The owner's main assistant is the front door. GRM helps it remember people, relationships, and events, and surface useful context. The owner decides whether to reach out.

## Interface

**GRM exposes an MCP server with multiple tools. Natural language is the primary interface between the main assistant and GRM, for both supplying information and asking questions.** The main assistant can tell GRM about a conversation or ask who it knows in a particular area without deciding whether the request is a semantic search, a graph traversal, a timeline lookup, or a combination. GRM handles that translation internally.

Use evidence associated with people and relationships to answer questions, not only explicit contact fields. For example, a question about who has discussed a topic can use the relevant part of a conversation, not merely a keyword in a person's profile. The exact MCP tools, schemas, and evidence storage layout are not decided.

## Entities, identity, and relationships

Use *entity* and *relationship* in product-facing language. *Node* and *edge* are implementation terms. GRM represents people, places, organizations, things, and events. A thing may be a product or an individual item, such as the particular copy of a book someone lent.

The owner is a person in the graph, not a special case outside it. A person's identifier is stable and independent of their name. Two people can share a name and a person can change names. Graph, semantic, and temporal records that refer to the same person use the same identifier. An update addressed to an existing person's identifier updates that record rather than creating another person. This does not specify the ID format or how ambiguous identities are resolved.

Relationships can connect people and other entities and carry information of their own. The graph can represent people the owner knows directly, contacts reached through other people, and people the owner is tracking without an established relationship. Attendance and interaction are separate: four people might attend an event, while only two are known to have met there. GRM should also be able to represent what happened, such as a book being lent and later returned, without treating that as evidence that everyone at the event participated.

## Time and evidence

Keep **when something happened or was true** distinct from **when it became known**. If GRM learns a month later about an interaction at an event, today's account of the event can include it, while an answer about what was known the next day cannot. The same principle applies to changing relationships and facts about people.

Ordinary responses can indicate that unrelated updates are waiting and let the main assistant retrieve them separately. Public monitoring through GroktoCrawl is optional, not required to use GRM. A public observation should be tied to a verified identity before it updates that person's record.

## Workflow terminology

**Workflow** is the canonical term for intentionally developing a relationship through stages. *Pipeline* and *funnel* are synonyms GRM can understand, not names for separate features. A workflow can serve a sales or a personal objective. Stage describes progress within that workflow, not an overall status of the person. One person can be relevant to several objectives; this does not require multiple person records.

This record does not settle ontology modularity, storage engines, MCP signatures, or internal agent topology. It does not authorize implementation under the [architecture gate](../CONTRIBUTING.md#architecture-first).
