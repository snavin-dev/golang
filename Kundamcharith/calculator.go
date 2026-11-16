package main

import (
	"errors"
	"fmt"
)

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

func main() {
	defer fmt.Println("Thank you for using the calculator")

	var a, b float64
	var op string

	fmt.Print("Enter first number: ")
	if _, err := fmt.Scan(&a); err != nil {
		fmt.Println("Error reading first number:", err)
		return
	}

	fmt.Print("Enter operator (+ - * /): ")
	if _, err := fmt.Scan(&op); err != nil {
		fmt.Println("Error reading operator:", err)
		return
	}

	fmt.Print("Enter second number: ")
	if _, err := fmt.Scan(&b); err != nil {
		fmt.Println("Error reading second number:", err)
		return
	}

	switch op {
	case "+":
		fmt.Println("Answer:", a+b)
	case "-":
		fmt.Println("Answer:", a-b)
	case "*":
		fmt.Println("Answer:", a*b)
	case "/":
		result, err := divide(a, b)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Answer:", result)
	default:
		fmt.Println("Error: unknown operator", op)
	}
}
