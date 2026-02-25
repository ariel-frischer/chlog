package main

import (
	"github.com/spf13/cobra"
)

const (
	defaultYAMLFile = "CHANGELOG.yaml"
	defaultMDFile   = "CHANGELOG.md"
)

var yamlFile string

var rootCmd = &cobra.Command{
	Use:     "chlog",
	Short:   "YAML-first changelog management",
	Long:    "Language-agnostic CLI for YAML-first changelog management.",
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&yamlFile, "file", "f", defaultYAMLFile, "path to CHANGELOG.yaml")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(scaffoldCmd)
}
