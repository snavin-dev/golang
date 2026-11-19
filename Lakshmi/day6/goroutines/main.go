package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello from a goroutine")
}

func work(name string) {
	for i := 1; i <= 3; i++ {
		fmt.Println(name, "step", i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// 1. main does not wait for goroutines: it carries on straight away.
	// sayHello only gets to print because main sleeps at the end.
	// Remove the time.Sleep below and it will almost never print.
	fmt.Println("--- main does not wait ---")
	go sayHello()
	fmt.Println("main carries on without waiting")

	// 2. Two goroutines running at the same time.
	// The order of the lines changes every run.
	fmt.Println("--- two goroutines ---")
	go work("A")
	go work("B")

	// 3. A goroutine from a function with no name
	go func() {
		fmt.Println("anonymous goroutine says hi")
	}()

	// Wait so the goroutines have time to finish.
	// time.Sleep is only for learning; Day 6 waitgroup shows the right way.
	time.Sleep(500 * time.Millisecond)
	fmt.Println("main finished")
}
