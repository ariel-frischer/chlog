package main

import (
	"path/filepath"
	"testing"

	"github.com/ariel-frischer/chlog/pkg/changelog"
)

func TestRunRemove_ExactMatch(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")

	u := changelog.Version{Version: "unreleased"}
	u.Public.Append("added", "Feature A")
	u.Public.Append("added", "Feature B")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{u},
	})

	removeVersion = "unreleased"
	removeInternal = false
	removeMatch = false

	if err := runRemove(nil, []string{"added", "Feature A"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := loadTestChangelog(t, yamlFile)
	entries := c.GetUnreleased().Public.Get("added")
	if len(entries) != 1 || entries[0] != "Feature B" {
		t.Errorf("entries = %v, want [Feature B]", entries)
	}
}

func TestRunRemove_SubstringMatch(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")

	u := changelog.Version{Version: "unreleased"}
	u.Public.Append("fixed", "Fix login timeout")
	u.Public.Append("fixed", "Fix signup error")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{u},
	})

	removeVersion = "unreleased"
	removeInternal = false
	removeMatch = true

	if err := runRemove(nil, []string{"fixed", "login"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := loadTestChangelog(t, yamlFile)
	entries := c.GetUnreleased().Public.Get("fixed")
	if len(entries) != 1 || entries[0] != "Fix signup error" {
		t.Errorf("entries = %v, want [Fix signup error]", entries)
	}
}

func TestRunRemove_SubstringMultipleMatchErrors(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")

	u := changelog.Version{Version: "unreleased"}
	u.Public.Append("fixed", "Fix login timeout")
	u.Public.Append("fixed", "Fix login redirect")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{u},
	})

	removeVersion = "unreleased"
	removeInternal = false
	removeMatch = true

	err := runRemove(nil, []string{"fixed", "login"})
	if err == nil {
		t.Fatal("expected error for multiple matches")
	}
}

func TestRunRemove_InternalEntry(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")

	u := changelog.Version{Version: "unreleased"}
	u.Internal.Append("changed", "Refactor auth")
	u.Internal.Append("changed", "Refactor DB")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{u},
	})

	removeVersion = "unreleased"
	removeInternal = true
	removeMatch = false

	if err := runRemove(nil, []string{"changed", "Refactor auth"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := loadTestChangelog(t, yamlFile)
	entries := c.GetUnreleased().Internal.Get("changed")
	if len(entries) != 1 || entries[0] != "Refactor DB" {
		t.Errorf("internal entries = %v, want [Refactor DB]", entries)
	}
}

func TestRunRemove_SpecificVersion(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")

	v := changelog.Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "Feature A")
	v.Public.Append("added", "Feature B")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{{Version: "unreleased"}, v},
	})

	removeVersion = "1.0.0"
	removeInternal = false
	removeMatch = false

	if err := runRemove(nil, []string{"added", "Feature A"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := loadTestChangelog(t, yamlFile)
	ver, _ := c.GetVersion("1.0.0")
	entries := ver.Public.Get("added")
	if len(entries) != 1 || entries[0] != "Feature B" {
		t.Errorf("entries = %v, want [Feature B]", entries)
	}
}

func TestRunRemove_VersionNotFound(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{{Version: "unreleased"}},
	})

	removeVersion = "9.9.9"
	removeInternal = false
	removeMatch = false

	err := runRemove(nil, []string{"added", "test"})
	if err == nil {
		t.Fatal("expected error for nonexistent version")
	}
}

func TestRunRemove_CategoryNotFound(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")

	u := changelog.Version{Version: "unreleased"}
	u.Public.Append("added", "Feature A")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{u},
	})

	removeVersion = "unreleased"
	removeInternal = false
	removeMatch = false

	err := runRemove(nil, []string{"removed", "anything"})
	if err == nil {
		t.Fatal("expected error for nonexistent category")
	}
}

func TestRunRemove_EntryNotFound(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")

	u := changelog.Version{Version: "unreleased"}
	u.Public.Append("added", "Feature A")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{u},
	})

	removeVersion = "unreleased"
	removeInternal = false
	removeMatch = false

	err := runRemove(nil, []string{"added", "Nonexistent"})
	if err == nil {
		t.Fatal("expected error for nonexistent entry")
	}
}

func TestRunRemove_CleansUpEmptyCategory(t *testing.T) {
	dir := t.TempDir()
	yamlFile = filepath.Join(dir, "CHANGELOG.yaml")

	u := changelog.Version{Version: "unreleased"}
	u.Public.Append("added", "Only entry")
	u.Public.Append("fixed", "Keep this")
	writeTestChangelog(t, yamlFile, &changelog.Changelog{
		Project:  "test",
		Versions: []changelog.Version{u},
	})

	removeVersion = "unreleased"
	removeInternal = false
	removeMatch = false

	if err := runRemove(nil, []string{"added", "Only entry"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := loadTestChangelog(t, yamlFile)
	cats := c.GetUnreleased().Public.CategoryNames()
	for _, cat := range cats {
		if cat == "added" {
			t.Error("empty 'added' category should have been removed")
		}
	}
	if len(cats) != 1 || cats[0] != "fixed" {
		t.Errorf("categories = %v, want [fixed]", cats)
	}
}
