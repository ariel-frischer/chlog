package changelog

import "fmt"

// Release promotes the unreleased version to a named release with the given
// version string and date (YYYY-MM-DD), then prepends a fresh unreleased block.
func (c *Changelog) Release(version, date string) error {
	unreleased := c.GetUnreleased()
	if unreleased == nil {
		return fmt.Errorf("no unreleased version found")
	}

	if unreleased.IsEmpty() && unreleased.Internal.IsEmpty() {
		return fmt.Errorf("unreleased version has no entries")
	}

	if _, err := c.GetVersion(version); err == nil {
		return fmt.Errorf("version %q already exists", version)
	}

	unreleased.Version = version
	unreleased.Date = date

	c.Versions = append([]Version{{
		Version: "unreleased",
	}}, c.Versions...)

	return nil
}
