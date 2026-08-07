# RepoWiki

**Persistent, Git-native project memory for any coding agent.**

OKF-compatible LLM Wiki with automatic (redacted) session capture via plugins/hooks for Claude Code, Codex, Cursor, Grok Build, Kimi, Herdr and more.

## Why

Coding agents have amnesia. Every session starts from zero: the agent re-reads your
code, re-derives your architecture, and re-asks questions you answered last week.
The *why* behind your code — decisions made, alternatives rejected, gotchas
discovered — lives in chat histories that expire, scattered across tools that
don't talk to each other.

RepoWiki gives every repository a **persistent memory that lives with the code**:

- **Survives sessions.** A decision captured today is consulted by the agent next month — automatically, because `AGENTS.md` / `CLAUDE.md` point to it.
- **Survives tools.** Captured in Claude Code, consulted by Codex or Cursor. It's just Markdown in `.llm-wiki/`, versioned and cloned with the repo.
- **Survives people.** A new teammate (human or agent) clones the repo and inherits the project's accumulated judgment, not just its code.
- **Portable by design.** Compiled knowledge is a valid [OKF](https://github.com/google/okf) v0.2 bundle, so it works with anything that speaks the format.

This is the project-memory layer that RAG over source code or raw chat logs can't
provide: small, intentional, high-signal knowledge — not transcript dumps.

## Example

```sh
cd my-project
repowiki init      # creates .llm-wiki/, seeds SCHEMA.md, injects AGENTS.md,
                   # installs the Claude Code skill
```

During work, you (or your agent, automatically) capture what matters:

```sh
echo "# Chose Postgres over DynamoDB

We need relational queries for billing reports; single-region is fine.
Rejected DynamoDB: access patterns too varied, team knows SQL." \
| repowiki ingest - --type decision --title "Chose Postgres over DynamoDB" --tags db,billing

repowiki index     # refresh the wiki's table of contents
```

Three weeks later — different session, maybe a different agent entirely:

> **You:** should we move billing to DynamoDB to cut costs?
>
> **Agent:** *(consults `.llm-wiki/` per AGENTS.md)* We decided against DynamoDB
> on 2026-08-07 — billing needs relational queries and the access patterns were
> too varied. If costs are the new driver, that decision could be revisited, but
> here's the context we'd be trading away…

The agent argues from your project's recorded judgment instead of re-deriving
(or contradicting) it. Finding things back works for humans too:

```sh
repowiki search dynamodb
repowiki get raw/decisions/2026-08-07-chose-postgres-over-dynamodb.md
```

## Status

Planning complete. Ready for Phase 1 (Foundation).

Full project plan → [PLAN.md](./PLAN.md)
CLI command reference → [docs/CLI.md](./docs/CLI.md)

### Core Goals
- Per-repository LLM Wiki that lives with the code (Git-versioned)
- OKF v0.2 compatible (Google Open Knowledge Format)
- Automatic / near-automatic redacted session capture via platform hooks & plugins
- Seamless integration with major coding agents (Claude Code first)
- Simple, reliable Go CLI (`repowiki`) as the backbone
- High-quality, intentional knowledge (not raw chat dumps)

### Phases (high level)
1. **Foundation** – Core Go CLI + SCHEMA.md + AGENTS.md injection + OKF structure
2. **Claude Code integration** – Real plugin + SessionEnd-style hooks for automatic capture
3. **Broader agent support** – Codex, Cursor, Grok Build, portable skills for the rest
4. **Polish & dogfooding**

---

Built incrementally. The plan lives in this repo and will evolve with the code.
