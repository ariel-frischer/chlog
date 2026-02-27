package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/spf13/cobra"
)

var (
	addVersion  string
	addInternal bool
)

var addCmd = &cobra.Command{
	Use:   "add <category> [entries...]",
	Short: "Add entries to the changelog",
	Long:  "Add one or more entries to a category in the changelog.",
	Example: `  chlog add added "Support dark mode"
  chlog add fixed --version 1.2.0 "Fix login timeout"
  chlog add changed --internal "Refactor auth middleware"
  chlog add added "Feature A" "Feature B"`,
	Args: cobra.MinimumNArgs(2),
	RunE: runAdd,
}

func init() {
	addCmd.Flags().StringVarP(&addVersion, "version", "v", "unreleased", "target version")
	addCmd.Flags().BoolVarP(&addInternal, "internal", "i", false, "add as internal entry")
}

func runAdd(cmd *cobra.Command, args []string) error {
	category := strings.ToLower(strings.TrimSpace(args[0]))
	entries := args[1:]
	if err := validateCategory(category); err != nil {
		return err
	}

	for _, text := range entries {
		if strings.TrimSpace(text) == "" {
			return fmt.Errorf("entry text must not be empty")
		}
	}

	c, err := changelog.Load(yamlFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s not found — run 'chlog init' first", yamlFile)
		}
		return err
	}

	v, err := resolveVersionForAdd(c, addVersion)
	if err != nil {
		return err
	}

	changes := &v.Public
	if addInternal {
		changes = &v.Internal
	}

	for _, text := range entries {
		changes.Append(category, text)
	}

	if err := changelog.Save(c, yamlFile); err != nil {
		return fmt.Errorf("saving %s: %w", yamlFile, err)
	}

	label := "public"
	if addInternal {
		label = "internal"
	}
	success("Added %d %s %s entr%s to %s", len(entries), label, categoryRef(category), pluralY(len(entries)), versionRef(v.Version))
	return nil
}

func resolveVersionForAdd(c *changelog.Changelog, version string) (*changelog.Version, error) {
	if strings.EqualFold(version, "unreleased") {
		u := c.GetUnreleased()
		if u != nil {
			return u, nil
		}
		// Auto-create unreleased block
		c.Versions = append([]changelog.Version{{Version: "unreleased"}}, c.Versions...)
		return &c.Versions[0], nil
	}
	v, err := c.GetVersion(version)
	if err != nil {
		return nil, fmt.Errorf("version %q not found — can only add to existing versions", version)
	}
	return v, nil
}

func validateCategory(category string) error {
	cfg := loadConfig()
	allowed := cfg.AllowedCategories()
	if allowed == nil {
		return nil // non-strict mode
	}
	for _, a := range allowed {
		if a == category {
			return nil
		}
	}
	return fmt.Errorf("unknown category %q (allowed: %s)", category, strings.Join(allowed, ", "))
}

func pluralY(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}
