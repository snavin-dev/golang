package main

import "fmt"

func main() {
    ages := map[string]int{
        "John": 22,
        "Alex": 25,
        "Sam": 20,
    }

    fmt.Println(ages["John"])

    age, ok := ages["Alex"]

    if ok {
        fmt.Println("Age:", age)
    } else {
        fmt.Println("Key not found")
    }
}