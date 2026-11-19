package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. select waits on several channels and handles whichever is ready first
	fmt.Println("--- select ---")
	fast := make(chan string)
	slow := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		fast <- "fast server"
	}()
	go func() {
		time.Sleep(300 * time.Millisecond)
		slow <- "slow server"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg := <-fast:
			fmt.Println("Got reply from", msg)
		case msg := <-slow:
			fmt.Println("Got reply from", msg)
		}
	}

	// 2. Timeout: stop waiting after 1 second
	fmt.Println("--- timeout ---")
	result := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second) // pretend this is slow work
		result <- "done"
	}()

	select {
	case r := <-result:
		fmt.Println(r)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout: took longer than 1 second")
	}

	// 3. default: do not wait at all
	fmt.Println("--- default ---")
	inbox := make(chan string)
	select {
	case m := <-inbox:
		fmt.Println(m)
	default:
		fmt.Println("no message yet, moving on")
	}
}
