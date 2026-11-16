package main

import "fmt"

func main() {
	// A slice has no size inside []
	fruits := []string{"apple"}

	// append adds items. Always save the result back.
	fruits = append(fruits, "banana")
	fruits = append(fruits, "mango", "grape")
	fmt.Println("Fruits:", fruits)

	// Slicing: [start:end] includes start, not end
	nums := []int{10, 20, 30, 40, 50}
	fmt.Println("nums[1:4] =", nums[1:4])
	fmt.Println("nums[:2]  =", nums[:2])
	fmt.Println("nums[2:]  =", nums[2:])

	// A part of a slice shares memory with the original
	part := nums[1:4]
	part[0] = 99
	fmt.Println("After changing part, nums =", nums)

	// len = items now, cap = room before it needs to grow
	fmt.Println("part len:", len(part), "cap:", cap(part))

	// make: length 3, capacity 5
	s := make([]int, 3, 5)
	fmt.Println("make:", s, "len:", len(s), "cap:", cap(s))

	s = append(s, 7)
	fmt.Println("after append:", s, "len:", len(s), "cap:", cap(s))

	// An empty slice is nil, and append still works on it
	var empty []int
	fmt.Println("Is empty nil?", empty == nil, "len:", len(empty))
	empty = append(empty, 1)
	fmt.Println("After append:", empty)
}
