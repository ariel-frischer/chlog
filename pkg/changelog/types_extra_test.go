package changelog

import (
	"errors"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	e := ValidationError{Field: "versions[0].date", Message: "invalid format"}
	want := "versions[0].date: invalid format"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestVersionNotFoundError_Error(t *testing.T) {
	e := VersionNotFoundError{Version: "3.0.0"}
	want := `version "3.0.0" not found`
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestVersionNotFoundError_IsError(t *testing.T) {
	_, err := (&Changelog{Project: "test"}).GetVersion("nope")
	var notFound VersionNotFoundError
	if !errors.As(err, &notFound) {
		t.Errorf("expected VersionNotFoundError, got %T", err)
	}
	if notFound.Version != "nope" {
		t.Errorf("Version = %q, want %q", notFound.Version, "nope")
	}
}

func TestMergedChanges_AllCategories(t *testing.T) {
	v := &Version{
		Changes: Changes{
			Added:      []string{"pub-add"},
			Changed:    []string{"pub-change"},
			Deprecated: []string{"pub-dep"},
			Removed:    []string{"pub-rem"},
			Fixed:      []string{"pub-fix"},
			Security:   []string{"pub-sec"},
		},
		Internal: Changes{
			Added:      []string{"int-add"},
			Changed:    []string{"int-change"},
			Deprecated: []string{"int-dep"},
			Removed:    []string{"int-rem"},
			Fixed:      []string{"int-fix"},
			Security:   []string{"int-sec"},
		},
	}
	merged := v.MergedChanges()

	if len(merged.Added) != 2 {
		t.Errorf("Added = %d, want 2", len(merged.Added))
	}
	if len(merged.Changed) != 2 {
		t.Errorf("Changed = %d, want 2", len(merged.Changed))
	}
	if len(merged.Deprecated) != 2 {
		t.Errorf("Deprecated = %d, want 2", len(merged.Deprecated))
	}
	if len(merged.Removed) != 2 {
		t.Errorf("Removed = %d, want 2", len(merged.Removed))
	}
	if len(merged.Fixed) != 2 {
		t.Errorf("Fixed = %d, want 2", len(merged.Fixed))
	}
	if len(merged.Security) != 2 {
		t.Errorf("Security = %d, want 2", len(merged.Security))
	}
	if merged.Count() != 12 {
		t.Errorf("Count() = %d, want 12", merged.Count())
	}
}

func TestMergedChanges_EmptyInternals(t *testing.T) {
	v := &Version{
		Changes:  Changes{Added: []string{"a", "b"}},
		Internal: Changes{},
	}
	merged := v.MergedChanges()
	if len(merged.Added) != 2 {
		t.Errorf("Added = %d, want 2 (internal empty)", len(merged.Added))
	}
}

func TestMergedChanges_EmptyPublic(t *testing.T) {
	v := &Version{
		Changes:  Changes{},
		Internal: Changes{Fixed: []string{"internal fix"}},
	}
	merged := v.MergedChanges()
	if len(merged.Fixed) != 1 {
		t.Errorf("Fixed = %d, want 1 (public empty)", len(merged.Fixed))
	}
}

func TestChanges_CategoryEntries_AllCategories(t *testing.T) {
	c := Changes{
		Added:      []string{"a"},
		Changed:    []string{"b"},
		Deprecated: []string{"c"},
		Removed:    []string{"d"},
		Fixed:      []string{"e"},
		Security:   []string{"f"},
	}
	for _, cat := range ValidCategories() {
		entries := c.CategoryEntries(cat)
		if len(entries) != 1 {
			t.Errorf("CategoryEntries(%q) = %d, want 1", cat, len(entries))
		}
	}
}

func TestVersion_IsUnreleased_MixedCase(t *testing.T) {
	cases := []string{"unreleased", "Unreleased", "UNRELEASED", "UnReleAsed"}
	for _, version := range cases {
		v := &Version{Version: version}
		if !v.IsUnreleased() {
			t.Errorf("IsUnreleased(%q) = false, want true", version)
		}
	}
}
