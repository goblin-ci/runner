// Container package provides docker container
// build / run functionality and executes build stacks
package container

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/goblin-ci/runner/stack"
)

// Container represents docker container
// strucure used for building proces
type Container struct {
	ID            string
	Stack         stack.Stack
	Stream        chan string
	BuildCommands []string
	WG            *sync.WaitGroup
}

// Write makes Container conform to io.Writer
func (c *Container) Write(p []byte) (n int, err error) {
	c.Stream <- string(p)
	return len(p), nil
}

// Observe listens for container stdout
// and sends it to string channel
func (c *Container) Observe() {
	defer c.WG.Done()
	select {
	case message := <-c.Stream:
		fmt.Println("Message received: ", message)
	// Container timeout
	case <-time.After(time.Second):
		fmt.Println("Observer timeout")
		return
	}
}

// Run starts up the container and
// runs build commands
func (c *Container) Run() {
	defer c.WG.Done()
	// Start up the container and get it's ID
	cmd := exec.Command("docker", "run", "-d", c.Stack.ImageName())
	result, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return
	}
	c.ID = string(result)

	fmt.Println("Container ID: ", c.ID)

	// Execute build commands and send data to stream
	if c.Stack.GetBuild() == nil {
		return
	}
}

// New creates and intializes new container
func New(s stack.Stack) *Container {
	return &Container{
		Stack:  s,
		Stream: make(chan string),
		WG:     &sync.WaitGroup{},
	}
}
