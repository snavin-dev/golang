# Error Handling and Packages in Go

A complete guide to two topics every Go program depends on: **how Go reports and handles errors**, and **how Go code is organised into packages**.

Each concept has a short explanation, when to use it, a small example, and the real output from running the example.

> **Go version:** every example works on Go 1.22 or newer, except `errors.AsType`, which needs **Go 1.26** or newer. Check yours with `go version`.

---

## Contents

**Part 1: Error Handling**

1. [How Go thinks about errors](#1-how-go-thinks-about-errors)
2. [The `error` type](#2-the-error-type)
3. [Returning and checking errors](#3-returning-and-checking-errors)
4. [Creating errors: `errors.New` and `fmt.Errorf`](#4-creating-errors-errorsnew-and-fmterrorf)
5. [Sentinel errors](#5-sentinel-errors)
6. [Custom error types](#6-custom-error-types)
7. [Wrapping errors: `%w` vs `%v`](#7-wrapping-errors-w-vs-v)
8. [`errors.Unwrap`](#8-errorsunwrap)
9. [`errors.Is`] [`errors.As`] (#9-errorsis)
10. [Several errors at once: `errors.Join`](#10-several-errors-at-once-errorsjoin-and-several-w)
11. [`errors.ErrUnsupported`](#11-errorserrunsupported)
12. [Errors and `defer`](#16-errors-and-defer)
13. [`panic` and `recover`](#17-panic-and-recover)

**Part 2: Packages**

1. [What a package is](#1-what-a-package-is)
2. [Why use packages`](#2-why-use-packages)
3. [How to create packages in Go](#3-how-to-create-packages-in-Go)
4. [Exported and unexported names](#4-exported-and-unexported-names)
5. [One package, several files](#5-one-package-several-files)
6. [Importing packages](#6-importing-packages)
7. [Third-party packages](#7-third-party-packages)


---

# Part 1: Error Handling

## 1. How Go thinks about errors

An error is a way to tell the program that something went wrong.

Go does not use try-catch like Java or Python. Instead, functions commonly return an error as a value.

In Go, **an error is a normal value that a function returns**. The caller checks it straight away, on the next line.

| | Exceptions (Java, Python, JS) | Go |
|---|---|---|
| How an error is reported | `throw` / `raise` | returned as the last value |
| How it's handled | `try` / `catch` somewhere above | `if err != nil` right after the call |
| Can you see which calls can fail? | Not from the code | Yes, the function returns `error` |
| Can you forget to handle it? | Easily | Harder, the value is right there |

The flow of a typical Go function looks like this:

```
call a function
      │
      ▼
 err != nil ? ──── yes ──▶ handle it or return it (with context)
      │
      no
      │
      ▼
 use the result and carry on
```

---

## 2. The `error` type

`error` is a built-in **interface** with a single method:

```go
type error interface {
    Error() string
}
```

Any type that has an `Error() string` method **is** an error. The value `nil` means **"no error"**.

| Value of `err` | Meaning |
|---|---|
| `nil` | everything worked |
| not `nil` | something failed; `err.Error()` describes what |

When you print an error with `fmt.Println(err)`, Go calls its `Error()` method for you.

---

## 3. Returning and checking errors

**The convention:** the error is the **last** return value. When there is an error, the other values are zero values. When there is no error, the error is `nil`.

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 2)
    fmt.Println(result, err)

    result, err = divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(result)
}
```

Output:

```
5 <nil>
Error: cannot divide by zero
```

**Handle the error first, then carry on.** Go code keeps the "happy path" on the left and returns early when something fails:

```go
data, err := os.ReadFile("config.txt")
if err != nil {
    return err          // stop here
}
// carry on with data
```

---

## 4. Creating errors: `errors.New` and `fmt.Errorf`

| Function | Use it when |
|---|---|
| `errors.New("text")` | the message is fixed |
| `fmt.Errorf("format", values...)` | the message needs values in it |

```go
err1 := errors.New("file is empty")
err2 := fmt.Errorf("user %q has %d unread messages", "ravi", 3)

fmt.Println(err1)
fmt.Println(err2)
fmt.Printf("%T\n", err1)
```

Output:

```
file is empty
user "ravi" has 3 unread messages
*errors.errorString
```

**Every call to `errors.New` creates a new, different error**, even with the same text:

```go
fmt.Println(errors.New("oops") == errors.New("oops"))
```

Output:

```
false
```

This is why you create an error once and keep it in a variable when you want to compare against it. That's the next section.

---

## 5. Sentinel errors

A **sentinel error** is an error created once, stored in a package-level variable, and returned whenever that particular problem happens. Callers compare against the variable.

By convention its name starts with `Err`.

```go
var ErrNotFound = errors.New("not found")

func findUser(name string) (string, error) {
    users := map[string]string{"asha": "Asha Rao"}
    full, ok := users[name]
    if !ok {
        return "", ErrNotFound
    }
    return full, nil
}

func main() {
    _, err := findUser("ravi")
    if err == ErrNotFound {
        fmt.Println("No such user")
    }
}
```

Output:

```
No such user
```

**When to use:** when the caller needs to know **which** problem happened, but no extra details are needed. Examples from the standard library: `io.EOF`, `fs.ErrNotExist`, `sql.ErrNoRows`.

> Once errors are wrapped (section 7), use `errors.Is(err, ErrNotFound)` instead of `==`.

---

## 6. Custom error types

When the caller needs **details** about the error (which field, what value, which code), make your own type with an `Error()` method.

```go
type ValidationError struct {
    Field  string
    Reason string
}

func (e *ValidationError) Error() string {
    return e.Field + ": " + e.Reason
}

func setAge(age int) error {
    if age <= 0 {
        return &ValidationError{Field: "age", Reason: "must be positive"}
    }
    return nil
}

func main() {
    err := setAge(-5)
    fmt.Println(err)

    if ve, ok := err.(*ValidationError); ok {
        fmt.Println("Field:", ve.Field)
    }
}
```

Output:

```
age: must be positive
Field: age
```

`err.(*ValidationError)` is a **type assertion**: "is this error a `*ValidationError`? If so, give it to me." It only checks the outer error. Once errors are wrapped, use `errors.As` or `errors.AsType` (sections 10 and 11).

**Sentinel error or custom type?**

| | Sentinel error | Custom error type |
|---|---|---|
| Example | `var ErrNotFound = errors.New(...)` | `type ValidationError struct{...}` |
| Carries extra data | No | Yes |
| Checked with | `errors.Is` | `errors.As` / `errors.AsType` |
| Use when | you only need to know *which* error | you need *details* about it |

---

## 7. Wrapping errors: `%w` 

As an error travels up through your functions, each one can **add context** so the final message says *where* and *why*:

`fmt.Errorf` does this. The verb you use matters:

| Verb | Message | Original error kept inside? | `errors.Is` finds it? |
|---|---|---|---|
| `%w` | adds context | **yes** (wrapped) | **yes** |
| `%v` | adds context | no, only its text | no |

```go
package main

import (
    "errors"
    "fmt"
)

func connectDatabase() error {
    return errors.New("connection timeout")
}

func main() {

    err := connectDatabase()

    if err != nil {
        wrappedErr := fmt.Errorf(
            "failed to initialize application: %w", err,
        )

        fmt.Println(wrappedErr)
    }
}
```

Output:

```
failed to initialize application: connection timeout
```

**Use `%w`** when callers might need to check the original error. **Use `%v`** when you deliberately want to hide it (for example, so callers don't depend on an internal detail).

---

## 8. `errors.Unwrap`

`errors.Unwrap(err)` returns the error **one layer inside** `err`, or `nil` if there's nothing inside. A wrapped error is like a set of boxes:

```
┌───────────────────────────────────────────┐
│ handle request: ...                       │  e3
│   ┌───────────────────────────────────┐   │
│   │ find user: ...                    │   │  e2
│   │   ┌───────────────────────────┐   │   │
│   │   │ not found  (ErrNotFound)  │   │   │  e1
│   │   └───────────────────────────┘   │   │
│   └───────────────────────────────────┘   │
└───────────────────────────────────────────┘
```

```go
e1 := ErrNotFound
e2 := fmt.Errorf("find user: %w", e1)
e3 := fmt.Errorf("handle request: %w", e2)

for e := error(e3); e != nil; e = errors.Unwrap(e) {
    fmt.Println("->", e)
}
```

Output:

```
-> handle request: find user: not found
-> find user: not found
-> not found
```

You rarely call `Unwrap` yourself. `errors.Is` and `errors.As` do this unwrapping for you.

---

## 9. `errors.Is`. and error.As

errors.Is() checks whether an error chain contains a particular error value.

errors.As() checks whether an error chain contains an error of a particular type.

```go
errors.Is(err, target) bool
```

**Question it answers:** "Is `target` anywhere inside `err`?" It checks `err`, then unwraps layer by layer and checks each one.

```go
errors.As(err, &target) bool
```

**Question it answers:** "Is there an error **of this type** anywhere inside `err`? If so, put it in `target`." Use it to get the **details** out of a custom error type, even when it's wrapped.


```go

package main

import (
    "errors"
    "fmt"
)

// Sentinel error
var ErrDatabaseConnection = errors.New("database connection failed")

// Custom error type
type DatabaseError struct {
    Code    int
    Message string
}

// Implement error interface
func (e *DatabaseError) Error() string {
    return e.Message
}

// Simulate database connection
func connectDatabase() error {

    dbErr := &DatabaseError{
        Code:    500,
        Message: "Connection timeout",
    }

    // Wrap sentinel error and custom error together
    return fmt.Errorf("%w: %w", ErrDatabaseConnection, dbErr)
}

func main() {

    err := connectDatabase()

    if err != nil {

        fmt.Println("Error:", err)

        // 1. Check error identity using errors.Is()
        if errors.Is(err, ErrDatabaseConnection) {
            fmt.Println("Identified: Database connection error")
        }

        // 2. Extract custom error using errors.As()
        var dbErr *DatabaseError

        if errors.As(err, &dbErr) {
            fmt.Println("Error Code:", dbErr.Code)
            fmt.Println("Error Message:", dbErr.Message)
        }
    }
}

```
Output:

```
Error: database connection failed: Connection timeout

Identified: Database connection error
Error Code: 500
Error Message: Connection timeout
```
---

## 10. Several errors at once: `errors.Join` and several `%w`

Sometimes more than one thing goes wrong. errors.Join() is used when you want to combine multiple errors into a single error.

It was introduced in Go 1.20.

Imagine an application that performs two operations:
1. Database connection.
2. Configuration loading.
Both operations fail, and we want to return both errors together instead of returning only one.

```go
package main

import (
    "errors"
    "fmt"
)

func main() {

    err1 := errors.New("database connection failed")
    err2 := errors.New("configuration file missing")

    // Combine multiple errors
    combinedErr := errors.Join(err1, err2)

    fmt.Println(combinedErr)
}
```

Output:

```
database connection failed
configuration file missing
```
---

## 11. `errors.ErrUnsupported`

A ready-made sentinel error (Go 1.21+) meaning **"this operation is not supported"**. The docs say not to return it directly, but to **wrap it** with context.

```go
func exportPDF() error {
    return fmt.Errorf("export pdf: %w", errors.ErrUnsupported)
}

func main() {
    err := exportPDF()
    fmt.Println(err)
    if errors.Is(err, errors.ErrUnsupported) {
        fmt.Println("PDF export is not available, try CSV")
    }
}
```

Output:

```
export pdf: unsupported operation
PDF export is not available, try CSV
```
---

## 12. Errors and `defer`

defer means execute this statement when the surrounding function is about to finish.

It is commonly used for cleanup operations such as:
- Closing files.
- Closing database connections.
- Releasing resources.

```go
package main

import "fmt"

func main() {
    defer fmt.Println("This executes at the end")

    fmt.Println("First statement")
    fmt.Println("Second statement")
}
```

Output:

```
First statement
Second statement
This executes at the end
```
Even though the defer statement appears first, its execution is postponed until the function exits.

Important: Multiple deferred statements execute in LIFO order (Last In, First Out).

```go
func main() {
    defer fmt.Println("First")
    defer fmt.Println("Second")
    defer fmt.Println("Third")

    fmt.Println("Main function")
}
```

Output:

```
Main function
Third
Second
First
```
---

## 13. `panic` and `recover`

A **`panic`** typically means something went unexpectedly wrong. Mostly we use it to fail fast on errors that shouldn’t occur during normal operation, or that we aren’t prepared to handle gracefully.

- The program starts normally.
- panic() is triggered.
- The remaining statements in the function are not executed.
- The program terminates unless the panic is recovered.
Important: Deferred functions still execute during panic unwinding.

```go
package main

import "fmt"

func main() {
    fmt.Println("Program started")

    panic("Something went seriously wrong")

    fmt.Println("Program finished")
}
```

Output:

```
Program started
panic: Something went seriously wrong
```

Go itself panics on mistakes like dividing an integer by zero or using an index past the end of a slice.

**`recover`** stops a panic. It only works inside a **deferred** function. You can use it to turn a panic into an ordinary error:

```go
package main

import "fmt"

func main() {

    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()

    fmt.Println("Program started")

    panic("Something went wrong")

    fmt.Println("Program finished")
}
```

Output:

```
Program started
Recovered from panic: Something went wrong
```
---

# Part 2: Packages

## 1. What a package is

A **package** in Go, is a way of organizing related Go source code files together or  say  a collection of related Go source files which helps developers organize code into reusable and manageable units.

A package can contain:
- Functions
- Variables
- Constants
- Structs
- Interfaces
- Types and methods

And

- Every Go file starts with `package <name>`.
- All files in one folder must use the **same** package name.
- **`package main`** is special: it builds into a program you can run, and it must have a `func main()`.
- Any other package (`package employee`, `package fmt`) is a **library**: other code imports it and uses it.

## 2. Why use packages

Imagine writing a 2,000-line application in a single main.go file.

What problems could we face?
1. Code becomes difficult to maintain.
2. Reusing functionality becomes difficult.
3. Multiple developers working on the same application may face conflicts.
4. Testing individual components becomes difficult.
5. Business logic becomes tightly coupled.
6. Changes in one area may affect unrelated functionality.

Packages solve these problems by providing logical separation.

```
employee-management/
│
├── main.go
├── employee/
│   ├── employee.go
│   └── salary.go
│
└── go.mod
```

Each package has its own responsibility.

This follows the principle of separation of concerns.

---

## 3. How to create packages in Go

These steps build the `employee-management` project from an empty folder.

### Step 1: Make the project folder and go into it

```bash
mkdir employee-management
cd employee-management
```

### Step 2: Create the module with `go mod init`

```bash
go mod init employee-management
```

Output:

```
go: creating new go.mod: module employee-management
```

This creates `go.mod`:

```go
module employee-management

go 1.22
```
### Step 3: Make a folder for the package

Each package gets its own folder. The folder name becomes the package name.

```bash
mkdir employee
```
### Step 4: Create the package files

employee/employee.go 
employee/salary.go

Both files must start with the same package name, matching the folder:

```go
package employee
```

**employee/employee.go**

```go
// Package employee stores employee details and works out salaries.
package employee

import "fmt"

var count int

func init() {
    fmt.Println("employee package is ready")
}

// Employee holds one employee's details.
type Employee struct {
    Name   string
    Role   string
    salary float64
}

// New creates an Employee and counts it.
func New(name, role string, salary float64) Employee {
    count++
    return Employee{Name: name, Role: role, salary: salary}
}

// Count returns how many employees were created.
func Count() int { return count }

// Salary returns the monthly salary.
func (e Employee) Salary() float64 { return e.salary }
```

**employee/salary.go**

```go
package employee

const bonusRate = 0.10

// YearlySalary returns the salary for 12 months.
func (e Employee) YearlySalary() float64 { return e.salary * 12 }

// Bonus returns 10% of the yearly salary.
func (e Employee) Bonus() float64 { return e.YearlySalary() * bonusRate }
```

**main.go**

```go
package main

import (
    "fmt"

    "employee-management/employee"
)

func init() {
    fmt.Println("main package is ready")
}

func main() {
    asha := employee.New("Asha", "Developer", 5000)
    ravi := employee.New("Ravi", "Tester", 4000)

    fmt.Println(asha.Name, "earns", asha.Salary(), "a month")
    fmt.Println(ravi.Name, "gets a bonus of", ravi.Bonus())
    fmt.Println("Employees created:", employee.Count())
}
```

Run it from the `employee-management` folder:

```bash
go run .
```

Output:

```
employee package is ready
main package is ready
Asha earns 5000 a month
Ravi gets a bonus of 4800
Employees created: 2
```

How the pieces connect:

```
                 go.mod: module employee-management
                               │
            ┌──────────────────┴──────────────────┐
            │                                     │
      package main                        package employee
        main.go   ── imports ──▶   employee.go  +  salary.go
                 "employee-management/employee"
```

The next sections explain each idea this project uses.

---

## 4. Exported and unexported names

Go has no `public` or `private` keywords. **The first letter of a name decides who can use it:**

| First letter | Called | Who can use it | Examples in the project |
|---|---|---|---|
| **Capital** | exported | any package that imports it | `Employee`, `New`, `Count`, `Salary`, `Bonus`, `Name`, `Role` |
| **lower-case** | unexported | only code in the same package | `count`, `salary`, `bonusRate`, `init` |

The rule applies to **everything**: functions, types, struct fields, methods, constants and variables.

```
package employee                        package main
┌──────────────────────────┐            ┌──────────────────────────┐
│ Employee   (exported) ───┼──────────▶ │ employee.Employee   ✅   │
│ New()      (exported) ───┼──────────▶ │ employee.New()      ✅   │
│ Name       (exported) ───┼──────────▶ │ asha.Name           ✅   │
│ salary     (unexported)  │     ✖      │ asha.salary         ❌   │
│ bonusRate  (unexported)  │     ✖      │ employee.bonusRate  ❌   │
└──────────────────────────┘            └──────────────────────────┘
```

---

## 5. One package, several files

A package can be split across as many files as you like, as long as they're in the **same folder** and use the **same package name**.

- Files in the same package **see each other's names** (exported or not) without importing anything. That's why `salary.go` can use the `Employee` type and the `salary` field from `employee.go`.
- Split files by topic to keep each one short and easy to read.

If two files in one folder use different package names, the build fails:


**`go run .` vs `go run main.go`.** `go run main.go` compiles **only** that one file. If `package main` is split across several files, use `go run .` to include them all.

---

## 6. Importing packages

### Single and grouped imports

```go
import "fmt"                          // one package

import (                              // several packages
    "fmt"                             // standard library first
    "strings"

    "employee-management/employee"    // then your own packages
)
```

`go fmt` sorts imports, and the convention is to put standard library packages in one group and your own or third-party packages in another, separated by a blank line.

You use an imported package through the **last part** of its path: `"employee-management/employee"` becomes `employee.New(...)`.

### Every import must be used

An unused import is a **compile error**, not a warning:

```
./main.go:5:2: "os" imported and not used
```

---

## 7. Third-party packages

Packages written by other people are downloaded by module path. You can search for them at [pkg.go.dev](https://pkg.go.dev).

```bash
go get github.com/google/uuid
```

Then import it like any other package:

```go
import "github.com/google/uuid"

id := uuid.New()
```

`go get` updates two files:

| File | What it holds |
|---|---|---|
| `go.mod` | the modules you depend on and their versions (`require` lines) |
| `go.sum` | checksums so everyone downloads exactly the same code |

A `go.mod` with a dependency looks like this:

```go
module employee-management

go 1.22

require github.com/google/uuid v1.6.0
```
---
