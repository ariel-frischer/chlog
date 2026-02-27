package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/spf13/cobra"
)

var (
	removeVersion  string
	removeInternal bool
	removeMatch    bool
)

var removeCmd = &cobra.Command{
	Use:   "remove <category> <entry>",
	Short: "Remove an entry from the changelog",
	Long:  "Remove a single entry from a category in the changelog.",
	Example: `  chlog remove added "Support dark mode"
  chlog remove added --match "dark mode"
  chlog remove fixed --version 1.2.0 "Fix login timeout"
  chlog remove changed --internal "Refactor auth"`,
	Args: cobra.ExactArgs(2),
	RunE: runRemove,
}

func init() {
	removeCmd.Flags().StringVarP(&removeVersion, "version", "v", "unreleased", "target version")
	removeCmd.Flags().BoolVarP(&removeInternal, "internal", "i", false, "remove from internal entries")
	removeCmd.Flags().BoolVarP(&removeMatch, "match", "m", false, "use case-insensitive substring matching")
}

func runRemove(cmd *cobra.Command, args []string) error {
	category := strings.ToLower(strings.TrimSpace(args[0]))
	text := args[1]
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("entry text must not be empty")
	}

	c, err := changelog.Load(yamlFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s not found — run 'chlog init' first", yamlFile)
		}
		return err
	}

	v, err := c.GetVersion(removeVersion)
	if err != nil {
		return err
	}

	changes := &v.Public
	if removeInternal {
		changes = &v.Internal
	}

	removed, err := changes.Remove(category, text, removeMatch)
	if err != nil {
		return formatRemoveError(err)
	}

	if err := changelog.Save(c, yamlFile); err != nil {
		return fmt.Errorf("saving %s: %w", yamlFile, err)
	}

	label := "public"
	if removeInternal {
		label = "internal"
	}
	success("Removed %s %s entry from %s: %s", label, categoryRef(category), versionRef(v.Version), removed)
	return nil
}

func formatRemoveError(err error) error {
	switch e := err.(type) {
	case changelog.MultipleMatchError:
		var b strings.Builder
		fmt.Fprintf(&b, "multiple entries match %q in %s:\n", e.Text, categoryRef(e.Category))
		for _, m := range e.Matches {
			fmt.Fprintf(&b, "  - %s\n", highlight(m))
		}
		b.WriteString("use exact text to remove a specific entry")
		return fmt.Errorf("%s", b.String())
	default:
		return err
	}
}
