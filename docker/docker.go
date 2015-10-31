// docker provides docker utility functions
package docker

import (
	"os/exec"
)

// Cmd returns exec.cmd docker command
func Cmd(args ...string) *exec.Cmd {
	return exec.Command("docker", args...)
}

// RunDetached runs container in detached mode
// Returns container ID, or an error
func RunDetached(arg string) ([]byte, error) {
	cmd := exec.Command("docker", append([]string{"run", "-d"}, arg)...)
	return cmd.Output()
}

// Run runs specific docker container
// Returns container ID or an error
func Run(arg string) ([]byte, error) {
	cmd := exec.Command("docker", append([]string{"run"}, arg)...)
	return cmd.Output()
}
