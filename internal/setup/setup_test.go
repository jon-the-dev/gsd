package setup

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func bufioReader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}

func TestLabel(t *testing.T) {
	cases := map[string]string{
		"agents/test-runner.md":     "agent test-runner",
		"agents/cloud-architect.md": "agent cloud-architect",
		"skills/merge/SKILL.md":     "skill merge",
		"skills/implement/SKILL.md": "skill implement",
	}
	for rel, want := range cases {
		if got := label(rel); got != want {
			t.Errorf("label(%q) = %q, want %q", rel, got, want)
		}
	}
}

func TestConfirmAutoNeverOverwrites(t *testing.T) {
	c := &confirmer{auto: true}
	if c.confirm("overwrite?") {
		t.Error("auto mode must answer no")
	}
}

func TestConfirmParsesAnswers(t *testing.T) {
	cases := map[string]bool{
		"y\n": true, "yes\n": true, "Y\n": true,
		"n\n": false, "\n": false, "anything\n": false,
	}
	for input, want := range cases {
		c := &confirmer{in: bufioReader(input)}
		if got := c.confirm("overwrite?"); got != want {
			t.Errorf("confirm(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestInstallTreeWritesMissingAndSkipsExisting(t *testing.T) {
	root := t.TempDir()

	// First run on an empty root installs everything bundled.
	c := &confirmer{auto: true}
	if err := installTree(root, c); err != nil {
		t.Fatalf("installTree: %v", err)
	}
	for _, rel := range []string{
		"agents/frontend-developer.md",
		"agents/test-runner.md",
		"skills/implement/SKILL.md",
		"skills/merge/SKILL.md",
	} {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Errorf("expected %s to be installed: %v", rel, err)
		}
	}

	// A pre-existing file must NOT be overwritten when the answer is no.
	target := filepath.Join(root, "agents", "test-runner.md")
	if err := os.WriteFile(target, []byte("CUSTOM"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := installTree(root, &confirmer{auto: true}); err != nil {
		t.Fatalf("installTree (second run): %v", err)
	}
	data, _ := os.ReadFile(target)
	if string(data) != "CUSTOM" {
		t.Errorf("auto mode overwrote an existing file: got %q", data)
	}

	// With an affirmative answer, the bundled content replaces it.
	if err := installTree(root, &confirmer{in: bufioReader(strings.Repeat("y\n", 8))}); err != nil {
		t.Fatalf("installTree (overwrite run): %v", err)
	}
	data, _ = os.ReadFile(target)
	if string(data) == "CUSTOM" {
		t.Error("expected bundled content to overwrite when answered yes")
	}
}
