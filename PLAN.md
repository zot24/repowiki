# RepoWiki — Full Project Plan

**Last updated:** 2026-08-06

## Vision

Every Git repository gets a persistent, versioned, agent-maintained knowledge base that lives *with the code*.

- Knowledge compounds across sessions, tools, and people.
- Agents automatically consult it when talking about the project.
- New decisions, planning insights, architecture notes, and session digests are easy to capture.
- A small, reliable CLI binary (`repowiki`) makes structural operations deterministic.
- The knowledge layer is **OKF-compatible**, so it is portable across agents, tools, and organizations.
- Session capture is automatic (or near-automatic) via platform hooks/plugins so the user barely notices it.

This is the “project memory” layer that pure RAG or chat history cannot provide.

## Core Principles

1. **Repo-local & versioned** — The wiki travels with `git clone`.
2. **Agent-native** — Any modern coding agent that respects `AGENTS.md` / `CLAUDE.md` becomes a wiki maintainer.
3. **Karpathy three layers, project-focused**:
   - `raw/` — immutable sources (sessions, notes, decisions, external docs)
   - `pages/` — compiled, interlinked, synthesized knowledge
   - `SCHEMA.md` — the protocol + project-specific conventions
4. **OKF-compatible by design** — Compiled knowledge forms a valid Open Knowledge Format (OKF) v0.2 bundle (Google).
5. **CLI for reliability, LLM for intelligence** — Binary handles structure, provenance, indexing, linting. Agent does smart synthesis.
6. **Session-first + automatic capture** — Redacted digests captured via hooks/plugins. Promotion to durable pages remains intentional/high-quality.
7. **Lightweight** — No global hub, no mandatory multi-agent research swarms, no required Obsidian, no servers.

## Directory Structure

```
.llm-wiki/                          # OKF knowledge bundle root (committed by default)
├── index.md                        # OKF root index (contains okf_version: "0.2")
├── log.md                          # OKF chronological activity log
├── SCHEMA.md                       # Agent protocol + project conventions
├── config.md                       # Optional project overrides
├── raw/                            # Karpathy-style immutable layer
│   ├── sessions/                   # Redacted digests from coding/planning sessions
│   ├── notes/
│   ├── sources/
│   └── decisions/
├── pages/                          # Compiled OKF concept documents
│   ├── overview.md
│   ├── architecture/
│   ├── decisions/
│   ├── plans/
│   ├── components/
│   ├── concepts/
│   └── open-questions.md
└── inbox/                          # Drop zone for quick captures
```

## Frontmatter (OKF-superset)

```yaml
---
# OKF required
type: decision                  # decision | architecture | plan | concept | component |
                                # overview | open-question | session | note | source | ...

# OKF recommended
title: "Authentication Decision"
description: "We chose JWT + refresh tokens for the public API"
tags: [auth, security, backend]
resource: null

# OKF provenance / trust / lifecycle
sources:
  - resource: "raw/sessions/2026-07-28-auth.md"
    title: "Session digest"
generated:
  by: "agent:claude-code/repowiki"
  at: "2026-07-28T20:15:00-03:00"
status: stable                  # draft | stable | deprecated

# Our useful extensions
created: 2026-07-28
updated: 2026-07-28
confidence: high
related:
  - "/pages/architecture/api.md"
---
```

Root `index.md` includes `okf_version: "0.2"`.

## CLI (`repowiki`)

Written in **Go** (Cobra) as a single static binary.

MVP commands:
- `repowiki init`
- `repowiki status`
- `repowiki ingest <file|url|-> [--type session|decision|note|source|auto]`
- `repowiki index`
- `repowiki lint [--fix]` (includes basic OKF conformance)
- `repowiki search` / `repowiki query`
- `repowiki get`
- `repowiki log`
- Session helpers

## Agent Integration & Automatic Session Capture

**Goal experience:** After `repowiki init` + one-time plugin install, knowledge accumulates almost invisibly.

### Universal core (works for every agent)
- SCHEMA.md + AGENTS.md / CLAUDE.md / .cursorrules injection
- Easy CLI ingest

### Platform plugins & hooks (seamless part)

| Priority | Agent          | Approach                                      | Capture level |
|----------|----------------|-----------------------------------------------|---------------|
| 1        | Claude Code    | Full plugin + SessionEnd / Stop hooks         | High          |
| 2        | Codex          | Plugin + SessionStart/SessionEnd hooks        | High          |
| 3        | Cursor         | Rules + skills + available hooks              | Medium–High   |
| 4        | Grok Build     | Native skill / AGENTS.md + best hooks         | Medium–High   |
| 5        | Kimi, Herdr... | Portable skill + documented patterns          | Good          |

**Privacy & quality:**
- Redacted digests only (no full transcripts by default)
- Easy opt-out
- Promotion of digests into durable `pages/` still benefits from LLM judgment

## Phases

1. **Foundation**  
   Core Go CLI, SCHEMA.md template, AGENTS.md injection, basic ingest/index/lint, OKF structure.

2. **Claude Code integration**  
   Real plugin with SessionEnd-style hooks for automatic redacted session capture.

3. **Broader agent support**  
   Codex, Cursor, Grok Build, portable skills for Kimi/Herdr/others.

4. **Polish & dogfooding**  
   The `repowiki` repo itself uses `.llm-wiki/` from early on. Refinements based on real usage.

## Success Criteria

- Agent in a `repowiki init`’d repo automatically consults the wiki.
- Session insights are captured with almost no user friction on major platforms.
- Knowledge compounds across tools and time.
- Compiled knowledge is valid OKF v0.2.
- Humans enjoy browsing the Markdown.
- The tool dogfoods itself successfully.

## Positioning

**Headline:** Persistent, Git-native project memory for any coding agent, with first-class (automatic) session capture.  
**Quiet standard:** OKF-compatible (Google Open Knowledge Format v0.2) so knowledge works everywhere.

---

This plan is living and will be updated as we implement.
