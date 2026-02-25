package changelog

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// CategoryStyle defines the icon and color for a category.
type CategoryStyle struct {
	Icon  string
	Color *color.Color
}

// FormatOptions controls terminal output formatting.
type FormatOptions struct {
	Plain    bool
	MaxWidth int
}

var categoryStyles = map[string]CategoryStyle{
	"added":      {Icon: "+", Color: color.New(color.FgGreen)},
	"changed":    {Icon: "~", Color: color.New(color.FgYellow)},
	"deprecated": {Icon: "!", Color: color.New(color.FgYellow)},
	"removed":    {Icon: "-", Color: color.New(color.FgRed)},
	"fixed":      {Icon: "x", Color: color.New(color.FgCyan)},
	"security":   {Icon: "🔒", Color: color.New(color.FgMagenta)},
}

func init() {
	// Disable colors when not writing to a terminal.
	if fi, err := os.Stdout.Stat(); err == nil {
		if fi.Mode()&os.ModeCharDevice == 0 {
			color.NoColor = true
		}
	}
}

// FormatTerminal formats the entire changelog for terminal output.
func FormatTerminal(c *Changelog, opts FormatOptions) string {
	if opts.MaxWidth == 0 {
		opts.MaxWidth = 80
	}
	var b strings.Builder
	bold := color.New(color.Bold)

	if opts.Plain {
		fmt.Fprintf(&b, "%s Changelog\n\n", c.Project)
	} else {
		fmt.Fprintf(&b, "%s\n\n", bold.Sprintf("%s Changelog", c.Project))
	}

	for i, v := range c.Versions {
		b.WriteString(FormatVersion(&v, opts))
		if i < len(c.Versions)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

// FormatVersion formats a single version for terminal output.
func FormatVersion(v *Version, opts FormatOptions) string {
	if opts.MaxWidth == 0 {
		opts.MaxWidth = 80
	}
	var b strings.Builder
	bold := color.New(color.Bold)

	header := formatVersionHeader(v)
	if opts.Plain {
		fmt.Fprintf(&b, "%s\n", header)
	} else {
		fmt.Fprintf(&b, "%s\n", bold.Sprint(header))
	}

	for _, cat := range ValidCategories() {
		entries := v.Changes.CategoryEntries(cat)
		if len(entries) == 0 {
			continue
		}

		style := categoryStyles[cat]
		catHeader := titleCase(cat)

		if opts.Plain {
			fmt.Fprintf(&b, "  %s %s\n", style.Icon, catHeader)
		} else {
			fmt.Fprintf(&b, "  %s\n", style.Color.Sprintf("%s %s", style.Icon, catHeader))
		}

		for _, entry := range entries {
			wrapped := wrapText(entry, opts.MaxWidth-6, "      ")
			fmt.Fprintf(&b, "    - %s\n", wrapped)
		}
	}

	return b.String()
}

func formatVersionHeader(v *Version) string {
	if v.IsUnreleased() {
		return "[Unreleased]"
	}
	return fmt.Sprintf("[%s] - %s", v.Version, v.Date)
}

// wrapText wraps text at word boundaries to fit within maxWidth, indenting continuation lines.
func wrapText(text string, maxWidth int, indent string) string {
	if maxWidth <= 0 || len(text) <= maxWidth {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	current := words[0]

	for _, word := range words[1:] {
		if len(current)+1+len(word) > maxWidth {
			lines = append(lines, current)
			current = word
		} else {
			current += " " + word
		}
	}
	lines = append(lines, current)

	return strings.Join(lines, "\n"+indent)
}
