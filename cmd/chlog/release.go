package main

import (
	"fmt"
	"time"

	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/spf13/cobra"
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
	ver := args[0]

	c, err := changelog.Load(yamlFile)
	if err != nil {
		return err
	}

	date := releaseDate
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	if err := c.Release(ver, date); err != nil {
		return err
	}

	if err := changelog.Save(c, yamlFile); err != nil {
		return fmt.Errorf("saving %s: %w", yamlFile, err)
	}

	success("Released %s (%s) — unreleased block reset", versionRef(ver), highlight(date))
	return nil
}
