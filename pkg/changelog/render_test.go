package changelog

import (
	"strings"
	"testing"
)

func TestRenderVersionMarkdown_Unreleased(t *testing.T) {
	v := &Version{Version: "unreleased"}
	v.Public.Append("added", "New feature")
	var b strings.Builder
	if err := RenderVersionMarkdown(v, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := b.String()

	if !strings.Contains(out, "## [Unreleased]") {
		t.Error("expected [Unreleased] header")
	}
	if !strings.Contains(out, "### Added") {
		t.Error("expected Added category")
	}
	if !strings.Contains(out, "- New feature") {
		t.Error("expected entry")
	}
}

func TestRenderVersionMarkdown_Released(t *testing.T) {
	v := &Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "Feature A")
	v.Public.Append("fixed", "Bug B")
	var b strings.Builder
	if err := RenderVersionMarkdown(v, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := b.String()

	if !strings.Contains(out, "## [1.0.0] - 2024-01-01") {
		t.Error("expected version header with date")
	}
	// Added should come before Fixed (order from YAML/Append)
	addedIdx := strings.Index(out, "### Added")
	fixedIdx := strings.Index(out, "### Fixed")
	if addedIdx == -1 || fixedIdx == -1 {
		t.Fatal("expected both Added and Fixed sections")
	}
	if addedIdx > fixedIdx {
		t.Error("Added should appear before Fixed")
	}
}

func TestRenderVersionMarkdown_SkipsEmptyCategories(t *testing.T) {
	v := &Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("fixed", "Bug fix")
	var b strings.Builder
	if err := RenderVersionMarkdown(v, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := b.String()

	if strings.Contains(out, "### Added") {
		t.Error("should not include empty Added section")
	}
	if !strings.Contains(out, "### Fixed") {
		t.Error("expected Fixed section")
	}
}

func TestRenderMarkdownString(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "Initial release")
	c := &Changelog{
		Project:  "test-project",
		Versions: []Version{v},
	}
	out, err := RenderMarkdownString(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "# Changelog") {
		t.Error("expected Changelog header")
	}
	if !strings.Contains(out, "test-project") {
		t.Error("expected project name")
	}
	if !strings.Contains(out, "Keep a Changelog") {
		t.Error("expected Keep a Changelog reference")
	}
}

func TestNormalizeGitURL(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"https":     {input: "https://github.com/org/repo.git", want: "https://github.com/org/repo"},
		"ssh":       {input: "git@github.com:org/repo.git", want: "https://github.com/org/repo"},
		"no_suffix": {input: "https://gitlab.com/org/repo", want: "https://gitlab.com/org/repo"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := normalizeGitURL(tc.input); got != tc.want {
				t.Errorf("normalizeGitURL(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
