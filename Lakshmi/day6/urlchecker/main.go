package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// result is what each goroutine sends back for one website.
type result struct {
	url    string
	status string        // for example "200 OK" or "404 Not Found"
	code   int           // for example 200 or 404
	err    error         // set when the site could not be reached at all
	took   time.Duration // how long the check took
}

// tally counts how many sites are up and down.
// Many goroutines update it at once, so a Mutex protects it.
type tally struct {
	mu   sync.Mutex
	up   int
	down int
}

func (t *tally) add(isUp bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if isUp {
		t.up++
	} else {
		t.down++
	}
}

// check visits one URL and returns its status code.
func check(client *http.Client, url string) result {
	start := time.Now()

	resp, err := client.Get(url)
	if err != nil {
		return result{url: url, err: err, took: time.Since(start)}
	}
	defer resp.Body.Close()

	return result{url: url, status: resp.Status, code: resp.StatusCode, took: time.Since(start)}
}

func main() {
	urls := []string{
		"https://go.dev",
		"https://pkg.go.dev",
		"https://github.com",
		"https://www.google.com",
		"https://go.dev/this-page-does-not-exist",
		"https://no-such-site.invalid",
	}

	// Each request gives up after 5 seconds.
	client := &http.Client{Timeout: 5 * time.Second}

	// Buffered, so a goroutine can always send its result and finish,
	// even if main has stopped waiting.
	results := make(chan result, len(urls))

	var wg sync.WaitGroup
	var stats tally

	start := time.Now()

	// Start one goroutine per website. They all run at the same time.
	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := check(client, url)
			stats.add(r.err == nil && r.code < 400)
			results <- r
		}()
	}

	// Print each result as it arrives, but stop waiting after 10 seconds.
	timeout := time.After(10 * time.Second)
	for i := 0; i < len(urls); i++ {
		select {
		case r := <-results:
			if r.err != nil {
				fmt.Printf("ERROR  %-42s %v\n", r.url, r.err)
			} else {
				fmt.Printf("%-6d %-42s %-15s %v\n", r.code, r.url, r.status, r.took.Round(time.Millisecond))
			}
		case <-timeout:
			fmt.Println("Stopped waiting: some sites took longer than 10 seconds")
			return
		}
	}

	wg.Wait() // every goroutine has finished, so the tally is complete

	fmt.Printf("\nChecked %d sites in %v: %d up, %d down\n",
		len(urls), time.Since(start).Round(time.Millisecond), stats.up, stats.down)
}
