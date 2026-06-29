---
name: implement
description: Resolve a single GitHub issue end to end — pull issue details, implement the change with developer subagents, run tests, then commit, push, and open a linked PR. Use when asked to implement, fix, or address a specific issue number.
---

# Implement a GitHub Issue

Work on the issue number(s) provided as arguments.

- Use the `gh` CLI to pull issue details and set the `wip` label (create it if missing).
- Scan the codebase to understand what the project does before making changes.
- Implement the requested change, feature, or bugfix using the appropriate subagents:
  - `frontend-developer` — React and frontend work
  - `backend-developer` — FastAPI and backend logic
  - `cloud-architect` — AWS, Terraform, or DevOps activities
  - `test-runner` — Python, JS, TS, TF, or other test execution
- Use a subagent to:
  - Run any tests from the Makefile that are appropriate.
  - If tests fail for reasons outside the scope of this work, either fix them in this run or open a GitHub issue so they're addressed next round.
  - Run pre-commit before committing to ensure the tree is clean before pushing.
- Commit, push, and open a PR. Ensure the original issue is linked to the PR and closed by it.
- If CI runner budget is exhausted, and there are no conflicts and local tests pass, merge in.
