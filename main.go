package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/redis.v3"

	"github.com/goblin-ci/dispatch"
	"github.com/goblin-ci/runner/docker"
	"github.com/goblin-ci/runner/github"
	"github.com/goblin-ci/runner/stack"
)

var mq dispatch.PubSuber

type fooWriter int

func (f fooWriter) Write(p []byte) (int, error) {
	fmt.Println(string(p))
	return len(p), nil
}

func main() {
	var fw fooWriter

	mq, err := dispatch.NewRedis("redis:6379")
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan bool, 2)
	duration := time.Second * 1
	recv, err := mq.Subscribe("github_webhook_push", stop, duration)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		stop <- true
	}()

	mq.Publish(
		"github_webhook_push",
		github.Push{
			RequestID:  "f28a10b9",
			Branch:     "master",
			OwnerName:  "aneshas",
			OwnerEmail: "anes.hasicic@gmail.com",
			CloneURL:   "https://github.com/aneshas/guinea-pig",
			FullName:   "aneshas/guinea-pig",
		})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case payload, ok := <-recv:
			if !ok {
				recv = nil
			}
			msg := payload.(*redis.Message)
			push := new(github.Push)
			err := json.Unmarshal([]byte(msg.Payload), &push)
			if err != nil {
				log.Println(err)
				break
			}

			// Instantiate golang stack
			goStack := stack.NewGolang("latest")

			// Create new contaner from stack
			cnt := docker.New(goStack, push)

			// Run the container
			cnt.WG.Add(2)
			go cnt.Run()
			go cnt.Observe(fw)

			log.Printf("Waiting for %s queue to finish", push.RequestID)
			// cnt.WG.Wait()

			log.Println(push)
		case <-sig:
			log.Println("OS Signal received, exiting...")
			recv = nil
		}

		if recv == nil {
			break
		}
	}

}
