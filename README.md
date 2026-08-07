# RepoWiki

**Persistent, Git-native project memory for any coding agent.**

OKF-compatible LLM Wiki with automatic (redacted) session capture via plugins/hooks for Claude Code, Codex, Cursor, Grok Build, Kimi, Herdr and more.

## Status

Planning complete. Ready for Phase 1 (Foundation).

Full project plan → [PLAN.md](./PLAN.md)

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
