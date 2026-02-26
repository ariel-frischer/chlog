package changelog

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type categoryStyle struct {
	Icon  string
	Color *color.Color
}

// FormatOptions controls terminal output formatting.
type FormatOptions struct {
	Plain           bool
	MaxWidth        int
	IncludeInternal bool
}

var categoryStyles = map[string]categoryStyle{
	"added":      {Icon: "+", Color: color.New(color.FgGreen)},
	"changed":    {Icon: "~", Color: color.New(color.FgYellow)},
	"deprecated": {Icon: "!", Color: color.New(color.FgYellow)},
	"removed":    {Icon: "-", Color: color.New(color.FgRed)},
	"fixed":      {Icon: "x", Color: color.New(color.FgCyan)},
	"security":   {Icon: "🔒", Color: color.New(color.FgMagenta)},
}

var defaultStyle = categoryStyle{Icon: "*", Color: color.New(color.FgWhite)}

func styleFor(category string) categoryStyle {
	if s, ok := categoryStyles[category]; ok {
		return s
	}
	return defaultStyle
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

	changes := v.Public
	if opts.IncludeInternal {
		changes = v.MergedChanges()
	}

	for _, cat := range changes.Categories {
		if len(cat.Entries) == 0 {
			continue
		}

		style := styleFor(cat.Name)
		catHeader := titleCase(cat.Name)

		if opts.Plain {
			fmt.Fprintf(&b, "  %s %s\n", style.Icon, catHeader)
		} else {
			fmt.Fprintf(&b, "  %s\n", style.Color.Sprintf("%s %s", style.Icon, catHeader))
		}

		for _, entry := range cat.Entries {
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
