# Release Summary Prompt

Use this when you're ready to cut a release and want the agent to summarize all changes since the last version.

## Prompt

```
Review all commits since the last release tag. Run `chlog scaffold --write` to
get a starting point from conventional commits, then curate the CHANGELOG.yaml
unreleased entries:

1. Merge related commits into single user-facing descriptions
2. Drop noise (typo fixes, formatting, CI tweaks) unless they affect users
3. Classify refactors and perf improvements as `internal`
4. Write entries from the user's perspective ("Added X" not "Implemented X handler")
5. Validate and sync when done

Keep entries concise — one line each, no implementation details in the public section.

Commands:

  chlog scaffold --write          # Auto-scaffold from conventional commits
  chlog validate                  # Check YAML schema is valid
  chlog sync                      # Regenerate CHANGELOG.md
  chlog check                     # Verify markdown matches YAML
  chlog release 1.2.0             # Promote unreleased → 1.2.0
  chlog extract 1.2.0 > notes.md  # Extract release notes for gh release

CHANGELOG.yaml schema:

  project: my-project
  versions:
    unreleased:
      added: []                       # Public / user-facing categories
      changed: []                     #   live directly on the version
      deprecated: []
      removed: []
      fixed: []
      security: []
      internal:                       # Implementation details (optional)
        added: []
        changed: []
        deprecated: []
        removed: []
        fixed: []
        security: []
    1.0.0:
      date: "2026-01-15"             # Required for released versions
      added:
        - "Initial release"
```

## Example result

Before (raw scaffold from 12 commits):

```yaml
added:
  - "feat: add timeout flag to request command"
  - "feat: add --json output flag"
fixed:
  - "fix: handle nil pointer in auth middleware"
  - "fix: adjust retry backoff timing"
  - "fix: correct header parsing for multipart"
internal:
  changed:
    - "refactor: extract http client helper"
    - "perf: cache DNS lookups"
```

After (agent-curated):

```yaml
added:
  - "Configurable request timeout via `--timeout` flag"
  - "JSON output format with `--json`"
fixed:
  - "Auth no longer crashes on expired tokens"
  - "Multipart uploads now work with custom headers"
internal:
  changed:
    - "Extracted shared HTTP client, added DNS caching"
```
