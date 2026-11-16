package main

import "fmt"

func celsiusToFahrenheit(celsius float64) float64 {
    fahrenheit := (celsius * 9 / 5) + 32
    return fahrenheit
}

func main() {
    var celsius float64

    fmt.Print("Enter temperature in Celsius: ")
    fmt.Scan(&celsius)

    result := celsiusToFahrenheit(celsius)

    fmt.Println("Temperature in Fahrenheit:", result)
}