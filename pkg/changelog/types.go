package changelog

import (
	"fmt"
	"strings"
)

// Changelog is the root structure for a YAML changelog file.
type Changelog struct {
	Project  string    `yaml:"project"`
	Versions []Version `yaml:"versions"`
}

// Version represents a single version entry in the changelog.
type Version struct {
	Version string  `yaml:"version"`
	Date    string  `yaml:"date,omitempty"`
	Changes Changes `yaml:"changes"`
}

// Changes groups entries by Keep a Changelog categories.
type Changes struct {
	Added      []string `yaml:"added,omitempty"`
	Changed    []string `yaml:"changed,omitempty"`
	Deprecated []string `yaml:"deprecated,omitempty"`
	Removed    []string `yaml:"removed,omitempty"`
	Fixed      []string `yaml:"fixed,omitempty"`
	Security   []string `yaml:"security,omitempty"`
}

// Entry is a flattened view of a single changelog entry with metadata.
type Entry struct {
	Text     string
	Category string
	Version  string
}

// ValidationError describes a validation failure with context.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// VersionNotFoundError indicates a requested version does not exist.
type VersionNotFoundError struct {
	Version string
}

func (e VersionNotFoundError) Error() string {
	return fmt.Sprintf("version %q not found", e.Version)
}

// IsUnreleased returns true if this version represents unreleased changes.
func (v *Version) IsUnreleased() bool {
	return strings.EqualFold(v.Version, "unreleased")
}

// IsEmpty returns true if all category lists are empty.
func (c Changes) IsEmpty() bool {
	return c.Count() == 0
}

// Count returns the total number of entries across all categories.
func (c Changes) Count() int {
	return len(c.Added) + len(c.Changed) + len(c.Deprecated) +
		len(c.Removed) + len(c.Fixed) + len(c.Security)
}

// ValidCategories returns category names in canonical order.
func ValidCategories() []string {
	return []string{"added", "changed", "deprecated", "removed", "fixed", "security"}
}

// CategoryEntries returns the entries for a given category name.
func (c Changes) CategoryEntries(category string) []string {
	switch category {
	case "added":
		return c.Added
	case "changed":
		return c.Changed
	case "deprecated":
		return c.Deprecated
	case "removed":
		return c.Removed
	case "fixed":
		return c.Fixed
	case "security":
		return c.Security
	default:
		return nil
	}
}
