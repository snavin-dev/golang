package main

import "fmt"

func main() {
	var celsius float64

	fmt.Print("Enter temperature in Celsius: ")
	fmt.Scan(&celsius) // reads the number you type into celsius

	fahrenheit := celsius*9/5 + 32

	fmt.Println(celsius, "°C is", fahrenheit, "°F")
}
