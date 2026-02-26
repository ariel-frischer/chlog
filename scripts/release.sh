#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-}"
if [[ -z "$VERSION" ]]; then
  echo "Usage: $0 <version>"
  echo "  e.g. $0 v0.1.0"
  exit 1
fi

# Strip leading v for chlog release (expects bare semver)
BARE_VERSION="${VERSION#v}"
# Ensure tag has v prefix
TAG="v${BARE_VERSION}"

echo "==> Pre-flight checks..."
if [[ -n "$(git status --porcelain)" ]]; then
  echo "Error: working tree is dirty"
  exit 1
fi

echo "==> Running tests..."
make test

echo "==> Running lint..."
make lint

echo "==> Building chlog..."
make build

echo "==> Checking unreleased entries..."
if ! ./bin/chlog show unreleased | grep -q .; then
  echo "Error: no unreleased entries in CHANGELOG.yaml"
  exit 1
fi

echo "==> Stamping changelog: ${BARE_VERSION}..."
./bin/chlog release "${BARE_VERSION}"

echo "==> Syncing CHANGELOG.md..."
./bin/chlog sync

echo "==> Committing changelog..."
git add CHANGELOG.yaml CHANGELOG.md
git commit -m "release: ${TAG}"

echo "==> Tagging ${TAG}..."
git tag -a "${TAG}" -m "Release ${TAG}"

echo "==> Pushing to origin..."
git push origin main
git push origin "${TAG}"

echo ""
echo "Done! ${TAG} tagged and pushed."
echo ""
echo "Next steps:"
echo "  1. Watch GitHub CI:  gh run watch"
echo "  2. Once green:       goreleaser release --clean"
