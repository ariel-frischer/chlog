package changelog

import "testing"

func TestCleanDescription(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"capitalize first letter": {
			input: "add feature",
			want:  "Add feature",
		},
		"already capitalized": {
			input: "Add feature",
			want:  "Add feature",
		},
		"truncate at period": {
			input: "add feature. more details here",
			want:  "Add feature",
		},
		"truncate at exclamation": {
			input: "fix crash! important details",
			want:  "Fix crash",
		},
		"truncate at question mark": {
			input: "why does this break? no one knows",
			want:  "Why does this break",
		},
		"trim whitespace": {
			input: "  add feature  ",
			want:  "Add feature",
		},
		"empty string": {
			input: "",
			want:  "",
		},
		"single word": {
			input: "feature",
			want:  "Feature",
		},
		"period at end": {
			input: "add feature.",
			want:  "Add feature",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cleanDescription(tc.input); got != tc.want {
				t.Errorf("cleanDescription(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestChanges_Append_NewAndExisting(t *testing.T) {
	tests := map[string]struct {
		category string
		entry    string
		check    func(*Changes) bool
	}{
		"added":      {category: "added", entry: "x", check: func(c *Changes) bool { return len(c.Get("added")) == 1 && c.Get("added")[0] == "x" }},
		"changed":    {category: "changed", entry: "x", check: func(c *Changes) bool { return len(c.Get("changed")) == 1 }},
		"deprecated": {category: "deprecated", entry: "x", check: func(c *Changes) bool { return len(c.Get("deprecated")) == 1 }},
		"removed":    {category: "removed", entry: "x", check: func(c *Changes) bool { return len(c.Get("removed")) == 1 }},
		"fixed":      {category: "fixed", entry: "x", check: func(c *Changes) bool { return len(c.Get("fixed")) == 1 }},
		"security":   {category: "security", entry: "x", check: func(c *Changes) bool { return len(c.Get("security")) == 1 }},
		"custom": {category: "performance", entry: "x", check: func(c *Changes) bool {
			return len(c.Get("performance")) == 1
		}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c := &Changes{}
			c.Append(tc.category, tc.entry)
			if !tc.check(c) {
				t.Errorf("Append(%q) did not produce expected result", tc.category)
			}
		})
	}
}

func TestScaffold_InternalNonBreakingRefactorAndPerf(t *testing.T) {
	commits := []GitCommit{
		{Hash: "a", Subject: "refactor: simplify handler"},
		{Hash: "b", Subject: "perf: cache results"},
	}
	v := Scaffold(commits, ScaffoldOptions{})

	if !v.IsEmpty() {
		t.Error("non-breaking refactor/perf should not appear in public changes")
	}
	if len(v.Internal.Get("changed")) != 2 {
		t.Errorf("internal changed = %d, want 2", len(v.Internal.Get("changed")))
	}
}

func TestScaffold_AllSkippedCommits(t *testing.T) {
	commits := []GitCommit{
		{Hash: "a", Subject: "chore: update deps"},
		{Hash: "b", Subject: "docs: update readme"},
		{Hash: "c", Subject: "test: add tests"},
		{Hash: "d", Subject: "ci: update pipeline"},
		{Hash: "e", Subject: "style: format code"},
		{Hash: "f", Subject: "build: update makefile"},
	}
	v := Scaffold(commits, ScaffoldOptions{})

	if !v.IsEmpty() {
		t.Error("all skipped types should produce empty public changes")
	}
	if !v.Internal.IsEmpty() {
		t.Error("all skipped types should produce empty internal changes")
	}
}

func TestScaffold_MixedCommitTypes(t *testing.T) {
	commits := []GitCommit{
		{Hash: "a", Subject: "feat: add dark mode"},
		{Hash: "b", Subject: "fix: resolve crash"},
		{Hash: "c", Subject: "deprecate: old API"},
		{Hash: "d", Subject: "remove: legacy module"},
		{Hash: "e", Subject: "chore: update deps"},
		{Hash: "f", Subject: "not a conventional commit"},
		{Hash: "g", Subject: "refactor: clean internals"},
	}
	v := Scaffold(commits, ScaffoldOptions{})

	if len(v.Public.Get("added")) != 1 {
		t.Errorf("added = %d, want 1", len(v.Public.Get("added")))
	}
	if len(v.Public.Get("fixed")) != 1 {
		t.Errorf("fixed = %d, want 1", len(v.Public.Get("fixed")))
	}
	if len(v.Public.Get("deprecated")) != 1 {
		t.Errorf("deprecated = %d, want 1", len(v.Public.Get("deprecated")))
	}
	if len(v.Public.Get("removed")) != 1 {
		t.Errorf("removed = %d, want 1", len(v.Public.Get("removed")))
	}
	if len(v.Internal.Get("changed")) != 1 {
		t.Errorf("internal changed = %d, want 1", len(v.Internal.Get("changed")))
	}
}

func TestParseConventionalCommit_EdgeCases(t *testing.T) {
	tests := map[string]struct {
		subject      string
		wantCat      string
		wantDesc     string
		wantBreaking bool
		wantInternal bool
	}{
		"extra spaces around colon": {
			subject:  "feat:   spaced description",
			wantCat:  "added",
			wantDesc: "Spaced description",
		},
		"scope with hyphen": {
			subject:  "feat(my-scope): add thing",
			wantCat:  "added",
			wantDesc: "Add thing",
		},
		"unknown type non-breaking": {
			subject: "random: something",
			wantCat: "",
		},
		"breaking unknown type": {
			subject:      "unknown!: drop support",
			wantCat:      "changed",
			wantDesc:     "BREAKING: Drop support",
			wantBreaking: true,
		},
		"empty subject after type": {
			subject: "",
			wantCat: "",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cat, desc, breaking, internal := ParseConventionalCommit(tc.subject)
			if cat != tc.wantCat {
				t.Errorf("category = %q, want %q", cat, tc.wantCat)
			}
			if tc.wantCat != "" && desc != tc.wantDesc {
				t.Errorf("description = %q, want %q", desc, tc.wantDesc)
			}
			if breaking != tc.wantBreaking {
				t.Errorf("breaking = %v, want %v", breaking, tc.wantBreaking)
			}
			if internal != tc.wantInternal {
				t.Errorf("internal = %v, want %v", internal, tc.wantInternal)
			}
		})
	}
}
