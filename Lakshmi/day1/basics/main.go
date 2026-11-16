package main

import "fmt"

func main() {
	// Way 1: var with a type
	var name string = "Asha"

	// Way 2: var without a type (Go works out the type)
	var age = 25

	// Way 3: short declaration with :=
	height := 1.65
	isStudent := true

	// A constant can never change
	const country = "India"

	fmt.Println("Name:", name)
	fmt.Println("Age:", age)
	fmt.Println("Height:", height)
	fmt.Println("Is student:", isStudent)
	fmt.Println("Country:", country)

	// Variables without a value get a "zero value"
	var count int
	var price float64
	var city string
	var done bool
	fmt.Println("Zero values:", count, price, city, done)
}
