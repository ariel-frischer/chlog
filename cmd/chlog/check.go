package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/spf13/cobra"
)

var (
	checkInternal bool
	checkSplit    bool
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Verify CHANGELOG.md matches CHANGELOG.yaml",
	Long:  "Exit 0 = in sync, exit 1 = out of sync, exit 2 = validation error.",
	RunE:  runCheck,
}

func init() {
	checkCmd.Flags().BoolVar(&checkInternal, "internal", false, "compare with internal entries included")
	checkCmd.Flags().BoolVar(&checkSplit, "split", false, "verify both public and internal changelogs")
}

func runCheck(cmd *cobra.Command, args []string) error {
	c, err := changelog.Load(yamlFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, errFmt.Sprintf("validation error: %v", err))
		os.Exit(2)
	}

	cfg := loadConfig()

	if checkSplit {
		return runCheckSplit(c, cfg)
	}

	internal := checkInternal || cfg.IncludeInternal
	return checkFile(c, cfg, internal, defaultMDFile)
}

func runCheckSplit(c *changelog.Changelog, cfg *changelog.Config) error {
	if err := checkFile(c, cfg, false, cfg.PublicFilePath()); err != nil {
		return err
	}
	return checkFile(c, cfg, true, cfg.InternalFilePath())
}

func checkFile(c *changelog.Changelog, cfg *changelog.Config, internal bool, path string) error {
	rendered, err := changelog.RenderMarkdownString(c, changelog.RenderOptions{
		IncludeInternal: internal,
		Config:          cfg,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, errFmt.Sprintf("render error: %v", err))
		os.Exit(2)
	}

	existing, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, errFmt.Sprintf("%s not found — run 'chlog sync' first", fileRef(path)))
		os.Exit(1)
	}

	if !bytes.Equal(existing, []byte(rendered)) {
		warn("%s is out of sync — run 'chlog sync'", fileRef(path))
		os.Exit(1)
	}

	success("%s is in sync", fileRef(path))
	return nil
}
