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
  - version: 1.0.0
    date: "2024-01-01"
    changes:
      added:
        - Init
`
	_, err := LoadFromReader(strings.NewReader(yaml))
	if err == nil {
		t.Fatal("expected error for unknown field (KnownFields enabled)")
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
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{Version: "", Date: "2024-01-01", Changes: Changes{Added: []string{"x"}}},
		},
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
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{Version: "1.0.0", Changes: Changes{Added: []string{"x"}}},
		},
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
	c := &Changelog{
		Project: "",
		Versions: []Version{
			{Version: "1.0.0", Date: "not-a-date", Changes: Changes{Added: []string{"x"}}},
			{Version: "1.0.0", Date: "2024-01-01", Changes: Changes{Added: []string{"y"}}},
		},
	}
	errs := Validate(c)
	if len(errs) < 3 {
		t.Errorf("expected at least 3 errors (project + date + duplicate), got %d: %v", len(errs), errs)
	}
}

func TestSave_RoundTrip(t *testing.T) {
	original := &Changelog{
		Project: "roundtrip-test",
		Versions: []Version{
			{
				Version: "unreleased",
				Changes: Changes{Added: []string{"WIP feature"}},
			},
			{
				Version: "1.0.0",
				Date:    "2024-01-01",
				Changes: Changes{
					Added:    []string{"Feature A", "Feature B"},
					Fixed:    []string{"Bug fix"},
					Security: []string{"CVE patch"},
				},
				Internal: Changes{Changed: []string{"Refactored handler"}},
			},
		},
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
	if loaded.Versions[1].Changes.Count() != original.Versions[1].Changes.Count() {
		t.Errorf("entry count mismatch: got %d, want %d",
			loaded.Versions[1].Changes.Count(), original.Versions[1].Changes.Count())
	}
	if len(loaded.Versions[1].Internal.Changed) != 1 {
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
  - version: unreleased
    changes:
      added:
        - Upcoming feature
  - version: 2.0.0
    date: "2024-06-01"
    changes:
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
  - version: 1.0.0
    date: "2024-01-01"
    changes:
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
	if v2.Changes.Count() != 6 {
		t.Errorf("v2.0.0 entry count = %d, want 6", v2.Changes.Count())
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
