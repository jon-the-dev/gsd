package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var repoRegex = regexp.MustCompile(
	`github\.com[:/]([^/]+/[^/.]+?)(?:\.git)?$`,
)

func DetectRepo() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", fmt.Errorf("no git remote 'origin' found")
	}

	url := strings.TrimSpace(string(out))
	matches := repoRegex.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not parse GitHub repo from remote: %s", url)
	}

	return matches[1], nil
}

func SprintWorktreeName() string {
	return fmt.Sprintf("mini-sprint-%s", time.Now().Format("20060102"))
}

func WorktreeName(issues []int) string {
	if len(issues) == 1 {
		return fmt.Sprintf("fix-issue-%d", issues[0])
	}
	parts := make([]string, len(issues))
	for i, num := range issues {
		parts[i] = fmt.Sprintf("%d", num)
	}
	return "fix-issue-" + strings.Join(parts, "-")
}

var slugRegex = regexp.MustCompile(`[^a-z0-9-]+`)

func BranchSlug(name string) string {
	lower := strings.ToLower(name)
	slug := slugRegex.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func PrepareWorktree(root, slug string) (string, error) {
	wtPath := filepath.Join(root, slug)

	gitDir := filepath.Join(wtPath, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		branch, err := currentBranch(wtPath)
		if err == nil && BranchSlug(branch) == slug {
			return wtPath, nil
		}
		return wtPath, nil
	}

	if err := os.MkdirAll(root, 0755); err != nil {
		return "", fmt.Errorf("failed to create worktree root: %w", err)
	}

	if branchExists(slug) {
		err := exec.Command("git", "worktree", "add", wtPath, slug).Run()
		if err != nil {
			return "", fmt.Errorf("git worktree add failed: %w", err)
		}
		return wtPath, nil
	}

	if remoteBranchExists(slug) {
		err := exec.Command(
			"git", "worktree", "add", wtPath, "-b", slug, "origin/"+slug,
		).Run()
		if err != nil {
			return "", fmt.Errorf("git worktree add (remote) failed: %w", err)
		}
		return wtPath, nil
	}

	err := exec.Command("git", "worktree", "add", "-b", slug, wtPath).Run()
	if err != nil {
		return "", fmt.Errorf("git worktree add (new branch) failed: %w", err)
	}
	return wtPath, nil
}

func currentBranch(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func branchExists(name string) bool {
	err := exec.Command("git", "rev-parse", "--verify", name).Run()
	return err == nil
}

func remoteBranchExists(name string) bool {
	err := exec.Command(
		"git", "rev-parse", "--verify", "origin/"+name,
	).Run()
	return err == nil
}
