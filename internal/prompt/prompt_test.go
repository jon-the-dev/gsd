package prompt

import (
	"strings"
	"testing"
)

func TestGoalSentence(t *testing.T) {
	tests := []struct {
		name  string
		count int
		label string
		want  string
	}{
		{
			"count only",
			3, "",
			"The goal is complete when the next 3 issues are implemented and verified.",
		},
		{
			"label only",
			0, "security",
			"Target every open github issue labeled `security`. The goal is complete when all open issues labeled `security` are implemented and verified.",
		},
		{
			"count and label",
			2, "bug",
			"Target open github issues labeled `bug`. The goal is complete when the next 2 such issues are implemented and verified.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GoalSentence(tt.count, tt.label)
			if got != tt.want {
				t.Errorf("GoalSentence(%d, %q) =\n  %q\nwant:\n  %q", tt.count, tt.label, got, tt.want)
			}
		})
	}
}

func TestBuildGoalContainsGoal(t *testing.T) {
	goal := GoalSentence(5, "")
	p := BuildGoal("owner/repo", "mini-sprint-20260624", goal)

	if !strings.Contains(p, goal) {
		t.Error("BuildGoal output does not contain the goal sentence")
	}
	if !strings.Contains(p, "owner/repo") {
		t.Error("BuildGoal output does not contain repo")
	}
	if !strings.Contains(p, "/implement") {
		t.Error("BuildGoal output does not contain /implement skill reference")
	}
	if !strings.Contains(p, "/code-audit") {
		t.Error("BuildGoal output does not contain /code-audit skill reference")
	}
}
