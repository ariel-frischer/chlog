package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildConfigTestCLI(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "chlog")
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("building CLI: %v\n%s", err, output)
	}
	return binary
}

func runConfigTestCLI(t *testing.T, binary, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("chlog %v: %v\n%s", args, err, output)
	}
}

func TestConfiguredCategoriesRoundTrip(t *testing.T) {
	binary := buildConfigTestCLI(t)
	for _, setting := range []string{"categories", "strict_categories"} {
		t.Run(setting, func(t *testing.T) {
			dir := t.TempDir()
			runConfigTestCLI(t, binary, dir, "init", "--project", "test")
			value := "added,changed,fixed,performance"
			if setting == "strict_categories" {
				value = "false"
			}
			runConfigTestCLI(t, binary, dir, "config", "set", setting, value)
			runConfigTestCLI(t, binary, dir, "add", "performance", "Faster public queries")
			runConfigTestCLI(t, binary, dir, "validate")
			runConfigTestCLI(t, binary, dir, "add", "performance", "Faster internal queries", "--internal")
			runConfigTestCLI(t, binary, dir, "sync", "--internal")
			runConfigTestCLI(t, binary, dir, "check", "--internal")
			data, err := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
			if err != nil {
				t.Fatal(err)
			}
			for _, entry := range []string{"Faster public queries", "Faster internal queries"} {
				if !strings.Contains(string(data), entry) {
					t.Errorf("rendered changelog missing %q: %s", entry, data)
				}
			}
		})
	}
}

func TestConfiguredPublicFileSyncCheck(t *testing.T) {
	binary := buildConfigTestCLI(t)
	dir := t.TempDir()
	runConfigTestCLI(t, binary, dir, "init", "--project", "test")
	if err := os.Mkdir(filepath.Join(dir, "docs"), 0755); err != nil {
		t.Fatal(err)
	}
	runConfigTestCLI(t, binary, dir, "config", "set", "public_file", "docs/CHANGELOG.md")
	defaultPath := filepath.Join(dir, "CHANGELOG.md")
	if err := os.WriteFile(defaultPath, []byte("leave this alone"), 0644); err != nil {
		t.Fatal(err)
	}
	runConfigTestCLI(t, binary, dir, "sync")
	target := filepath.Join(dir, "docs", "CHANGELOG.md")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Initial project setup") {
		t.Fatalf("configured output missing entry: %s", data)
	}
	runConfigTestCLI(t, binary, dir, "check")
	data, err = os.ReadFile(defaultPath)
	if err != nil || string(data) != "leave this alone" {
		t.Fatalf("default output changed: %q, %v", data, err)
	}
	if err := os.WriteFile(target, []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(binary, "check")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != 1 || !strings.Contains(string(output), "docs/CHANGELOG.md is out of sync") {
		t.Fatalf("stale configured output: error %v, output %s", err, output)
	}
}
