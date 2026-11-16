package main

import "fmt"

// divide returns two values:
// the answer, and true/false to say if the division worked
func divide(a float64, b float64) (float64, bool) {
	if b == 0 {
		return 0, false
	}
	return a / b, true
}

func main() {
	defer fmt.Println("Thank you for using the calculator")

	var a, b float64
	var op string

	fmt.Print("Enter first number: ")
	fmt.Scan(&a)
	fmt.Print("Enter operator (+ - * /): ")
	fmt.Scan(&op)
	fmt.Print("Enter second number: ")
	fmt.Scan(&b)

	switch op {
	case "+":
		fmt.Println("Answer:", a+b)
	case "-":
		fmt.Println("Answer:", a-b)
	case "*":
		fmt.Println("Answer:", a*b)
	case "/":
		result, ok := divide(a, b)
		if ok {
			fmt.Println("Answer:", result)
		} else {
			fmt.Println("Error: cannot divide by zero")
		}
	default:
		fmt.Println("Error: unknown operator", op)
	}
}
