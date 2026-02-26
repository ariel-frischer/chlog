package changelog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Missing(t *testing.T) {
	cfg, err := LoadConfig("/nonexistent/.chlog.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RepoURL != "" {
		t.Errorf("expected empty RepoURL, got %q", cfg.RepoURL)
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".chlog.yaml")
	if err := os.WriteFile(path, []byte("repo_url: https://github.com/example/repo\n"), 0644); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RepoURL != "https://github.com/example/repo" {
		t.Errorf("RepoURL = %q, want %q", cfg.RepoURL, "https://github.com/example/repo")
	}
}

func TestLoadConfig_IncludeInternal(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".chlog.yaml")
	if err := os.WriteFile(path, []byte("include_internal: true\n"), 0644); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.IncludeInternal {
		t.Error("expected IncludeInternal=true, got false")
	}
}

func TestLoadConfig_IncludeInternalDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".chlog.yaml")
	if err := os.WriteFile(path, []byte("repo_url: https://example.com\n"), 0644); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.IncludeInternal {
		t.Error("expected IncludeInternal=false by default, got true")
	}
}

func TestResolveRepoURL_ConfigOverridesGit(t *testing.T) {
	cfg := &Config{RepoURL: "https://example.com/my/repo"}
	got := ResolveRepoURL(cfg)
	if got != "https://example.com/my/repo" {
		t.Errorf("ResolveRepoURL = %q, want config value", got)
	}
}

func TestResolveRepoURL_NilConfig(t *testing.T) {
	// Should not panic with nil config
	_ = ResolveRepoURL(nil)
}

func TestLoadConfig_FileOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".chlog.yaml")
	if err := os.WriteFile(path, []byte("public_file: public.md\ninternal_file: docs/internal.md\n"), 0644); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PublicFile != "public.md" {
		t.Errorf("PublicFile = %q, want %q", cfg.PublicFile, "public.md")
	}
	if cfg.InternalFile != "docs/internal.md" {
		t.Errorf("InternalFile = %q, want %q", cfg.InternalFile, "docs/internal.md")
	}
}

func TestConfig_FilePathDefaults(t *testing.T) {
	cfg := &Config{}
	if got := cfg.PublicFilePath(); got != DefaultPublicFile {
		t.Errorf("PublicFilePath() = %q, want %q", got, DefaultPublicFile)
	}
	if got := cfg.InternalFilePath(); got != DefaultInternalFile {
		t.Errorf("InternalFilePath() = %q, want %q", got, DefaultInternalFile)
	}
}

func TestConfig_FilePathOverrides(t *testing.T) {
	cfg := &Config{PublicFile: "custom.md", InternalFile: "docs/custom-internal.md"}
	if got := cfg.PublicFilePath(); got != "custom.md" {
		t.Errorf("PublicFilePath() = %q, want %q", got, "custom.md")
	}
	if got := cfg.InternalFilePath(); got != "docs/custom-internal.md" {
		t.Errorf("InternalFilePath() = %q, want %q", got, "docs/custom-internal.md")
	}
}

func TestSaveConfig(t *testing.T) {
	tests := map[string]struct {
		cfg  *Config
		want map[string]string
	}{
		"all fields": {
			cfg: &Config{
				RepoURL:         "https://github.com/example/repo",
				IncludeInternal: true,
				PublicFile:      "docs/CHANGELOG.md",
				InternalFile:    "docs/internal.md",
			},
			want: map[string]string{
				"repo_url":         "https://github.com/example/repo",
				"include_internal": "true",
				"public_file":      "docs/CHANGELOG.md",
				"internal_file":    "docs/internal.md",
			},
		},
		"repo url only": {
			cfg: &Config{RepoURL: "https://github.com/example/repo"},
			want: map[string]string{
				"repo_url": "https://github.com/example/repo",
			},
		},
		"empty config": {
			cfg:  &Config{},
			want: map[string]string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, ".chlog.yaml")

			if err := SaveConfig(tc.cfg, path); err != nil {
				t.Fatalf("SaveConfig: %v", err)
			}

			got, err := LoadConfig(path)
			if err != nil {
				t.Fatalf("LoadConfig: %v", err)
			}

			if got.RepoURL != tc.cfg.RepoURL {
				t.Errorf("RepoURL = %q, want %q", got.RepoURL, tc.cfg.RepoURL)
			}
			if got.IncludeInternal != tc.cfg.IncludeInternal {
				t.Errorf("IncludeInternal = %v, want %v", got.IncludeInternal, tc.cfg.IncludeInternal)
			}
			if got.PublicFile != tc.cfg.PublicFile {
				t.Errorf("PublicFile = %q, want %q", got.PublicFile, tc.cfg.PublicFile)
			}
			if got.InternalFile != tc.cfg.InternalFile {
				t.Errorf("InternalFile = %q, want %q", got.InternalFile, tc.cfg.InternalFile)
			}
		})
	}
}

func TestSaveConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".chlog.yaml")

	original := &Config{
		RepoURL:         "https://github.com/test/project",
		IncludeInternal: true,
	}

	if err := SaveConfig(original, path); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	if loaded.RepoURL != original.RepoURL {
		t.Errorf("RepoURL mismatch: got %q, want %q", loaded.RepoURL, original.RepoURL)
	}
	if loaded.IncludeInternal != original.IncludeInternal {
		t.Errorf("IncludeInternal mismatch: got %v, want %v", loaded.IncludeInternal, original.IncludeInternal)
	}
}

func TestAllowedCategories(t *testing.T) {
	// Default: returns DefaultCategories
	cfg := &Config{}
	allowed := cfg.AllowedCategories()
	if len(allowed) != 6 {
		t.Errorf("default allowed = %d, want 6", len(allowed))
	}

	// Custom categories
	cfg2 := &Config{Categories: []string{"added", "performance"}}
	allowed2 := cfg2.AllowedCategories()
	if len(allowed2) != 2 {
		t.Errorf("custom allowed = %d, want 2", len(allowed2))
	}

	// Non-strict: returns nil
	strictFalse := false
	cfg3 := &Config{StrictCategories: &strictFalse}
	if cfg3.AllowedCategories() != nil {
		t.Error("non-strict should return nil")
	}

	// Strict true explicitly: returns default
	strictTrue := true
	cfg4 := &Config{StrictCategories: &strictTrue}
	if len(cfg4.AllowedCategories()) != 6 {
		t.Errorf("strict=true allowed = %d, want 6", len(cfg4.AllowedCategories()))
	}
}

func TestSaveConfig_BadPath(t *testing.T) {
	cfg := &Config{RepoURL: "https://example.com"}
	err := SaveConfig(cfg, "/nonexistent/dir/.chlog.yaml")
	if err == nil {
		t.Fatal("expected error for bad path, got nil")
	}
}
