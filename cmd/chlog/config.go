package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/spf13/cobra"
)

var configForce bool

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage .chlog.yaml configuration",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold a .chlog.yaml config file",
	RunE:  runConfigInit,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show resolved configuration",
	RunE:  runConfigShow,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Long: `Set a config value in .chlog.yaml. Valid keys:
  repo_url          URL for commit links in CHANGELOG.md
  changelog_file    Path to CHANGELOG.yaml source (default: CHANGELOG.yaml)
  public_file       Output path for public CHANGELOG.md (default: CHANGELOG.md)
  internal_file     Output path for internal CHANGELOG (default: CHANGELOG-internal.md)
  include_internal  Include internal entries in public output (true/false)
  strict_categories Enforce allowed categories (true/false)
  categories        Comma-separated list of allowed categories`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open config in $EDITOR",
	RunE:  runConfigEdit,
}

func init() {
	configInitCmd.Flags().BoolVar(&configForce, "force", false, "overwrite existing config")
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configEditCmd)
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	path := configFile
	if _, err := os.Stat(path); err == nil && !configForce {
		return fmt.Errorf("%s already exists (use --force to overwrite)", path)
	}

	repoLine := "# repo_url: https://github.com/owner/repo\n"
	if url, err := changelog.DetectRepoURL(); err == nil {
		repoLine = fmt.Sprintf("repo_url: %s\n", url)
		fmt.Printf("Detected repo URL: %s\n", highlight(url))
	}

	content := "# chlog configuration — all fields are optional\n\n" +
		repoLine +
		"# changelog_file: CHANGELOG.yaml\n" +
		"# public_file: CHANGELOG.md\n" +
		"# internal_file: CHANGELOG-internal.md\n" +
		"# include_internal: false\n" +
		"# categories: [added, changed, deprecated, removed, fixed, security]\n" +
		"# strict_categories: true\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	success("Created %s", fileRef(path))
	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg := loadConfig()

	repoURL, repoSource := cfg.RepoURL, "file"
	if repoURL == "" {
		if url, err := changelog.DetectRepoURL(); err == nil {
			repoURL, repoSource = url, "detected"
		} else {
			repoURL, repoSource = "(none)", "default"
		}
	}

	strictStr, strictSource := "true", "default"
	if cfg.StrictCategories != nil {
		strictSource = "file"
		if !*cfg.StrictCategories {
			strictStr = "false"
		}
	}

	cats := cfg.AllowedCategories()
	catStr := strings.Join(cats, ", ")
	if len(cats) == 0 {
		catStr = "(any)"
	}

	changelogFile := cfg.ChangelogFile
	if changelogFile == "" {
		changelogFile = defaultYAMLFile
	}

	printConfigRow("repo_url", repoURL, repoSource)
	printConfigRow("changelog_file", changelogFile, sourceLabel(cfg.ChangelogFile != ""))
	printConfigRow("public_file", cfg.PublicFilePath(), sourceLabel(cfg.PublicFile != ""))
	printConfigRow("internal_file", cfg.InternalFilePath(), sourceLabel(cfg.InternalFile != ""))
	printConfigRow("include_internal", fmt.Sprintf("%v", cfg.IncludeInternal), sourceLabel(cfg.IncludeInternal))
	printConfigRow("strict_categories", strictStr, strictSource)
	printConfigRow("categories", catStr, sourceLabel(len(cfg.Categories) > 0))
	return nil
}

func printConfigRow(key, value, source string) {
	fmt.Printf("%-20s %-40s (%s)\n", highlight(key+":")+" ", value, source)
}

func sourceLabel(isCustom bool) string {
	if isCustom {
		return "file"
	}
	return "default"
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key, value := args[0], args[1]

	cfg, err := changelog.LoadConfig(configFile)
	if err != nil {
		return err
	}

	switch key {
	case "repo_url":
		cfg.RepoURL = value
	case "changelog_file":
		cfg.ChangelogFile = value
	case "public_file":
		cfg.PublicFile = value
	case "internal_file":
		cfg.InternalFile = value
	case "include_internal":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("include_internal expects true/false, got %q", value)
		}
		cfg.IncludeInternal = b
	case "strict_categories":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("strict_categories expects true/false, got %q", value)
		}
		cfg.StrictCategories = &b
	case "categories":
		parts := strings.Split(value, ",")
		cats := make([]string, 0, len(parts))
		for _, p := range parts {
			if t := strings.TrimSpace(p); t != "" {
				cats = append(cats, t)
			}
		}
		cfg.Categories = cats
	default:
		return fmt.Errorf("unknown key %q\nvalid keys: repo_url, changelog_file, public_file, internal_file, include_internal, strict_categories, categories", key)
	}

	if err := changelog.SaveConfig(cfg, configFile); err != nil {
		return err
	}
	success("Set %s = %s in %s", highlight(key), highlight(value), fileRef(configFile))
	return nil
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		warn("%s not found — run `chlog config init` first", fileRef(configFile))
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	c := exec.Command(editor, configFile)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
