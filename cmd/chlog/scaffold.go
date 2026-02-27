package main

import (
	"fmt"
	"os"

	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/spf13/cobra"
)

var (
	scaffoldWrite   bool
	scaffoldVersion string
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Generate changelog entries from conventional commits",
	RunE:  runScaffold,
}

func init() {
	scaffoldCmd.Flags().BoolVar(&scaffoldWrite, "write", false, "merge into existing CHANGELOG.yaml")
	scaffoldCmd.Flags().StringVar(&scaffoldVersion, "version", "", "version string (default: unreleased)")
}

func runScaffold(cmd *cobra.Command, args []string) error {
	sinceTag, err := changelog.LatestTag()
	if err != nil {
		sinceTag = ""
	}

	commits, err := changelog.GitLog(sinceTag)
	if err != nil {
		return fmt.Errorf("reading git log: %w", err)
	}

	if len(commits) == 0 {
		warn("No commits found")
		return nil
	}

	v := changelog.Scaffold(commits, changelog.ScaffoldOptions{Version: scaffoldVersion})
	if v.IsEmpty() && v.Internal.IsEmpty() {
		warn("No conventional commits found")
		return nil
	}

	if !scaffoldWrite {
		data, err := changelog.MarshalVersionEntry(v)
		if err != nil {
			return fmt.Errorf("marshaling YAML: %w", err)
		}
		fmt.Print(string(data))
		return nil
	}

	return writeScaffold(v)
}

func writeScaffold(v *changelog.Version) error {
	c, err := changelog.Load(yamlFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s not found — run 'chlog init' first", yamlFile)
		}
		return err
	}

	existing := c.GetUnreleased()
	if existing != nil && v.Version == "unreleased" {
		existing.Public.Merge(v.Public)
		existing.Internal.Merge(v.Internal)
	} else {
		c.Versions = append([]changelog.Version{*v}, c.Versions...)
	}

	if err := changelog.Save(c, yamlFile); err != nil {
		return fmt.Errorf("saving %s: %w", yamlFile, err)
	}

	success("Updated %s with %d entries", fileRef(yamlFile), v.Count()+v.Internal.Count())
	return nil
}
