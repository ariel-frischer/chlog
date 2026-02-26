package changelog

import "testing"

func TestParseGitLog(t *testing.T) {
	tests := map[string]struct {
		input string
		want  int
		first GitCommit
	}{
		"single commit": {
			input: "abc123 feat: add feature",
			want:  1,
			first: GitCommit{Hash: "abc123", Subject: "feat: add feature"},
		},
		"multiple commits": {
			input: "abc123 feat: add feature\ndef456 fix: bug fix\nghi789 chore: update deps",
			want:  3,
			first: GitCommit{Hash: "abc123", Subject: "feat: add feature"},
		},
		"empty input": {
			input: "",
			want:  0,
		},
		"subject with spaces": {
			input: "abc123 feat: add multi word feature description here",
			want:  1,
			first: GitCommit{Hash: "abc123", Subject: "feat: add multi word feature description here"},
		},
		"hash only no space": {
			input: "abc123",
			want:  0,
		},
		"full sha": {
			input: "abc123def456abc123def456abc123def456abc12 fix: something",
			want:  1,
			first: GitCommit{Hash: "abc123def456abc123def456abc123def456abc12", Subject: "fix: something"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := parseGitLog(tc.input)
			if len(got) != tc.want {
				t.Fatalf("parseGitLog() returned %d commits, want %d", len(got), tc.want)
			}
			if tc.want > 0 {
				if got[0].Hash != tc.first.Hash {
					t.Errorf("first.Hash = %q, want %q", got[0].Hash, tc.first.Hash)
				}
				if got[0].Subject != tc.first.Subject {
					t.Errorf("first.Subject = %q, want %q", got[0].Subject, tc.first.Subject)
				}
			}
		})
	}
}

func TestNormalizeGitURL_Extended(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"https with .git": {
			input: "https://github.com/org/repo.git",
			want:  "https://github.com/org/repo",
		},
		"ssh": {
			input: "git@github.com:org/repo.git",
			want:  "https://github.com/org/repo",
		},
		"https no suffix": {
			input: "https://gitlab.com/org/repo",
			want:  "https://gitlab.com/org/repo",
		},
		"ssh no suffix": {
			input: "git@gitlab.com:group/project",
			want:  "https://gitlab.com/group/project",
		},
		"ssh nested group": {
			input: "git@gitlab.com:group/subgroup/project.git",
			want:  "https://gitlab.com/group/subgroup/project",
		},
		"https already clean": {
			input: "https://github.com/user/project",
			want:  "https://github.com/user/project",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := normalizeGitURL(tc.input); got != tc.want {
				t.Errorf("normalizeGitURL(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
