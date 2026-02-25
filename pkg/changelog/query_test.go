package changelog

import (
	"testing"
)

func loadTestChangelog(t *testing.T) *Changelog {
	t.Helper()
	c, err := Load("testdata/valid.yaml")
	if err != nil {
		t.Fatalf("loading test fixture: %v", err)
	}
	return c
}

func TestGetVersion(t *testing.T) {
	c := loadTestChangelog(t)
	tests := map[string]struct {
		version string
		wantErr bool
	}{
		"exact":         {version: "1.0.0", wantErr: false},
		"v_prefix":      {version: "v1.0.0", wantErr: false},
		"unreleased":    {version: "unreleased", wantErr: false},
		"not_found":     {version: "9.9.9", wantErr: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v, err := c.GetVersion(tc.version)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				if _, ok := err.(VersionNotFoundError); !ok {
					t.Errorf("expected VersionNotFoundError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if v == nil {
					t.Error("expected non-nil version")
				}
			}
		})
	}
}

func TestGetUnreleased(t *testing.T) {
	c := loadTestChangelog(t)
	v := c.GetUnreleased()
	if v == nil {
		t.Fatal("expected unreleased version")
	}
	if !v.IsUnreleased() {
		t.Errorf("expected unreleased, got %q", v.Version)
	}
}

func TestGetLatestRelease(t *testing.T) {
	c := loadTestChangelog(t)
	v := c.GetLatestRelease()
	if v == nil {
		t.Fatal("expected latest release")
	}
	if v.Version != "1.1.0" {
		t.Errorf("expected 1.1.0, got %q", v.Version)
	}
}

func TestListVersions(t *testing.T) {
	c := loadTestChangelog(t)
	versions := c.ListVersions()
	if len(versions) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(versions))
	}
	if versions[0] != "unreleased" {
		t.Errorf("first version = %q, want unreleased", versions[0])
	}
}

func TestGetLastN(t *testing.T) {
	c := loadTestChangelog(t)
	entries := c.GetLastN(2)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestAllEntries(t *testing.T) {
	c := loadTestChangelog(t)
	entries := c.AllEntries()
	// unreleased: 1 added, 1.1.0: 2 added + 1 fixed, 1.0.0: 1 added + 1 security = 6
	if len(entries) != 6 {
		t.Errorf("expected 6 entries, got %d", len(entries))
	}
}

func TestGetVersionCount(t *testing.T) {
	c := loadTestChangelog(t)
	if got := c.GetVersionCount(); got != 3 {
		t.Errorf("GetVersionCount() = %d, want 3", got)
	}
}

func TestGetEntryCount(t *testing.T) {
	c := loadTestChangelog(t)
	if got := c.GetEntryCount(); got != 6 {
		t.Errorf("GetEntryCount() = %d, want 6", got)
	}
}

func TestHasUnreleased(t *testing.T) {
	c := loadTestChangelog(t)
	if !c.HasUnreleased() {
		t.Error("expected HasUnreleased() = true")
	}
}
