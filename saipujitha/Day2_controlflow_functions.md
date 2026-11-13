# Day 2 -- Go: Control Flow & Functions

A practical study guide for conditionals, loops, functions, `defer`

## Learning roadmap

1.  Conditionals: `if`, `else`, short initialization, and `switch`

2.  Loops: standard `for`, while-style `for`, infinite loops, `break`,

    `continue`, and `range`

3.  Functions: parameters, return values, multiple returns, named

    returns, and errors

4.  `defer`: delayed execution and resource cleanup

------------------------------------------------------------------------

# 1. Conditionals in Go

Conditionals let a program make decisions based on whether an expression

is true or false. They are useful for validating input, checking

application status, and handling different outcomes.

## 1.1 Simple `if`

``` go

age := 23

if age >= 18 {

    fmt.Println("Eligible")

}

```

## 1.2 `if-else`

``` go

temperature := 35

if temperature > 30 {

    fmt.Println("It is hot")

} else {

    fmt.Println("Weather is pleasant")

}

```

## 1.3 `if-else if-else`

Go evaluates conditions from top to bottom and executes the first

matching branch.

``` go

marks := 85

if marks >= 90 {

    fmt.Println("Grade A")

} else if marks >= 75 {

    fmt.Println("Grade B")

} else if marks >= 50 {

    fmt.Println("Grade C")

} else {

    fmt.Println("Fail")

}

```

## 1.4 Short initialization statement

Go allows a variable to be declared and initialized inside an `if`

statement.

``` go

if age := 23; age >= 18 {

    fmt.Println("Eligible")

} else {

    fmt.Println("Not eligible")

}

```

The variable is in scope within the `if` and its associated `else`

block, not after the complete conditional.

A common Go error-handling pattern is:

``` go

if err := someFunction(); err != nil {

    fmt.Println("Error:", err)

    return

}

```


## 1.5 `switch`

A `switch` selects a code block based on a matching case.

``` go

day := 3

switch day {

case 1:

    fmt.Println("Monday")

case 2:

    fmt.Println("Tuesday")

case 3:

    fmt.Println("Wednesday")

default:

    fmt.Println("Invalid day")

}

```

Go automatically exits a matching case. An explicit `break` is generally

unnecessary.

Multiple values can appear in one case:

``` go

day := "Saturday"

switch day {

case "Saturday", "Sunday":

    fmt.Println("Weekend")

default:

    fmt.Println("Weekday")

}

```

A switch can also omit its expression, making each case a condition:

``` go

temperature := 35

switch {

case temperature > 40:

    fmt.Println("Very hot")

case temperature > 30:

    fmt.Println("Hot")

case temperature > 20:

    fmt.Println("Normal")

default:

    fmt.Println("Cold")

}

```

A switch can include initialization:

``` go

switch status := 503; status {

case 200:

    fmt.Println("Success")

case 404:

    fmt.Println("Not found")

case 503:

    fmt.Println("Service unavailable")

default:

    fmt.Println("Unknown status")

}

```

### `if-else` vs `switch`

  -----------------------------------------------------------------------

  `if-else`                           `switch`

  ----------------------------------- -----------------------------------

  Useful for ranges and complex       Useful for matching one value

  logical expressions                 against several cases

  Supports `&&`, `||`, and            Supports multiple values in a case

  comparisons                         

  General decision-making             Selecting an operation or status

  -----------------------------------------------------------------------

------------------------------------------------------------------------

# 2. Loops in Go

Loops repeat a block of code. They are useful for processing records,

iterating over servers, retrying operations, and periodically checking

application health.

Go has one looping keyword: **`for`**. It supports standard,

while-style, and infinite loops.

## 2.1 Standard `for` loop

A standard loop has three parts: initialization, condition, and post

statement.

``` go

for i := 1; i <= 5; i++ {

    fmt.Println(i)

}

```

Output: `1 2 3 4 5`

-   Initialization runs once.

-   The condition is checked before each iteration.

-   The post statement runs after each iteration.

## 2.2 While-style loop

Go has no separate `while` keyword. Use `for` with only a condition.

``` go

i := 1

for i <= 5 {

    fmt.Println(i)

    i++

}

```

Initialize before the loop, check the condition before each iteration,

and update the variable inside the loop.

## 2.3 Infinite loop

``` go

for {

    fmt.Println("Checking application health...")

    time.Sleep(5 * time.Second)

}

```

This continues until stopped by `break`, `return`, or an external

interruption.

## 2.4 `break`

`break` terminates the loop immediately.

``` go

for i := 1; i <= 10; i++ {

    if i == 5 {

        break

    }

    fmt.Println(i)

}

```

Output: `1 2 3 4`

## 2.5 `continue`

`continue` skips the rest of the current iteration and proceeds to the

next one.

``` go

for i := 1; i <= 5; i++ {

    if i == 3 {

        continue

    }

    fmt.Println(i)

}

```

Output: `1 2 4 5`

### `break` vs `continue`

  -----------------------------------------------------------------------

  `break`                             `continue`

  ----------------------------------- -----------------------------------

  Terminates the loop                 Skips the current iteration

  Control moves outside the loop      Control proceeds to the next

                                      iteration

  -----------------------------------------------------------------------

## 2.6 Nested loops

A nested loop is a loop inside another loop. The inner loop completes

its iterations for each iteration of the outer loop.

``` go

for i := 1; i <= 3; i++ {

    for j := 1; j <= 2; j++ {

        fmt.Println("i:", i, "j:", j)

    }

}

```

## 2.7 `range` loop

`range` iterates over arrays, slices, maps, strings, and channels.

``` go

servers := []string{"server-1", "server-2", "server-3"}

for index, server := range servers {

    fmt.Println(index, server)

}

```

If the index is not needed, use the blank identifier `_`:

``` go

for _, server := range servers {

    fmt.Println(server)

}

```

### Loop patterns summary

  Pattern       Example                      Purpose

  ------------- ---------------------------- ---------------------------

  Standard      `for i := 0; i < 5; i++`     Counted iterations

  While-style   `for i < 5`                  Condition-based iteration

  Infinite      `for {}`                     Continuous execution

  Range         `for _, v := range values`   Iterate over collections

------------------------------------------------------------------------

# 3. Functions in Go

A function is a reusable block of code that performs a specific task.

Functions divide a program into smaller parts and avoid repeating logic.

## 3.1 Basic function

``` go

func greet() {

    fmt.Println("Welcome to Go programming")

}

func main() {

    greet()

}

```

`func` declares a function. `main()` is the entry point of an executable

Go program.

## 3.2 Function with a parameter

``` go

func greet(name string) {

    fmt.Println("Hello", name)

}

func main() {

    greet("Pujitha")

}

```

A parameter is declared in the function signature. An argument is the

actual value passed during the call.

## 3.3 Multiple parameters

``` go

func add(a, b int) {

    result := a + b

    fmt.Println("Sum:", result)

}

```

When adjacent parameters share a type, Go allows `a, b int` instead of

`a int, b int`.

## 3.4 Function with one return value

``` go

func add(a, b int) int {

    result := a + b

    return result

}

func main() {

    sum := add(10, 20)

    fmt.Println("Sum:", sum)

}

```

The `int` after the parameter list specifies the return type. `return`

sends the result to the caller.

## 3.5 Multiple return values

Go allows a function to return more than one value. This is commonly

used to return a result and an error.

``` go

func calculate(a, b int) (int, int) {

    sum := a + b

    difference := a - b

    return sum, difference

}

func main() {

    sum, difference := calculate(20, 10)

    fmt.Println("Sum:", sum)

    fmt.Println("Difference:", difference)

}

```

If a returned value is not needed, use `_`:

``` go

sum, _ := calculate(20, 10)

```

## 3.6 Named return values

Go lets us name return variables in the function signature.

``` go

func calculate(a, b int) (result int) {
    result = a + b
    return
}
```

The return variables are declared by the signature. A bare `return` returns their current values. 


## Normal return values

``` go
func add(a int, b int) int {
    return a+b
}
```

we specify only the retun type.
we calculate the result and directly return it.
no named return variable is created.

### Normal returns vs named returns

  -----------------------------------------------------------------------

  Normal returns                      Named returns

  ----------------------------------- -----------------------------------

  Return types are specified in the   Return variables are named in the

  signature                           signature

  Usually returns explicit            Can use a bare `return`

  expressions                         

  Straightforward for short functions Can help when result names add

                                      clarity

  -----------------------------------------------------------------------

------------------------------------------------------------------------

# 4. `defer` in Go

The `defer` keyword schedules a function call to run when the

surrounding function is about to return.

Common cleanup uses: - Closing files. - Closing HTTP response bodies. -

Closing database rows or connections. - Releasing locks.

## 4.1 Basic `defer`

``` go

func main() {

    fmt.Println("Start")

    defer fmt.Println("Deferred statement")

    fmt.Println("End")

}

```

Output:

``` text

Start

End

Deferred statement

```

The deferred call runs as the surrounding function is returning.

## 4.2 Multiple defer statements: LIFO

Deferred calls execute in **Last In, First Out (LIFO)** order.

``` go

func main() {

    defer fmt.Println("First defer")

    defer fmt.Println("Second defer")

    defer fmt.Println("Third defer")

    fmt.Println("Main function")

}

```

Output:

``` text

Main function

Third defer

Second defer

First defer

```

## 4.3 `defer` for file cleanup

``` go

package main

import (

    "fmt"

    "os"

)

func main() {

    file, err := os.Open("sample.txt")

    if err != nil {

        fmt.Println("Error opening file:", err)

        return

    }

    defer file.Close()

    fmt.Println("File opened successfully")

    // Perform file operations here.

}

```
 