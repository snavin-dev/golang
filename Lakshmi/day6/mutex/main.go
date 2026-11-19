package main

import (
	"fmt"
	"sync"
)

// countWithoutLock has a data race: 1000 goroutines change count at once,
// so some updates are lost and the answer is usually less than 1000.
func countWithoutLock() int {
	count := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count++
		}()
	}

	wg.Wait()
	return count
}

// countWithLock uses a Mutex so only one goroutine changes count at a time.
func countWithLock() int {
	count := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}

	wg.Wait()
	return count
}

func main() {
	fmt.Println("without Mutex:", countWithoutLock(), "(expected 1000)")
	fmt.Println("with Mutex:   ", countWithLock(), "(always 1000)")
}
