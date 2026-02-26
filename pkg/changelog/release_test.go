package changelog

import "testing"

func TestRelease(t *testing.T) {
	tests := map[string]struct {
		changelog *Changelog
		version   string
		date      string
		wantErr   string
	}{
		"basic release": {
			changelog: func() *Changelog {
				v := Version{Version: "unreleased"}
				v.Public.Append("added", "New feature")
				v2 := Version{Version: "1.0.0", Date: "2024-01-01"}
				v2.Public.Append("added", "Init")
				return &Changelog{Project: "test", Versions: []Version{v, v2}}
			}(),
			version: "2.0.0",
			date:    "2024-06-01",
		},
		"no unreleased": {
			changelog: func() *Changelog {
				v := Version{Version: "1.0.0", Date: "2024-01-01"}
				v.Public.Append("added", "Init")
				return &Changelog{Project: "test", Versions: []Version{v}}
			}(),
			version: "2.0.0",
			date:    "2024-06-01",
			wantErr: "no unreleased version found",
		},
		"empty unreleased": {
			changelog: &Changelog{
				Project:  "test",
				Versions: []Version{{Version: "unreleased"}},
			},
			version: "1.0.0",
			date:    "2024-06-01",
			wantErr: "unreleased version has no entries",
		},
		"duplicate version": {
			changelog: func() *Changelog {
				v := Version{Version: "unreleased"}
				v.Public.Append("fixed", "Bug fix")
				v2 := Version{Version: "1.0.0", Date: "2024-01-01"}
				v2.Public.Append("added", "Init")
				return &Changelog{Project: "test", Versions: []Version{v, v2}}
			}(),
			version: "1.0.0",
			date:    "2024-06-01",
			wantErr: `version "1.0.0" already exists`,
		},
		"duplicate version with v prefix": {
			changelog: func() *Changelog {
				v := Version{Version: "unreleased"}
				v.Public.Append("fixed", "Bug fix")
				v2 := Version{Version: "1.0.0", Date: "2024-01-01"}
				v2.Public.Append("added", "Init")
				return &Changelog{Project: "test", Versions: []Version{v, v2}}
			}(),
			version: "v1.0.0",
			date:    "2024-06-01",
			wantErr: `version "v1.0.0" already exists`,
		},
		"internal only entries": {
			changelog: func() *Changelog {
				v := Version{Version: "unreleased"}
				v.Internal.Append("changed", "Refactored internals")
				return &Changelog{Project: "test", Versions: []Version{v}}
			}(),
			version: "1.0.0",
			date:    "2024-06-01",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.changelog.Release(tc.version, tc.date)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErr)
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected error %q, got %q", tc.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// After release, first version should be unreleased
			if len(tc.changelog.Versions) < 2 {
				t.Fatal("expected at least 2 versions after release")
			}
			if !tc.changelog.Versions[0].IsUnreleased() {
				t.Errorf("first version should be unreleased, got %q", tc.changelog.Versions[0].Version)
			}
			if !tc.changelog.Versions[0].IsEmpty() {
				t.Error("new unreleased block should have empty changes")
			}

			// Second version should be the newly released one
			released := tc.changelog.Versions[1]
			if released.Version != tc.version {
				t.Errorf("released version = %q, want %q", released.Version, tc.version)
			}
			if released.Date != tc.date {
				t.Errorf("released date = %q, want %q", released.Date, tc.date)
			}
		})
	}
}

func TestRelease_PreservesExistingVersions(t *testing.T) {
	unreleased := Version{Version: "unreleased"}
	unreleased.Public.Append("added", "Feature A")
	unreleased.Public.Append("added", "Feature B")
	v11 := Version{Version: "1.1.0", Date: "2024-03-01"}
	v11.Public.Append("fixed", "Bug")
	v10 := Version{Version: "1.0.0", Date: "2024-01-01"}
	v10.Public.Append("added", "Init")

	c := &Changelog{
		Project:  "test",
		Versions: []Version{unreleased, v11, v10},
	}

	if err := c.Release("2.0.0", "2024-06-01"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(c.Versions) != 4 {
		t.Fatalf("expected 4 versions, got %d", len(c.Versions))
	}

	want := []string{"unreleased", "2.0.0", "1.1.0", "1.0.0"}
	for i, w := range want {
		if c.Versions[i].Version != w {
			t.Errorf("versions[%d] = %q, want %q", i, c.Versions[i].Version, w)
		}
	}

	// Verify the released version kept its entries
	released := c.Versions[1]
	if len(released.Public.Get("added")) != 2 {
		t.Errorf("released added = %d entries, want 2", len(released.Public.Get("added")))
	}
}
