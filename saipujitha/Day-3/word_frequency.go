package main

import (
    "fmt"
    "strings"
)

func main() {
    paragraph := "I work on go and go is powerful"

    words := strings.Fields(
        strings.ToLower(paragraph),
    )

    wordCount := make(map[string]int)

    for _, word := range words {
        wordCount[word]++
    }

    fmt.Println("Word Frequency:")

    for word, count := range wordCount {
        fmt.Println(word, ":", count)
    }
}