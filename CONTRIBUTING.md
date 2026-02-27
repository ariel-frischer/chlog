# Contributing to chlog

Thanks for your interest in contributing!

## Getting Started

```bash
git clone https://github.com/ariel-frischer/chlog.git
cd chlog
make install   # Download dependencies
make build     # Build binary
make test      # Run tests
```

## Development

```bash
make build     # Build to bin/chlog
make test      # Run all tests
make lint      # Run golangci-lint
make fmt       # Format code
```

## Pull Requests

1. Fork the repo and create your branch from `main`
2. Add tests for any new functionality
3. Ensure `make test` and `make lint` pass
4. Update `CHANGELOG.yaml` with your changes (`chlog add added "Description"` or `chlog scaffold --write`)
5. Run `chlog sync` to regenerate `CHANGELOG.md`

## Reporting Issues

Use [GitHub Issues](https://github.com/ariel-frischer/chlog/issues). Include:
- What you expected vs what happened
- Steps to reproduce
- `chlog --version` output
- OS and architecture

## Code Style

- Functions under 40 lines
- Errors wrapped with context: `fmt.Errorf("doing X: %w", err)`
- Table-driven tests with `map[string]struct{}`

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
