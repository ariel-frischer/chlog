package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ariel-frischer/chlog/pkg/changelog"
)

var syncInternal bool

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Generate CHANGELOG.md from CHANGELOG.yaml",
	RunE:  runSync,
}

func init() {
	syncCmd.Flags().BoolVar(&syncInternal, "internal", false, "include internal entries")
}

func runSync(cmd *cobra.Command, args []string) error {
	c, err := changelog.Load(yamlFile)
	if err != nil {
		return err
	}

	rendered, err := changelog.RenderMarkdownString(c, changelog.RenderOptions{
		IncludeInternal: syncInternal,
		Config:          loadConfig(),
	})
	if err != nil {
		return fmt.Errorf("rendering markdown: %w", err)
	}

	existing, _ := os.ReadFile(defaultMDFile)
	if bytes.Equal(existing, []byte(rendered)) {
		fmt.Println("CHANGELOG.md is up to date")
		return nil
	}

	if err := os.WriteFile(defaultMDFile, []byte(rendered), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", defaultMDFile, err)
	}

	fmt.Printf("Generated %s\n", defaultMDFile)
	return nil
}
