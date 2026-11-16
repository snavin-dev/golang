
package main

import (
	"fmt"
	"strings"
)

func main() {
	var paragraph string

	fmt.Println("Enter a paragraph:")
	fmt.Scanln(&paragraph)

	words := strings.Fields(paragraph)

	wordCount := make(map[string]int)

	for _, word := range words {
		wordCount[word]++
	}

	fmt.Println("\nWord Frequency:")

	for word, count := range wordCount {
		fmt.Println(word, ":", count)
	}
}
