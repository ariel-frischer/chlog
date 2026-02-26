package changelog

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultCategories are the standard Keep a Changelog categories.
var DefaultCategories = []string{"added", "changed", "deprecated", "removed", "fixed", "security"}

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
		keyNode := versionKeyNode(v.Version)

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

	keyNode := versionKeyNode(v.Version)
	var valNode yaml.Node
	if err := valNode.Encode(*v); err != nil {
		return nil, fmt.Errorf("encoding version %s: %w", v.Version, err)
	}
	wrapper.Content = append(wrapper.Content, keyNode, &valNode)

	return yaml.Marshal(wrapper)
}

// versionKeyNode creates a YAML scalar node for a version key.
// Only forces !!str tag when the value would be misinterpreted by YAML
// (e.g. "1.0" as float, "1" as int, "true" as bool).
func versionKeyNode(version string) *yaml.Node {
	node := &yaml.Node{Kind: yaml.ScalarNode, Value: version}
	var parsed interface{}
	if err := yaml.Unmarshal([]byte(version), &parsed); err == nil {
		if _, ok := parsed.(string); !ok {
			node.Tag = "!!str"
		}
	}
	return node
}

// CategoryEntry holds entries for a single changelog category.
type CategoryEntry struct {
	Name    string
	Entries []string
}

// Changes is an ordered collection of changelog categories.
type Changes struct {
	Categories []CategoryEntry
}

// Get returns the entries for a category, or nil if not found.
func (c Changes) Get(category string) []string {
	for _, cat := range c.Categories {
		if cat.Name == category {
			return cat.Entries
		}
	}
	return nil
}

// Append adds an entry to a category, creating the category if needed.
func (c *Changes) Append(category, entry string) {
	for i := range c.Categories {
		if c.Categories[i].Name == category {
			c.Categories[i].Entries = append(c.Categories[i].Entries, entry)
			return
		}
	}
	c.Categories = append(c.Categories, CategoryEntry{Name: category, Entries: []string{entry}})
}

// Remove removes an entry from the given category.
// If substring is true, performs case-insensitive substring matching.
// Returns the removed entry text on success.
func (c *Changes) Remove(category, text string, substring bool) (string, error) {
	catIdx := -1
	for i := range c.Categories {
		if c.Categories[i].Name == category {
			catIdx = i
			break
		}
	}
	if catIdx == -1 {
		return "", CategoryNotFoundError{Category: category}
	}

	entries := c.Categories[catIdx].Entries

	if substring {
		return c.removeSubstring(catIdx, entries, category, text)
	}
	return c.removeExact(catIdx, entries, category, text)
}

func (c *Changes) removeExact(catIdx int, entries []string, category, text string) (string, error) {
	for i, e := range entries {
		if e == text {
			c.Categories[catIdx].Entries = append(entries[:i], entries[i+1:]...)
			c.cleanupEmpty(catIdx)
			return e, nil
		}
	}
	return "", EntryNotFoundError{Category: category, Text: text}
}

func (c *Changes) removeSubstring(catIdx int, entries []string, category, text string) (string, error) {
	lower := strings.ToLower(text)
	var matches []string
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e), lower) {
			matches = append(matches, e)
		}
	}

	switch len(matches) {
	case 0:
		return "", EntryNotFoundError{Category: category, Text: text}
	case 1:
		return c.removeExact(catIdx, entries, category, matches[0])
	default:
		return "", MultipleMatchError{Category: category, Text: text, Matches: matches}
	}
}

func (c *Changes) cleanupEmpty(catIdx int) {
	if len(c.Categories[catIdx].Entries) == 0 {
		c.Categories = append(c.Categories[:catIdx], c.Categories[catIdx+1:]...)
	}
}

// Merge appends all entries from other into c, preserving order.
func (c *Changes) Merge(other Changes) {
	for _, cat := range other.Categories {
		for _, entry := range cat.Entries {
			c.Append(cat.Name, entry)
		}
	}
}

// Clone returns a deep copy of the Changes.
func (c Changes) Clone() Changes {
	clone := Changes{Categories: make([]CategoryEntry, len(c.Categories))}
	for i, cat := range c.Categories {
		entries := make([]string, len(cat.Entries))
		copy(entries, cat.Entries)
		clone.Categories[i] = CategoryEntry{Name: cat.Name, Entries: entries}
	}
	return clone
}

// IsEmpty returns true if there are no entries in any category.
func (c Changes) IsEmpty() bool {
	return c.Count() == 0
}

// Count returns the total number of entries across all categories.
func (c Changes) Count() int {
	n := 0
	for _, cat := range c.Categories {
		n += len(cat.Entries)
	}
	return n
}

// CategoryNames returns the names of all categories in order.
func (c Changes) CategoryNames() []string {
	names := make([]string, len(c.Categories))
	for i, cat := range c.Categories {
		names[i] = cat.Name
	}
	return names
}

// UnmarshalYAML parses a YAML mapping where each key is a category name.
func (c *Changes) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("changes: expected mapping, got %d", value.Kind)
	}
	for i := 0; i < len(value.Content)-1; i += 2 {
		key := value.Content[i].Value
		var entries []string
		if err := value.Content[i+1].Decode(&entries); err != nil {
			return fmt.Errorf("changes.%s: %w", key, err)
		}
		c.Categories = append(c.Categories, CategoryEntry{Name: key, Entries: entries})
	}
	return nil
}

// MarshalYAML emits a YAML mapping with each category as a key.
func (c Changes) MarshalYAML() (interface{}, error) {
	node := &yaml.Node{Kind: yaml.MappingNode}
	for _, cat := range c.Categories {
		if len(cat.Entries) == 0 {
			continue
		}
		var entriesNode yaml.Node
		if err := entriesNode.Encode(cat.Entries); err != nil {
			return nil, fmt.Errorf("encoding %s: %w", cat.Name, err)
		}
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: cat.Name},
			&entriesNode,
		)
	}
	return node, nil
}

// Version represents a single version entry in the changelog.
type Version struct {
	Version  string  `yaml:"-"`
	Date     string  `yaml:"-"`
	Public   Changes `yaml:"-"`
	Internal Changes `yaml:"-"`
}

// MergedChanges returns Changes with internal entries merged into a clone of public.
func (v *Version) MergedChanges() Changes {
	merged := v.Public.Clone()
	merged.Merge(v.Internal)
	return merged
}

// IsEmpty returns true if all public category lists are empty.
func (v *Version) IsEmpty() bool {
	return v.Public.IsEmpty()
}

// Count returns the total number of public entries.
func (v *Version) Count() int {
	return v.Public.Count()
}

// IsUnreleased returns true if this version represents unreleased changes.
func (v *Version) IsUnreleased() bool {
	return strings.EqualFold(v.Version, "unreleased")
}

// UnmarshalYAML parses a version node where "date" and "internal" are special keys,
// and everything else is a public category.
func (v *Version) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("version: expected mapping, got %d", value.Kind)
	}

	for i := 0; i < len(value.Content)-1; i += 2 {
		key := value.Content[i].Value
		val := value.Content[i+1]

		switch key {
		case "date":
			v.Date = val.Value
		case "internal":
			if err := val.Decode(&v.Internal); err != nil {
				return fmt.Errorf("version.internal: %w", err)
			}
		default:
			// Everything else is a public category
			var entries []string
			if err := val.Decode(&entries); err != nil {
				return fmt.Errorf("version.%s: %w", key, err)
			}
			v.Public.Categories = append(v.Public.Categories, CategoryEntry{
				Name:    key,
				Entries: entries,
			})
		}
	}
	return nil
}

// MarshalYAML emits date, then public categories, then internal if non-empty.
func (v Version) MarshalYAML() (interface{}, error) {
	node := &yaml.Node{Kind: yaml.MappingNode}

	if v.Date != "" {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "date"},
			&yaml.Node{Kind: yaml.ScalarNode, Value: v.Date},
		)
	}

	// Public categories
	for _, cat := range v.Public.Categories {
		if len(cat.Entries) == 0 {
			continue
		}
		var entriesNode yaml.Node
		if err := entriesNode.Encode(cat.Entries); err != nil {
			return nil, fmt.Errorf("encoding %s: %w", cat.Name, err)
		}
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: cat.Name},
			&entriesNode,
		)
	}

	// Internal
	if !v.Internal.IsEmpty() {
		var internalNode yaml.Node
		if err := internalNode.Encode(v.Internal); err != nil {
			return nil, fmt.Errorf("encoding internal: %w", err)
		}
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "internal"},
			&internalNode,
		)
	}

	return node, nil
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

// CategoryNotFoundError indicates a requested category does not exist.
type CategoryNotFoundError struct {
	Category string
}

func (e CategoryNotFoundError) Error() string {
	return fmt.Sprintf("category %q not found", e.Category)
}

// EntryNotFoundError indicates a requested entry does not exist in a category.
type EntryNotFoundError struct {
	Category string
	Text     string
}

func (e EntryNotFoundError) Error() string {
	return fmt.Sprintf("entry %q not found in category %q", e.Text, e.Category)
}

// MultipleMatchError indicates multiple entries matched a substring search.
type MultipleMatchError struct {
	Category string
	Text     string
	Matches  []string
}

func (e MultipleMatchError) Error() string {
	return fmt.Sprintf("multiple entries in %q match %q: %v", e.Category, e.Text, e.Matches)
}
