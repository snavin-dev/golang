package main

import "fmt"

// double gets a copy of n, so the caller's variable does not change
func double(n int) {
	n = n * 2
}

// doublePtr gets the address of n, so it can change the original
func doublePtr(n *int) {
	*n = *n * 2
}

func main() {
	x := 10
	p := &x // & gives the address of x

	fmt.Println("Address in p:", p)
	fmt.Println("Value at p:", *p) // * gives the value at the address

	*p = 20 // change x through the pointer
	fmt.Println("x is now:", x)

	y := 5
	double(y)
	fmt.Println("After double(y):", y)

	doublePtr(&y)
	fmt.Println("After doublePtr(&y):", y)
}
