
package main

import (
	"errors"
	"fmt"
)

func calculate(a, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("division by zero is not allowed")
		}
		return a / b, nil
	default:
		return 0, errors.New("invalid operator")
	}
}

func main() {
	var num1, num2 float64
	var operator string

	fmt.Print("Enter first number: ")
	_, err := fmt.Scan(&num1)
	if err != nil {
		fmt.Println("Invalid first number:", err)
		return
	}

	fmt.Print("Enter operator (+, -, *, /): ")
	_, err = fmt.Scan(&operator)
	if err != nil {
		fmt.Println("Invalid operator input:", err)
		return
	}

	fmt.Print("Enter second number: ")
	_, err = fmt.Scan(&num2)
	if err != nil {
		fmt.Println("Invalid second number:", err)
		return
	}

	result, err := calculate(num1, num2, operator)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Result:", result)
}
