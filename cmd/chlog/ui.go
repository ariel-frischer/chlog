package main

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	successFmt = color.New(color.FgGreen)
	warnFmt    = color.New(color.FgYellow)
	errFmt     = color.New(color.FgRed)
	boldFmt    = color.New(color.Bold)
	fileFmt    = color.New(color.FgCyan)
	versionFmt = color.New(color.FgMagenta, color.Bold)
)

// success prints a green success message.
func success(format string, a ...any) {
	fmt.Println(successFmt.Sprintf(format, a...))
}

// warn prints a yellow warning message.
func warn(format string, a ...any) {
	fmt.Println(warnFmt.Sprintf(format, a...))
}

// highlight returns a bold string.
func highlight(s string) string {
	return boldFmt.Sprint(s)
}

// fileRef returns a cyan-colored file path.
func fileRef(s string) string {
	return fileFmt.Sprint(s)
}

// versionRef returns a magenta bold version string.
func versionRef(s string) string {
	return versionFmt.Sprint(s)
}

// categoryRef returns a category name colored to match show output.
func categoryRef(category string) string {
	styles := map[string]*color.Color{
		"added":      color.New(color.FgGreen),
		"changed":    color.New(color.FgYellow),
		"deprecated": color.New(color.FgYellow),
		"removed":    color.New(color.FgRed),
		"fixed":      color.New(color.FgCyan),
		"security":   color.New(color.FgMagenta),
	}
	if c, ok := styles[category]; ok {
		return c.Sprint(category)
	}
	return category
}
