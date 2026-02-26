# CLAUDE.md

## Project: chlog

Language-agnostic CLI for YAML-first changelog management. Single YAML source of truth, auto-generated CHANGELOG.md, CI validation, release notes extraction.

## Commands

```bash
make build          # Build binary
make test           # Run tests
make lint           # Run linters
make fmt            # Format code
```

## Architecture

```
cmd/chlog/          # CLI entry point (cobra)
pkg/changelog/      # Library (importable by other Go projects)
  types.go          # Changelog, Version, Changes structs + Remove method + error types
  parser.go         # YAML loading + validation
  query.go          # Query methods (GetVersion, GetLastN, etc.)
  render.go         # CHANGELOG.md generation
  format.go         # Terminal rendering (colors, icons)
```

## Coding Standards

- Functions under 40 lines
- Errors wrapped with context: `fmt.Errorf("doing X: %w", err)`
- Map-based table tests: `map[string]struct{}`
- Accept interfaces, return concrete types
