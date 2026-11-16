
# Go Language – KT Notes

These notes cover the Go topics discussed during the first three KT sessions.

---

# Day 1 – Introduction to Go

## 1. What is Go?

- Go (Golang) is a statically typed, compiled programming language.
- It is designed for simple, readable, and efficient software development.
- It has built-in support for concurrency through goroutines and channels.
- Go is commonly used for backend services, APIs, cloud tools, and infrastructure automation.

## 2. Key Features

- **Statically typed:** Types are checked by the compiler.
- **Compiled:** Go source code is compiled into an executable.
- **Garbage collected:** Memory is managed by the Go runtime.
- **Concurrency support:** Goroutines and channels are part of the language.
- **Simple syntax:** The language keeps common programming tasks straightforward.
- **Standard library:** Go includes packages for HTTP, JSON, files, networking, testing, and more.

## 3. Basic Go Program Structure

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

- `package main`: Declares the package. An executable Go program uses `main`.
- `import`: Makes another package available to the current file.
- `func main()`: Entry point of an executable Go program.
- `fmt`: Standard library package used for formatted input and output.
- `fmt.Println()`: Prints a line to standard output.

## 4. Packages and the Standard Library

- A **package** groups related Go source files and functionality.
- Packages help organize code and make it reusable.
- The Go standard library provides commonly used functionality.

Examples:

| Package | Purpose |
|---|---|
| `fmt` | Formatted input and output |
| `net/http` | HTTP clients and servers |
| `os` | Operating-system functionality |
| `encoding/json` | JSON encoding and decoding |
| `time` | Dates, times, and durations |
| `log` | Logging |

## 5. Basic Go Commands

| Command | Purpose |
|---|---|
| `go run .` | Compiles and runs the current package |
| `go build` | Builds the package |
| `go test ./...` | Runs tests in the current module and its packages |
| `go fmt ./...` | Formats Go source files |
| `go vet ./...` | Reports suspicious code patterns |
| `go version` | Displays the installed Go version |

---

# Day 2 – Go Modules and Dependency Management

## 1. What is a Go Module?

- A module is a collection of related Go packages managed and versioned together.
- A module is identified by a module path.
- Modules help manage external dependencies and make builds more consistent across local machines and CI/CD environments.
- The `go.mod` file is located at the root of a module.

## 2. Creating a Module

```bash
mkdir go-module-demo
cd go-module-demo
go mod init example.com/go-module-demo
```

This creates a `go.mod` file.

## 3. The `go.mod` File

The `go.mod` file describes the module and its dependencies.

Example:

```go
module example.com/go-module-demo

go 1.23

require github.com/google/uuid v1.6.0
```

Common entries:

- `module`: The module's import path.
- `go`: The Go language version or module semantics declared by the module.
- `require`: Lists required modules and their versions.
- `replace`: Replaces a module with another module source or version.
- `exclude`: Excludes a specific module version from being used.

## 4. The `go.sum` File

- `go.sum` contains checksums used to verify downloaded module content.
- It helps Go verify that dependency content matches the expected content.
- It is not a traditional lock file and does not independently select dependency versions.
- Go may add checksum entries for module versions and their `go.mod` files.
- Do not manually edit checksum entries.

## 5. Direct and Indirect Dependencies

- **Direct dependency:** A module imported by the current module's packages.
- **Indirect dependency:** A dependency required by another dependency, or otherwise needed by the module graph but not directly imported by the current module's packages.
- Go records dependency requirements in `go.mod`.

## 6. Dependency Management Commands

| Command | Purpose |
|---|---|
| `go mod init <module-path>` | Creates a new module |
| `go get <module>@<version>` | Adds or updates a dependency |
| `go mod tidy` | Adds missing module requirements and removes unused ones |
| `go mod download` | Downloads module dependencies |
| `go list -m all` | Lists modules in the build list |
| `go mod graph` | Displays the module requirement graph |
| `go mod verify` | Verifies downloaded module cache content |
| `go mod why -m <module>` | Explains why a module is needed |

## 7. Basic Dependency Demo

Create a Go file:

```go
package main

import (
    "fmt"

    "github.com/google/uuid"
)

func main() {
    id := uuid.New()
    fmt.Println("Generated UUID:", id)
}
```

Run:

```bash
go mod tidy
go run .
go build
```

After `go mod tidy`, inspect the `go.mod` and `go.sum` files.

## 8. Important Points

- Commit both `go.mod` and `go.sum` to version control.
- Use `go mod tidy` to keep dependency declarations in sync with the code.
- Review dependency changes before merging them.
- `go mod verify` checks module cache content; it does not assess whether a dependency is secure.
- Avoid upgrading dependencies without checking compatibility and release notes.

---

# Day 3 – Variables, Data Types, Zero Values, Constants, and Pointers

## 1. Variables in Go

A variable stores a value.

Go supports different ways to declare variables.

### Using `var`

```go
var name string = "Go"
var age int = 10
```

Type inference can also be used with `var`:

```go
var language = "Go"
```

### Using the Short Declaration Operator `:=`

```go
name := "Go"
age := 10
```

- `:=` declares and initializes a variable.
- It can only be used inside functions.
- Go infers the variable's type from the assigned value.
- `var` can be used at package level and inside functions.

### Multiple Variable Declarations

```go
var (
    name    string = "Go"
    version int    = 1
)
```

## 2. `var` vs `:=`

| `var` | `:=` |
|---|---|
| Can be used inside and outside functions | Can only be used inside functions |
| Can declare a variable without an initial value | Requires an initial value |
| Type can be specified explicitly or inferred | Type is inferred |
| Useful for package-level variables | Commonly used for local variables |

Example:

```go
var count int
count = 10

total := 20
```

## 3. Common Data Types

| Type | Description | Example |
|---|---|---|
| `int` | Integer | `var count int = 10` |
| `int32`, `int64` | Integers of specific sizes | `var id int64 = 100` |
| `float32`, `float64` | Floating-point numbers | `price := 10.5` |
| `string` | Text | `name := "Go"` |
| `bool` | Boolean value | `active := true` |
| `byte` | Alias for `uint8` | `var b byte = 'A'` |
| `rune` | Alias for `int32`, commonly used for Unicode code points | `r := '界'` |

Notes:

- Go is statically typed.
- A variable's type is determined at declaration.
- Go does not automatically convert between different numeric types.
- A string is an immutable sequence of bytes, commonly used to represent UTF-8 text.

## 4. Zero Values

When a variable is declared without an explicit initial value, Go assigns it a zero value.

| Type | Zero Value |
|---|---|
| Integer types | `0` |
| Floating-point types | `0` |
| `bool` | `false` |
| `string` | `""` |
| Pointer | `nil` |
| Slice | `nil` |
| Map | `nil` |
| Channel | `nil` |
| Interface | `nil` |
| Function | `nil` |

Example:

```go
var count int
var name string
var active bool

fmt.Println(count)  // 0
fmt.Println(name)   // ""
fmt.Println(active) // false
```

### Nil Slice vs Empty Slice

```go
var first []int
second := []int{}
```

- `first` is a nil slice.
- `second` is an empty but non-nil slice.
- Both have length zero.
- A nil slice can be appended to.

## 5. Constants

- Constants are declared using the `const` keyword.
- Their values cannot be changed after declaration.
- Constants can be typed or untyped.

```go
const appName = "MonitoringService"
const maxRetries int = 3
```

Example:

```go
const timeout = 30

fmt.Println(timeout)
```

## 6. Type Conversion

Go requires explicit type conversion when converting between different types.

```go
var count int = 10

var price float64 = float64(count)

fmt.Println(price)
```

Example of converting a floating-point number to an integer:

```go
value := 10.8
result := int(value)

fmt.Println(result) // 10
```

- Conversion from `float64` to `int` truncates the fractional part.
- Type conversion can result in loss of information.
- Go does not perform implicit numeric conversions.

## 7. Pointers

A pointer stores the memory address of another variable.

- `&` is used to get the address of a variable.
- `*` is used to declare a pointer type or dereference a pointer.
- A pointer that does not point to a value has the zero value `nil`.

Example:

```go
package main

import "fmt"

func main() {
    number := 10

    ptr := &number

    fmt.Println("Value:", number)
    fmt.Println("Address:", ptr)
    fmt.Println("Value through pointer:", *ptr)
}
```

## 8. Modifying a Value Using a Pointer

```go
package main

import "fmt"

func updateValue(number *int) {
    *number = 100
}

func main() {
    value := 10

    updateValue(&value)

    fmt.Println(value) // 100
}
```

- The function receives a pointer to the original variable.
- Dereferencing the pointer allows the function to modify the original value.

## 9. Pass by Value vs Passing a Pointer

Go always passes arguments by value.

### Pass by Value

```go
func update(number int) {
    number = 100
}

func main() {
    value := 10
    update(value)

    fmt.Println(value) // 10
}
```

- The function receives a copy of the value.
- Changes to the copy do not modify the original variable.

### Passing a Pointer

```go
func update(number *int) {
    *number = 100
}

func main() {
    value := 10
    update(&value)

    fmt.Println(value) // 100
}
```

- The function receives a copy of the pointer value.
- Both pointer values refer to the same underlying variable.
- Dereferencing the pointer allows the function to modify that variable.

---

# Quick Revision

| Day | Topics |
|---|---|
| Day 1 | Introduction to Go, key features, program structure, packages, standard library, and basic Go commands |
| Day 2 | Go modules, `go.mod`, `go.sum`, dependencies, and dependency management commands |
| Day 3 | Variables, declarations, data types, zero values, constants, type conversion, pointers, and pass by value |

---

# Official Learning Resources

- Go Getting Started: https://go.dev/doc/tutorial/getting-started
- Go Language Tour: https://go.dev/tour/
- Go Specification: https://go.dev/ref/spec
- Go Modules Reference: https://go.dev/ref/mod
- Go Modules Guide: https://go.dev/doc/modules/
- Go Standard Library: https://pkg.go.dev/std
- Go by Example: https://gobyexample.com/
- Go Playground: https://go.dev/play/
