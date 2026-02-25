# chlog

Language-agnostic CLI for YAML-first changelog management.

Single `CHANGELOG.yaml` source of truth → auto-generated `CHANGELOG.md` → CI validation.

## Install

```bash
go install gitlab.com/ariel-frischer/chlog@latest
```

## Usage

```bash
chlog init                  # Create CHANGELOG.yaml in current directory
chlog sync                  # Generate CHANGELOG.md from CHANGELOG.yaml
chlog check                 # CI gate — verify markdown matches YAML
chlog extract 0.3.0         # Output release notes for a version
chlog validate              # Validate YAML schema
chlog scaffold              # Auto-scaffold entry from conventional commits
chlog show                  # View changelog in terminal
chlog show 0.3.0            # View specific version
chlog show --last 10        # View last 10 entries
```

## Why?

Commit-based changelog tools (git-cliff, semantic-release) dump git logs. That's fine for human-curated commits, but with AI agents generating dozens of implementation commits, you need **curated, user-facing entries**.

`chlog` separates product communication from implementation history.

## Schema

```yaml
project: my-tool
versions:
  - version: unreleased
    changes:
      added:
        - "New feature description"

  - version: 0.1.0
    date: "2026-02-24"
    changes:
      added:
        - "Initial release"
      fixed:
        - "Bug fix description"
```

Six categories from [Keep a Changelog](https://keepachangelog.com/): `added`, `changed`, `deprecated`, `removed`, `fixed`, `security`.

## CI

Add `.github/workflows/changelog.yml` to validate your changelog on PRs:

```yaml
name: Changelog Check

on:
  pull_request:
    paths:
      - "CHANGELOG.yaml"
      - "CHANGELOG.md"

jobs:
  changelog-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - run: go install gitlab.com/ariel-frischer/chlog@latest

      - name: Validate CHANGELOG.yaml
        run: chlog validate

      - name: Verify CHANGELOG.md is in sync
        run: chlog check
```

Exit codes: `0` = in sync, `1` = out of sync, `2` = validation error.
