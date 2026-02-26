package changelog

import (
	"strings"
	"testing"
)

func TestLoad_Valid(t *testing.T) {
	c, err := Load("testdata/valid.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Project != "test-project" {
		t.Errorf("project = %q, want %q", c.Project, "test-project")
	}
	if len(c.Versions) != 3 {
		t.Errorf("version count = %d, want 3", len(c.Versions))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("testdata/nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFromReader_Valid(t *testing.T) {
	yaml := `project: myproject
versions:
  1.0.0:
    date: "2024-01-01"
    added:
      - Initial release
`
	c, err := LoadFromReader(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Project != "myproject" {
		t.Errorf("project = %q, want %q", c.Project, "myproject")
	}
}

func TestLoadFromReader_InvalidYAML(t *testing.T) {
	_, err := LoadFromReader(strings.NewReader("{{invalid"))
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		fixture    string
		wantErrors []string
	}{
		"empty_project": {
			fixture:    "testdata/empty_project.yaml",
			wantErrors: []string{"project"},
		},
		"duplicate_versions": {
			fixture:    "testdata/duplicate_versions.yaml",
			wantErrors: []string{"duplicate"},
		},
		"invalid_date": {
			fixture:    "testdata/invalid_date.yaml",
			wantErrors: []string{"invalid date"},
		},
		"multiple_unreleased": {
			fixture:    "testdata/multiple_unreleased.yaml",
			wantErrors: []string{"only one unreleased"},
		},
		"empty_entries": {
			fixture:    "testdata/empty_entries.yaml",
			wantErrors: []string{"must not be empty"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Load(tc.fixture)
			if err == nil {
				t.Fatal("expected validation error")
			}
			for _, want := range tc.wantErrors {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q should contain %q", err.Error(), want)
				}
			}
		})
	}
}

func TestValidate_EmptyUnreleasedAllowed(t *testing.T) {
	yaml := `project: test
versions:
  unreleased: {}
  1.0.0:
    date: "2024-01-01"
    added:
      - Initial release
`
	c, err := LoadFromReader(strings.NewReader(yaml))
	if err != nil {
		t.Fatalf("empty unreleased should be valid, got: %v", err)
	}
	if !c.HasUnreleased() {
		t.Error("expected unreleased version")
	}
}

func TestValidate_EmptyReleasedRejected(t *testing.T) {
	yaml := `project: test
versions:
  1.0.0:
    date: "2024-01-01"
`
	_, err := LoadFromReader(strings.NewReader(yaml))
	if err == nil {
		t.Fatal("empty released version should fail validation")
	}
	if !strings.Contains(err.Error(), "at least one entry") {
		t.Errorf("error %q should mention 'at least one entry'", err.Error())
	}
}

func TestNormalizeVersion(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"plain":      {input: "1.0.0", want: "1.0.0"},
		"v_prefix":   {input: "v1.0.0", want: "1.0.0"},
		"uppercase":  {input: "V1.0.0", want: "1.0.0"},
		"unreleased": {input: "Unreleased", want: "unreleased"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := NormalizeVersion(tc.input); got != tc.want {
				t.Errorf("NormalizeVersion(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestSave(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{
				Version: "1.0.0",
				Date:    "2024-01-01",
				Added:   []string{"Initial release"},
			},
		},
	}
	path := t.TempDir() + "/out.yaml"
	if err := Save(c, path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() after Save() error: %v", err)
	}
	if loaded.Project != "test" {
		t.Errorf("project = %q, want %q", loaded.Project, "test")
	}
}
