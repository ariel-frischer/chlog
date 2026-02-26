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
	Version    string   `yaml:"version"`
	Date       string   `yaml:"date,omitempty"`
	Added      []string `yaml:"added,omitempty"`
	Changed    []string `yaml:"changed,omitempty"`
	Deprecated []string `yaml:"deprecated,omitempty"`
	Removed    []string `yaml:"removed,omitempty"`
	Fixed      []string `yaml:"fixed,omitempty"`
	Security   []string `yaml:"security,omitempty"`
	Internal   Changes  `yaml:"internal,omitempty"`
}

// Changes returns a Changes value built from the version's direct category fields.
func (v *Version) Changes() Changes {
	return Changes{
		Added:      v.Added,
		Changed:    v.Changed,
		Deprecated: v.Deprecated,
		Removed:    v.Removed,
		Fixed:      v.Fixed,
		Security:   v.Security,
	}
}

// MergedChanges returns Changes with internal entries merged in.
func (v *Version) MergedChanges() Changes {
	return Changes{
		Added:      append(append([]string{}, v.Added...), v.Internal.Added...),
		Changed:    append(append([]string{}, v.Changed...), v.Internal.Changed...),
		Deprecated: append(append([]string{}, v.Deprecated...), v.Internal.Deprecated...),
		Removed:    append(append([]string{}, v.Removed...), v.Internal.Removed...),
		Fixed:      append(append([]string{}, v.Fixed...), v.Internal.Fixed...),
		Security:   append(append([]string{}, v.Security...), v.Internal.Security...),
	}
}

// IsEmpty returns true if all category lists on this version are empty.
func (v *Version) IsEmpty() bool {
	return v.Changes().IsEmpty()
}

// Count returns the total number of entries across all categories on this version.
func (v *Version) Count() int {
	return v.Changes().Count()
}

// CategoryEntries returns the entries for a given category name on this version.
func (v *Version) CategoryEntries(category string) []string {
	return v.Changes().CategoryEntries(category)
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
