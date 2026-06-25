package cmd

import (
	"fmt"
	"os"

	"github.com/jon-the-dev/gsd/internal/agent"
	gitutil "github.com/jon-the-dev/gsd/internal/git"
	"github.com/jon-the-dev/gsd/internal/ui"
	"github.com/spf13/cobra"
)

var termCmd = &cobra.Command{
	Use:   "term [branch-name]",
	Short: "Create or reuse a git worktree and launch an agent",
	Args:  cobra.ExactArgs(1),
	RunE:  runTerm,
}

func runTerm(cmd *cobra.Command, args []string) error {
	selectedAgent, err := resolveTermAgent(cmd)
	if err != nil {
		return err
	}

	branchName := args[0]
	slug := gitutil.BranchSlug(branchName)

	root := termWorktreeRoot(selectedAgent)
	wtPath, err := gitutil.PrepareWorktree(root, slug)
	if err != nil {
		return fmt.Errorf("worktree preparation failed: %w", err)
	}

	ui.PrintTermBanner(selectedAgent, slug, wtPath)

	if !ciMode {
		if err := ui.Countdown(5); err != nil {
			return err
		}
	}

	return agent.LaunchInDir(selectedAgent, wtPath, ciMode)
}

func resolveTermAgent(cmd *cobra.Command) (string, error) {
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

	if env := os.Getenv("GSD_TERM_AGENT"); env != "" {
		return agent.Validate(env)
	}
	if env := os.Getenv("GSD_AGENT"); env != "" {
		return agent.Validate(env)
	}

	return "claude", nil
}

func termWorktreeRoot(selectedAgent string) string {
	if env := os.Getenv("GSD_TERM_WORKTREE_ROOT"); env != "" {
		return env
	}

	switch selectedAgent {
	case "codex":
		if env := os.Getenv("GSD_CODEX_WORKTREE_ROOT"); env != "" {
			return env
		}
		return ".codex/worktrees"
	case "kiro":
		if env := os.Getenv("GSD_KIRO_WORKTREE_ROOT"); env != "" {
			return env
		}
		return ".kiro/worktrees"
	default:
		return ".gsd-term/worktrees"
	}
}
