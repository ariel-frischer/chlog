package main

import (
	"path/filepath"
	"testing"

	"github.com/ariel-frischer/chlog/pkg/changelog"
)

func TestRunAdd_BasicPublicEntry(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{{Version: "unreleased"}},
	})

	addVersion = "unreleased"
	addInternal = false

	if err := runAdd(nil, []string{"added", "New feature"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := loadTestChangelog(t, yamlFile)
	u := c.GetUnreleased()
	if u == nil {
		t.Fatal("expected unreleased version")
	}
	entries := u.Public.Get("added")
	if len(entries) != 1 || entries[0] != "New feature" {
		t.Errorf("entries = %v, want [New feature]", entries)
	}
}

func TestRunAdd_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{{Version: "unreleased"}},
	})

	addVersion = "unreleased"
	addInternal = false

	if err := runAdd(nil, []string{"added", "Feature A", "Feature B", "Feature C"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := loadTestChangelog(t, yamlFile)
	entries := c.GetUnreleased().Public.Get("added")
	if len(entries) != 3 {
		t.Fatalf("entries count = %d, want 3", len(entries))
	}
	want := []string{"Feature A", "Feature B", "Feature C"}
	for i, w := range want {
		if entries[i] != w {
			t.Errorf("entries[%d] = %q, want %q", i, entries[i], w)
		}
	}
}

func TestRunAdd_InternalEntry(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{{Version: "unreleased"}},
	})

	addVersion = "unreleased"
	addInternal = true

	if err := runAdd(nil, []string{"changed", "Refactor auth"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := loadTestChangelog(t, yamlFile)
	u := c.GetUnreleased()
	if u.Public.Get("changed") != nil {
		t.Error("entry should not be in public changes")
	}
	entries := u.Internal.Get("changed")
	if len(entries) != 1 || entries[0] != "Refactor auth" {
		t.Errorf("internal entries = %v, want [Refactor auth]", entries)
	}
}

func TestRunAdd_AutoCreatesUnreleased(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")

	v := changelog.Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "Init")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{v},
	})

	addVersion = "unreleased"
	addInternal = false

	if err := runAdd(nil, []string{"fixed", "Bug fix"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := loadTestChangelog(t, yamlFile)
	if len(c.Versions) != 2 {
		t.Fatalf("versions = %d, want 2", len(c.Versions))
	}
	if !c.Versions[0].IsUnreleased() {
		t.Error("first version should be unreleased")
	}
	entries := c.Versions[0].Public.Get("fixed")
	if len(entries) != 1 || entries[0] != "Bug fix" {
		t.Errorf("entries = %v, want [Bug fix]", entries)
	}
}

func TestRunAdd_SpecificVersion(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")

	v := changelog.Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "Init")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{{Version: "unreleased"}, v},
	})

	addVersion = "1.0.0"
	addInternal = false

	if err := runAdd(nil, []string{"fixed", "Hotfix"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := loadTestChangelog(t, yamlFile)
	ver, err := c.GetVersion("1.0.0")
	if err != nil {
		t.Fatalf("version not found: %v", err)
	}
	entries := ver.Public.Get("fixed")
	if len(entries) != 1 || entries[0] != "Hotfix" {
		t.Errorf("entries = %v, want [Hotfix]", entries)
	}
}

func TestRunAdd_NonexistentVersionErrors(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{{Version: "unreleased"}},
	})

	addVersion = "9.9.9"
	addInternal = false

	err := runAdd(nil, []string{"added", "test"})
	if err == nil {
		t.Fatal("expected error for nonexistent version")
	}
}

func TestRunAdd_AppendsToExisting(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")

	u := changelog.Version{Version: "unreleased"}
	u.Public.Append("added", "Existing feature")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{u},
	})

	addVersion = "unreleased"
	addInternal = false

	if err := runAdd(nil, []string{"added", "New feature"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := loadTestChangelog(t, yamlFile)
	entries := c.GetUnreleased().Public.Get("added")
	if len(entries) != 2 {
		t.Fatalf("entries count = %d, want 2", len(entries))
	}
	if entries[0] != "Existing feature" || entries[1] != "New feature" {
		t.Errorf("entries = %v", entries)
	}
}

func TestPluralY(t *testing.T) {
	tests := map[string]struct {
		n    int
		want string
	}{
		"singular": {n: 1, want: "y"},
		"plural":   {n: 2, want: "ies"},
		"zero":     {n: 0, want: "ies"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := pluralY(tc.n); got != tc.want {
				t.Errorf("pluralY(%d) = %q, want %q", tc.n, got, tc.want)
			}
		})
	}
}

// writeTestChangelog saves a changelog to a temp file.
func writeTestChangelog(t *testing.T, path string, c *changelog.Changelog) {
	t.Helper()
	if err := changelog.Save(c, path); err != nil {
		t.Fatalf("saving test changelog: %v", err)
	}
}

// loadTestChangelog loads a changelog from a temp file.
func loadTestChangelog(t *testing.T, path string) *changelog.Changelog {
	t.Helper()
	c, err := changelog.Load(path)
	if err != nil {
		t.Fatalf("loading test changelog: %v", err)
	}
	return c
}
