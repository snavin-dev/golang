package main

import "fmt"

// This program crashes on purpose to show a deadlock.
// An unbuffered send waits for a receiver, but nobody else is running,
// so Go stops with: fatal error: all goroutines are asleep - deadlock!
func main() {
	ch := make(chan int)
	ch <- 1 // waits forever: no other goroutine will receive
	fmt.Println(<-ch)
}
