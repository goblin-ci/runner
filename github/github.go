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
func CloneRepo(repo, branch, containerID string, w io.Writer) (err error) {
	args := []string{
		"exec",
		containerID,
		"git",
		"clone",
		"--depth",
		"10",
		"--branch",
		fmt.Sprintf("%s", branch),
		repo,
		"/go/src/goblin/app",
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = w
	cmd.Stderr = w

	return cmd.Run()
}
