package changelog

import (
	"regexp"
	"strings"
	"unicode"
)

// ScaffoldOptions controls scaffold behavior.
type ScaffoldOptions struct {
	Version string
}

var conventionalPattern = regexp.MustCompile(
	`^(\w+)(?:\([\w-]+\))?(!)?\s*:\s*(.+)$`,
)

// commitTypeMap maps conventional commit types to changelog categories.
var commitTypeMap = map[string]string{
	"feat":      "added",
	"fix":       "fixed",
	"refactor":  "changed",
	"perf":      "changed",
	"deprecate": "deprecated",
	"remove":    "removed",
}

// internalTypes are commit types routed to internal changes by default.
var internalTypes = map[string]bool{
	"refactor": true,
	"perf":     true,
}

// skippedTypes are commit types that don't belong in changelogs.
var skippedTypes = map[string]bool{
	"chore": true, "docs": true, "style": true,
	"test": true, "ci": true, "build": true,
}

// ParseConventionalCommit extracts category, description, breaking flag, and internal flag from a commit subject.
func ParseConventionalCommit(subject string) (category, description string, breaking, internal bool) {
	m := conventionalPattern.FindStringSubmatch(subject)
	if m == nil {
		return "", "", false, false
	}

	commitType := strings.ToLower(m[1])
	breaking = m[2] == "!"
	description = cleanDescription(m[3])

	if skippedTypes[commitType] && !breaking {
		return "", "", false, false
	}

	category = commitTypeMap[commitType]
	if category == "" && !breaking {
		return "", "", false, false
	}

	internal = internalTypes[commitType]

	if breaking {
		category = "changed"
		description = "BREAKING: " + description
		internal = false
	}

	return category, description, breaking, internal
}

// Scaffold creates a Version from a list of git commits.
func Scaffold(commits []GitCommit, opts ScaffoldOptions) *Version {
	version := opts.Version
	if version == "" {
		version = "unreleased"
	}

	v := &Version{Version: version}

	for _, c := range commits {
		cat, desc, _, isInternal := ParseConventionalCommit(c.Subject)
		if cat == "" {
			continue
		}
		if isInternal {
			appendToCategory(&v.Internal, cat, desc)
		} else {
			appendToCategory(&v.Changes, cat, desc)
		}
	}

	return v
}

func appendToCategory(c *Changes, category, entry string) {
	switch category {
	case "added":
		c.Added = append(c.Added, entry)
	case "changed":
		c.Changed = append(c.Changed, entry)
	case "deprecated":
		c.Deprecated = append(c.Deprecated, entry)
	case "removed":
		c.Removed = append(c.Removed, entry)
	case "fixed":
		c.Fixed = append(c.Fixed, entry)
	case "security":
		c.Security = append(c.Security, entry)
	}
}

func cleanDescription(s string) string {
	s = strings.TrimSpace(s)
	// Take first sentence only
	if idx := strings.IndexAny(s, ".!?"); idx != -1 {
		s = s[:idx]
	}
	s = strings.TrimSpace(s)
	// Capitalize first letter
	if len(s) > 0 {
		runes := []rune(s)
		runes[0] = unicode.ToUpper(runes[0])
		s = string(runes)
	}
	return s
}
