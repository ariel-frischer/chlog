# Changelog

All notable changes to chlog will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

## [0.0.4] - 2026-02-26

### Added

- Version command with build info display
- ASCII branding for CLI
- Production CI workflow
- Open source GitHub release preparation (LICENSE, .goreleaser config)

## [0.0.3] - 2026-02-25

### Added

- Support for internal changelog entries (excluded from public CHANGELOG.md)
- `chlog release` command to promote unreleased changes to a versioned release
- Configuration section and `.chlog.yaml` file support
- Tests for config loading and repo URL resolution

### Changed

- README restructured with Quickstart section, new command examples, and clearer usage instructions

### Fixed

- Include config loading in check command for proper validation
- Allow empty unreleased versions while enforcing non-empty released versions

## [0.0.2] - 2026-02-25

### Added

- Core CLI with init, check, extract, sync, validate, and show commands
- YAML-based changelog parsing and validation
- Markdown rendering from CHANGELOG.yaml
- Terminal output with colors and icons

### Changed

- Makefile install target alias and go-install target for streamlined dependency management

## [0.0.1] - 2026-02-24

### Added

- Initial project scaffolding (Go module, directory structure, Makefile)

### Fixed

- Use correct gitlab.com module path instead of github.com

[Unreleased]: https://github.com/ariel-frischer/chlog/compare/v0.0.4...HEAD
[0.0.4]: https://github.com/ariel-frischer/chlog/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/ariel-frischer/chlog/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/ariel-frischer/chlog/compare/v0.0.1...v0.0.2
