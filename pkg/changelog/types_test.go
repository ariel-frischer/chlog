package changelog

import "testing"

func TestVersion_IsUnreleased(t *testing.T) {
	tests := map[string]struct {
		version string
		want    bool
	}{
		"lowercase":  {version: "unreleased", want: true},
		"titlecase":  {version: "Unreleased", want: true},
		"uppercase":  {version: "UNRELEASED", want: true},
		"semver":     {version: "1.0.0", want: false},
		"empty":      {version: "", want: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := &Version{Version: tc.version}
			if got := v.IsUnreleased(); got != tc.want {
				t.Errorf("IsUnreleased() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestChanges_IsEmpty(t *testing.T) {
	tests := map[string]struct {
		changes Changes
		want    bool
	}{
		"empty":     {changes: Changes{}, want: true},
		"has_added": {changes: Changes{Added: []string{"x"}}, want: false},
		"has_fixed": {changes: Changes{Fixed: []string{"y"}}, want: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.changes.IsEmpty(); got != tc.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestChanges_Count(t *testing.T) {
	tests := map[string]struct {
		changes Changes
		want    int
	}{
		"empty":    {changes: Changes{}, want: 0},
		"one":      {changes: Changes{Added: []string{"a"}}, want: 1},
		"multiple": {changes: Changes{Added: []string{"a", "b"}, Fixed: []string{"c"}}, want: 3},
		"all_categories": {
			changes: Changes{
				Added: []string{"a"}, Changed: []string{"b"},
				Deprecated: []string{"c"}, Removed: []string{"d"},
				Fixed: []string{"e"}, Security: []string{"f"},
			},
			want: 6,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.changes.Count(); got != tc.want {
				t.Errorf("Count() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestChanges_CategoryEntries(t *testing.T) {
	c := Changes{
		Added:   []string{"a"},
		Fixed:   []string{"f"},
		Changed: []string{"c"},
	}
	tests := map[string]struct {
		category string
		wantLen  int
	}{
		"added":   {category: "added", wantLen: 1},
		"fixed":   {category: "fixed", wantLen: 1},
		"changed": {category: "changed", wantLen: 1},
		"empty":   {category: "removed", wantLen: 0},
		"unknown": {category: "bogus", wantLen: 0},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := c.CategoryEntries(tc.category)
			if len(got) != tc.wantLen {
				t.Errorf("CategoryEntries(%q) len = %d, want %d", tc.category, len(got), tc.wantLen)
			}
		})
	}
}

func TestValidCategories(t *testing.T) {
	cats := ValidCategories()
	if len(cats) != 6 {
		t.Fatalf("expected 6 categories, got %d", len(cats))
	}
	expected := []string{"added", "changed", "deprecated", "removed", "fixed", "security"}
	for i, cat := range cats {
		if cat != expected[i] {
			t.Errorf("category[%d] = %q, want %q", i, cat, expected[i])
		}
	}
}
