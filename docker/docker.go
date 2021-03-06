// Package docker provides docker utility functions
package docker

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/goblin-ci/runner/github"
	"github.com/goblin-ci/runner/stack"
)

// Docker represents docker container
// strucure used for building proces
type Docker struct {
	ID         string
	Stack      stack.Stack
	Stream     chan string
	Done       chan bool
	WG         *sync.WaitGroup
	GithubPush *github.Push
}

// runDetached runs container in detached mode
// Returns container ID, or an error
func (c *Docker) runDetached() ([]byte, error) {
	cmd := exec.Command("docker", []string{"run", "-d", c.Stack.ImageName()}...)
	cmd.Stderr = os.Stderr
	r, err := cmd.Output()
	if len(r) > 1 {
		return r[:16], nil
	}
	log.Println(string(r))
	return nil, err
}

// run runs specific docker container
// Returns container ID or an error
func (c *Docker) run() ([]byte, error) {
	cmd := exec.Command("docker", []string{"run", c.Stack.ImageName()}...)
	cmd.Stderr = os.Stderr
	r, err := cmd.Output()
	log.Println(string(r))
	return r, err
}

// execInteractive executes command inside of docker container
// with inearctive output sent to provided io.Writer
func (c *Docker) execInteractive(args []string) error {
	cm := append([]string{"exec", c.ID}, args...)
	if strings.ToLower(args[0]) == "cd" {
		cm = append([]string{"exec", c.ID, "bash", "-c"}, "'"+strings.Join(args, " ")+"'")
	}
	cmd := exec.Command("docker", cm...)
	cmd.Stdout = c
	cmd.Stderr = c
	return cmd.Run()
}

// stop stops docker container with provided ID
func (c *Docker) stop() error {
	cmd := exec.Command("docker", []string{"kill", c.ID}...)
	log.Printf("Killing container %s", c.ID)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Docker) remove() error {
	cmd := exec.Command("docker", []string{"rm", "-f", c.ID}...)
	log.Printf("Removing container %s", c.ID)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Write makes Docker conform to io.Writer
func (c *Docker) Write(p []byte) (n int, err error) {
	c.Stream <- string(p)
	return len(p), nil
}

// Observe listens for container stdout
// and sends it to string channel
func (c *Docker) Observe(w io.Writer) {
	defer func() {
		err := c.stop()
		if err != nil {
			log.Println(err)
		}
		err = c.remove()
		if err != nil {
			// TODO Save to dangling container table
			log.Println(err)
		}
		c.WG.Done()
	}()

	for {
		select {
		case message := <-c.Stream:
			fmt.Fprintf(w, message)
		case <-time.After(time.Minute * 10):
			log.Println("Observer timeout")
			return
		case <-c.Done:
			log.Println("Build complete.")
			return
		}
	}
}

// Run starts up the container and
// runs build commands
func (c *Docker) Run() {
	defer func() {
		c.Done <- true
		c.WG.Done()
	}()

	// Start up the container and get it's ID
	ID, err := c.runDetached()
	if err != nil {
		log.Println(err)
		log.Println("Exiting...")
		// TODO Close channel
		return
	}
	c.ID = string(ID)

	log.Println("New container up and running: ", c.ID)

	// TODO
	// Setup .ssh keys for private repos (priority low)
	// Clone the repo
	cloneCmd := c.GithubPush.CloneCmd()
	fmt.Fprintf(c, strings.Join(cloneCmd, " "))
	err = c.execInteractive(cloneCmd)
	if err != nil {
		log.Println(err)
		return
	}

	// Parse yml file
	// Determine go version and set proper ENV acordingly
	// Check for build commands and set them if any

	/*
		c.Stack.SetBuild([]string{
			"pwd",
			"ls -l",
			"cd /go/bin",
			"ls -1",
			"ls -l /go/src",
		})
	*/

	// Execute build commands and send data to stream
	// Close channel on build error
	build := c.Stack.GetBuild()
	if build == nil {
		build = c.Stack.DefaultBuild()
	}

	for _, cmd := range build {
		cmdSlice := strings.Split(cmd, " ")
		// TODO Send json event with new command
		fmt.Fprintf(c, ">> "+cmd)
		err = c.execInteractive(cmdSlice)
		if err != nil {
			fmt.Fprintf(c, "\nBUILD FAILED!")
			log.Println("BUILD FAILED!")
			return
		}
	}

	fmt.Fprintf(c, "\nBUILD SUCCESSFUL")
}

// New creates and intializes new container
func New(s stack.Stack, push *github.Push) *Docker {
	return &Docker{
		Stack:      s,
		Stream:     make(chan string),
		Done:       make(chan bool),
		WG:         &sync.WaitGroup{},
		GithubPush: push,
	}
}
