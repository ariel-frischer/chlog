package changelog

import (
	"os"
	"strings"
	"testing"
)

func TestLoadFromReader_UnknownField(t *testing.T) {
	yaml := `project: test
unknown_field: should_fail
versions:
  1.0.0:
    date: "2024-01-01"
    added:
      - Init
`
	_, err := LoadFromReader(strings.NewReader(yaml))
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestValidate_WhitespaceOnlyProject(t *testing.T) {
	c := &Changelog{
		Project:  "   ",
		Versions: []Version{},
	}
	errs := Validate(c)
	found := false
	for _, e := range errs {
		if e.Field == "project" {
			found = true
		}
	}
	if !found {
		t.Error("expected validation error for whitespace-only project name")
	}
}

func TestValidate_EmptyVersionString(t *testing.T) {
	v := Version{Version: "", Date: "2024-01-01"}
	v.Public.Append("added", "x")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}
	errs := Validate(c)
	if len(errs) == 0 {
		t.Fatal("expected validation error for empty version string")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "must not be empty") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'must not be empty' error for version")
	}
}

func TestValidate_DateMissing(t *testing.T) {
	v := Version{Version: "1.0.0"}
	v.Public.Append("added", "x")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}
	errs := Validate(c)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "date required") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'date required' error for released version without date")
	}
}

func TestValidate_NoVersions(t *testing.T) {
	c := &Changelog{
		Project:  "test",
		Versions: nil,
	}
	errs := Validate(c)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid project with no versions, got %d: %v", len(errs), errs)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	v1 := Version{Version: "1.0.0", Date: "not-a-date"}
	v1.Public.Append("added", "x")
	v2 := Version{Version: "1.0.0", Date: "2024-01-01"}
	v2.Public.Append("added", "y")
	c := &Changelog{
		Project:  "",
		Versions: []Version{v1, v2},
	}
	errs := Validate(c)
	if len(errs) < 3 {
		t.Errorf("expected at least 3 errors (project + date + duplicate), got %d: %v", len(errs), errs)
	}
}

func TestSave_RoundTrip(t *testing.T) {
	unreleased := Version{Version: "unreleased"}
	unreleased.Public.Append("added", "WIP feature")

	released := Version{Version: "1.0.0", Date: "2024-01-01"}
	released.Public.Append("added", "Feature A")
	released.Public.Append("added", "Feature B")
	released.Public.Append("fixed", "Bug fix")
	released.Public.Append("security", "CVE patch")
	released.Internal.Append("changed", "Refactored handler")

	original := &Changelog{
		Project:  "roundtrip-test",
		Versions: []Version{unreleased, released},
	}

	path := t.TempDir() + "/roundtrip.yaml"
	if err := Save(original, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Project != original.Project {
		t.Errorf("project = %q, want %q", loaded.Project, original.Project)
	}
	if len(loaded.Versions) != len(original.Versions) {
		t.Fatalf("version count = %d, want %d", len(loaded.Versions), len(original.Versions))
	}
	if loaded.Versions[1].Count() != original.Versions[1].Count() {
		t.Errorf("entry count mismatch: got %d, want %d",
			loaded.Versions[1].Count(), original.Versions[1].Count())
	}
	if len(loaded.Versions[1].Internal.Get("changed")) != 1 {
		t.Error("internal changes not preserved in round-trip")
	}
}

func TestSave_BadPath(t *testing.T) {
	c := &Changelog{Project: "test"}
	err := Save(c, "/nonexistent/dir/file.yaml")
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}

func TestLoadFromReader_MultipleVersionsWithAllCategories(t *testing.T) {
	yaml := `project: full-test
versions:
  unreleased:
    added:
      - Upcoming feature
  2.0.0:
    date: "2024-06-01"
    added:
      - Major feature
    changed:
      - Updated API
    deprecated:
      - Old endpoint
    removed:
      - Legacy code
    fixed:
      - Critical bug
    security:
      - Patched CVE
  1.0.0:
    date: "2024-01-01"
    added:
      - Initial release
`
	c, err := LoadFromReader(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Versions) != 3 {
		t.Fatalf("version count = %d, want 3", len(c.Versions))
	}

	v2 := c.Versions[1]
	if v2.Count() != 6 {
		t.Errorf("v2.0.0 entry count = %d, want 6", v2.Count())
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/.chlog.yaml"
	if err := os.WriteFile(path, []byte("{{invalid yaml"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML config")
	}
}

func TestValidate_CustomCategories(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("performance", "Improved query speed")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}

	// Default strict mode should reject "performance"
	errs := Validate(c)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, `unknown category "performance"`) {
			found = true
		}
	}
	if !found {
		t.Error("expected unknown category error in strict mode")
	}

	// With custom categories including "performance"
	cfg := &Config{Categories: []string{"performance", "added"}}
	errs = Validate(c, cfg)
	for _, e := range errs {
		if strings.Contains(e.Message, "unknown category") {
			t.Errorf("unexpected unknown category error with custom config: %v", e)
		}
	}

	// Non-strict mode accepts anything
	strictFalse := false
	cfg2 := &Config{StrictCategories: &strictFalse}
	errs = Validate(c, cfg2)
	for _, e := range errs {
		if strings.Contains(e.Message, "unknown category") {
			t.Errorf("unexpected unknown category error in non-strict mode: %v", e)
		}
	}
}
