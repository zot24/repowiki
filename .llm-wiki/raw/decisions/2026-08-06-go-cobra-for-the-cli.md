---
type: decision
title: Go + Cobra for the CLI
tags:
  - cli
  - architecture
status: draft
created: "2026-08-06"
updated: "2026-08-06"
---

# Go + Cobra for the CLI

Decided in planning: the repowiki CLI is a single static Go binary built with Cobra.
Rationale: reliability and determinism for structural operations (init, ingest, index, lint),
easy cross-platform distribution, no runtime dependencies. LLM agents handle synthesis;
the binary handles structure.
