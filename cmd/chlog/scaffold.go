package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/ariel-frischer/chlog/pkg/changelog"
	"gopkg.in/yaml.v3"
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
		fmt.Println("No commits found")
		return nil
	}

	v := changelog.Scaffold(commits, changelog.ScaffoldOptions{Version: scaffoldVersion})
	if v.Changes.IsEmpty() && v.Internal.IsEmpty() {
		fmt.Println("No conventional commits found")
		return nil
	}

	if !scaffoldWrite {
		data, err := yaml.Marshal(v)
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
		mergeChanges(&existing.Changes, &v.Changes)
		mergeChanges(&existing.Internal, &v.Internal)
	} else {
		c.Versions = append([]changelog.Version{*v}, c.Versions...)
	}

	if err := changelog.Save(c, yamlFile); err != nil {
		return fmt.Errorf("saving %s: %w", yamlFile, err)
	}

	fmt.Printf("Updated %s with %d entries\n", yamlFile, v.Changes.Count()+v.Internal.Count())
	return nil
}

func mergeChanges(dst, src *changelog.Changes) {
	dst.Added = append(dst.Added, src.Added...)
	dst.Changed = append(dst.Changed, src.Changed...)
	dst.Deprecated = append(dst.Deprecated, src.Deprecated...)
	dst.Removed = append(dst.Removed, src.Removed...)
	dst.Fixed = append(dst.Fixed, src.Fixed...)
	dst.Security = append(dst.Security, src.Security...)
}
