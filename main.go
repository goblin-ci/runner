package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goblin-ci/dispatch"
	"github.com/goblin-ci/runner/docker"
	"github.com/goblin-ci/runner/github"
	"github.com/goblin-ci/runner/stack"
)

type dummyWriter int

func (d dummyWriter) Write(p []byte) (int, error) {
	fmt.Println(string(p))
	return len(p), nil
}

var d dummyWriter
var mq dispatch.PubSuber

func main() {
	mq, err := dispatch.NewRedis("redis:6379")
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan bool, 1)
	duration := time.Second * 1
	recv, err := mq.Subscribe("foo", stop, duration)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		//stop <- true
	}()

	mq.Publish("foo",
		github.Push{
			RequestID: "askdfjsjf",
			Branch:    "master",
		})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case payload, ok := <-recv:
			if !ok {
				recv = nil
			}
			log.Println("New payload")
			log.Println(payload)
		case <-sig:
			log.Println("OS Signal received, exiting...")
			stop <- true
			recv = nil
		}

		if recv == nil {
			break
		}
	}

	log.Println("EOP")

	//http.HandleFunc("/", handle)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {

	// Receive json Repo request via MQ
	// and decode it to Repo struct
	push := github.Push{
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
	cnt := docker.New(goStack, &push)

	// Run the container
	cnt.WG.Add(2)
	go cnt.Run()
	go cnt.Observe(w)

	log.Printf("Waiting for %s queue to finish", push.RequestID)
	cnt.WG.Wait()
}
