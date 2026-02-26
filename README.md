<div align="center">

<pre>
█▀▀ █ █ █   █▀█ █▀▀
█▄▄ █▀█ █▄▄ █▄█ █▄█
</pre>

**YAML-First Changelog Management — CLI + Go Library**

[![CI](https://github.com/ariel-frischer/chlog/actions/workflows/ci.yml/badge.svg)](https://github.com/ariel-frischer/chlog/actions/workflows/ci.yml)
[![GitHub Release](https://img.shields.io/github/v/release/ariel-frischer/chlog)](https://github.com/ariel-frischer/chlog/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/ariel-frischer/chlog)](https://goreportcard.com/report/github.com/ariel-frischer/chlog)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Single `CHANGELOG.yaml` source of truth → auto-generated `CHANGELOG.md` → CI validation.

</div>

## Install

**Quick install** (Linux/macOS):

```bash
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/chlog/main/install.sh | sh
```

**Go install**:

```bash
go install github.com/ariel-frischer/chlog@latest
```

**From source**:

```bash
git clone https://github.com/ariel-frischer/chlog.git
cd chlog
make build    # Binary at bin/chlog
```

## Quickstart

```bash
chlog init                          # Create CHANGELOG.yaml + .chlog.yaml (auto-detects repo URL)
chlog init --project myapp          # Skip project name prompt
chlog sync                          # Generate CHANGELOG.md from YAML
chlog check                         # CI gate — verify markdown matches YAML
chlog validate                      # Validate YAML schema
chlog show                          # View changelog in terminal
chlog show 0.3.0                    # View specific version
chlog show --last 5                 # View last 5 entries
chlog extract 0.3.0                 # Output release notes (for gh release)
chlog scaffold                      # Auto-scaffold from conventional commits
chlog scaffold --write              # Scaffold and merge into CHANGELOG.yaml
chlog scaffold --version 1.2.0      # Scaffold with explicit version string
chlog release 1.0.0                 # Promote unreleased → 1.0.0 with today's date
chlog release 1.0.0 --date 2026-03-01  # Promote with explicit date
```

## Why?

Commit-based changelog tools (git-cliff, semantic-release) dump git logs. That's fine for human-curated commits, but with AI agents generating dozens of implementation commits, you want **curated public & internal entries**.

`chlog` separates product communication from implementation history:

```
# git log (what tools like git-cliff give you)
a4f2c1 fix: adjust retry backoff timing
b92e0a refactor: extract http client helper
c7d31b fix: handle nil pointer in auth middleware
d1a8ef feat: add timeout flag
e53f90 chore: update deps

# CHANGELOG.yaml (what you write, or your agent summarizes)
added:
  - "Configurable request timeout via --timeout flag"
fixed:
  - "Auth no longer crashes on expired tokens"
```

## Schema

```yaml
project: my-tool
versions:
  - version: unreleased
    changes:
      added:
        - "New feature description"
    internal:
      changed:
        - "Refactored auth middleware"

  - version: 0.1.0
    date: "2026-02-24"
    changes:
      added:
        - "Initial release"
      fixed:
        - "Bug fix description"
```

Six categories from [Keep a Changelog](https://keepachangelog.com/): `added`, `changed`, `deprecated`, `removed`, `fixed`, `security`.

### Internal entries

chlog supports a two-tier model: **public** entries (customer-facing release notes) and **internal** entries (implementation details like refactors, perf improvements, dependency updates). Public entries live under `changes`, internal entries under `internal` — same categories, separate audiences.

By default, internal entries are excluded from output. Include them with `--internal`:

```bash
chlog sync --internal       # Render internal entries in CHANGELOG.md
chlog show --internal       # Show internal entries in terminal
chlog extract 1.0 --internal
chlog check --internal      # Compare with internal entries included
```

To always include internal entries, set `include_internal: true` in `.chlog.yaml` (see [Config](#config)). The `--internal` flag and config option are OR'd together — config provides the team default, the flag always adds them.

`chlog scaffold` auto-classifies `refactor`/`perf` conventional commits as internal, so `scaffold --write` populates both tiers automatically.

## Config

Optional `.chlog.yaml` in your project root. Created automatically by `chlog init` with auto-detected values:

```yaml
repo_url: https://github.com/myorg/myproject
include_internal: true
```

| Field | Default | Description |
|-------|---------|-------------|
| `repo_url` | auto-detect from `git remote origin` | Used for version comparison links in `CHANGELOG.md` |
| `include_internal` | `false` | Include internal entries in all commands (`sync`, `show`, `extract`, `check`) |

## CI

GitHub Actions example:

```yaml
name: Changelog Check
on:
  pull_request:
    paths: ["CHANGELOG.yaml", "CHANGELOG.md"]
jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: go install github.com/ariel-frischer/chlog@latest
      - run: chlog validate
      - run: chlog check
```

Exit codes: `0` in sync, `1` out of sync, `2` validation error.

## Library

`chlog` is both a CLI tool and an importable Go library. Use `pkg/changelog` directly in your own Go projects for programmatic changelog management — no subprocess needed.

```go
import "github.com/ariel-frischer/chlog/pkg/changelog"

// Load and query
c, err := changelog.Load("CHANGELOG.yaml")
latest := c.GetLatestRelease()
entries := c.GetLastN(5)

// Programmatic release
c.Release("2.0.0", "2024-06-01")
changelog.Save(c, "CHANGELOG.yaml")

// Render to Markdown
md := changelog.RenderMarkdown(c, changelog.RenderOptions{})

// Parse from any io.Reader
c, err = changelog.LoadFromReader(reader)
```

See the [package documentation](https://pkg.go.dev/github.com/ariel-frischer/chlog/pkg/changelog) for the full API.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

[MIT](LICENSE)
