package changelog

import (
	"os/exec"
	"strings"
)

type GitServcie struct {
	workingDir string
}

func NewGitService(workingDir string) GitServcie {
	return GitServcie{
		workingDir: workingDir,
	}
}

func (g GitServcie) CurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = g.workingDir

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
