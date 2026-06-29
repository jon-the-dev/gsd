// Package setup installs the Claude subagents and skills that gsd's prompt
// templates depend on into the user's Claude config dir, prompting before it
// overwrites anything that already exists.
package setup

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jon-the-dev/gsd/internal/ui"
)

//go:embed all:assets
var assets embed.FS

// externalSkills are referenced by gsd's prompts but maintained elsewhere, so
// setup only verifies their presence rather than bundling a copy that drifts.
var externalSkills = map[string]string{
	"next-issue": "jons-ai-toolkit marketplace",
	"code-audit": "gstack",
}

// Run writes bundled agents/skills into the Claude config dir and reports on
// the external skills it can only verify. reader supplies overwrite answers
// (os.Stdin in normal use); auto is true when prompts should be skipped and
// existing files left untouched (CI / non-interactive).
func Run(reader io.Reader, auto bool) error {
	root, err := claudeRoot()
	if err != nil {
		return err
	}
	ui.Info("Claude config dir: " + root)

	c := &confirmer{in: bufio.NewReader(reader), auto: auto}
	if err := installTree(root, c); err != nil {
		return err
	}

	verifyExternal(root)
	ui.Success("Setup complete — restart any running claude session to load changes")
	return nil
}

// installTree walks the embedded assets and mirrors them under root, treating
// "assets/agents/x.md" -> "<root>/agents/x.md" and likewise for skills.
func installTree(root string, c *confirmer) error {
	return fs.WalkDir(assets, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel := strings.TrimPrefix(path, "assets/")
		dest := filepath.Join(root, rel)
		return installFile(path, dest, label(rel), c)
	})
}

func installFile(srcPath, dest, label string, c *confirmer) error {
	if _, err := os.Stat(dest); err == nil {
		if !c.confirm(fmt.Sprintf("%s already exists at %s — overwrite?", label, dest)) {
			ui.Info("Kept existing " + label)
			return nil
		}
	}

	data, err := assets.ReadFile(srcPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", dest, err)
	}
	ui.Success("Installed " + label)
	return nil
}

func verifyExternal(root string) {
	for skill, source := range externalSkills {
		dir := filepath.Join(root, "skills", skill)
		if _, err := os.Stat(dir); err == nil {
			ui.Success("/" + skill + " skill present")
			continue
		}
		ui.Warn(fmt.Sprintf(
			"/%s skill not found — it's maintained by %s, not bundled with gsd; install it separately",
			skill, source,
		))
	}
}

// label turns an embedded path into a human reference, e.g.
// "agents/test-runner.md" -> "agent test-runner",
// "skills/merge/SKILL.md" -> "skill merge".
func label(rel string) string {
	parts := strings.Split(rel, "/")
	switch parts[0] {
	case "agents":
		return "agent " + strings.TrimSuffix(parts[len(parts)-1], ".md")
	case "skills":
		return "skill " + parts[1]
	default:
		return rel
	}
}

// claudeRoot resolves the Claude config dir, honoring CLAUDE_CONFIG_DIR.
func claudeRoot() (string, error) {
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".claude"), nil
}

type confirmer struct {
	in   *bufio.Reader
	auto bool
}

// confirm asks a yes/no question defaulting to no. In auto mode it answers no
// (never clobbers existing files unattended) and records the skip.
func (c *confirmer) confirm(question string) bool {
	if c.auto {
		ui.Warn(question + " [auto: no]")
		return false
	}
	fmt.Printf("%s [y/N]: ", question)
	line, err := c.in.ReadString('\n')
	if err != nil && line == "" {
		return false
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes"
}
