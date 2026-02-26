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
	unreleased := Version{Version: "unreleased"}
	unreleased.Public.Append("added", "WIP")
	v2 := Version{Version: "2.0.0", Date: "2024-06-01"}
	v2.Public.Append("added", "Feature")
	v1 := Version{Version: "1.0.0", Date: "2024-01-01"}
	v1.Public.Append("added", "Init")

	c := &Changelog{
		Project:  "test-project",
		Versions: []Version{unreleased, v2, v1},
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
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "Init")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
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
	v := &Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "New thing")
	v.Public.Append("changed", "Updated thing")
	v.Public.Append("deprecated", "Old thing")
	v.Public.Append("removed", "Gone thing")
	v.Public.Append("fixed", "Fixed thing")
	v.Public.Append("security", "Secure thing")

	var b strings.Builder
	if err := RenderVersionMarkdown(v, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := b.String()

	// Verify ordering matches append order (which is canonical)
	categories := []string{"Added", "Changed", "Deprecated", "Removed", "Fixed", "Security"}
	lastIdx := -1
	for _, cat := range categories {
		idx := strings.Index(out, "### "+cat)
		if idx == -1 {
			t.Errorf("missing ### %s section", cat)
			continue
		}
		if idx < lastIdx {
			t.Errorf("### %s appeared out of order", cat)
		}
		lastIdx = idx
	}
}

func TestRenderVersionMarkdown_EmptyChanges(t *testing.T) {
	v := &Version{
		Version: "unreleased",
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
