package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"leohetsch.com/simulation/consumer"
)

func main() {

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// TODO run a consumer with specific topic, group ID and logic

	c := make(chan string) // TODO maybe don't use producer domain/struct
	go consumer.NewConsumer(c, sigchan)

	// TODO better way to do this??
loop:
	for {
		select {
		case ev := <-c:
			fmt.Println(ev)
		case <-sigchan:
			break loop
		}
	}
}
