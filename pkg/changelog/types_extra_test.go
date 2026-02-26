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
		Public: makeChangesMulti(map[string][]string{
			"added": {"pub-add"}, "changed": {"pub-change"},
			"deprecated": {"pub-dep"}, "removed": {"pub-rem"},
			"fixed": {"pub-fix"}, "security": {"pub-sec"},
		}),
		Internal: makeChangesMulti(map[string][]string{
			"added": {"int-add"}, "changed": {"int-change"},
			"deprecated": {"int-dep"}, "removed": {"int-rem"},
			"fixed": {"int-fix"}, "security": {"int-sec"},
		}),
	}
	merged := v.MergedChanges()

	if len(merged.Get("added")) != 2 {
		t.Errorf("Added = %d, want 2", len(merged.Get("added")))
	}
	if len(merged.Get("changed")) != 2 {
		t.Errorf("Changed = %d, want 2", len(merged.Get("changed")))
	}
	if len(merged.Get("deprecated")) != 2 {
		t.Errorf("Deprecated = %d, want 2", len(merged.Get("deprecated")))
	}
	if len(merged.Get("removed")) != 2 {
		t.Errorf("Removed = %d, want 2", len(merged.Get("removed")))
	}
	if len(merged.Get("fixed")) != 2 {
		t.Errorf("Fixed = %d, want 2", len(merged.Get("fixed")))
	}
	if len(merged.Get("security")) != 2 {
		t.Errorf("Security = %d, want 2", len(merged.Get("security")))
	}
	if merged.Count() != 12 {
		t.Errorf("Count() = %d, want 12", merged.Count())
	}
}

func TestMergedChanges_EmptyInternals(t *testing.T) {
	v := &Version{
		Public: makeChangesMulti(map[string][]string{"added": {"a", "b"}}),
	}
	merged := v.MergedChanges()
	if len(merged.Get("added")) != 2 {
		t.Errorf("Added = %d, want 2 (internal empty)", len(merged.Get("added")))
	}
}

func TestMergedChanges_EmptyPublic(t *testing.T) {
	v := &Version{
		Internal: makeChanges("fixed", "internal fix"),
	}
	merged := v.MergedChanges()
	if len(merged.Get("fixed")) != 1 {
		t.Errorf("Fixed = %d, want 1 (public empty)", len(merged.Get("fixed")))
	}
}

func TestChanges_Get_AllCategories(t *testing.T) {
	c := makeChangesMulti(map[string][]string{
		"added": {"a"}, "changed": {"b"},
		"deprecated": {"c"}, "removed": {"d"},
		"fixed": {"e"}, "security": {"f"},
	})
	for _, cat := range DefaultCategories {
		entries := c.Get(cat)
		if len(entries) != 1 {
			t.Errorf("Get(%q) = %d, want 1", cat, len(entries))
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

func TestChanges_Append(t *testing.T) {
	var c Changes
	c.Append("added", "first")
	c.Append("added", "second")
	c.Append("fixed", "bug")

	if len(c.Get("added")) != 2 {
		t.Errorf("added = %d, want 2", len(c.Get("added")))
	}
	if len(c.Get("fixed")) != 1 {
		t.Errorf("fixed = %d, want 1", len(c.Get("fixed")))
	}
}

func TestChanges_Clone(t *testing.T) {
	orig := makeChanges("added", "x")
	clone := orig.Clone()
	clone.Append("added", "y")

	if len(orig.Get("added")) != 1 {
		t.Error("Clone mutated original")
	}
	if len(clone.Get("added")) != 2 {
		t.Error("Clone didn't append")
	}
}

func TestChanges_CategoryNames(t *testing.T) {
	c := makeChangesMulti(map[string][]string{
		"added": {"a"}, "fixed": {"b"},
	})
	names := c.CategoryNames()
	if len(names) != 2 {
		t.Fatalf("CategoryNames() = %d, want 2", len(names))
	}
}
