# gsd — Get Stuff Done with AI agents

![Version](https://img.shields.io/badge/version-0.2.0-blue) ![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white) ![License](https://img.shields.io/github/license/jon-the-dev/gsd) ![Issues](https://img.shields.io/github/issues/jon-the-dev/gsd) ![Last Commit](https://img.shields.io/github/last-commit/jon-the-dev/gsd)

`gsd` is a CLI that wraps AI coding agents (Claude, Codex, Kiro) to autonomously
resolve GitHub issues. Point it at issue numbers — or just tell it how many to
pick up — and it builds a contextual prompt, prepares an isolated git worktree,
and launches the agent to do the work.

## How it works

For every run, `gsd`:

1. **Detects the target repo** from your `origin` remote (or `GSD_ISSUE_REPO`).
2. **Builds a contextual prompt** tailored to the mode you invoked (explicit
   issues, goal-based, label-based, or priority auto-select).
3. **Names a git worktree/branch** derived from the issue numbers or the date.
4. **Syncs the repo** (`git pull` when an upstream is tracked) and launches the
   selected agent with that prompt.

The agent then pulls issue details with the `gh` CLI, implements the change,
runs the project's tests, and opens a PR linked to the issue.

## Requirements

- [Go](https://go.dev/) 1.26+ (to build from source)
- `git` and the [GitHub CLI](https://cli.github.com/) (`gh`), authenticated
- At least one agent CLI on your `PATH`:
  - [`claude`](https://github.com/anthropics/claude-code) (default)
  - [`codex`](https://github.com/openai/codex)
  - `kiro-cli`

## Install

```bash
# Build and install to ~/.local/bin
make install

# …or install directly with Go
go install github.com/jon-the-dev/gsd@latest

# …or build the binary locally
make build        # produces ./gsd
```

Make sure `~/.local/bin` (or your `GOBIN`) is on your `PATH`.

## Quick start

```bash
gsd --setup        # one-time: install the Claude agents/skills gsd's prompts use
gsd "#123"         # fix issue #123 with Claude
gsd 5              # work the next 5 open issues until done
```

## Usage

```text
gsd [issues | count] [flags]
gsd term <branch-name> [flags]
gsd version
```

### Explicit issues

Pass one or more issue references. Commas, ranges, spaces, and `and` are all
accepted, and duplicates are de-duplicated.

```bash
gsd "#123"                 # single issue
gsd --codex "#123,#456"    # multiple issues with Codex
gsd "12-15"                # a range (12, 13, 14, 15)
gsd "1 and 7"              # mixed forms
```

### Goal-based (iterate until done)

Give a count and `gsd` runs in goal mode, working issues one by one until the
target is met.

```bash
gsd 5                      # the next 5 open issues
gsd --label security       # every open issue with the `security` label
gsd 3 --label bug          # the next 3 issues labeled `bug`
```

### Auto mode

Let `gsd` pick the work for you, priority-sorted by the P0–P7 label rubric
(blocker → critical/security → bug+high → … → chore/docs).

```bash
gsd --auto 5               # auto-select up to 5 non-conflicting issues
```

### Terminal mode

Create (or reuse) a worktree for a branch and drop the agent straight into it —
no issue required.

```bash
gsd term feature/my-branch
gsd term hotfix/login --codex
```

### Selecting an agent

Claude is the default. Override per-run with a flag or the `--agent` value, or
set a default with `GSD_AGENT`.

```bash
gsd --codex "#42"
gsd --kiro "#42"
gsd --agent codex "#42"
```

### Other flags

| Flag        | Description                                                      |
| ----------- | ---------------------------------------------------------------- |
| `--claude`  | Use the Claude agent (default).                                  |
| `--codex`   | Use the Codex agent.                                             |
| `--kiro`    | Use the Kiro agent.                                              |
| `--agent`   | Agent to use (`claude`, `codex`, `kiro`).                        |
| `--auto N`  | Auto-select `N` issues by priority.                              |
| `--label L` | Target open issues carrying label `L`.                           |
| `--yolo`    | Skip the agent's permission prompts (dangerous).                 |
| `--setup`   | Install the Claude agents/skills `gsd` depends on, then exit.    |

Before launching, `gsd` prints a briefing and a short countdown (skipped when
`GSD_CI` is set).

## Setup (`gsd --setup`)

`gsd`'s prompts reference a set of Claude subagents and skills. `--setup`
installs the bundled ones into your Claude config dir (`$CLAUDE_CONFIG_DIR`, or
`~/.claude`):

- **Agents:** `frontend-developer`, `backend-developer`, `cloud-architect`,
  `test-runner`
- **Skills:** `implement`, `merge`

Missing files are written; existing ones prompt `overwrite? [y/N]` per item
(auto-declined under `GSD_CI`). Two further skills the prompts use —
`next-issue` (jons-ai-toolkit) and `code-audit` (gstack) — are maintained
elsewhere, so setup only verifies their presence and warns if absent. The
bundled definitions live in [`internal/setup/assets/`](internal/setup/assets).

## Environment variables

| Variable                    | Purpose                                                        |
| --------------------------- | -------------------------------------------------------------- |
| `GSD_AGENT`                 | Default agent (`claude`, `codex`, `kiro`).                     |
| `GSD_TERM_AGENT`            | Agent override for `gsd term`.                                 |
| `GSD_ISSUE_REPO`            | Explicit `owner/repo` (otherwise auto-detected from `origin`). |
| `GSD_CI`                    | Skip the countdown and add permission-bypass flags.            |
| `CLAUDE_CONFIG_DIR`         | Where `--setup` installs agents/skills (default `~/.claude`).  |
| `GSD_TERM_WORKTREE_ROOT`    | Worktree root for `gsd term`.                                  |
| `GSD_CODEX_WORKTREE_ROOT`   | Worktree root for Codex (default `.codex/worktrees`).          |
| `GSD_KIRO_WORKTREE_ROOT`    | Worktree root for Kiro (default `.kiro/worktrees`).            |

## Project layout

```text
main.go              # Entry point
cmd/                 # Cobra commands (root, run, term, version)
internal/
  agent/             # Agent launch logic (claude, codex, kiro)
  git/               # Repo detection, worktree management, branch slugs
  issue/             # Issue-number parsing from flexible input
  prompt/            # Prompt template generation
  setup/             # `--setup` installer + bundled agents/skills
  ui/                # Terminal output (banner, briefing, countdown, colors)
```

## Development

```bash
make build       # build the binary
make test        # run the test suite
make lint        # golangci-lint
make install     # install to ~/.local/bin
make clean       # remove the built binary
```

## Contributing

1. Open (or pick up) an issue describing the change.
2. Branch, implement, and keep the tests green (`make test`).
3. Open a PR that links the issue.

See [`CHANGELOG.md`](CHANGELOG.md) for notable changes.

## License

Released under the [GNU Affero General Public License v3.0](LICENSE).
