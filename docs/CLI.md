# `repowiki` CLI reference

A single static binary. All commands operate on the nearest `.llm-wiki/` found by
walking up from the current directory (except `init`, which creates one).

## Commands

| Command | What it does |
|---------|--------------|
| `repowiki init [dir]` | Create the `.llm-wiki/` bundle: full directory structure, `SCHEMA.md`, `index.md` (OKF v0.2), `log.md`, `config.md`, starter pages. Injects a managed block into `AGENTS.md` (created if missing) and `CLAUDE.md` / `.cursorrules` (only if present), and installs the Claude Code project skill at `.claude/skills/repowiki/SKILL.md`. Idempotent — never overwrites existing files, refreshes managed blocks. |
| `repowiki status` | Overview: page/raw counts, pending inbox items, recent activity. |
| `repowiki ingest <file\|url\|->` | Capture content into `raw/` with frontmatter + provenance. Reads a file, fetches a URL (stored as `source`), or reads stdin with `-`. Appends to `log.md`. |
| `repowiki index` | Regenerate `index.md` from `pages/` and `raw/` (titles + descriptions from frontmatter, grouped by section). |
| `repowiki lint [--fix]` | Check structure and OKF conformance: frontmatter presence, required `type`/`title`, valid `status`, root `okf_version`, broken relative links. `--fix` adds frontmatter stubs where missing. |
| `repowiki search <query>` | Case-insensitive full-text search across the wiki (`query` is an alias). Prints `path:line: match`. |
| `repowiki get <wiki-path>` | Print a document by wiki-relative path, e.g. `repowiki get pages/overview.md`. |
| `repowiki log [-n N]` | Show the last N activity entries (default 20). |
| `repowiki skill [install]` | Print the agent skill to stdout, or `install` it to `.claude/skills/repowiki/SKILL.md`. |

## `ingest` flags

| Flag | Meaning |
|------|---------|
| `--type` | `session` \| `decision` \| `note` \| `source` \| `auto` (default). `auto` uses the content's own frontmatter `type`, else `note`; URLs default to `source`. |
| `--title` | Override the title (otherwise: frontmatter title → first `#` heading → filename). |
| `--tags` | Comma-separated tags merged into frontmatter. |

Destination: `raw/<type-dir>/YYYY-MM-DD-<slug>.md` (suffixes `-2`, `-3`… on collision).

## `init` flags

| Flag | Meaning |
|------|---------|
| `--no-inject` | Skip `AGENTS.md` / `CLAUDE.md` / `.cursorrules` injection. |
| `--no-skill` | Skip installing the Claude Code project skill. |

## Typical loop

```sh
repowiki init                                   # once per repo
echo "insight..." | repowiki ingest - --type note --title "Gotcha with X"
repowiki index                                  # keep the index fresh
repowiki lint                                   # before committing
repowiki search "auth"                          # when you need something back
```
