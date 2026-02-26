package changelog

import "testing"

func TestVersion_IsUnreleased(t *testing.T) {
	tests := map[string]struct {
		version string
		want    bool
	}{
		"lowercase": {version: "unreleased", want: true},
		"titlecase": {version: "Unreleased", want: true},
		"uppercase": {version: "UNRELEASED", want: true},
		"semver":    {version: "1.0.0", want: false},
		"empty":     {version: "", want: false},
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
		"has_added": {changes: makeChanges("added", "x"), want: false},
		"has_fixed": {changes: makeChanges("fixed", "y"), want: false},
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
		"one":      {changes: makeChanges("added", "a"), want: 1},
		"multiple": {changes: makeChangesMulti(map[string][]string{"added": {"a", "b"}, "fixed": {"c"}}), want: 3},
		"all_categories": {
			changes: makeChangesMulti(map[string][]string{
				"added": {"a"}, "changed": {"b"},
				"deprecated": {"c"}, "removed": {"d"},
				"fixed": {"e"}, "security": {"f"},
			}),
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

func TestChanges_Get(t *testing.T) {
	c := makeChangesMulti(map[string][]string{
		"added": {"a"}, "fixed": {"f"}, "changed": {"c"},
	})
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
			got := c.Get(tc.category)
			if len(got) != tc.wantLen {
				t.Errorf("Get(%q) len = %d, want %d", tc.category, len(got), tc.wantLen)
			}
		})
	}
}

func TestDefaultCategories(t *testing.T) {
	cats := DefaultCategories
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

func TestChanges_Remove(t *testing.T) {
	tests := map[string]struct {
		changes   Changes
		category  string
		text      string
		substring bool
		wantText  string
		wantErr   interface{} // nil, or error type
		wantCats  int         // expected category count after removal
	}{
		"exact_match": {
			changes:  makeChangesMulti(map[string][]string{"added": {"Feature A", "Feature B"}}),
			category: "added", text: "Feature A",
			wantText: "Feature A", wantCats: 1,
		},
		"exact_no_match": {
			changes:  makeChanges("added", "Feature A"),
			category: "added", text: "Feature X",
			wantErr: EntryNotFoundError{},
		},
		"category_not_found": {
			changes:  makeChanges("added", "Feature A"),
			category: "removed", text: "anything",
			wantErr: CategoryNotFoundError{},
		},
		"substring_match": {
			changes:   makeChangesMulti(map[string][]string{"fixed": {"Fix login timeout", "Fix signup error"}}),
			category:  "fixed", text: "login",
			substring: true,
			wantText:  "Fix login timeout", wantCats: 1,
		},
		"substring_case_insensitive": {
			changes:   makeChanges("added", "Support Dark Mode"),
			category:  "added", text: "dark mode",
			substring: true,
			wantText:  "Support Dark Mode", wantCats: 0,
		},
		"substring_multiple_matches": {
			changes:   makeChangesMulti(map[string][]string{"fixed": {"Fix login timeout", "Fix login redirect"}}),
			category:  "fixed", text: "login",
			substring: true,
			wantErr:   MultipleMatchError{},
		},
		"empty_category_cleanup": {
			changes:  makeChanges("added", "only entry"),
			category: "added", text: "only entry",
			wantText: "only entry", wantCats: 0,
		},
		"order_preserved": {
			changes:  makeChangesMulti(map[string][]string{"added": {"A", "B", "C"}}),
			category: "added", text: "B",
			wantText: "B", wantCats: 1,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.changes.Remove(tc.category, tc.text, tc.substring)

			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				switch tc.wantErr.(type) {
				case CategoryNotFoundError:
					if _, ok := err.(CategoryNotFoundError); !ok {
						t.Fatalf("expected CategoryNotFoundError, got %T", err)
					}
				case EntryNotFoundError:
					if _, ok := err.(EntryNotFoundError); !ok {
						t.Fatalf("expected EntryNotFoundError, got %T", err)
					}
				case MultipleMatchError:
					if _, ok := err.(MultipleMatchError); !ok {
						t.Fatalf("expected MultipleMatchError, got %T", err)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantText {
				t.Errorf("Remove() = %q, want %q", got, tc.wantText)
			}
			if len(tc.changes.Categories) != tc.wantCats {
				t.Errorf("categories count = %d, want %d", len(tc.changes.Categories), tc.wantCats)
			}

			// For order_preserved, verify remaining entries
			if name == "order_preserved" {
				entries := tc.changes.Get("added")
				if len(entries) != 2 || entries[0] != "A" || entries[1] != "C" {
					t.Errorf("remaining entries = %v, want [A C]", entries)
				}
			}
		})
	}
}

// makeChanges creates a Changes with a single category and single entry.
func makeChanges(category, entry string) Changes {
	return Changes{Categories: []CategoryEntry{{Name: category, Entries: []string{entry}}}}
}

// makeChangesMulti creates a Changes from a map. Order follows DefaultCategories for known
// categories, then any unknown in map iteration order.
func makeChangesMulti(m map[string][]string) Changes {
	var c Changes
	// First, add in default order for deterministic output
	for _, cat := range DefaultCategories {
		if entries, ok := m[cat]; ok {
			c.Categories = append(c.Categories, CategoryEntry{Name: cat, Entries: entries})
		}
	}
	// Then any non-default categories
	for cat, entries := range m {
		found := false
		for _, dc := range DefaultCategories {
			if cat == dc {
				found = true
				break
			}
		}
		if !found {
			c.Categories = append(c.Categories, CategoryEntry{Name: cat, Entries: entries})
		}
	}
	return c
}
