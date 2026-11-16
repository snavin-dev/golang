# Practice Go: Day 1 and Day 2

My notes and practice programs from the first two days of learning Go.

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

## What I learned

1. Every Go program starts in `func main()` inside `package main`.
2. `:=` is the short way to create a variable.
3. Every type has a zero value, so a variable always has a value.
4. `for` is the only loop in Go, and `switch` doesn't need `break`.
5. A function can return more than one value. The calculator uses this to check for division by zero.
6. `defer` runs a line when the function ends.
