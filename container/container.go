// Container package provides docker container
// build / run functionality and executes build stacks
package container

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/goblin-ci/runner/docker"
	"github.com/goblin-ci/runner/stack"
)

// Container represents docker container
// strucure used for building proces
type Container struct {
	ID            string
	Stack         stack.Stack
	Stream        chan string
	RepoURL       string
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
	ID, err := docker.RunDetached(c.Stack.ImageName())
	if err != nil {
		log.Println(err)
		return
	}
	c.ID = string(ID)

	fmt.Println("Container ID: ", c.ID)

	// TODO
	// Setup .ssh keys for private repos (priority low)
	// Clone the repo
	// Parse yml file
	// Determine go version and set proper ENV acordingly
	// Checko for build commands and set them if any

	// Execute build commands and send data to stream
	if c.Stack.GetBuild() == nil {
		return
	}
}

// New creates and intializes new container
func New(s stack.Stack, repoURL string) *Container {
	return &Container{
		Stack:   s,
		Stream:  make(chan string),
		WG:      &sync.WaitGroup{},
		RepoURL: repoURL,
	}
}
