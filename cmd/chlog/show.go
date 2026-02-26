package main

import (
	"fmt"

	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/spf13/cobra"
)

var (
	showLast     int
	showPlain    bool
	showInternal bool
)

var showCmd = &cobra.Command{
	Use:   "show [version]",
	Short: "Display changelog in terminal",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runShow,
}

func init() {
	showCmd.Flags().IntVarP(&showLast, "last", "n", 0, "show last N entries")
	showCmd.Flags().BoolVar(&showPlain, "plain", false, "disable colors and icons")
	showCmd.Flags().BoolVar(&showInternal, "internal", false, "include internal entries")
}

func runShow(cmd *cobra.Command, args []string) error {
	c, err := changelog.Load(yamlFile)
	if err != nil {
		return err
	}

	cfg := loadConfig()
	internal := showInternal || cfg.IncludeInternal
	opts := changelog.FormatOptions{Plain: showPlain, IncludeInternal: internal}

	if len(args) == 1 {
		v, err := c.GetVersion(args[0])
		if err != nil {
			return err
		}
		fmt.Print(changelog.FormatVersion(v, opts))
		return nil
	}

	if showLast > 0 {
		entries := c.GetLastN(showLast, changelog.QueryOptions{IncludeInternal: internal})
		for _, e := range entries {
			fmt.Printf("[%s] %s: %s\n", e.Version, e.Category, e.Text)
		}
		return nil
	}

	fmt.Print(changelog.FormatTerminal(c, opts))
	return nil
}
