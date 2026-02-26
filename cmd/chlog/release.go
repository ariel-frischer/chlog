package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/ariel-frischer/chlog/pkg/changelog"
)

var releaseDate string

var releaseCmd = &cobra.Command{
	Use:   "release <version>",
	Short: "Promote unreleased changes to a versioned release",
	Long:  "Stamps the unreleased block with the given version and today's date, then creates a fresh unreleased block.",
	Args:  cobra.ExactArgs(1),
	RunE:  runRelease,
}

func init() {
	releaseCmd.Flags().StringVar(&releaseDate, "date", "", "release date in YYYY-MM-DD format (default: today)")
}

func runRelease(cmd *cobra.Command, args []string) error {
	version := args[0]

	c, err := changelog.Load(yamlFile)
	if err != nil {
		return err
	}

	unreleased := c.GetUnreleased()
	if unreleased == nil {
		return fmt.Errorf("no unreleased version found")
	}

	if unreleased.Changes.IsEmpty() && unreleased.Internal.IsEmpty() {
		return fmt.Errorf("unreleased version has no entries")
	}

	// Check for duplicate version
	if _, err := c.GetVersion(version); err == nil {
		return fmt.Errorf("version %q already exists", version)
	}

	date := releaseDate
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// Stamp unreleased → versioned
	unreleased.Version = version
	unreleased.Date = date

	// Prepend fresh unreleased block
	c.Versions = append([]changelog.Version{{
		Version: "unreleased",
		Changes: changelog.Changes{},
	}}, c.Versions...)

	if err := changelog.Save(c, yamlFile); err != nil {
		return fmt.Errorf("saving %s: %w", yamlFile, err)
	}

	fmt.Printf("Released %s (%s) — unreleased block reset\n", version, date)
	return nil
}
