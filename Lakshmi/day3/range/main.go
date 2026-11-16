package main

import "fmt"

func main() {
	fruits := []string{"apple", "banana", "mango"}

	// Over a slice: index and value
	for i, fruit := range fruits {
		fmt.Println(i, fruit)
	}

	// Use _ when you don't need the index
	for _, fruit := range fruits {
		fmt.Println("Fruit:", fruit)
	}

	// Over a map: key and value (order is random)
	ages := map[string]int{"Asha": 25, "Ravi": 30}
	for name, age := range ages {
		fmt.Println(name, "is", age)
	}

	// Over a string: position and character
	for i, ch := range "Go!" {
		fmt.Println(i, ch, string(ch))
	}
}
