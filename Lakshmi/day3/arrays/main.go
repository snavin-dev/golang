package main

import "fmt"

func main() {
	// An array has a fixed size
	var scores [3]int
	fmt.Println("Empty array:", scores)

	scores[0] = 90
	scores[1] = 85
	scores[2] = 70
	fmt.Println("Scores:", scores)

	// Create an array with values
	colors := [3]string{"red", "green", "blue"}
	fmt.Println("Second color:", colors[1])
	fmt.Println("Length:", len(colors))

	// Copying an array copies all the values
	arr := [3]int{1, 2, 3}
	copyArr := arr
	copyArr[0] = 100
	fmt.Println("Original:", arr)
	fmt.Println("Copy:", copyArr)
}
