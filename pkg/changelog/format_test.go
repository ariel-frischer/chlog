package changelog

import (
	"strings"
	"testing"
)

func TestFormatVersion_Plain(t *testing.T) {
	v := &Version{
		Version: "1.0.0",
		Date:    "2024-01-01",
		Added:   []string{"Feature A"},
		Fixed:   []string{"Bug B"},
	}
	out := FormatVersion(v, FormatOptions{Plain: true})

	if !strings.Contains(out, "[1.0.0] - 2024-01-01") {
		t.Error("expected version header")
	}
	if !strings.Contains(out, "+ Added") {
		t.Error("expected Added with icon")
	}
	if !strings.Contains(out, "x Fixed") {
		t.Error("expected Fixed with icon")
	}
	if !strings.Contains(out, "- Feature A") {
		t.Error("expected entry A")
	}
}

func TestFormatVersion_Unreleased(t *testing.T) {
	v := &Version{
		Version: "unreleased",
		Added:   []string{"WIP"},
	}
	out := FormatVersion(v, FormatOptions{Plain: true})
	if !strings.Contains(out, "[Unreleased]") {
		t.Error("expected [Unreleased] header")
	}
}

func TestFormatTerminal_Plain(t *testing.T) {
	c := &Changelog{
		Project: "myproject",
		Versions: []Version{
			{
				Version: "1.0.0",
				Date:    "2024-01-01",
				Added:   []string{"Init"},
			},
		},
	}
	out := FormatTerminal(c, FormatOptions{Plain: true})
	if !strings.Contains(out, "myproject Changelog") {
		t.Error("expected project header")
	}
}

func TestWrapText(t *testing.T) {
	tests := map[string]struct {
		text     string
		maxWidth int
		wantNL   bool
	}{
		"short": {text: "Short text", maxWidth: 80, wantNL: false},
		"long":  {text: "This is a very long text that should wrap at some point in the output", maxWidth: 30, wantNL: true},
		"exact": {text: "fits", maxWidth: 10, wantNL: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := wrapText(tc.text, tc.maxWidth, "  ")
			hasNL := strings.Contains(result, "\n")
			if hasNL != tc.wantNL {
				t.Errorf("wrapText(%q, %d) newline = %v, want %v\nresult: %q",
					tc.text, tc.maxWidth, hasNL, tc.wantNL, result)
			}
		})
	}
}
