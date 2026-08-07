# AGENTS

<!-- repowiki:begin — managed by `repowiki init`, do not edit inside this block -->

## Project memory (.llm-wiki)

This repository has a persistent knowledge base in `.llm-wiki/` maintained with the `repowiki` CLI.

**Before answering questions about this project**, consult:
- `.llm-wiki/index.md` — map of all compiled knowledge
- `.llm-wiki/pages/` — architecture, decisions, plans, components, concepts
- `.llm-wiki/SCHEMA.md` — the protocol for reading and writing wiki content

**When you learn something durable** (a decision, an architectural insight, a gotcha,
a plan), capture it:
- Quick capture: drop a markdown note in `.llm-wiki/inbox/`
- Structured capture: `repowiki ingest <file> --type decision|note|session|source`
- Then keep the index fresh: `repowiki index`

Follow the conventions in `.llm-wiki/SCHEMA.md`. Prefer updating an existing page
over creating a near-duplicate.

<!-- repowiki:end -->
