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
	os.WriteFile(path, []byte("repo_url: https://github.com/example/repo\n"), 0644)

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
	os.WriteFile(path, []byte("include_internal: true\n"), 0644)

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
	os.WriteFile(path, []byte("repo_url: https://example.com\n"), 0644)

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
	os.WriteFile(path, []byte("public_file: public.md\ninternal_file: docs/internal.md\n"), 0644)

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
