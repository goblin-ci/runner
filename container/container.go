// Package container package provides docker container
// build / run functionality and executes build stacks
package container

import (
	"io"
)

// Container interface provides common
// functionality for container engines
type ContainerWriter interface {
	io.Writer

	// Observe listens for container stdout
	// and sends it to string channel
	Observe(w io.Writer)
	// Run starts up the container
	// runs build commands, kills and removes
	// container upon completion
	Run()
}
