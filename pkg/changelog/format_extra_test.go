package changelog

import (
	"strings"
	"testing"
)

func TestWrapText_EdgeCases(t *testing.T) {
	tests := map[string]struct {
		text     string
		maxWidth int
		indent   string
		want     string
	}{
		"zero width returns as-is": {
			text: "hello world", maxWidth: 0, indent: "  ",
			want: "hello world",
		},
		"negative width returns as-is": {
			text: "hello world", maxWidth: -1, indent: "  ",
			want: "hello world",
		},
		"single long word": {
			text: "superlongword", maxWidth: 5, indent: "  ",
			want: "superlongword",
		},
		"empty text": {
			text: "", maxWidth: 80, indent: "  ",
			want: "",
		},
		"whitespace only": {
			text: "   ", maxWidth: 80, indent: "  ",
			want: "   ",
		},
		"exact fit no wrap": {
			text: "abc def", maxWidth: 7, indent: "  ",
			want: "abc def",
		},
		"wraps at word boundary": {
			text: "abc def ghi", maxWidth: 7, indent: ">>",
			want: "abc def\n>>ghi",
		},
		"multiple wraps": {
			text: "one two three four", maxWidth: 9, indent: "  ",
			want: "one two\n  three\n  four",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := wrapText(tc.text, tc.maxWidth, tc.indent)
			if got != tc.want {
				t.Errorf("wrapText(%q, %d, %q) =\n  %q\nwant:\n  %q", tc.text, tc.maxWidth, tc.indent, got, tc.want)
			}
		})
	}
}

func TestFormatVersionHeader(t *testing.T) {
	tests := map[string]struct {
		version *Version
		want    string
	}{
		"released": {
			version: &Version{Version: "1.2.3", Date: "2024-03-15"},
			want:    "[1.2.3] - 2024-03-15",
		},
		"unreleased": {
			version: &Version{Version: "unreleased"},
			want:    "[Unreleased]",
		},
		"unreleased titlecase": {
			version: &Version{Version: "Unreleased"},
			want:    "[Unreleased]",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := formatVersionHeader(tc.version); got != tc.want {
				t.Errorf("formatVersionHeader() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestFormatTerminal_MultipleVersions(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{Version: "unreleased", Added: []string{"WIP"}},
			{Version: "1.0.0", Date: "2024-01-01", Added: []string{"Init"}},
		},
	}
	out := FormatTerminal(c, FormatOptions{Plain: true})

	if !strings.Contains(out, "[Unreleased]") {
		t.Error("missing unreleased header")
	}
	if !strings.Contains(out, "[1.0.0]") {
		t.Error("missing 1.0.0 header")
	}
	// Verify separator between versions
	unreleasedIdx := strings.Index(out, "[Unreleased]")
	releasedIdx := strings.Index(out, "[1.0.0]")
	if unreleasedIdx > releasedIdx {
		t.Error("unreleased should appear before released versions")
	}
}

func TestFormatTerminal_DefaultMaxWidth(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{Version: "1.0.0", Date: "2024-01-01", Added: []string{"Feature"}},
		},
	}
	// MaxWidth=0 should use default of 80
	out := FormatTerminal(c, FormatOptions{Plain: true, MaxWidth: 0})
	if !strings.Contains(out, "Feature") {
		t.Error("expected feature in output")
	}
}

func TestFormatVersion_EmptyChanges(t *testing.T) {
	v := &Version{
		Version: "unreleased",
	}
	out := FormatVersion(v, FormatOptions{Plain: true})
	if !strings.Contains(out, "[Unreleased]") {
		t.Error("expected unreleased header")
	}
	// Should only have the header line
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line for empty version, got %d: %v", len(lines), lines)
	}
}

func TestFormatVersion_CategoryOrder(t *testing.T) {
	v := &Version{
		Version:  "1.0.0",
		Date:     "2024-01-01",
		Security: []string{"Patched vuln"},
		Added:    []string{"New thing"},
		Fixed:    []string{"Bug fix"},
	}
	out := FormatVersion(v, FormatOptions{Plain: true})

	addedIdx := strings.Index(out, "Added")
	fixedIdx := strings.Index(out, "Fixed")
	secIdx := strings.Index(out, "Security")

	if addedIdx == -1 || fixedIdx == -1 || secIdx == -1 {
		t.Fatalf("missing categories in output:\n%s", out)
	}
	if addedIdx > fixedIdx || fixedIdx > secIdx {
		t.Error("categories should appear in canonical order: added < fixed < security")
	}
}
