package changelog

import (
	"strings"
	"testing"
)

func TestMergedChanges(t *testing.T) {
	v := &Version{
		Changes:  Changes{Added: []string{"public"}, Fixed: []string{"pub-fix"}},
		Internal: Changes{Added: []string{"internal"}, Changed: []string{"refactor"}},
	}
	merged := v.MergedChanges()

	if len(merged.Added) != 2 {
		t.Errorf("merged.Added = %d, want 2", len(merged.Added))
	}
	if len(merged.Fixed) != 1 {
		t.Errorf("merged.Fixed = %d, want 1", len(merged.Fixed))
	}
	if len(merged.Changed) != 1 {
		t.Errorf("merged.Changed = %d, want 1", len(merged.Changed))
	}
}

func TestMergedChanges_DoesNotMutateOriginal(t *testing.T) {
	v := &Version{
		Changes:  Changes{Added: []string{"public"}},
		Internal: Changes{Added: []string{"internal"}},
	}
	_ = v.MergedChanges()
	if len(v.Changes.Added) != 1 {
		t.Error("MergedChanges mutated original Changes")
	}
}

func TestLoad_WithInternal(t *testing.T) {
	c, err := Load("testdata/with_internal.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Versions) != 2 {
		t.Fatalf("version count = %d, want 2", len(c.Versions))
	}
	if len(c.Versions[0].Internal.Changed) != 1 {
		t.Errorf("unreleased internal.changed = %d, want 1", len(c.Versions[0].Internal.Changed))
	}
	if len(c.Versions[1].Internal.Changed) != 1 {
		t.Errorf("v1.0.0 internal.changed = %d, want 1", len(c.Versions[1].Internal.Changed))
	}
	if len(c.Versions[1].Internal.Fixed) != 1 {
		t.Errorf("v1.0.0 internal.fixed = %d, want 1", len(c.Versions[1].Internal.Fixed))
	}
}

func TestLoad_InternalOnly(t *testing.T) {
	c, err := Load("testdata/internal_only.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Versions[0].Changes.IsEmpty() != true {
		t.Error("expected empty public changes")
	}
	if c.Versions[0].Internal.IsEmpty() {
		t.Error("expected non-empty internal changes")
	}
}

func TestValidate_InternalEmptyEntry(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{
				Version:  "1.0.0",
				Date:     "2024-01-01",
				Changes:  Changes{Added: []string{"public"}},
				Internal: Changes{Changed: []string{""}},
			},
		},
	}
	errs := Validate(c)
	if len(errs) == 0 {
		t.Fatal("expected validation error for empty internal entry")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Field, "internal.changed") {
			found = true
		}
	}
	if !found {
		t.Error("expected error referencing internal.changed field")
	}
}

func TestRenderVersionMarkdown_PublicOnly(t *testing.T) {
	v := &Version{
		Version:  "1.0.0",
		Date:     "2024-01-01",
		Changes:  Changes{Added: []string{"Public feature"}},
		Internal: Changes{Changed: []string{"Refactored internals"}},
	}
	var b strings.Builder
	if err := RenderVersionMarkdown(v, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := b.String()
	if !strings.Contains(out, "Public feature") {
		t.Error("expected public entry")
	}
	if strings.Contains(out, "Refactored internals") {
		t.Error("should not include internal entry without IncludeInternal")
	}
}

func TestRenderVersionMarkdown_WithInternal(t *testing.T) {
	v := &Version{
		Version:  "1.0.0",
		Date:     "2024-01-01",
		Changes:  Changes{Added: []string{"Public feature"}},
		Internal: Changes{Changed: []string{"Refactored internals"}},
	}
	var b strings.Builder
	if err := RenderVersionMarkdown(v, &b, RenderOptions{IncludeInternal: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := b.String()
	if !strings.Contains(out, "Public feature") {
		t.Error("expected public entry")
	}
	if !strings.Contains(out, "Refactored internals") {
		t.Error("expected internal entry when IncludeInternal is true")
	}
}

func TestFormatVersion_PublicOnly(t *testing.T) {
	v := &Version{
		Version:  "1.0.0",
		Date:     "2024-01-01",
		Changes:  Changes{Added: []string{"Public"}},
		Internal: Changes{Changed: []string{"Internal"}},
	}
	out := FormatVersion(v, FormatOptions{Plain: true})
	if strings.Contains(out, "Internal") {
		t.Error("should not show internal entries without IncludeInternal")
	}
}

func TestFormatVersion_WithInternal(t *testing.T) {
	v := &Version{
		Version:  "1.0.0",
		Date:     "2024-01-01",
		Changes:  Changes{Added: []string{"Public"}},
		Internal: Changes{Changed: []string{"Internal change"}},
	}
	out := FormatVersion(v, FormatOptions{Plain: true, IncludeInternal: true})
	if !strings.Contains(out, "Internal change") {
		t.Error("expected internal entry when IncludeInternal is true")
	}
}

func TestAllEntries_PublicOnly(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{
				Version:  "1.0.0",
				Date:     "2024-01-01",
				Changes:  Changes{Added: []string{"public"}},
				Internal: Changes{Changed: []string{"internal"}},
			},
		},
	}
	entries := c.AllEntries()
	if len(entries) != 1 {
		t.Errorf("AllEntries() = %d, want 1 (public only)", len(entries))
	}
}

func TestAllEntries_WithInternal(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{
				Version:  "1.0.0",
				Date:     "2024-01-01",
				Changes:  Changes{Added: []string{"public"}},
				Internal: Changes{Changed: []string{"internal"}},
			},
		},
	}
	entries := c.AllEntries(QueryOptions{IncludeInternal: true})
	if len(entries) != 2 {
		t.Errorf("AllEntries(internal) = %d, want 2", len(entries))
	}
}

func TestGetEntryCount_WithInternal(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{
				Version:  "1.0.0",
				Date:     "2024-01-01",
				Changes:  Changes{Added: []string{"a"}},
				Internal: Changes{Changed: []string{"b", "c"}},
			},
		},
	}
	if got := c.GetEntryCount(); got != 1 {
		t.Errorf("GetEntryCount() = %d, want 1", got)
	}
	if got := c.GetEntryCount(QueryOptions{IncludeInternal: true}); got != 3 {
		t.Errorf("GetEntryCount(internal) = %d, want 3", got)
	}
}

func TestScaffold_InternalRouting(t *testing.T) {
	commits := []GitCommit{
		{Hash: "a", Subject: "feat: add feature"},
		{Hash: "b", Subject: "refactor: clean up handler"},
		{Hash: "c", Subject: "perf: speed up queries"},
		{Hash: "d", Subject: "fix: fix bug"},
	}
	v := Scaffold(commits, ScaffoldOptions{})
	if len(v.Changes.Added) != 1 {
		t.Errorf("public added = %d, want 1", len(v.Changes.Added))
	}
	if len(v.Changes.Fixed) != 1 {
		t.Errorf("public fixed = %d, want 1", len(v.Changes.Fixed))
	}
	if len(v.Changes.Changed) != 0 {
		t.Errorf("public changed = %d, want 0 (refactor/perf should be internal)", len(v.Changes.Changed))
	}
	if len(v.Internal.Changed) != 2 {
		t.Errorf("internal changed = %d, want 2 (refactor + perf)", len(v.Internal.Changed))
	}
}

func TestScaffold_BreakingRefactorIsPublic(t *testing.T) {
	commits := []GitCommit{
		{Hash: "a", Subject: "refactor!: new config format"},
	}
	v := Scaffold(commits, ScaffoldOptions{})
	if len(v.Changes.Changed) != 1 {
		t.Errorf("public changed = %d, want 1 (breaking refactor)", len(v.Changes.Changed))
	}
	if len(v.Internal.Changed) != 0 {
		t.Errorf("internal changed = %d, want 0 (breaking is always public)", len(v.Internal.Changed))
	}
}
