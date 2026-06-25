package agent

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var validAgents = []string{"claude", "codex", "kiro"}

func Validate(name string) (string, error) {
	lower := strings.ToLower(strings.TrimSpace(name))
	for _, a := range validAgents {
		if lower == a {
			return a, nil
		}
	}
	return "", fmt.Errorf(
		"unknown agent %q (valid: %s)",
		name, strings.Join(validAgents, ", "),
	)
}

func Launch(name, prompt, worktree string, yolo bool) error {
	switch name {
	case "claude":
		return launchClaude(prompt, worktree, yolo)
	case "codex":
		return launchCodex(prompt, worktree, yolo)
	case "kiro":
		return launchKiro(prompt, worktree, yolo)
	default:
		return fmt.Errorf("unsupported agent: %s", name)
	}
}

func LaunchInDir(name, dir string, yolo bool) error {
	switch name {
	case "claude":
		return launchClaudeInDir(dir, yolo)
	case "codex":
		return launchCodexInDir(dir, yolo)
	case "kiro":
		return launchKiroInDir(dir, yolo)
	default:
		return fmt.Errorf("unsupported agent: %s", name)
	}
}

func launchClaude(prompt, worktree string, yolo bool) error {
	args := []string{"-w", worktree, prompt}
	if yolo {
		args = append([]string{"--dangerously-skip-permissions"}, args...)
	}
	return run("claude", args, "")
}

func launchClaudeInDir(dir string, yolo bool) error {
	args := []string{}
	if yolo {
		args = append(args, "--dangerously-skip-permissions")
	}
	return run("claude", args, dir)
}

func launchCodex(prompt, worktree string, yolo bool) error {
	root := os.Getenv("GSD_CODEX_WORKTREE_ROOT")
	if root == "" {
		root = ".codex/worktrees"
	}
	dir := root + "/" + worktree

	args := []string{prompt}
	if yolo {
		args = append([]string{"--dangerously-bypass-approvals-and-sandbox"}, args...)
	}
	return run("codex", args, dir)
}

func launchCodexInDir(dir string, yolo bool) error {
	args := []string{}
	if yolo {
		args = append(args, "--dangerously-bypass-approvals-and-sandbox")
	}
	return run("codex", args, dir)
}

func launchKiro(prompt, worktree string, yolo bool) error {
	root := os.Getenv("GSD_KIRO_WORKTREE_ROOT")
	if root == "" {
		root = ".kiro/worktrees"
	}
	dir := root + "/" + worktree

	args := []string{"chat", prompt}
	if yolo {
		args = append([]string{"-a"}, args...)
	}
	return run("kiro-cli", args, dir)
}

func launchKiroInDir(dir string, yolo bool) error {
	args := []string{"chat"}
	if yolo {
		args = append([]string{"-a"}, args...)
	}
	return run("kiro-cli", args, dir)
}

func run(binary string, args []string, dir string) error {
	cmd := exec.Command(binary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Run()
}
