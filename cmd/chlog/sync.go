package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/spf13/cobra"
)

var (
	syncInternal bool
	syncSplit    bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Generate CHANGELOG.md from CHANGELOG.yaml",
	RunE:  runSync,
}

func init() {
	syncCmd.Flags().BoolVar(&syncInternal, "internal", false, "include internal entries")
	syncCmd.Flags().BoolVar(&syncSplit, "split", false, "generate both public and internal changelogs")
}

func runSync(cmd *cobra.Command, args []string) error {
	c, err := changelog.Load(yamlFile)
	if err != nil {
		return err
	}

	cfg := loadConfig()

	if syncSplit {
		return runSyncSplit(c, cfg)
	}

	internal := syncInternal || cfg.IncludeInternal
	return syncFile(c, cfg, internal, defaultMDFile)
}

func runSyncSplit(c *changelog.Changelog, cfg *changelog.Config) error {
	if err := syncFile(c, cfg, false, cfg.PublicFilePath()); err != nil {
		return err
	}
	return syncFile(c, cfg, true, cfg.InternalFilePath())
}

func syncFile(c *changelog.Changelog, cfg *changelog.Config, internal bool, path string) error {
	rendered, err := changelog.RenderMarkdownString(c, changelog.RenderOptions{
		IncludeInternal: internal,
		Config:          cfg,
	})
	if err != nil {
		return fmt.Errorf("rendering markdown: %w", err)
	}

	existing, _ := os.ReadFile(path)
	if bytes.Equal(existing, []byte(rendered)) {
		success("%s is up to date", fileRef(path))
		return nil
	}

	if err := os.WriteFile(path, []byte(rendered), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	success("Generated %s", fileRef(path))
	return nil
}
