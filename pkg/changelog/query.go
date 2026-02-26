package changelog

// GetVersion returns a version by its identifier, normalizing "v" prefix.
func (c *Changelog) GetVersion(version string) (*Version, error) {
	normalized := NormalizeVersion(version)
	for i := range c.Versions {
		if NormalizeVersion(c.Versions[i].Version) == normalized {
			return &c.Versions[i], nil
		}
	}
	return nil, VersionNotFoundError{Version: version}
}

// GetUnreleased returns the unreleased version, or nil if none exists.
func (c *Changelog) GetUnreleased() *Version {
	for i := range c.Versions {
		if c.Versions[i].IsUnreleased() {
			return &c.Versions[i]
		}
	}
	return nil
}

// GetLatestRelease returns the most recent released version (skipping unreleased).
func (c *Changelog) GetLatestRelease() *Version {
	for i := range c.Versions {
		if !c.Versions[i].IsUnreleased() {
			return &c.Versions[i]
		}
	}
	return nil
}

// ListVersions returns all version identifiers in order.
func (c *Changelog) ListVersions() []string {
	versions := make([]string, len(c.Versions))
	for i, v := range c.Versions {
		versions[i] = v.Version
	}
	return versions
}

// QueryOptions controls which entries are included in query results.
type QueryOptions struct {
	IncludeInternal bool
}

// GetLastN returns the first n entries across all versions, newest first.
func (c *Changelog) GetLastN(n int, opts ...QueryOptions) []Entry {
	all := c.AllEntries(opts...)
	if n >= len(all) {
		return all
	}
	return all[:n]
}

// AllEntries returns all entries flattened across all versions, newest first.
func (c *Changelog) AllEntries(opts ...QueryOptions) []Entry {
	var opt QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	var entries []Entry
	for _, v := range c.Versions {
		entries = append(entries, flattenChanges(&v, opt.IncludeInternal)...)
	}
	return entries
}

// GetVersionCount returns the total number of versions.
func (c *Changelog) GetVersionCount() int {
	return len(c.Versions)
}

// GetEntryCount returns the total number of entries across all versions.
func (c *Changelog) GetEntryCount(opts ...QueryOptions) int {
	var opt QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	count := 0
	for _, v := range c.Versions {
		if opt.IncludeInternal {
			count += v.MergedChanges().Count()
		} else {
			count += v.Count()
		}
	}
	return count
}

// HasUnreleased returns true if there is an unreleased version.
func (c *Changelog) HasUnreleased() bool {
	return c.GetUnreleased() != nil
}

// flattenChanges converts a version's changes into a flat entry slice.
func flattenChanges(v *Version, includeInternal bool) []Entry {
	changes := v.Public
	if includeInternal {
		changes = v.MergedChanges()
	}
	var entries []Entry
	for _, cat := range changes.Categories {
		for _, text := range cat.Entries {
			entries = append(entries, Entry{
				Text:     text,
				Category: cat.Name,
				Version:  v.Version,
			})
		}
	}
	return entries
}
