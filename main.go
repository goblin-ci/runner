package main

import (
	"fmt"
	"sync"

	"github.com/goblin-ci/runner/container"
	"github.com/goblin-ci/runner/stack"
)

var wg sync.WaitGroup

func main() {
	var goStack stack.Go
	cnt := container.New(&goStack)

	cnt.WG.Add(2)
	go cnt.Run()
	go cnt.Observe()

	fmt.Println("Waiting for goroutines to finish")
	cnt.WG.Wait()
}
