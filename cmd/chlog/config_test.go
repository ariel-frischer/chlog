package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ariel-frischer/chlog/pkg/changelog"
)

func TestRunConfigInit_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	configFile = filepath.Join(dir, ".chlog.yaml")

	if err := runConfigInit(nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("config file not created: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "chlog configuration") {
		t.Error("expected config header comment")
	}
	if !strings.Contains(content, "changelog_file") {
		t.Error("expected changelog_file comment")
	}
}

func TestRunConfigInit_FailsIfExists(t *testing.T) {
	dir := t.TempDir()
	configFile = filepath.Join(dir, ".chlog.yaml")

	if err := os.WriteFile(configFile, []byte("existing"), 0644); err != nil {
		t.Fatal(err)
	}

	err := runConfigInit(nil, nil)
	if err == nil {
		t.Fatal("expected error when file already exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error = %q, want 'already exists'", err)
	}
}

func TestRunConfigInit_ForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	configFile = filepath.Join(dir, ".chlog.yaml")
	configForce = true
	t.Cleanup(func() { configForce = false })

	if err := os.WriteFile(configFile, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := runConfigInit(nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(configFile)
	if strings.Contains(string(data), "old content") {
		t.Error("expected old content to be overwritten")
	}
}

func TestRunConfigSet(t *testing.T) {
	trueVal := true

	tests := map[string]struct {
		key     string
		value   string
		check   func(*testing.T, *changelog.Config)
		wantErr bool
	}{
		"repo_url": {
			key: "repo_url", value: "https://github.com/test/repo",
			check: func(t *testing.T, c *changelog.Config) {
				if c.RepoURL != "https://github.com/test/repo" {
					t.Errorf("RepoURL = %q", c.RepoURL)
				}
			},
		},
		"changelog_file": {
			key: "changelog_file", value: "changelogs/CHANGELOG.yaml",
			check: func(t *testing.T, c *changelog.Config) {
				if c.ChangelogFile != "changelogs/CHANGELOG.yaml" {
					t.Errorf("ChangelogFile = %q", c.ChangelogFile)
				}
			},
		},
		"public_file": {
			key: "public_file", value: "docs/CHANGELOG.md",
			check: func(t *testing.T, c *changelog.Config) {
				if c.PublicFile != "docs/CHANGELOG.md" {
					t.Errorf("PublicFile = %q", c.PublicFile)
				}
			},
		},
		"internal_file": {
			key: "internal_file", value: "docs/CHANGELOG-internal.md",
			check: func(t *testing.T, c *changelog.Config) {
				if c.InternalFile != "docs/CHANGELOG-internal.md" {
					t.Errorf("InternalFile = %q", c.InternalFile)
				}
			},
		},
		"include_internal true": {
			key: "include_internal", value: "true",
			check: func(t *testing.T, c *changelog.Config) {
				if !c.IncludeInternal {
					t.Error("IncludeInternal should be true")
				}
			},
		},
		"include_internal false": {
			key: "include_internal", value: "false",
			check: func(t *testing.T, c *changelog.Config) {
				if c.IncludeInternal {
					t.Error("IncludeInternal should be false")
				}
			},
		},
		"strict_categories true": {
			key: "strict_categories", value: "true",
			check: func(t *testing.T, c *changelog.Config) {
				if c.StrictCategories == nil || *c.StrictCategories != trueVal {
					t.Error("StrictCategories should be true")
				}
			},
		},
		"strict_categories false": {
			key: "strict_categories", value: "false",
			check: func(t *testing.T, c *changelog.Config) {
				if c.StrictCategories == nil || *c.StrictCategories {
					t.Error("StrictCategories should be false")
				}
			},
		},
		"categories": {
			key: "categories", value: "added, changed, fixed",
			check: func(t *testing.T, c *changelog.Config) {
				if len(c.Categories) != 3 {
					t.Fatalf("Categories len = %d, want 3", len(c.Categories))
				}
				want := []string{"added", "changed", "fixed"}
				for i, w := range want {
					if c.Categories[i] != w {
						t.Errorf("Categories[%d] = %q, want %q", i, c.Categories[i], w)
					}
				}
			},
		},
		"unknown key": {
			key: "bad_key", value: "whatever",
			wantErr: true,
		},
		"include_internal bad value": {
			key: "include_internal", value: "yes",
			wantErr: true,
		},
		"strict_categories bad value": {
			key: "strict_categories", value: "yes",
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			configFile = filepath.Join(dir, ".chlog.yaml")

			err := runConfigSet(nil, []string{tc.key, tc.value})
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			cfg, err := changelog.LoadConfig(configFile)
			if err != nil {
				t.Fatalf("loading config: %v", err)
			}
			tc.check(t, cfg)
		})
	}
}

func TestRunConfigSet_UpdatesExistingFile(t *testing.T) {
	dir := t.TempDir()
	configFile = filepath.Join(dir, ".chlog.yaml")

	// Pre-populate with a value.
	initial := &changelog.Config{RepoURL: "https://github.com/old/repo"}
	if err := changelog.SaveConfig(initial, configFile); err != nil {
		t.Fatal(err)
	}

	if err := runConfigSet(nil, []string{"public_file", "out/CHANGELOG.md"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := changelog.LoadConfig(configFile)
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}
	if cfg.RepoURL != "https://github.com/old/repo" {
		t.Errorf("existing RepoURL was overwritten, got %q", cfg.RepoURL)
	}
	if cfg.PublicFile != "out/CHANGELOG.md" {
		t.Errorf("PublicFile = %q, want out/CHANGELOG.md", cfg.PublicFile)
	}
}

func TestRunConfigShow_NoError(t *testing.T) {
	dir := t.TempDir()
	configFile = filepath.Join(dir, ".chlog.yaml")

	// Should work with missing config (all defaults).
	if err := runConfigShow(nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should work with an existing config too.
	cfg := &changelog.Config{
		RepoURL:    "https://github.com/test/repo",
		PublicFile: "out/CHANGELOG.md",
	}
	if err := changelog.SaveConfig(cfg, configFile); err != nil {
		t.Fatal(err)
	}
	if err := runConfigShow(nil, nil); err != nil {
		t.Fatalf("unexpected error with config file: %v", err)
	}
}
