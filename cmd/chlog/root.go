package main

import (
	"os"

	"github.com/ariel-frischer/chlog/internal/version"
	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	defaultYAMLFile = "CHANGELOG.yaml"
	defaultMDFile   = "CHANGELOG.md"
)

var (
	yamlFile   string
	configFile string
)

var rootCmd = &cobra.Command{
	Use:     "chlog",
	Short:   "YAML-first changelog management",
	Long:    "Language-agnostic CLI for YAML-first changelog management.",
	Version: version.Version,
}

func init() {
	// Disable colors when not writing to a terminal.
	if fi, err := os.Stdout.Stat(); err == nil {
		if fi.Mode()&os.ModeCharDevice == 0 {
			color.NoColor = true
		}
	}

	var noColor bool
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if noColor {
			color.NoColor = true
		}
	}

	rootCmd.PersistentFlags().StringVarP(&yamlFile, "file", "f", defaultYAMLFile, "path to CHANGELOG.yaml")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", changelog.DefaultConfigFile, "path to config file")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(scaffoldCmd)
	rootCmd.AddCommand(releaseCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(versionCmd)
}

func loadConfig() *changelog.Config {
	cfg, err := changelog.LoadConfig(configFile)
	if err != nil {
		return &changelog.Config{}
	}
	return cfg
}
