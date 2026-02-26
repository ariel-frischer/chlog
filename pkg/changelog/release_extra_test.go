package changelog

import "testing"

func TestRelease_PreservesEntriesInPromotedVersion(t *testing.T) {
	v := Version{Version: "unreleased"}
	v.Public.Append("added", "Feature A")
	v.Public.Append("added", "Feature B")
	v.Public.Append("fixed", "Bug X")
	v.Public.Append("changed", "API update")
	v.Internal.Append("changed", "Refactored handler")

	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
	}

	if err := c.Release("1.0.0", "2024-06-01"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	released := c.Versions[1]
	if len(released.Public.Get("added")) != 2 {
		t.Errorf("added = %d, want 2", len(released.Public.Get("added")))
	}
	if len(released.Public.Get("fixed")) != 1 {
		t.Errorf("fixed = %d, want 1", len(released.Public.Get("fixed")))
	}
	if len(released.Public.Get("changed")) != 1 {
		t.Errorf("changed = %d, want 1", len(released.Public.Get("changed")))
	}
	if len(released.Internal.Get("changed")) != 1 {
		t.Errorf("internal changed = %d, want 1", len(released.Internal.Get("changed")))
	}
}

func TestRelease_NewUnreleasedIsClean(t *testing.T) {
	v := Version{Version: "unreleased"}
	v.Public.Append("added", "x")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
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
	unreleased := Version{Version: "unreleased"}
	unreleased.Public.Append("added", "new")
	v2 := Version{Version: "2.0.0", Date: "2024-06-01"}
	v2.Public.Append("added", "v2")
	v1 := Version{Version: "1.0.0", Date: "2024-01-01"}
	v1.Public.Append("added", "v1")

	c := &Changelog{
		Project:  "test",
		Versions: []Version{unreleased, v2, v1},
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
	v := Version{Version: "unreleased"}
	v.Public.Append("added", "feature")
	c := &Changelog{
		Project:  "test",
		Versions: []Version{v},
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
