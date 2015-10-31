package main

import (
	"fmt"
	"sync"

	"github.com/goblin-ci/runner/container"
	"github.com/goblin-ci/runner/stack"
)

var wg sync.WaitGroup

func main() {
	// New webhook triggered, got repo url etc... from MQ
	// Instantiate golang stack
	goStack := stack.NewGolang("latest")

	// Create new contaner from stack
	cnt := container.New(goStack, "")

	// Run the container
	cnt.WG.Add(2)
	go cnt.Run()
	go cnt.Observe()

	fmt.Println("Waiting for goroutines to finish")
	cnt.WG.Wait()
}
