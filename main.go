package main

import (
	"fmt"
	"log"

	"github.com/goblin-ci/runner/container"
	"github.com/goblin-ci/runner/stack"
)

type dummyWriter int

func (d dummyWriter) Write(p []byte) (int, error) {
	fmt.Println("--- START ---")
	fmt.Println(string(p))
	fmt.Println("--- END ---")
	return len(p), nil
}

var d dummyWriter

func main() {
	// New webhook triggered, got repo url etc... from MQ
	// Instantiate golang stack
	goStack := stack.NewGolang("latest")

	// Create new contaner from stack
	cnt := container.New(goStack, "")

	// Run the container
	cnt.WG.Add(2)
	go cnt.Run()
	go cnt.Observe(d)

	log.Println("Waiting for build queue to finish")
	cnt.WG.Wait()
}
