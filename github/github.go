// Package github provides utility functions
// for github cloning etc...
package github

import (
	"fmt"
	"io"
	"os/exec"
)

// CloneRepo clones github repo
// provided with internal repo ID
// TODO take an ID as an argument instead
func CloneRepo(repo, branch string, w io.Writer) (err error) {
	args := []string{
		"clone",
		"--depth 10",
		fmt.Sprintf("--branch %s", branch),
		repo,
		"/go/src/goblin/app",
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = w
	cmd.Stderr = w

	return cmd.Run()
}
