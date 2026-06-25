package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/jon-the-dev/gsd/internal/agent"
	gitutil "github.com/jon-the-dev/gsd/internal/git"
	"github.com/jon-the-dev/gsd/internal/issue"
	"github.com/jon-the-dev/gsd/internal/prompt"
	"github.com/jon-the-dev/gsd/internal/ui"
	"github.com/spf13/cobra"
)

func runGSD(cmd *cobra.Command, args []string) error {
	selectedAgent, err := resolveAgent(cmd)
	if err != nil {
		return err
	}

	repo, err := gitutil.DetectRepo()
	if err != nil {
		return fmt.Errorf("failed to detect GitHub repo: %w", err)
	}
	if envRepo := os.Getenv("GSD_ISSUE_REPO"); envRepo != "" {
		repo = envRepo
	}

	ui.PrintBanner()

	if autoFlag > 0 {
		return runAutoMode(selectedAgent, repo, autoFlag)
	}

	count, isCount := parseCountArg(args)
	if isCount || labelFlag != "" {
		return runGoalMode(selectedAgent, repo, count, labelFlag)
	}

	if len(args) == 0 {
		return fmt.Errorf(
			"no target specified\n\n" +
				"Usage:\n" +
				"  gsd \"#123\"              Fix specific issues\n" +
				"  gsd 5                   Next 5 open issues\n" +
				"  gsd --label security    All issues with label\n" +
				"  gsd 3 --label bug       Next 3 with label\n" +
				"  gsd --auto 5            Priority-sorted auto-select",
		)
	}

	issues, err := issue.Parse(args)
	if err != nil {
		return fmt.Errorf("failed to parse issues: %w", err)
	}

	worktree := gitutil.WorktreeName(issues)
	p := prompt.BuildStandard(issues, repo, worktree)

	ui.PrintBriefing(selectedAgent, issues, repo, worktree)

	if !ciMode {
		if err := ui.Countdown(10); err != nil {
			return err
		}
	}

	if err := gitPull(); err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	return agent.Launch(selectedAgent, p, worktree, ciMode)
}

func runGoalMode(selectedAgent, repo string, count int, label string) error {
	worktree := gitutil.SprintWorktreeName()
	goal := prompt.GoalSentence(count, label)
	p := prompt.BuildGoal(repo, worktree, goal)

	ui.PrintGoalBriefing(selectedAgent, repo, count, label, worktree)

	if !ciMode {
		if err := ui.Countdown(10); err != nil {
			return err
		}
	}

	if err := gitPull(); err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	return agent.Launch(selectedAgent, p, worktree, ciMode)
}

func runAutoMode(selectedAgent, repo string, count int) error {
	worktree := "fix-auto-batch"
	p := prompt.BuildAuto(repo, count)

	ui.PrintAutoBriefing(selectedAgent, repo, count)

	if !ciMode {
		if err := ui.Countdown(10); err != nil {
			return err
		}
	}

	if err := gitPull(); err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	return agent.Launch(selectedAgent, p, worktree, ciMode)
}

func parseCountArg(args []string) (int, bool) {
	if len(args) != 1 {
		return 0, false
	}
	n, err := strconv.Atoi(args[0])
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}

func resolveAgent(cmd *cobra.Command) (string, error) {
	if agentFlag != "" {
		return agent.Validate(agentFlag)
	}

	claude, _ := cmd.Flags().GetBool("claude")
	codex, _ := cmd.Flags().GetBool("codex")
	kiro, _ := cmd.Flags().GetBool("kiro")

	switch {
	case claude:
		return "claude", nil
	case codex:
		return "codex", nil
	case kiro:
		return "kiro", nil
	}

	if env := os.Getenv("GSD_AGENT"); env != "" {
		return agent.Validate(env)
	}

	return "claude", nil
}

func gitPull() error {
	out, err := exec.Command(
		"git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}",
	).Output()
	if err != nil || len(out) == 0 {
		ui.Warn("No upstream tracking branch — skipping pull")
		return nil
	}

	ui.Info("Syncing repository...")
	pullCmd := exec.Command("git", "pull")
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr
	if err := pullCmd.Run(); err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	return nil
}
