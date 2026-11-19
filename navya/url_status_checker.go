
package main

import (
    "fmt"
    "net/http"
    "sync"
)

func checkWebsite(url string, wg *sync.WaitGroup) {
    defer wg.Done()

    response, err := http.Get(url)

    if err != nil {
        fmt.Println(url, "Error:", err)
        return
    }

    defer response.Body.Close()

    fmt.Println(url, "Status:", response.StatusCode)
}

func main() {
    websites := []string{
        "https://www.google.com",
        "https://www.github.com",
        "https://www.amazon.com",
        "https://www.wikipedia.org",
    }

    var wg sync.WaitGroup

    for _, url := range websites {
        wg.Add(1)
        go checkWebsite(url, &wg)
    }

    wg.Wait()

    fmt.Println("All websites checked!")
}
