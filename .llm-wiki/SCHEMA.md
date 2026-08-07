---
type: concept
title: "SCHEMA — RepoWiki protocol"
description: "How agents and humans read and write this wiki"
status: stable
okf_version: "0.2"
---

# SCHEMA — how this wiki works

This directory is an OKF v0.2 knowledge bundle managed with the `repowiki` CLI.

## Layers

| Layer | Path | Rule |
|-------|------|------|
| Raw | `raw/` | Immutable sources: session digests, notes, decisions, external docs. Append-only — never rewrite. |
| Compiled | `pages/` | Synthesized, interlinked knowledge. Freely edited and refactored. |
| Protocol | `SCHEMA.md` | This file. Conventions and project-specific rules. |
| Inbox | `inbox/` | Drop zone for quick captures awaiting ingestion. |

## Frontmatter

Every document needs YAML frontmatter with at least `type` and `title`:

- `type`: decision | architecture | plan | concept | component | overview | open-question | session | note | source
- `title`, `description`, `tags`
- `status`: draft | stable | deprecated
- `sources`: provenance pointers back to `raw/` files
- `created` / `updated` (YYYY-MM-DD), `confidence`, `related`

## Rules for agents

1. **Consult before answering.** Check `index.md` and relevant `pages/` first.
2. **Capture intentionally.** Redacted digests, decisions, and insights — not raw chat dumps.
3. **Cite provenance.** Compiled pages list their `sources` from `raw/`.
4. **Never rewrite `raw/`.** Promote and synthesize into `pages/` instead.
5. **Keep structure via the CLI.** `repowiki ingest`, `repowiki index`, `repowiki lint`.
6. **No secrets.** Redact credentials, tokens, PII before anything lands here.

## Project-specific conventions

_Add project conventions below this line._
