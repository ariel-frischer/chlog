package changelog

import (
	"strings"
	"testing"
)

func TestTitleCase(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"lowercase":     {input: "added", want: "Added"},
		"already upper": {input: "Added", want: "Added"},
		"empty":         {input: "", want: ""},
		"single char":   {input: "a", want: "A"},
		"unicode":       {input: "änderung", want: "Änderung"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := titleCase(tc.input); got != tc.want {
				t.Errorf("titleCase(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestRenderComparisonLinks_GitHub(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{Version: "unreleased"},
			{Version: "2.0.0", Date: "2024-06-01"},
			{Version: "1.0.0", Date: "2024-01-01"},
		},
	}
	var b strings.Builder
	renderComparisonLinks(c, &b, "https://github.com/org/repo")
	out := b.String()

	if !strings.Contains(out, "[Unreleased]: https://github.com/org/repo/compare/v2.0.0...HEAD") {
		t.Errorf("expected unreleased comparison link, got:\n%s", out)
	}
	if !strings.Contains(out, "[2.0.0]: https://github.com/org/repo/compare/v1.0.0...v2.0.0") {
		t.Errorf("expected 2.0.0 comparison link, got:\n%s", out)
	}
}

func TestRenderComparisonLinks_GitLab(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{Version: "unreleased"},
			{Version: "1.0.0", Date: "2024-01-01"},
		},
	}
	var b strings.Builder
	renderComparisonLinks(c, &b, "https://gitlab.com/org/repo")
	out := b.String()

	if !strings.Contains(out, "/-/compare/") {
		t.Errorf("expected GitLab compare path, got:\n%s", out)
	}
}

func TestRenderComparisonLinks_EmptyVersions(t *testing.T) {
	c := &Changelog{Project: "test", Versions: nil}
	var b strings.Builder
	renderComparisonLinks(c, &b, "https://github.com/org/repo")
	if b.Len() != 0 {
		t.Errorf("expected no output for empty versions, got: %q", b.String())
	}
}

func TestRenderComparisonLinks_SingleVersion(t *testing.T) {
	c := &Changelog{
		Project:  "test",
		Versions: []Version{{Version: "1.0.0", Date: "2024-01-01"}},
	}
	var b strings.Builder
	renderComparisonLinks(c, &b, "https://github.com/org/repo")
	// Single release with no predecessor — no comparison link possible
	if b.Len() != 0 {
		t.Errorf("expected no output for single version, got: %q", b.String())
	}
}

func TestRenderComparisonLinks_UnreleasedOnly(t *testing.T) {
	c := &Changelog{
		Project:  "test",
		Versions: []Version{{Version: "unreleased"}},
	}
	var b strings.Builder
	renderComparisonLinks(c, &b, "https://github.com/org/repo")
	// Unreleased with no previous release — no comparison link
	if b.Len() != 0 {
		t.Errorf("expected no output for unreleased-only, got: %q", b.String())
	}
}

func TestRenderMarkdown_WithComparisonLinks(t *testing.T) {
	c := &Changelog{
		Project: "test-project",
		Versions: []Version{
			{Version: "unreleased", Changes: Changes{Added: []string{"WIP"}}},
			{Version: "2.0.0", Date: "2024-06-01", Changes: Changes{Added: []string{"Feature"}}},
			{Version: "1.0.0", Date: "2024-01-01", Changes: Changes{Added: []string{"Init"}}},
		},
	}
	cfg := &Config{RepoURL: "https://github.com/org/repo"}
	out, err := RenderMarkdownString(c, RenderOptions{Config: cfg})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "# Changelog") {
		t.Error("missing Changelog header")
	}
	if !strings.Contains(out, "[Unreleased]: https://github.com/org/repo/compare/v2.0.0...HEAD") {
		t.Error("missing unreleased comparison link")
	}
	if !strings.Contains(out, "[2.0.0]: https://github.com/org/repo/compare/v1.0.0...v2.0.0") {
		t.Error("missing 2.0.0 comparison link")
	}
}

func TestRenderMarkdown_NoLinksWithoutConfig(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{Version: "1.0.0", Date: "2024-01-01", Changes: Changes{Added: []string{"Init"}}},
		},
	}
	out, err := RenderMarkdownString(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "/compare/") {
		t.Error("should not contain comparison links without config")
	}
}

func TestRenderVersionMarkdown_AllCategories(t *testing.T) {
	v := &Version{
		Version: "1.0.0",
		Date:    "2024-01-01",
		Changes: Changes{
			Added:      []string{"New thing"},
			Changed:    []string{"Updated thing"},
			Deprecated: []string{"Old thing"},
			Removed:    []string{"Gone thing"},
			Fixed:      []string{"Fixed thing"},
			Security:   []string{"Secure thing"},
		},
	}
	var b strings.Builder
	if err := RenderVersionMarkdown(v, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := b.String()

	// Verify canonical ordering
	categories := []string{"Added", "Changed", "Deprecated", "Removed", "Fixed", "Security"}
	lastIdx := -1
	for _, cat := range categories {
		idx := strings.Index(out, "### "+cat)
		if idx == -1 {
			t.Errorf("missing ### %s section", cat)
			continue
		}
		if idx < lastIdx {
			t.Errorf("### %s appeared out of canonical order", cat)
		}
		lastIdx = idx
	}
}

func TestRenderVersionMarkdown_EmptyChanges(t *testing.T) {
	v := &Version{
		Version: "unreleased",
		Changes: Changes{},
	}
	var b strings.Builder
	if err := RenderVersionMarkdown(v, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := b.String()

	if !strings.Contains(out, "## [Unreleased]") {
		t.Error("expected unreleased header even with empty changes")
	}
	if strings.Contains(out, "###") {
		t.Error("should not render any category sections for empty changes")
	}
}
