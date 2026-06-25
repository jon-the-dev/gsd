# gsd

A CLI that wraps AI coding agents (Claude, Codex, Kiro) to autonomously resolve GitHub issues. Point it at issue numbers and it builds a contextual prompt, prepares a git worktree, and launches the agent to fix them.

## Build & Test

```bash
make build       # Build binary
make test        # Run tests
make lint        # golangci-lint
make install     # Install to ~/.local/bin
```

## Architecture

```
main.go              # Entry point
cmd/                 # Cobra commands (root, term, version)
internal/
  agent/             # Agent launch logic (claude, codex, kiro)
  git/               # Repo detection, worktree management, branch slugs
  issue/             # Issue number parsing from flexible input
  prompt/            # Prompt template generation
  ui/                # Terminal output (banner, briefing, countdown, colors)
```

## Usage

```bash
# Explicit issues
gsd "#123"                    # Fix issue with Claude (default)
gsd --codex "#123,#456"      # Fix issues with Codex

# Goal-based (iterates until done)
gsd 5                         # Next 5 open issues
gsd --label security          # All open issues with label
gsd 3 --label bug             # Next 3 issues with label

# Other modes
gsd --auto 5                  # Priority-sorted auto-select (P0-P7 rubric)
gsd term feature/my-branch    # Launch agent in a worktree
```

## Environment Variables

- `GSD_AGENT` — default agent (claude/codex/kiro)
- `GSD_TERM_AGENT` — override for `gsd term`
- `GSD_ISSUE_REPO` — explicit owner/repo (auto-detected from git remote)
- `GSD_CODEX_WORKTREE_ROOT` / `GSD_KIRO_WORKTREE_ROOT` / `GSD_TERM_WORKTREE_ROOT`
- `GSD_CI` — skip countdown, add permission-bypass flags
