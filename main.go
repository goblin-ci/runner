package main

import (
	"fmt"
	"sync"

	"github.com/goblin-ci/runner/container"
)

var wg sync.WaitGroup

func main() {
	cnt := container.New("ubuntu")

	cnt.WG.Add(2)
	go cnt.Run()
	go cnt.Observe()

	fmt.Println("Waiting for goroutines to finish")
	cnt.WG.Wait()
}
