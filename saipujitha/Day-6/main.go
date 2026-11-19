
package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func checkWebsite(url string, wg *sync.WaitGroup, results chan string) {
	defer wg.Done()

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)

	if err != nil {
		results <- fmt.Sprintf("%s -> Error: %v", url, err)
		return
	}

	defer resp.Body.Close()

	results <- fmt.Sprintf("%s -> Status: %s", url, resp.Status)
}

func main() {

	websites := []string{
		"https://google.com",
		"https://github.com",
		"https://golang.org",
		"https://stackoverflow.com",
		"https://amazon.com",
		"https://example.com",
	}

	var wg sync.WaitGroup

	results := make(chan string, len(websites))

	start := time.Now()

	// Launch goroutines
	for _, url := range websites {
		wg.Add(1)
		go checkWebsite(url, &wg, results)
	}

	// Wait for all goroutines
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and print results
	fmt.Println("Checking websites...\n")

	for result := range results {
		fmt.Println(result)
	}

	fmt.Println("\nAll websites checked!")
	fmt.Println("Total Time:", time.Since(start))
}
