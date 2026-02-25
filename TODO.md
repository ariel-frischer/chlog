# TODO

## Core Library (`pkg/changelog/`)

- [ ] `types.go` — Changelog, Version, Changes, Entry structs
- [ ] `parser.go` — YAML loading + validation (project not empty, no duplicate versions, at most one unreleased, date format, no empty entries)
- [ ] `query.go` — GetVersion, GetUnreleased, GetLastN, AllEntries, ListVersions, GetLatestRelease
- [ ] `render.go` — CHANGELOG.md generation (Keep a Changelog compliant, comparison links)
- [ ] `format.go` — Terminal rendering with colors/icons, text wrapping, `--plain` mode

## CLI Commands (`cmd/chlog/`)

- [ ] `chlog init` — Create CHANGELOG.yaml with project name prompt
- [ ] `chlog sync` — Regenerate CHANGELOG.md from YAML (idempotent)
- [ ] `chlog check` — CI gate: byte-compare rendered vs actual CHANGELOG.md
- [ ] `chlog validate` — Schema validation only (no sync)
- [ ] `chlog extract <version>` — Output single version as markdown (for `gh release create --notes-file`)
- [ ] `chlog show [version]` — Terminal display with `--last N` and `--plain` flags
- [ ] `chlog scaffold` — Auto-scaffold YAML block from conventional commits since last tag

## Scaffold System (pm-style)

- [ ] Parse conventional commits (feat→added, fix→fixed, refactor→changed, etc.)
- [ ] Skip chore/docs/style/test/ci prefixes
- [ ] Detect breaking changes via `!` in subject
- [ ] Clean subjects (strip prefix, capitalize, first sentence only)
- [ ] `--write` flag to insert into CHANGELOG.yaml
- [ ] `--version` flag to override version string

## CI / Integration

- [ ] GitHub Actions example in README
- [ ] `chlog check` exit codes (0 = in sync, 1 = out of sync, 2 = validation error)
- [ ] goreleaser config for cross-platform binaries

## Stretch

- [ ] `chlog release <version>` — Move unreleased→versioned, add date, create fresh unreleased block
- [ ] Config file (`.chlog.yaml`) for project-specific settings (repo URL for comparison links, custom categories)
- [ ] Importable as Go library: `import "github.com/ariel-frischer/chlog/pkg/changelog"`
