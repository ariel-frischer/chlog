package main

import (
	"fmt"

	"github.com/ariel-frischer/chlog/pkg/changelog"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate CHANGELOG.yaml schema",
	RunE:  runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	_, err := changelog.Load(yamlFile)
	if err != nil {
		return err
	}
	fmt.Println("CHANGELOG.yaml is valid")
	return nil
}
