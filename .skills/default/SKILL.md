---
name: chlog
description: >
  YAML-first changelog management CLI. Use when creating, updating, validating,
  or rendering changelogs. Covers init, sync, check, show, extract, scaffold,
  release commands, internal entries, and CHANGELOG.yaml schema.
license: MIT
compatibility:
  - Claude Code
  - Cursor
  - Codex
  - Gemini CLI
  - VS Code
metadata:
  author: ariel-frischer
  version: 0.0.4
  tags: changelog, yaml, go, cli, ci
allowed-tools: Bash Read Write Edit
---

# chlog

YAML-first changelog CLI. Single `CHANGELOG.yaml` → auto-generated `CHANGELOG.md` → CI validation.

## Commands

```bash
chlog init                      # Create CHANGELOG.yaml + .chlog.yaml (auto-detects repo URL)
chlog init --project myapp      # Skip project name prompt
chlog validate                  # Validate YAML schema
chlog sync                      # Generate CHANGELOG.md from YAML
chlog check                     # CI gate: exit 0=sync, 1=stale, 2=invalid
chlog show                      # Terminal display (colors + icons)
chlog show 1.2.0                # Single version
chlog show --last 5             # Last N entries
chlog show --plain              # No ANSI
chlog extract 1.0.0             # Markdown for one version (pipe to gh release)
chlog scaffold                  # Dry-run: conventional commits → YAML
chlog scaffold --write          # Merge into CHANGELOG.yaml
chlog release 1.0.0             # Promote unreleased → 1.0.0 (today's date)
chlog release 1.0.0 --date 2026-03-01
```

Global flags: `-f` (CHANGELOG.yaml path), `--config` (.chlog.yaml path), `--internal` (include internal entries).

## Schema

```yaml
project: my-project
versions:
  unreleased:
    added:
      - "New feature"
    internal:                   # Excluded by default, use --internal
      changed:
        - "Refactored auth"
  "1.0.0":
    date: "2024-01-01"         # Required for released versions
    added: []
    changed: []
    deprecated: []
    removed: []
    fixed: []
    security: []
```

Six categories from [Keep a Changelog](https://keepachangelog.com/) live directly on each version entry.

## Config (.chlog.yaml)

```yaml
repo_url: https://github.com/org/repo   # Auto-detected from git remote
include_internal: false                  # Include internal entries by default
```

## Scaffold Mapping

| Commit Type | Category | Tier |
|---|---|---|
| `feat` | added | public |
| `fix` | fixed | public |
| `refactor`, `perf` | changed | internal |
| `deprecate` | deprecated | public |
| `remove` | removed | public |
| `chore`, `docs`, `style`, `test`, `ci`, `build` | skipped | — |

Breaking changes (`feat!:`) → `changed`, prefixed `BREAKING: `.

## Workflows

**Setup:** `chlog init` → edit CHANGELOG.yaml → `chlog sync`

**Release:** `chlog scaffold --write` → curate → `chlog release 1.2.0` → `chlog sync` → `chlog extract 1.2.0 > notes.md`

**CI:** `chlog validate && chlog check`

## Go Library

```go
import "github.com/ariel-frischer/chlog/pkg/changelog"

c, _ := changelog.Load("CHANGELOG.yaml")
v, _ := c.GetVersion("1.0.0")
md, _ := changelog.RenderMarkdownString(c)
c.Release("2.0.0", "2024-06-01")
changelog.Save(c, "CHANGELOG.yaml")
```
