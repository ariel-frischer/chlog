# Per-Change Prompt

Add this to your task prompt so the agent updates the changelog alongside the code.

## Prompt

```
After implementing the changes, add a changelog entry to CHANGELOG.yaml under
the `unreleased` version. Place user-facing entries directly on the version
(added, changed, fixed, etc.) and use the `internal` section for refactors,
perf, or implementation details.

Write entries as concise, user-facing descriptions (not commit messages).
You can add entries via CLI or by editing YAML directly:

  chlog add added "New feature"             # Add to unreleased
  chlog add changed -i "Refactored X"       # Add internal entry
  chlog remove added "Reverted feature"     # Remove an entry

After editing, validate and sync:

  chlog validate        # Check YAML schema is valid
  chlog sync            # Regenerate CHANGELOG.md
  chlog check           # Verify markdown matches YAML

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
        changed: []                   #   same categories, separate audience
    1.0.0:
      date: "2026-01-15"             # Required for released versions
      added:
        - "Initial release"

Categories are arbitrary YAML keys on each version. The default 6 from Keep
a Changelog are enforced by default (strict mode). To allow custom categories
like `performance` or `infrastructure`, add to .chlog.yaml:

  categories: [added, changed, fixed, performance]  # custom allowlist
  strict_categories: false                           # or disable validation entirely
```

## Example result

```yaml
versions:
  unreleased:
    added:
      - "Export to CSV from the dashboard"
    internal:
      changed:
        - "Switched CSV serialization to encoding/csv for streaming"
```
