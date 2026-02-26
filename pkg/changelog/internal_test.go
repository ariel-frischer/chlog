package changelog

import (
	"strings"
	"testing"
)

func TestMergedChanges(t *testing.T) {
	v := &Version{}
	v.Public.Append("added", "public")
	v.Public.Append("fixed", "pub-fix")
	v.Internal.Append("added", "internal")
	v.Internal.Append("changed", "refactor")

	merged := v.MergedChanges()

	if len(merged.Get("added")) != 2 {
		t.Errorf("merged.added = %d, want 2", len(merged.Get("added")))
	}
	if len(merged.Get("fixed")) != 1 {
		t.Errorf("merged.fixed = %d, want 1", len(merged.Get("fixed")))
	}
	if len(merged.Get("changed")) != 1 {
		t.Errorf("merged.changed = %d, want 1", len(merged.Get("changed")))
	}
}

func TestMergedChanges_DoesNotMutateOriginal(t *testing.T) {
	v := &Version{}
	v.Public.Append("added", "public")
	v.Internal.Append("added", "internal")

	_ = v.MergedChanges()
	if len(v.Public.Get("added")) != 1 {
		t.Error("MergedChanges mutated original Public")
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
	if len(c.Versions[0].Internal.Get("changed")) != 1 {
		t.Errorf("unreleased internal.changed = %d, want 1", len(c.Versions[0].Internal.Get("changed")))
	}
	if len(c.Versions[1].Internal.Get("changed")) != 1 {
		t.Errorf("v1.0.0 internal.changed = %d, want 1", len(c.Versions[1].Internal.Get("changed")))
	}
	if len(c.Versions[1].Internal.Get("fixed")) != 1 {
		t.Errorf("v1.0.0 internal.fixed = %d, want 1", len(c.Versions[1].Internal.Get("fixed")))
	}
}

func TestLoad_InternalOnly(t *testing.T) {
	c, err := Load("testdata/internal_only.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.Versions[0].IsEmpty() {
		t.Error("expected empty public changes")
	}
	if c.Versions[0].Internal.IsEmpty() {
		t.Error("expected non-empty internal changes")
	}
}

func TestValidate_InternalEmptyEntry(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "public")
	v.Internal.Append("changed", "")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
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
	v := &Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "Public feature")
	v.Internal.Append("changed", "Refactored internals")

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
	v := &Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "Public feature")
	v.Internal.Append("changed", "Refactored internals")

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
	v := &Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "Public")
	v.Internal.Append("changed", "Internal")

	out := FormatVersion(v, FormatOptions{Plain: true})
	if strings.Contains(out, "Internal") {
		t.Error("should not show internal entries without IncludeInternal")
	}
}

func TestFormatVersion_WithInternal(t *testing.T) {
	v := &Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "Public")
	v.Internal.Append("changed", "Internal change")

	out := FormatVersion(v, FormatOptions{Plain: true, IncludeInternal: true})
	if !strings.Contains(out, "Internal change") {
		t.Error("expected internal entry when IncludeInternal is true")
	}
}

func TestAllEntries_PublicOnly(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "public")
	v.Internal.Append("changed", "internal")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}
	entries := c.AllEntries()
	if len(entries) != 1 {
		t.Errorf("AllEntries() = %d, want 1 (public only)", len(entries))
	}
}

func TestAllEntries_WithInternal(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "public")
	v.Internal.Append("changed", "internal")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}
	entries := c.AllEntries(QueryOptions{IncludeInternal: true})
	if len(entries) != 2 {
		t.Errorf("AllEntries(internal) = %d, want 2", len(entries))
	}
}

func TestGetEntryCount_WithInternal(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "a")
	v.Internal.Append("changed", "b")
	v.Internal.Append("changed", "c")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
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
	if len(v.Public.Get("added")) != 1 {
		t.Errorf("public added = %d, want 1", len(v.Public.Get("added")))
	}
	if len(v.Public.Get("fixed")) != 1 {
		t.Errorf("public fixed = %d, want 1", len(v.Public.Get("fixed")))
	}
	if len(v.Public.Get("changed")) != 0 {
		t.Errorf("public changed = %d, want 0 (refactor/perf should be internal)", len(v.Public.Get("changed")))
	}
	if len(v.Internal.Get("changed")) != 2 {
		t.Errorf("internal changed = %d, want 2 (refactor + perf)", len(v.Internal.Get("changed")))
	}
}

func TestScaffold_BreakingRefactorIsPublic(t *testing.T) {
	commits := []GitCommit{
		{Hash: "a", Subject: "refactor!: new config format"},
	}
	v := Scaffold(commits, ScaffoldOptions{})
	if len(v.Public.Get("changed")) != 1 {
		t.Errorf("public changed = %d, want 1 (breaking refactor)", len(v.Public.Get("changed")))
	}
	if len(v.Internal.Get("changed")) != 0 {
		t.Errorf("internal changed = %d, want 0 (breaking is always public)", len(v.Internal.Get("changed")))
	}
}
