// Package docker provides docker utility functions
package docker

import (
	"os/exec"
)

// Cmd returns exec.cmd docker command
func Cmd(args ...string) *exec.Cmd {
	return exec.Command("docker", args...)
}
