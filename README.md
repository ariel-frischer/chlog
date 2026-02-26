<div align="center">

<pre>
█▀▀ █ █ █   █▀█ █▀▀
█▄▄ █▀█ █▄▄ █▄█ █▄█
</pre>

**YAML-First Changelog Management**

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
chlog init                          # Create CHANGELOG.yaml
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

Entries under `internal` are excluded by default. Pass `--internal` to include them:

```bash
chlog sync --internal       # Render internal entries in CHANGELOG.md
chlog show --internal       # Show internal entries in terminal
chlog extract 1.0 --internal
chlog check --internal      # Compare with internal entries included
```

Scaffold auto-classifies `refactor`/`perf` commits as internal.

## Config

Optional `.chlog.yaml` in your project root:

```yaml
repo_url: https://github.com/myorg/myproject
```

The `repo_url` is used for version comparison links in `CHANGELOG.md`. If omitted, chlog auto-detects from `git remote origin`.

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

Importable as a Go library:

```go
import "github.com/ariel-frischer/chlog/pkg/changelog"
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

[MIT](LICENSE)
