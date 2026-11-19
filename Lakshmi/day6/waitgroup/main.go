package main

import (
	"fmt"
	"sync"
	"time"
)

func cook(dish string, wg *sync.WaitGroup) {
	defer wg.Done() // counter - 1 when cook finishes
	time.Sleep(100 * time.Millisecond)
	fmt.Println(dish, "is ready")
}

func main() {
	var wg sync.WaitGroup
	dishes := []string{"rice", "dal", "curry"}

	start := time.Now()
	for _, d := range dishes {
		wg.Add(1)       // counter + 1, before starting the goroutine
		go cook(d, &wg) // pass a pointer, not a copy
	}

	wg.Wait() // wait until the counter is 0
	fmt.Println("All dishes are ready")

	// All 3 cooked at the same time, so this is about 100ms, not 300ms
	fmt.Println("Took about", time.Since(start).Round(100*time.Millisecond))
}
