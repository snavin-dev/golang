package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

func main() {
	fmt.Print("Enter a paragraph: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading paragraph:", err)
			return
		}
		fmt.Fprintln(os.Stderr, "No paragraph entered.")
		return
	}

	counts := make(map[string]int)
	for _, word := range strings.FieldsFunc(scanner.Text(), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}) {
		counts[strings.ToLower(word)]++
	}

	words := make([]string, 0, len(counts))
	for word := range counts {
		words = append(words, word)
	}
	sort.Strings(words)

	for _, word := range words {
		fmt.Printf("%s: %d\n", word, counts[word])
	}
}