# Per-Change Prompt

Add this to your task prompt so the agent updates the changelog alongside the code.

## Prompt

```
After implementing the changes, add a changelog entry to CHANGELOG.yaml under
the `unreleased` version. Place user-facing entries directly on the version
(added, changed, fixed, etc.) and use the `internal` section for refactors,
perf, or implementation details.

Write entries as concise, user-facing descriptions (not commit messages).
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
        added: []
        changed: []
        deprecated: []
        removed: []
        fixed: []
        security: []
    "1.0.0":
      date: "2026-01-15"             # Required for released versions
      added:
        - "Initial release"
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
