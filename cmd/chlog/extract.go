package main

import (
	"os"

	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/spf13/cobra"
)

var extractInternal bool

var extractCmd = &cobra.Command{
	Use:   "extract <version>",
	Short: "Extract a single version as markdown",
	Args:  cobra.ExactArgs(1),
	RunE:  runExtract,
}

func init() {
	extractCmd.Flags().BoolVar(&extractInternal, "internal", false, "include internal entries")
}

func runExtract(cmd *cobra.Command, args []string) error {
	c, err := changelog.Load(yamlFile)
	if err != nil {
		return err
	}

	v, err := c.GetVersion(args[0])
	if err != nil {
		return err
	}

	cfg := loadConfig()
	internal := extractInternal || cfg.IncludeInternal

	return changelog.RenderVersionMarkdown(v, os.Stdout, changelog.RenderOptions{IncludeInternal: internal})
}
