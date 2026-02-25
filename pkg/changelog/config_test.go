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
