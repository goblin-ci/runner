// Package github provides utility functions
// for github cloning etc...
package github

import (
	"io"
	"os/exec"
)

// CloneRepo clones github repo
// provided with internal repo ID
// TODO tak ID as an argument instead
func CloneRepo(repo string, w io.Writer) (err error) {
	args := []string{
		"clone",
		"--depth 10",
		repo,
		"/go/src/github.com/app",
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = w
	cmd.Stderr = w

	return cmd.Run()
}
