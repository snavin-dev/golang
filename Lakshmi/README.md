# Practice Go: Day 1 to Day 3

My notes and practice programs from the first three days of learning Go.

```
practice-go/
├── go.mod
├── README.md
├── day1/
│   ├── basics/main.go          variables, constants, types
│   └── tempconverter/main.go   Day 1 goal: Celsius to Fahrenheit
└── day2/
    ├── basics/main.go          if, switch, loops, functions, defer
    └── calculator/main.go      Day 2 goal: calculator
└── day3/
    ├── arrays/main.go          fixed-size lists
    ├── slices/main.go          lists that can grow
    ├── maps/main.go            key -> value pairs
    ├── range/main.go           looping with range
    ├── pointers/main.go        & and *
    └── wordcounter/main.go     Day 3 goal: word-frequency counter
```

---

## Day 1: Setup and Basics

### Setup

1. Install Go and check it works:
   ```bash
   go version
   ```
2. Create a module (a Go project) inside the project folder:
   ```bash
   go mod init practice-go
   ```
   This creates a `go.mod` file with the project name.
3. Run a program:
   ```bash
   go run ./day1/basics
   ```

### Structure of a Go program

```go
package main          // this file belongs to package main

import "fmt"          // fmt is used for printing

func main() {         // the program starts here
    fmt.Println("Hello, Go")
}
```

- `package main` means this is a program you can run.
- `import` brings in code from other packages. `fmt` is used to print.
- `func main()` is where the program starts.

### Variables

There are three ways to create a variable:

```go
var name string = "Asha"   // with a type
var age = 25               // Go works out the type
height := 1.65             // short way, only inside functions
```

### Constants

A constant's value can't be changed after it's set.

```go
const country = "India"
```

### Basic types

| Type | Holds | Example |
|---|---|---|
| `int` | whole numbers | `25` |
| `float64` | decimal numbers | `1.65` |
| `string` | text | `"Asha"` |
| `bool` | true or false | `true` |

### Zero values

A variable created without a value gets a default value:

| Type | Default |
|---|---|
| `int` | `0` |
| `float64` | `0` |
| `string` | `""` (empty) |
| `bool` | `false` |

### Day 1 goal: temperature converter

Formula: **Fahrenheit = Celsius × 9 / 5 + 32**

```go
var celsius float64

fmt.Print("Enter temperature in Celsius: ")
fmt.Scan(&celsius)

fahrenheit := celsius*9/5 + 32

fmt.Println(celsius, "°C is", fahrenheit, "°F")
```

`fmt.Scan(&celsius)` waits for you to type a number and stores it in `celsius`. Write the `&` before the variable name for now; Day 3 explains what it means.

Run it:

```bash
go run ./day1/tempconverter
```

Check these inputs:

| Enter | Expected output |
|---|---|
| `100` | `100 °C is 212 °F` |
| `0` | `0 °C is 32 °F` |
| `37` | `37 °C is 98.6 °F` |

---

## Day 2: Control Flow and Functions

### if / else

```go
marks := 75
if marks >= 50 {
    fmt.Println("Result: Pass")
} else {
    fmt.Println("Result: Fail")
}
```

You can create a variable inside the `if`. It only exists inside that `if` block.

```go
if age := 20; age >= 18 {
    fmt.Println("Adult")
}
```

### switch

`switch` checks a value against several cases. Go stops after the first match, so you don't need `break`.

```go
switch day {
case 1:
    fmt.Println("Monday")
case 6, 7:
    fmt.Println("Weekend")
default:
    fmt.Println("Weekday")
}
```

### Loops

Go has only one loop keyword, `for`, used in three ways:

```go
// 1. Standard loop
for i := 1; i <= 3; i++ {
    fmt.Println("Count:", i)
}

// 2. While-style loop
n := 1
for n <= 3 {
    fmt.Println(n)
    n++
}

// 3. Infinite loop, stopped with break
x := 1
for {
    if x > 3 {
        break
    }
    fmt.Println(x)
    x++
}
```

### Functions

**Parameters and a return value:**

```go
func add(a int, b int) int {
    return a + b
}
```

**Multiple return values:**

```go
func divideAndRemainder(a int, b int) (int, int) {
    return a / b, a % b
}

q, r := divideAndRemainder(17, 5)   // q = 3, r = 2
```

**Named return values.** The return values get names, and a plain `return` sends them back:

```go
func rectangle(width int, height int) (area int, perimeter int) {
    area = width * height
    perimeter = 2 * (width + height)
    return
}
```

### defer

`defer` runs a line at the end of the function, just before it finishes.

```go
func main() {
    defer fmt.Println("Program finished")   // prints last
    fmt.Println("Hello")                    // prints first
}
```

It's mostly used for cleanup, such as closing a file when you're done with it.

Run all the Day 2 examples:

```bash
go run ./day2/basics
```

### Day 2 goal: calculator

The `divide` function returns **two values**: the answer, and `true` or `false` to say if the division worked. If the second number is 0, it returns `false` instead of dividing.

```go
func divide(a float64, b float64) (float64, bool) {
    if b == 0 {
        return 0, false
    }
    return a / b, true
}
```

In `main`, a `switch` picks the operation, and the result of `divide` is checked with `if ok`:

```go
case "/":
    result, ok := divide(a, b)
    if ok {
        fmt.Println("Answer:", result)
    } else {
        fmt.Println("Error: cannot divide by zero")
    }
```

Run it:

```bash
go run ./day2/calculator
```

Check these inputs:

| First number | Operator | Second number | Expected output |
|---|---|---|---|
| `10` | `+` | `5` | `Answer: 15` |
| `10` | `-` | `5` | `Answer: 5` |
| `6` | `*` | `7` | `Answer: 42` |
| `10` | `/` | `4` | `Answer: 2.5` |
| `10` | `/` | `0` | `Error: cannot divide by zero` |
| `2` | `%` | `3` | `Error: unknown operator %` |

Every run ends with `Thank you for using the calculator`, printed by `defer`.

---

## Day 3: Data Structures and Pointers

### Arrays

An array is a list with a **fixed size**. It can't grow or shrink.

```go
var scores [3]int                         // [0 0 0]
scores[0] = 90                            // positions start at 0

colors := [3]string{"red", "green", "blue"}
fmt.Println(colors[1])                    // green
fmt.Println(len(colors))                  // 3
```

Copying an array copies all its values, so changing the copy doesn't change the original:

```go
arr := [3]int{1, 2, 3}
copyArr := arr
copyArr[0] = 100
// arr = [1 2 3], copyArr = [100 2 3]
```

Run it:

```bash
go run ./day3/arrays
```

### Slices

A slice is a list that **can grow**. It looks like an array but has no size inside `[]`.

```go
fruits := []string{"apple"}
fruits = append(fruits, "banana")          // add one item
fruits = append(fruits, "mango", "grape")  // add several
// [apple banana mango grape]
```

Always save the result of `append` back into the variable.

**Taking part of a slice** with `[start:end]` includes `start` but not `end`:

```go
nums := []int{10, 20, 30, 40, 50}

nums[1:4]   // [20 30 40]
nums[:2]    // [10 20]
nums[2:]    // [30 40 50]
```

A part of a slice **shares memory** with the original:

```go
part := nums[1:4]
part[0] = 99
// nums is now [10 99 30 40 50]
```

**len and cap:**
- `len` is how many items are in the slice now.
- `cap` is how many fit before Go has to make more room.

**make** creates a slice with a length and a capacity:

```go
s := make([]int, 3, 5)   // [0 0 0], len 3, cap 5
s = append(s, 7)         // [0 0 0 7], len 4, cap 5
```

| | Array | Slice |
|---|---|---|
| Written as | `[3]int{1, 2, 3}` | `[]int{1, 2, 3}` |
| Size | Fixed | Can grow |
| Used | Rarely | Almost always |

Run it:

```bash
go run ./day3/slices
```

### Maps

A map stores **key → value** pairs, like a dictionary.

```go
ages := map[string]int{
    "Asha": 25,
    "Ravi": 30,
}

ages["Mia"] = 22       // add
ages["Asha"] = 26      // change
delete(ages, "Ravi")   // remove
len(ages)              // 2
```

Create an empty map with `{}`:

```go
counts := map[string]int{}
```

**Comma ok.** A missing key gives the zero value (`0` for int). To check whether the key really exists, ask for a second value:

```go
age, ok := ages["Tom"]
// age = 0, ok = false → Tom is not in the map

if age, ok := ages["Asha"]; ok {
    fmt.Println("Asha is", age)
}
```

The order of items in a map is **random**.

Run it:

```bash
go run ./day3/maps
```

### Range

`for ... range` loops over slices, maps and strings.

```go
fruits := []string{"apple", "banana", "mango"}

for i, fruit := range fruits {   // index and value
    fmt.Println(i, fruit)
}

for _, fruit := range fruits {   // _ skips the index
    fmt.Println(fruit)
}

for name, age := range ages {    // map: key and value
    fmt.Println(name, age)
}

for i, ch := range "Go!" {       // string: position and character
    fmt.Println(i, string(ch))
}
```

| Looping over | First value | Second value |
|---|---|---|
| slice | index | item |
| map | key | value |
| string | position | character |

Run it:

```bash
go run ./day3/range
```

### Pointers

A pointer holds the **address** of a variable, the place in memory where it's stored.

| Symbol | Meaning | Example |
|---|---|---|
| `&` | get the address of a variable | `p := &x` |
| `*` | use the value at that address | `*p = 20` |

```go
x := 10
p := &x          // p holds the address of x
fmt.Println(*p)  // 10
*p = 20
fmt.Println(x)   // 20, x changed through the pointer
```

**Why use pointers?** A function gets a **copy** of what you pass in. To let the function change your variable, pass its address:

```go
func double(n int) {
    n = n * 2            // changes only the copy
}

func doublePtr(n *int) {
    *n = *n * 2          // changes the original
}

y := 5
double(y)       // y is still 5
doublePtr(&y)   // y is now 10
```

This is why Day 1 used `fmt.Scan(&celsius)`. `Scan` needs the address so it can put the typed number into the variable.

Run it:

```bash
go run ./day3/pointers
```

### Day 3 goal: word-frequency counter

This program uses a slice, a map, range and comma ok together.

```go
text := "Go is fun. Go is fast, and Go is simple!"

counts := map[string]int{}   // word -> how many times
order := []string{}          // words in the order they appear

for _, w := range strings.Fields(text) {
    w = strings.Trim(w, ".,!?")   // remove punctuation
    w = strings.ToLower(w)        // "Go" and "go" are the same

    if _, seen := counts[w]; !seen {
        order = append(order, w)
    }

    counts[w]++
}

for _, w := range order {
    fmt.Println(w, counts[w])
}
```

How it works:
1. `strings.Fields` splits the text into a slice of words.
2. `range` goes through each word.
3. `Trim` removes punctuation and `ToLower` makes every word lower-case.
4. `counts[w]++` adds 1. A new word starts at 0, so it becomes 1.
5. The map's order is random, so the `order` slice keeps the words in the order they first appeared.

Run it:

```bash
go run ./day3/wordcounter
```

Expected output:

```
go 3
is 3
fun 1
fast 1
and 1
simple 1
```

---

## What I learned

**Short note*

1. Every Go program starts in `func main()` inside `package main`.
2. `:=` is the short way to create a variable.
3. Every type has a zero value, so a variable always has a value.
4. `for` is the only loop in Go, and `switch` doesn't need `break`.
5. A function can return more than one value. The calculator uses this to check for division by zero.
6. `defer` runs a line when the function ends.
7. Arrays have a fixed size. Slices can grow with `append`, and slices are what you use most.
8. Part of a slice shares memory with the original slice.
9. Maps store key → value pairs. Use comma ok to check whether a key exists.
10. `range` is the easy way to loop over slices, maps and strings.
11. `&` gets a variable's address and `*` uses the value at that address. Pass a pointer when a function needs to change your variable.
