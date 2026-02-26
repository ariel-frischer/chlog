package changelog

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Changelog is the root structure for a YAML changelog file.
type Changelog struct {
	Project  string    `yaml:"project"`
	Versions []Version `yaml:"-"`
}

// UnmarshalYAML implements custom YAML unmarshaling for map-keyed versions format.
func (c *Changelog) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping node, got %d", value.Kind)
	}

	for i := 0; i < len(value.Content)-1; i += 2 {
		key := value.Content[i].Value
		val := value.Content[i+1]

		switch key {
		case "project":
			c.Project = val.Value
		case "versions":
			if val.Kind != yaml.MappingNode {
				return fmt.Errorf("versions: expected mapping, got %d", val.Kind)
			}
			seen := map[string]bool{}
			for j := 0; j < len(val.Content)-1; j += 2 {
				versionKey := val.Content[j].Value
				versionVal := val.Content[j+1]

				if seen[versionKey] {
					return fmt.Errorf("versions: duplicate key %q", versionKey)
				}
				seen[versionKey] = true

				var v Version
				if err := versionVal.Decode(&v); err != nil {
					return fmt.Errorf("versions.%s: %w", versionKey, err)
				}
				v.Version = versionKey
				c.Versions = append(c.Versions, v)
			}
		default:
			return fmt.Errorf("unknown field %q", key)
		}
	}
	return nil
}

// MarshalYAML implements custom YAML marshaling for map-keyed versions format.
func (c Changelog) MarshalYAML() (interface{}, error) {
	root := &yaml.Node{Kind: yaml.MappingNode}

	// project key
	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "project"},
		&yaml.Node{Kind: yaml.ScalarNode, Value: c.Project},
	)

	// versions key
	versionsMap := &yaml.Node{Kind: yaml.MappingNode}
	for _, v := range c.Versions {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: v.Version, Tag: "!!str"}

		var valNode yaml.Node
		if err := valNode.Encode(v); err != nil {
			return nil, fmt.Errorf("encoding version %s: %w", v.Version, err)
		}
		versionsMap.Content = append(versionsMap.Content, keyNode, &valNode)
	}

	root.Content = append(root.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "versions"},
		versionsMap,
	)

	return root, nil
}

// MarshalVersionEntry marshals a single version as a map-keyed YAML entry.
// Used by scaffold dry-run to produce output like "unreleased:\n  added:\n    - ...".
func MarshalVersionEntry(v *Version) ([]byte, error) {
	wrapper := &yaml.Node{Kind: yaml.MappingNode}

	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: v.Version, Tag: "!!str"}
	var valNode yaml.Node
	if err := valNode.Encode(*v); err != nil {
		return nil, fmt.Errorf("encoding version %s: %w", v.Version, err)
	}
	wrapper.Content = append(wrapper.Content, keyNode, &valNode)

	return yaml.Marshal(wrapper)
}

// Version represents a single version entry in the changelog.
type Version struct {
	Version    string   `yaml:"-"`
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
