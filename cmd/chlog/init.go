package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/spf13/cobra"
)

var projectName string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new CHANGELOG.yaml",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().StringVar(&projectName, "project", "", "project name (default: directory name)")
}

func runInit(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(yamlFile); err == nil {
		return fmt.Errorf("%s already exists", yamlFile)
	}

	name := projectName
	if name == "" {
		name = promptProjectName()
	}

	c := &changelog.Changelog{
		Project: name,
		Versions: []changelog.Version{
			{
				Version: "unreleased",
				Added:   []string{"Initial project setup"},
			},
		},
	}

	if err := changelog.Save(c, yamlFile); err != nil {
		return fmt.Errorf("creating %s: %w", yamlFile, err)
	}
	fmt.Printf("Created %s for project %q\n", yamlFile, name)

	// Create .chlog.yaml config with auto-detected repo URL.
	configPath := changelog.DefaultConfigFile
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("%s already exists, skipping config creation\n", configPath)
		return nil
	}

	cfg := &changelog.Config{}
	if url, err := changelog.DetectRepoURL(); err == nil {
		cfg.RepoURL = url
		fmt.Printf("Detected repo URL: %s\n", url)
	}

	if err := changelog.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("creating %s: %w", configPath, err)
	}
	fmt.Printf("Created %s\n", configPath)
	return nil
}

func promptProjectName() string {
	dir, _ := os.Getwd()
	defaultName := filepath.Base(dir)

	fmt.Printf("Project name [%s]: ", defaultName)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultName
	}
	return input
}
