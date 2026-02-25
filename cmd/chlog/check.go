package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/ariel-frischer/chlog/pkg/changelog"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Verify CHANGELOG.md matches CHANGELOG.yaml",
	Long:  "Exit 0 = in sync, exit 1 = out of sync, exit 2 = validation error.",
	RunE:  runCheck,
}

func runCheck(cmd *cobra.Command, args []string) error {
	c, err := changelog.Load(yamlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "validation error: %v\n", err)
		os.Exit(2)
	}

	rendered, err := changelog.RenderMarkdownString(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render error: %v\n", err)
		os.Exit(2)
	}

	existing, err := os.ReadFile(defaultMDFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s not found — run 'chlog sync' first\n", defaultMDFile)
		os.Exit(1)
	}

	if !bytes.Equal(existing, []byte(rendered)) {
		fmt.Fprintf(os.Stderr, "%s is out of sync — run 'chlog sync'\n", defaultMDFile)
		os.Exit(1)
	}

	fmt.Println("CHANGELOG.md is in sync")
	return nil
}
