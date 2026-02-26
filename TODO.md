# TODO

## Core Library (`pkg/changelog/`)

- [x] `types.go` — Changelog, Version, Changes, Entry structs
- [x] `parser.go` — YAML loading + validation (project not empty, no duplicate versions, at most one unreleased, date format, no empty entries)
- [x] `query.go` — GetVersion, GetUnreleased, GetLastN, AllEntries, ListVersions, GetLatestRelease
- [x] `render.go` — CHANGELOG.md generation (Keep a Changelog compliant, comparison links)
- [x] `format.go` — Terminal rendering with colors/icons, text wrapping, `--plain` mode
- [x] `scaffold.go` — Conventional commit parsing + scaffold generation
- [x] `git.go` — DetectRepoURL, GitLog, LatestTag

## CLI Commands (`cmd/chlog/`)

- [x] `chlog init` — Create CHANGELOG.yaml with project name prompt
- [x] `chlog sync` — Regenerate CHANGELOG.md from YAML (idempotent)
- [x] `chlog check` — CI gate: byte-compare rendered vs actual CHANGELOG.md
- [x] `chlog validate` — Schema validation only (no sync)
- [x] `chlog extract <version>` — Output single version as markdown (for `gh release create --notes-file`)
- [x] `chlog show [version]` — Terminal display with `--last N` and `--plain` flags
- [x] `chlog scaffold` — Auto-scaffold YAML block from conventional commits since last tag

## Scaffold System (pm-style)

- [x] Parse conventional commits (feat→added, fix→fixed, refactor→changed, etc.)
- [x] Skip chore/docs/style/test/ci prefixes
- [x] Detect breaking changes via `!` in subject
- [x] Clean subjects (strip prefix, capitalize, first sentence only)
- [x] `--write` flag to insert into CHANGELOG.yaml
- [x] `--version` flag to override version string

## Public / Internal Visibility

Support marking entries as public-facing vs internal-only, so `sync`/`extract` can generate public release notes while `show` can include everything.

### Option A: Entry-level flag (mixed scalar/map YAML)

```yaml
added:
  - "Public feature"
  - text: "Internal refactor detail"
    internal: true
```

- Granular per-entry control
- Changes `[]string` → custom type that unmarshals both string and `{text, internal}` map
- Every consumer (render, format, query, scaffold) needs to handle the new type

### Option B: Version-level split (separate `internal` block)

```yaml
versions:
  1.0.0:
    date: "2024-03-15"
    changes:
      added:
        - Public feature
    internal:
      changed:
        - Refactored auth middleware
```

- Clean separation, entries stay `[]string`
- Add `Internal Changes` field to `Version`
- Add `--internal` flag to render/show/extract to include internal entries
- Smaller code change

### Implementation tasks

- [x] Pick approach (A or B) — chose B (version-level split)
- [x] Update `Version` struct + YAML tags
- [x] Update parser validation (internal entries follow same rules)
- [x] Add `--internal` flag to `show`, `extract`, `sync`
- [x] `check` should compare with same internal setting as `sync` produced
- [x] Update `scaffold` to classify commits (`refactor`/`perf` → internal by default)
- [x] Tests for both public-only and full rendering

## CI / Integration

- [x] GitHub Actions example in README
- [x] `chlog check` exit codes (0 = in sync, 1 = out of sync, 2 = validation error)
- [x] goreleaser config for cross-platform binaries

## Stretch

- [x] `chlog release <version>` — Move unreleased→versioned, add date, create fresh unreleased block
- [x] Config file (`.chlog.yaml`) for project-specific settings (repo URL for comparison links)
- [x] Importable as Go library: `import "github.com/ariel-frischer/chlog/pkg/changelog"`
