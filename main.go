package main

import (
	"fmt"
	"log"

	"github.com/goblin-ci/runner/container"
	"github.com/goblin-ci/runner/github"
	"github.com/goblin-ci/runner/stack"
)

type dummyWriter int

func (d dummyWriter) Write(p []byte) (int, error) {
	fmt.Println(string(p))
	return len(p), nil
}

var d dummyWriter

func main() {

	// Receive json Repo request via MQ
	// and decode it to Repo struct
	repo := github.Repo{
		RequestID:  "f28a10b9",
		Branch:     "master",
		OwnerName:  "aneshas",
		OwnerEmail: "anes.hasicic@gmail.com",
		CloneURL:   "https://github.com/aneshas/guinea-pig",
		FullName:   "aneshas/guinea-pig",
	}

	// Instantiate golang stack
	goStack := stack.NewGolang("latest")

	// Create new contaner from stack
	cnt := container.New(goStack, &repo)

	// Run the container
	cnt.WG.Add(2)
	go cnt.Run()
	go cnt.Observe(d)

	log.Printf("Waiting for %s queue to finish", repo.RequestID)
	cnt.WG.Wait()
}
