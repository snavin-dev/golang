package main

import "fmt"

func main() {
    numbers := [3]int{10, 20, 30}
    fmt.Println("Array:", numbers)

    fruits := []string{"Apple", "Banana"}
    fruits = append(fruits, "Mango")

    fmt.Println("Slice:", fruits)
    fmt.Println("Length:", len(fruits))
    fmt.Println("Capacity:", cap(fruits))
}