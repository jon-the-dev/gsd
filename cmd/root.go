package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.2.0"

var (
	agentFlag string
	autoFlag  int
	labelFlag string
	yoloFlag  bool
	ciMode    bool
)

var rootCmd = &cobra.Command{
	Use:   "gsd [issues | count]",
	Short: "GSD — Get Stuff Done with AI agents",
	Long: `A CLI that wraps AI coding agents (Claude, Codex, Kiro) to autonomously resolve GitHub issues.

Modes:
  gsd "#123,#456"          Explicit issues
  gsd 5                    Next 5 open issues (goal-based)
  gsd --label security     All open issues with label (goal-based)
  gsd 3 --label bug        Next 3 issues with label (goal-based)
  gsd --auto 5             Priority-sorted auto-select`,
	Args: cobra.ArbitraryArgs,
	RunE: runGSD,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&agentFlag, "agent", "", "Agent to use (claude, codex, kiro)")
	rootCmd.Flags().BoolP("claude", "", false, "Use Claude agent")
	rootCmd.Flags().BoolP("codex", "", false, "Use Codex agent")
	rootCmd.Flags().BoolP("kiro", "", false, "Use Kiro agent")
	rootCmd.Flags().IntVar(&autoFlag, "auto", 0, "Auto-select N issues by priority")
	rootCmd.Flags().StringVar(&labelFlag, "label", "", "Target issues with this GitHub label")
	rootCmd.Flags().BoolVar(&yoloFlag, "yolo", false, "Skip permission prompts (dangerous)")
	rootCmd.AddCommand(termCmd)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gsd %s\n", version)
	},
}

func Execute() error {
	ciMode = os.Getenv("GSD_CI") != ""
	return rootCmd.Execute()
}
