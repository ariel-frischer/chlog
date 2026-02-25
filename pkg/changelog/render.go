package changelog

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

// RenderMarkdown writes a full Keep a Changelog-compliant markdown document.
func RenderMarkdown(c *Changelog, w io.Writer) error {
	fmt.Fprintf(w, "# Changelog\n\n")
	fmt.Fprintf(w, "All notable changes to %s will be documented in this file.\n\n", c.Project)
	fmt.Fprintf(w, "The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).\n\n")

	for _, v := range c.Versions {
		if err := RenderVersionMarkdown(&v, w); err != nil {
			return err
		}
	}

	repoURL, err := DetectRepoURL()
	if err == nil && repoURL != "" {
		renderComparisonLinks(c, w, repoURL)
	}

	return nil
}

// RenderMarkdownString renders the changelog to a string.
func RenderMarkdownString(c *Changelog) (string, error) {
	var b strings.Builder
	if err := RenderMarkdown(c, &b); err != nil {
		return "", err
	}
	return b.String(), nil
}

// RenderVersionMarkdown writes a single version as markdown.
func RenderVersionMarkdown(v *Version, w io.Writer) error {
	if v.IsUnreleased() {
		fmt.Fprintf(w, "## [Unreleased]\n\n")
	} else {
		fmt.Fprintf(w, "## [%s] - %s\n\n", v.Version, v.Date)
	}

	for _, cat := range ValidCategories() {
		entries := v.Changes.CategoryEntries(cat)
		if len(entries) == 0 {
			continue
		}
		fmt.Fprintf(w, "### %s\n\n", titleCase(cat))
		for _, entry := range entries {
			fmt.Fprintf(w, "- %s\n", entry)
		}
		fmt.Fprintln(w)
	}

	return nil
}

func titleCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func renderComparisonLinks(c *Changelog, w io.Writer, repoURL string) {
	if len(c.Versions) == 0 {
		return
	}

	isGitLab := strings.Contains(repoURL, "gitlab")
	comparePath := "/compare/"
	if isGitLab {
		comparePath = "/-/compare/"
	}

	for i, v := range c.Versions {
		if v.IsUnreleased() {
			if i+1 < len(c.Versions) {
				prev := c.Versions[i+1].Version
				fmt.Fprintf(w, "[Unreleased]: %s%sv%s...HEAD\n", repoURL, comparePath, prev)
			}
		} else if i+1 < len(c.Versions) && !c.Versions[i+1].IsUnreleased() {
			prev := c.Versions[i+1].Version
			fmt.Fprintf(w, "[%s]: %s%sv%s...v%s\n", v.Version, repoURL, comparePath, prev, v.Version)
		}
	}
}
