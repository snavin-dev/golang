package main

import "fmt"

func square(n int, results chan int) {
	results <- n * n // send the answer back
}

// producer can only send on out (chan<-)
func producer(out chan<- string) {
	out <- "hello from producer"
}

// consumer can only receive from in (<-chan)
func consumer(in <-chan string) {
	fmt.Println("consumer got:", <-in)
}

func main() {
	// 1. Unbuffered channel: send and receive wait for each other
	fmt.Println("--- unbuffered channel ---")
	messages := make(chan string)
	go func() {
		messages <- "ping" // send
	}()
	msg := <-messages // receive: waits until "ping" arrives
	fmt.Println("Received:", msg)

	// 2. Getting results back from goroutines.
	// Results arrive in the order the goroutines finish.
	fmt.Println("--- results from goroutines ---")
	results := make(chan int)
	for i := 1; i <= 3; i++ {
		go square(i, results)
	}
	for i := 1; i <= 3; i++ {
		fmt.Println("square:", <-results)
	}

	// 3. Buffered channel: room for 3 values, sending only waits when full
	fmt.Println("--- buffered channel ---")
	orders := make(chan string, 3)
	orders <- "tea"
	orders <- "coffee"
	fmt.Println("len:", len(orders), "cap:", cap(orders))
	fmt.Println(<-orders) // first in, first out
	fmt.Println(<-orders)

	// 4. close and range
	fmt.Println("--- close and range ---")
	nums := make(chan int)
	go func() {
		for i := 1; i <= 3; i++ {
			nums <- i
		}
		close(nums) // only the sender closes
	}()
	for n := range nums { // stops when nums is closed
		fmt.Println("got", n)
	}
	v, ok := <-nums // closed channel: zero value and false
	fmt.Println("after close:", v, ok)

	// 5. Send-only and receive-only channels
	fmt.Println("--- send-only and receive-only ---")
	pipe := make(chan string)
	go producer(pipe)
	consumer(pipe)
}
