package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "Go is fun. Go is fast, and Go is simple!"

	counts := map[string]int{} // word -> how many times
	order := []string{}        // words in the order they first appear

	// strings.Fields splits the text into words at the spaces
	for _, w := range strings.Fields(text) {
		w = strings.Trim(w, ".,!?") // remove punctuation
		w = strings.ToLower(w)      // "Go" and "go" are the same word

		// First time we see this word? Remember its order.
		if _, seen := counts[w]; !seen {
			order = append(order, w)
		}

		counts[w]++ // a missing key starts at 0
	}

	for _, w := range order {
		fmt.Println(w, counts[w])
	}
}
