package main

import "fmt"

// A function with two parameters and one return value
func add(a int, b int) int {
	return a + b
}

// A function with two return values
func divideAndRemainder(a int, b int) (int, int) {
	return a / b, a % b
}

// A function with named return values
func rectangle(width int, height int) (area int, perimeter int) {
	area = width * height
	perimeter = 2 * (width + height)
	return
}

func main() {
	// defer: this line runs at the very end of main
	defer fmt.Println("Program finished")

	// if / else
	marks := 75
	if marks >= 50 {
		fmt.Println("Result: Pass")
	} else {
		fmt.Println("Result: Fail")
	}

	// if with a short statement
	if age := 20; age >= 18 {
		fmt.Println("Adult")
	}

	// switch
	day := 6
	switch day {
	case 1:
		fmt.Println("Monday")
	case 6, 7:
		fmt.Println("Weekend")
	default:
		fmt.Println("Weekday")
	}

	// Loop 1: standard for loop
	for i := 1; i <= 3; i++ {
		fmt.Println("Count:", i)
	}

	// Loop 2: while-style loop
	n := 1
	for n <= 3 {
		fmt.Println("While-style:", n)
		n++
	}

	// Loop 3: infinite loop, stopped with break
	x := 1
	for {
		if x > 3 {
			break
		}
		fmt.Println("Infinite loop:", x)
		x++
	}

	// Calling functions
	fmt.Println("add(2, 3) =", add(2, 3))

	q, r := divideAndRemainder(17, 5)
	fmt.Println("17 / 5 =", q, "remainder", r)

	area, perimeter := rectangle(4, 3)
	fmt.Println("Area:", area, "Perimeter:", perimeter)
}
