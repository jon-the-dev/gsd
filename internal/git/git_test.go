package git

import "testing"

func TestBranchSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"feature/my-branch", "feature-my-branch"},
		{"Feature/MY_Branch", "feature-my-branch"},
		{"fix #123: broken stuff", "fix-123-broken-stuff"},
		{"simple", "simple"},
		{"--leading-trailing--", "leading-trailing"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := BranchSlug(tt.input)
			if got != tt.want {
				t.Errorf("BranchSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestWorktreeName(t *testing.T) {
	tests := []struct {
		issues []int
		want   string
	}{
		{[]int{123}, "fix-issue-123"},
		{[]int{1, 2, 3}, "fix-issue-1-2-3"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := WorktreeName(tt.issues)
			if got != tt.want {
				t.Errorf("WorktreeName(%v) = %q, want %q", tt.issues, got, tt.want)
			}
		})
	}
}
