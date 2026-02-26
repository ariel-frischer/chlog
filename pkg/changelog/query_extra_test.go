package changelog

import "testing"

func TestGetVersion_EmptyChangelog(t *testing.T) {
	c := &Changelog{Project: "empty"}
	_, err := c.GetVersion("1.0.0")
	if err == nil {
		t.Fatal("expected error for empty changelog")
	}
	if _, ok := err.(VersionNotFoundError); !ok {
		t.Errorf("expected VersionNotFoundError, got %T", err)
	}
}

func TestGetUnreleased_None(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "x")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}
	if v := c.GetUnreleased(); v != nil {
		t.Errorf("expected nil, got version %q", v.Version)
	}
}

func TestGetLatestRelease_UnreleasedOnly(t *testing.T) {
	c := &Changelog{
		Project:  "test",
		Versions: []Version{{Version: "unreleased"}},
	}
	if v := c.GetLatestRelease(); v != nil {
		t.Errorf("expected nil, got version %q", v.Version)
	}
}

func TestGetLatestRelease_Empty(t *testing.T) {
	c := &Changelog{Project: "test"}
	if v := c.GetLatestRelease(); v != nil {
		t.Errorf("expected nil for empty changelog, got %q", v.Version)
	}
}

func TestListVersions_Empty(t *testing.T) {
	c := &Changelog{Project: "test"}
	versions := c.ListVersions()
	if len(versions) != 0 {
		t.Errorf("expected 0 versions, got %d", len(versions))
	}
}

func TestGetLastN_MoreThanAvailable(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "only one")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}
	entries := c.GetLastN(100)
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestGetLastN_Zero(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "a")
	v.Public.Append("added", "b")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}
	entries := c.GetLastN(0)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestGetLastN_ExactCount(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "a")
	v.Public.Append("added", "b")
	v.Public.Append("added", "c")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}
	entries := c.GetLastN(3)
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestAllEntries_Order(t *testing.T) {
	unreleased := Version{Version: "unreleased"}
	unreleased.Public.Append("added", "newest")
	released := Version{Version: "1.0.0", Date: "2024-01-01"}
	released.Public.Append("added", "oldest")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{unreleased, released},
	}
	entries := c.AllEntries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Version != "unreleased" {
		t.Errorf("first entry should be from unreleased, got %q", entries[0].Version)
	}
	if entries[1].Version != "1.0.0" {
		t.Errorf("second entry should be from 1.0.0, got %q", entries[1].Version)
	}
}

func TestAllEntries_CategoryOrder(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "add")
	v.Public.Append("fixed", "fix")
	v.Public.Append("security", "sec")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}
	entries := c.AllEntries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	// Order follows append order: added, fixed, security
	if entries[0].Category != "added" {
		t.Errorf("first category = %q, want 'added'", entries[0].Category)
	}
	if entries[1].Category != "fixed" {
		t.Errorf("second category = %q, want 'fixed'", entries[1].Category)
	}
	if entries[2].Category != "security" {
		t.Errorf("third category = %q, want 'security'", entries[2].Category)
	}
}

func TestHasUnreleased_False(t *testing.T) {
	v := Version{Version: "1.0.0", Date: "2024-01-01"}
	v.Public.Append("added", "x")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}
	if c.HasUnreleased() {
		t.Error("expected false when no unreleased version")
	}
}

func TestGetVersionCount_Empty(t *testing.T) {
	c := &Changelog{Project: "test"}
	if got := c.GetVersionCount(); got != 0 {
		t.Errorf("GetVersionCount() = %d, want 0", got)
	}
}

func TestGetEntryCount_Empty(t *testing.T) {
	c := &Changelog{Project: "test"}
	if got := c.GetEntryCount(); got != 0 {
		t.Errorf("GetEntryCount() = %d, want 0", got)
	}
}

func TestGetEntryCount_MultipleVersions(t *testing.T) {
	unreleased := Version{Version: "unreleased"}
	unreleased.Public.Append("added", "a")
	released := Version{Version: "1.0.0", Date: "2024-01-01"}
	released.Public.Append("added", "b")
	released.Public.Append("fixed", "c")
	released.Public.Append("fixed", "d")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{unreleased, released},
	}
	if got := c.GetEntryCount(); got != 4 {
		t.Errorf("GetEntryCount() = %d, want 4", got)
	}
}

func TestGetVersion_NormalizesInput(t *testing.T) {
	v := Version{Version: "2.0.0", Date: "2024-06-01"}
	v.Public.Append("added", "x")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}
	tests := map[string]string{
		"exact":    "2.0.0",
		"v-prefix": "v2.0.0",
		"V-prefix": "V2.0.0",
	}
	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			v, err := c.GetVersion(input)
			if err != nil {
				t.Fatalf("GetVersion(%q) error: %v", input, err)
			}
			if v.Version != "2.0.0" {
				t.Errorf("got version %q, want 2.0.0", v.Version)
			}
		})
	}
}
