package changelog

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

// RenderOptions controls markdown rendering behavior.
type RenderOptions struct {
	IncludeInternal bool
	Config          *Config
}

// RenderMarkdown writes a full Keep a Changelog-compliant markdown document.
func RenderMarkdown(c *Changelog, w io.Writer, opts ...RenderOptions) error {
	var opt RenderOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	if _, err := fmt.Fprintf(w, "# Changelog\n\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "All notable changes to %s will be documented in this file.\n\n", c.Project); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).\n\n"); err != nil {
		return err
	}

	for _, v := range c.Versions {
		if err := RenderVersionMarkdown(&v, w, opt); err != nil {
			return err
		}
	}

	if repoURL := ResolveRepoURL(opt.Config); repoURL != "" {
		renderComparisonLinks(c, w, repoURL)
	}

	return nil
}

// RenderMarkdownString renders the changelog to a string.
func RenderMarkdownString(c *Changelog, opts ...RenderOptions) (string, error) {
	var b strings.Builder
	if err := RenderMarkdown(c, &b, opts...); err != nil {
		return "", err
	}
	return b.String(), nil
}

// RenderVersionMarkdown writes a single version as markdown.
func RenderVersionMarkdown(v *Version, w io.Writer, opts ...RenderOptions) error {
	var opt RenderOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	if v.IsUnreleased() {
		if _, err := fmt.Fprintf(w, "## [Unreleased]\n\n"); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(w, "## [%s] - %s\n\n", v.Version, v.Date); err != nil {
			return err
		}
	}

	changes := v.Public
	if opt.IncludeInternal {
		changes = v.MergedChanges()
	}

	for _, cat := range changes.Categories {
		if len(cat.Entries) == 0 {
			continue
		}
		if _, err := fmt.Fprintf(w, "### %s\n\n", titleCase(cat.Name)); err != nil {
			return err
		}
		for _, entry := range cat.Entries {
			if _, err := fmt.Fprintf(w, "- %s\n", entry); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
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
				_, _ = fmt.Fprintf(w, "[Unreleased]: %s%sv%s...HEAD\n", repoURL, comparePath, prev)
			}
		} else if i+1 < len(c.Versions) && !c.Versions[i+1].IsUnreleased() {
			prev := c.Versions[i+1].Version
			_, _ = fmt.Fprintf(w, "[%s]: %s%sv%s...v%s\n", v.Version, repoURL, comparePath, prev, v.Version)
		}
	}
}
