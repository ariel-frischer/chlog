package changelog

import "testing"

func TestParseConventionalCommit(t *testing.T) {
	tests := map[string]struct {
		subject      string
		wantCat      string
		wantDesc     string
		wantBreaking bool
		wantInternal bool
	}{
		"feat": {
			subject:  "feat: add user authentication",
			wantCat:  "added",
			wantDesc: "Add user authentication",
		},
		"fix": {
			subject:  "fix: resolve login redirect",
			wantCat:  "fixed",
			wantDesc: "Resolve login redirect",
		},
		"refactor": {
			subject:      "refactor: simplify auth middleware",
			wantCat:      "changed",
			wantDesc:     "Simplify auth middleware",
			wantInternal: true,
		},
		"perf": {
			subject:      "perf: optimize database queries",
			wantCat:      "changed",
			wantDesc:     "Optimize database queries",
			wantInternal: true,
		},
		"deprecate": {
			subject:  "deprecate: old API endpoint",
			wantCat:  "deprecated",
			wantDesc: "Old API endpoint",
		},
		"remove": {
			subject:  "remove: legacy auth system",
			wantCat:  "removed",
			wantDesc: "Legacy auth system",
		},
		"breaking": {
			subject:      "feat!: new auth system",
			wantCat:      "changed",
			wantDesc:     "BREAKING: New auth system",
			wantBreaking: true,
		},
		"scoped": {
			subject:  "feat(auth): add OAuth2 support",
			wantCat:  "added",
			wantDesc: "Add OAuth2 support",
		},
		"scoped_breaking": {
			subject:      "refactor(api)!: change response format",
			wantCat:      "changed",
			wantDesc:     "BREAKING: Change response format",
			wantBreaking: true,
		},
		"skip_chore": {
			subject: "chore: update deps",
			wantCat: "",
		},
		"skip_docs": {
			subject: "docs: update readme",
			wantCat: "",
		},
		"skip_test": {
			subject: "test: add unit tests",
			wantCat: "",
		},
		"skip_ci": {
			subject: "ci: update pipeline",
			wantCat: "",
		},
		"skip_style": {
			subject: "style: format code",
			wantCat: "",
		},
		"skip_build": {
			subject: "build: update makefile",
			wantCat: "",
		},
		"non_conventional": {
			subject: "just a regular commit message",
			wantCat: "",
		},
		"first_sentence": {
			subject:  "feat: add feature. This is extra detail",
			wantCat:  "added",
			wantDesc: "Add feature",
		},
		"breaking_chore": {
			subject:      "chore!: drop node 14 support",
			wantCat:      "changed",
			wantDesc:     "BREAKING: Drop node 14 support",
			wantBreaking: true,
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

func TestScaffold(t *testing.T) {
	commits := []GitCommit{
		{Hash: "aaa", Subject: "feat: add dark mode"},
		{Hash: "bbb", Subject: "fix: resolve crash on startup"},
		{Hash: "ccc", Subject: "chore: update deps"},
		{Hash: "ddd", Subject: "refactor!: new config format"},
	}

	v := Scaffold(commits, ScaffoldOptions{})
	if v.Version != "unreleased" {
		t.Errorf("version = %q, want unreleased", v.Version)
	}
	if len(v.Public.Get("added")) != 1 {
		t.Errorf("added count = %d, want 1", len(v.Public.Get("added")))
	}
	if len(v.Public.Get("fixed")) != 1 {
		t.Errorf("fixed count = %d, want 1", len(v.Public.Get("fixed")))
	}
	if len(v.Public.Get("changed")) != 1 {
		t.Errorf("changed count = %d, want 1 (breaking refactor)", len(v.Public.Get("changed")))
	}
}

func TestScaffold_CustomVersion(t *testing.T) {
	commits := []GitCommit{
		{Hash: "aaa", Subject: "feat: add feature"},
	}
	v := Scaffold(commits, ScaffoldOptions{Version: "2.0.0"})
	if v.Version != "2.0.0" {
		t.Errorf("version = %q, want 2.0.0", v.Version)
	}
}

func TestScaffold_EmptyCommits(t *testing.T) {
	v := Scaffold(nil, ScaffoldOptions{})
	if !v.IsEmpty() {
		t.Error("expected empty changes for nil commits")
	}
}
