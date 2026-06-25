package ui

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
)

func PrintBanner() {
	fmt.Printf(`%s%s
  ██████  ███████ ██████
 ██       ██      ██   ██
 ██   ███ ███████ ██   ██
 ██    ██      ██ ██   ██
  ██████  ███████ ██████
%s`, colorGreen, colorBold, colorReset)
	fmt.Println()
	fmt.Printf("%s  Get Stuff Done — AI Agent Launcher%s\n\n", colorCyan, colorReset)
}

func PrintBriefing(agent string, issues []int, repo, worktree string) {
	issueStrs := make([]string, len(issues))
	for i, num := range issues {
		issueStrs[i] = fmt.Sprintf("#%d", num)
	}

	fmt.Printf("%s┌─────────────────────────────────────────┐%s\n", colorCyan, colorReset)
	fmt.Printf("%s│  MISSION BRIEFING                       │%s\n", colorCyan, colorReset)
	fmt.Printf("%s├─────────────────────────────────────────┤%s\n", colorCyan, colorReset)
	fmt.Printf("%s│%s  Agent:     %-28s%s│%s\n", colorCyan, colorReset, agent, colorCyan, colorReset)
	fmt.Printf("%s│%s  Issues:    %-28s%s│%s\n", colorCyan, colorReset, strings.Join(issueStrs, ", "), colorCyan, colorReset)
	fmt.Printf("%s│%s  Repo:      %-28s%s│%s\n", colorCyan, colorReset, repo, colorCyan, colorReset)
	fmt.Printf("%s│%s  Worktree:  %-28s%s│%s\n", colorCyan, colorReset, worktree, colorCyan, colorReset)
	fmt.Printf("%s└─────────────────────────────────────────┘%s\n\n", colorCyan, colorReset)
}

func PrintAutoBriefing(agent, repo string, count int) {
	fmt.Printf("%s┌─────────────────────────────────────────┐%s\n", colorCyan, colorReset)
	fmt.Printf("%s│  MISSION BRIEFING (AUTO)                │%s\n", colorCyan, colorReset)
	fmt.Printf("%s├─────────────────────────────────────────┤%s\n", colorCyan, colorReset)
	fmt.Printf("%s│%s  Agent:     %-28s%s│%s\n", colorCyan, colorReset, agent, colorCyan, colorReset)
	fmt.Printf("%s│%s  Mode:      auto (top %d issues)        %s│%s\n", colorCyan, colorReset, count, colorCyan, colorReset)
	fmt.Printf("%s│%s  Repo:      %-28s%s│%s\n", colorCyan, colorReset, repo, colorCyan, colorReset)
	fmt.Printf("%s└─────────────────────────────────────────┘%s\n\n", colorCyan, colorReset)
}

func PrintGoalBriefing(agent, repo string, count int, label, worktree string) {
	mode := "goal"
	if count > 0 && label != "" {
		mode = fmt.Sprintf("next %d × label:%s", count, label)
	} else if label != "" {
		mode = fmt.Sprintf("all label:%s", label)
	} else if count > 0 {
		mode = fmt.Sprintf("next %d issues", count)
	}

	fmt.Printf("%s┌─────────────────────────────────────────┐%s\n", colorCyan, colorReset)
	fmt.Printf("%s│  MISSION BRIEFING (GOAL)                │%s\n", colorCyan, colorReset)
	fmt.Printf("%s├─────────────────────────────────────────┤%s\n", colorCyan, colorReset)
	fmt.Printf("%s│%s  Agent:     %-28s%s│%s\n", colorCyan, colorReset, agent, colorCyan, colorReset)
	fmt.Printf("%s│%s  Mode:      %-28s%s│%s\n", colorCyan, colorReset, mode, colorCyan, colorReset)
	fmt.Printf("%s│%s  Repo:      %-28s%s│%s\n", colorCyan, colorReset, repo, colorCyan, colorReset)
	fmt.Printf("%s│%s  Worktree:  %-28s%s│%s\n", colorCyan, colorReset, worktree, colorCyan, colorReset)
	fmt.Printf("%s└─────────────────────────────────────────┘%s\n\n", colorCyan, colorReset)
}

func PrintTermBanner(agent, slug, path string) {
	fmt.Printf("%s┌─────────────────────────────────────────┐%s\n", colorCyan, colorReset)
	fmt.Printf("%s│  GSD-TERM                               │%s\n", colorCyan, colorReset)
	fmt.Printf("%s├─────────────────────────────────────────┤%s\n", colorCyan, colorReset)
	fmt.Printf("%s│%s  Agent:     %-28s%s│%s\n", colorCyan, colorReset, agent, colorCyan, colorReset)
	fmt.Printf("%s│%s  Branch:    %-28s%s│%s\n", colorCyan, colorReset, slug, colorCyan, colorReset)
	fmt.Printf("%s│%s  Worktree:  %-28s%s│%s\n", colorCyan, colorReset, path, colorCyan, colorReset)
	fmt.Printf("%s└─────────────────────────────────────────┘%s\n\n", colorCyan, colorReset)
}

func Countdown(seconds int) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	fmt.Printf("%sLaunching in %d seconds (Ctrl+C to abort)...%s\n", colorYellow, seconds, colorReset)

	for i := seconds; i > 0; i-- {
		select {
		case <-sigCh:
			fmt.Printf("\n%sAborted.%s\n", colorRed, colorReset)
			return fmt.Errorf("aborted by user")
		case <-time.After(1 * time.Second):
			fmt.Printf("\r%s%d...%s  ", colorYellow, i-1, colorReset)
		}
	}
	fmt.Println()
	return nil
}

func Info(msg string) {
	fmt.Printf("%s▸ %s%s\n", colorCyan, msg, colorReset)
}

func Success(msg string) {
	fmt.Printf("%s✓ %s%s\n", colorGreen, msg, colorReset)
}

func Warn(msg string) {
	fmt.Printf("%s⚠ %s%s\n", colorYellow, msg, colorReset)
}

func Error(msg string) {
	fmt.Printf("%s✗ %s%s\n", colorRed, msg, colorReset)
}
