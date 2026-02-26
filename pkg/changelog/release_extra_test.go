package changelog

import "testing"

func TestRelease_PreservesEntriesInPromotedVersion(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{
				Version:  "unreleased",
				Added:    []string{"Feature A", "Feature B"},
				Fixed:    []string{"Bug X"},
				Changed:  []string{"API update"},
				Internal: Changes{Changed: []string{"Refactored handler"}},
			},
		},
	}

	if err := c.Release("1.0.0", "2024-06-01"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	released := c.Versions[1]
	if len(released.Added) != 2 {
		t.Errorf("added = %d, want 2", len(released.Added))
	}
	if len(released.Fixed) != 1 {
		t.Errorf("fixed = %d, want 1", len(released.Fixed))
	}
	if len(released.Changed) != 1 {
		t.Errorf("changed = %d, want 1", len(released.Changed))
	}
	if len(released.Internal.Changed) != 1 {
		t.Errorf("internal changed = %d, want 1", len(released.Internal.Changed))
	}
}

func TestRelease_NewUnreleasedIsClean(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{Version: "unreleased", Added: []string{"x"}},
		},
	}
	if err := c.Release("1.0.0", "2024-01-01"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unreleased := c.Versions[0]
	if !unreleased.IsUnreleased() {
		t.Error("first version should be unreleased")
	}
	if !unreleased.IsEmpty() {
		t.Error("new unreleased should have empty public changes")
	}
	if !unreleased.Internal.IsEmpty() {
		t.Error("new unreleased should have empty internal changes")
	}
}

func TestRelease_MultiplePreviousVersions(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{Version: "unreleased", Added: []string{"new"}},
			{Version: "2.0.0", Date: "2024-06-01", Added: []string{"v2"}},
			{Version: "1.0.0", Date: "2024-01-01", Added: []string{"v1"}},
		},
	}
	if err := c.Release("3.0.0", "2024-12-01"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(c.Versions) != 4 {
		t.Fatalf("expected 4 versions, got %d", len(c.Versions))
	}

	want := []string{"unreleased", "3.0.0", "2.0.0", "1.0.0"}
	for i, w := range want {
		if c.Versions[i].Version != w {
			t.Errorf("versions[%d] = %q, want %q", i, c.Versions[i].Version, w)
		}
	}
}

func TestRelease_VersionIsCorrectlyStamped(t *testing.T) {
	c := &Changelog{
		Project: "test",
		Versions: []Version{
			{Version: "unreleased", Added: []string{"feature"}},
		},
	}
	if err := c.Release("1.2.3", "2024-07-15"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	released := c.Versions[1]
	if released.Version != "1.2.3" {
		t.Errorf("version = %q, want 1.2.3", released.Version)
	}
	if released.Date != "2024-07-15" {
		t.Errorf("date = %q, want 2024-07-15", released.Date)
	}
}
