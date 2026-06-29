---
name: merge
description: Watch a pull request's CI checks and merge it once they pass. Use after opening or pushing a PR when the goal is to land it automatically on green CI. Triggers on "merge the PR", "land it when CI passes", or "auto-merge".
---

# Merge on Green CI

Watch the pull request's CI/CD checks and merge it once they pass.

1. Identify the PR for the current branch:

   ```bash
   gh pr view --json number,headRefName,mergeable,statusCheckRollup
   ```

2. Watch the checks until they complete:

   ```bash
   gh pr checks --watch
   ```

3. Decide based on the outcome:
   - **All checks pass** — merge with a short comment noting CI passed, then delete the branch:

     ```bash
     gh pr merge --squash --delete-branch \
       --subject "Merge: <PR title>" \
       --body "CI passed — merging."
     ```

   - **Checks fail** — do not merge. Summarize which checks failed and hand back to the appropriate developer subagent to fix.

4. If GitHub Actions is unavailable or the runner budget is exhausted, fall back to: confirm there are no merge conflicts (`mergeable: MERGEABLE`) and local tests pass, then merge with a comment explaining CI was skipped.

Never force-merge over failing required checks unless explicitly told to bypass them.
