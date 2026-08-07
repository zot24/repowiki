package wiki

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	injectBegin = "<!-- repowiki:begin — managed by `repowiki init`, do not edit inside this block -->"
	injectEnd   = "<!-- repowiki:end -->"
)

// AgentsSnippet is the block injected into AGENTS.md / CLAUDE.md.
const AgentsSnippet = `## Project memory (.llm-wiki)

This repository has a persistent knowledge base in ` + "`.llm-wiki/`" + ` maintained with the ` + "`repowiki`" + ` CLI.

**Before answering questions about this project**, consult:
- ` + "`.llm-wiki/index.md`" + ` — map of all compiled knowledge
- ` + "`.llm-wiki/pages/`" + ` — architecture, decisions, plans, components, concepts
- ` + "`.llm-wiki/SCHEMA.md`" + ` — the protocol for reading and writing wiki content

**When you learn something durable** (a decision, an architectural insight, a gotcha,
a plan), capture it:
- Quick capture: drop a markdown note in ` + "`.llm-wiki/inbox/`" + `
- Structured capture: ` + "`repowiki ingest <file> --type decision|note|session|source`" + `
- Then keep the index fresh: ` + "`repowiki index`" + `

Follow the conventions in ` + "`.llm-wiki/SCHEMA.md`" + `. Prefer updating an existing page
over creating a near-duplicate.`

func schemaTemplate() string {
	return `---
type: concept
title: "SCHEMA — RepoWiki protocol"
description: "How agents and humans read and write this wiki"
status: stable
okf_version: "` + OKFVersion + `"
---

# SCHEMA — how this wiki works

This directory is an OKF v` + OKFVersion + ` knowledge bundle managed with the ` + "`repowiki`" + ` CLI.

## Layers

| Layer | Path | Rule |
|-------|------|------|
| Raw | ` + "`raw/`" + ` | Immutable sources: session digests, notes, decisions, external docs. Append-only — never rewrite. |
| Compiled | ` + "`pages/`" + ` | Synthesized, interlinked knowledge. Freely edited and refactored. |
| Protocol | ` + "`SCHEMA.md`" + ` | This file. Conventions and project-specific rules. |
| Inbox | ` + "`inbox/`" + ` | Drop zone for quick captures awaiting ingestion. |

## Frontmatter

Every document needs YAML frontmatter with at least ` + "`type`" + ` and ` + "`title`" + `:

- ` + "`type`" + `: decision | architecture | plan | concept | component | overview | open-question | session | note | source
- ` + "`title`" + `, ` + "`description`" + `, ` + "`tags`" + `
- ` + "`status`" + `: draft | stable | deprecated
- ` + "`sources`" + `: provenance pointers back to ` + "`raw/`" + ` files
- ` + "`created`" + ` / ` + "`updated`" + ` (YYYY-MM-DD), ` + "`confidence`" + `, ` + "`related`" + `

## Rules for agents

1. **Consult before answering.** Check ` + "`index.md`" + ` and relevant ` + "`pages/`" + ` first.
2. **Capture intentionally.** Redacted digests, decisions, and insights — not raw chat dumps.
3. **Cite provenance.** Compiled pages list their ` + "`sources`" + ` from ` + "`raw/`" + `.
4. **Never rewrite ` + "`raw/`" + `.** Promote and synthesize into ` + "`pages/`" + ` instead.
5. **Keep structure via the CLI.** ` + "`repowiki ingest`" + `, ` + "`repowiki index`" + `, ` + "`repowiki lint`" + `.
6. **No secrets.** Redact credentials, tokens, PII before anything lands here.

## Project-specific conventions

_Add project conventions below this line._
`
}

func indexTemplate() string {
	return `---
type: overview
title: "Wiki index"
description: "Root index of this OKF knowledge bundle"
okf_version: "` + OKFVersion + `"
status: stable
---

# Wiki index

_Empty wiki. Run ` + "`repowiki index`" + ` after adding content to regenerate this file._
`
}

func logTemplate() string {
	return `---
type: note
title: "Activity log"
description: "Chronological log of wiki activity"
status: stable
---

# Activity log

`
}

func configTemplate() string {
	return `---
type: note
title: "RepoWiki config"
description: "Optional project overrides for repowiki"
status: draft
---

# Config

No overrides yet. Future keys: capture opt-out, redaction rules, custom types.
`
}

// InitBundle creates the .llm-wiki structure under repoRoot.
// Existing files are never overwritten. Returns the created wiki.
func InitBundle(repoRoot string) (*Wiki, []string, error) {
	w := &Wiki{RepoRoot: repoRoot, Root: filepath.Join(repoRoot, DirName)}
	var created []string

	dirs := []string{w.Root, filepath.Join(w.Root, "inbox"), filepath.Join(w.Root, "pages")}
	for _, s := range RawSubdirs {
		dirs = append(dirs, filepath.Join(w.Root, "raw", s))
	}
	for _, s := range PageSubdirs {
		dirs = append(dirs, filepath.Join(w.Root, "pages", s))
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return nil, nil, err
		}
		// Keep empty directories present after git clone.
		keep := filepath.Join(d, ".gitkeep")
		if entries, err := os.ReadDir(d); err == nil && len(entries) == 0 {
			if err := os.WriteFile(keep, nil, 0o644); err != nil {
				return nil, nil, err
			}
		}
	}

	files := map[string]string{
		"index.md":  indexTemplate(),
		"log.md":    logTemplate(),
		"SCHEMA.md": schemaTemplate(),
		"config.md": configTemplate(),
		filepath.Join("pages", "overview.md"): `---
type: overview
title: "Project overview"
description: "What this project is and how it fits together"
status: draft
---

# Project overview

_Not yet written. Synthesize from raw/ material as it accumulates._
`,
		filepath.Join("pages", "open-questions.md"): `---
type: open-question
title: "Open questions"
description: "Unresolved questions and pending decisions"
status: draft
---

# Open questions

- _None recorded yet._
`,
	}
	for rel, content := range files {
		p := filepath.Join(w.Root, rel)
		if _, err := os.Stat(p); err == nil {
			continue
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			return nil, nil, err
		}
		created = append(created, filepath.Join(DirName, rel))
	}
	return w, created, nil
}

// InjectAgentsBlock inserts (or refreshes) the managed repowiki block in the
// given file, creating the file when create is true and it doesn't exist.
// Returns what happened: "created", "updated", "unchanged", or "skipped".
func InjectAgentsBlock(path string, create bool) (string, error) {
	block := injectBegin + "\n\n" + AgentsSnippet + "\n\n" + injectEnd
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		if !create {
			return "skipped", nil
		}
		header := fmt.Sprintf("# %s\n\n", strings.TrimSuffix(filepath.Base(path), ".md"))
		content := header + block + "\n"
		return "created", os.WriteFile(path, []byte(content), 0o644)
	}
	if err != nil {
		return "", err
	}
	content := string(raw)
	if b := strings.Index(content, injectBegin); b >= 0 {
		if e := strings.Index(content, injectEnd); e > b {
			updated := content[:b] + block + content[e+len(injectEnd):]
			if updated == content {
				return "unchanged", nil
			}
			return "updated", os.WriteFile(path, []byte(updated), 0o644)
		}
	}
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += "\n" + block + "\n"
	return "updated", os.WriteFile(path, []byte(content), 0o644)
}

// Today returns the current date as YYYY-MM-DD.
func Today() string { return time.Now().Format("2006-01-02") }
