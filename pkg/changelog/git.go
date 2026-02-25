package changelog

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitCommit represents a single git commit.
type GitCommit struct {
	Hash    string
	Subject string
}

// DetectRepoURL returns the remote origin URL, normalized to HTTPS without .git suffix.
func DetectRepoURL() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", fmt.Errorf("getting remote URL: %w", err)
	}
	return normalizeGitURL(strings.TrimSpace(string(out))), nil
}

// GitLog returns commits since the given tag (or all commits if sinceTag is empty).
func GitLog(sinceTag string) ([]GitCommit, error) {
	args := []string{"log", "--format=%H %s"}
	if sinceTag != "" {
		args = append(args, sinceTag+"..HEAD")
	}

	out, err := exec.Command("git", args...).Output()
	if err != nil {
		// No commits yet or invalid range — return empty
		return nil, nil
	}

	return parseGitLog(strings.TrimSpace(string(out))), nil
}

// LatestTag returns the most recent tag reachable from HEAD.
func LatestTag() (string, error) {
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		return "", fmt.Errorf("getting latest tag: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func parseGitLog(output string) []GitCommit {
	if output == "" {
		return nil
	}
	lines := strings.Split(output, "\n")
	commits := make([]GitCommit, 0, len(lines))
	for _, line := range lines {
		if hash, subject, ok := strings.Cut(line, " "); ok {
			commits = append(commits, GitCommit{Hash: hash, Subject: subject})
		}
	}
	return commits
}

func normalizeGitURL(url string) string {
	// SSH → HTTPS: git@host:org/repo.git → https://host/org/repo
	if strings.HasPrefix(url, "git@") {
		url = strings.TrimPrefix(url, "git@")
		url = strings.Replace(url, ":", "/", 1)
		url = "https://" + url
	}
	url = strings.TrimSuffix(url, ".git")
	return url
}
