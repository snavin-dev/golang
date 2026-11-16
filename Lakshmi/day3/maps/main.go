package main

import "fmt"

func main() {
	// A map stores key -> value pairs
	ages := map[string]int{
		"Asha": 25,
		"Ravi": 30,
	}

	ages["Mia"] = 22  // add
	ages["Asha"] = 26 // change
	fmt.Println("Ravi's age:", ages["Ravi"])

	delete(ages, "Ravi") // remove
	fmt.Println("Number of people:", len(ages))

	// A missing key gives the zero value
	fmt.Println("Tom's age:", ages["Tom"])

	// Comma ok: check if the key exists
	age, ok := ages["Tom"]
	fmt.Println("Tom:", age, ok)

	if age, ok := ages["Asha"]; ok {
		fmt.Println("Asha is", age)
	} else {
		fmt.Println("Asha not found")
	}

	// Start with an empty map
	counts := map[string]int{}
	counts["go"] = 1
	fmt.Println("counts:", counts)
}
