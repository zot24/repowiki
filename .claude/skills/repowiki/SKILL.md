---
name: repowiki
description: Consult and maintain this repository's persistent project memory in .llm-wiki/ using the repowiki CLI. Use when answering questions about this project's architecture, decisions, plans, or history; when the user makes a decision or shares a durable insight worth remembering; when asked to "check the wiki", "remember this", or capture a session digest; and at the end of substantial work sessions to record what was learned.
---

# RepoWiki — project memory

This repository keeps a persistent, Git-versioned knowledge base in `.llm-wiki/`
(OKF v0.2 bundle) maintained with the `repowiki` CLI.

## Consult before answering

For any question about this project's architecture, decisions, plans, or history:

1. `repowiki status` — what exists
2. Read `.llm-wiki/index.md` — map of compiled knowledge
3. `repowiki search <term>` then `repowiki get <path>` — drill in
4. Conventions live in `.llm-wiki/SCHEMA.md`

Prefer wiki answers over guessing from code alone; the wiki records *why*.

## Capture when you learn something durable

Decisions, architectural insights, gotchas, plans — not routine chatter:

    echo "..." | repowiki ingest - --type decision --title "Chose X over Y" --tags a,b
    repowiki ingest path/to/notes.md --type note
    repowiki index    # always regenerate after ingesting

Types: `decision` | `note` | `session` | `source`.

## End-of-session digest

After substantial work, ingest a **redacted digest** (never a raw transcript;
never secrets, tokens, or PII):

    echo "# Session digest — <topic>
    - what was built/changed and why
    - decisions made, alternatives rejected
    - gotchas discovered
    - open questions" | repowiki ingest - --type session --title "<topic> session"

## Rules

- `raw/` is append-only — never rewrite it. Synthesize into `pages/` instead.
- Compiled pages cite their `sources` from `raw/`.
- Prefer updating an existing page over creating a near-duplicate.
- Run `repowiki lint` before committing wiki changes.
- Full command reference: `repowiki --help` or docs/CLI.md in the repowiki repo.
