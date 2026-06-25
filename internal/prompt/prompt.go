package prompt

import (
	"fmt"
	"strings"
)

func BuildStandard(issues []int, repo, worktree string) string {
	issueStr := formatIssues(issues)

	return fmt.Sprintf(`# Address GitHub Issue

Work on issue %s.

Target GitHub repository: %s
Use branch/worktree name: %s

## Instructions

- Use gh CLI to pull issue details and set wip label
- Scan the codebase to understand what the project does
- Implement the requested change, feature, or bugfix
- Use available subagents (frontend-developer, backend-developer, cloud-architect, test-runner)
- Run tests from Makefile
- Run pre-commit if configured
- Commit, push, and open a PR
- Link the original issue and ensure it's closed by the PR`, issueStr, repo, worktree)
}

func BuildGoal(repo, worktree, goal string) string {
	return fmt.Sprintf(`/goal Please review the open github issues and /implement them one by one.
Use /next-issue to identify the ideal next issue then use the /implement skill to implement, deploy, validate and /merge in the associated PRs.
Ensure we are using subagents listed in the /implement skill to minimize context bloat.

%s

Target GitHub repository: %s
Use branch/worktree name: %s

Use the /code-audit skill and ensure the code passes with a rating of 8 or greater.
Pass the code back to the developer subagents to fix.
If a web feature use playwright mcp server and tools. Tests passed etc.`, goal, repo, worktree)
}

func GoalSentence(count int, label string) string {
	switch {
	case count > 0 && label != "":
		return fmt.Sprintf(
			"Target open github issues labeled `%s`. "+
				"The goal is complete when the next %d such issues are implemented and verified.",
			label, count,
		)
	case label != "":
		return fmt.Sprintf(
			"Target every open github issue labeled `%s`. "+
				"The goal is complete when all open issues labeled `%s` are implemented and verified.",
			label, label,
		)
	default:
		return fmt.Sprintf(
			"The goal is complete when the next %d issues are implemented and verified.",
			count,
		)
	}
}

func BuildAuto(repo string, count int) string {
	return fmt.Sprintf(`# Auto-select and resolve GitHub issues

Target GitHub repository: %s
Select up to %d non-conflicting issues to resolve.

## Priority Rubric (highest first)

- P0: blocker
- P1: critical, security
- P2: bug + high
- P3: bug + medium
- P4: bug + low
- P5: enhancement + high
- P6: enhancement, feature
- P7: chore, tech-debt, documentation

## Instructions

1. Run: gh issue list --repo %s --state open --limit 50
2. Categorize issues by priority using the label rubric above
3. Select the top %d issues by priority
4. Verify selected issues don't conflict on shared files
5. For each selected issue:
   - Use gh CLI to pull issue details and set wip label
   - Scan the codebase to understand what the project does
   - Implement the requested change, feature, or bugfix
   - Use available subagents (frontend-developer, backend-developer, cloud-architect, test-runner)
   - Run tests from Makefile
   - Run pre-commit if configured
   - Commit, push, and open a PR
   - Link the original issue and ensure it's closed by the PR`, repo, count, repo, count)
}

func formatIssues(issues []int) string {
	parts := make([]string, len(issues))
	for i, num := range issues {
		parts[i] = fmt.Sprintf("#%d", num)
	}
	return strings.Join(parts, ", ")
}
