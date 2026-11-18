# Day 6: Concurrency in Go — Goroutines, Channels, Select, and Sync

**Presentation duration:**
**Goal:** Understand how Go runs tasks concurrently, how goroutines communicate, and how to coordinate and protect them safely.

---

## 1. What is concurrency?

Concurrency means structuring a program so that multiple tasks can make progress during overlapping periods of time.

For example, an application might:
- Read data from an API.
- Process that data.
- Write the result to a file.

These tasks may be handled concurrently instead of making the program wait for each task to finish before starting the next one.

**Concurrency vs. parallelism**
- **Concurrency:** Multiple tasks are in progress during the same period. Their execution may take turns.
- **Parallelism:** Multiple tasks are literally running at the same time, usually on different CPU cores.

Go makes concurrency convenient with goroutines and channels.

---

## 2. Goroutines — running functions concurrently

A **goroutine** is a lightweight unit of execution managed by the Go runtime. We start one by putting the `go` keyword before a function call.

Goroutines are cheaper to create than traditional operating-system threads, so a Go program can use many of them.

### Example: Start a goroutine

Create a file named `goroutine.go`:

```go
package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello from goroutine")
}

func main() {
	go sayHello()

	fmt.Println("Hello from main")

	time.Sleep( time.Second)
}
```

### Explanation

- `go sayHello()` starts `sayHello` in a new goroutine.
- The main function continues without waiting for that function to finish.
- `time.Sleep` is used here only to keep the example simple and give the goroutine time to print.

**Important:** `time.Sleep` is not the right way to coordinate goroutines in real applications. The program can exit when `main` finishes, even if other goroutines are still running. We will use `sync.WaitGroup` for proper coordination later.

**Try it:** Run the program several times. The order of the two printed messages can vary.

---

## 3. Channels — communication between goroutines

A **channel** is a typed way for goroutines to send and receive values.

A useful way to think about it is: one goroutine produces a value, and another goroutine receives it.

Create a channel with `make`:

```go
ch := make(chan int)
```

### 3.1 Unbuffered channel

An **unbuffered channel** has no storage capacity. A send waits until another goroutine is ready to receive, and a receive waits until a sender is ready.

Create `unbuffered_channel.go`:

```go
package main

import "fmt"

func sendData(ch chan int) {
	ch <- 100 // Send value into the channel
}

func main() {
	ch := make(chan int)

	go sendData(ch)

	value := <-ch // Receive value from the channel
	fmt.Println("Received:", value)
}
```

### Explanation

- `make(chan int)` creates an unbuffered integer channel.
- `ch <- 100` sends `100`.
- `<-ch` receives the value.
- The send and receive synchronize with each other.

### 3.2 Buffered channel

A **buffered channel** can hold a fixed number of values without requiring a receiver immediately.

Create `buffered_channel.go`:

```go
package main

import "fmt"

func main() {
	ch := make(chan string, 2)

	ch <- "Go"
	ch <- "Concurrency"

	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
```

### Explanation

- `make(chan string, 2)` creates a channel with capacity 2.
- Both values can be sent before receiving them.
- Sending another value while the buffer is full blocks until space becomes available.

**Remember:**
- Unbuffered channel: sender and receiver must meet.
- Buffered channel: values can wait in the channel buffer, up to its capacity.

---

## 4. Select — working with multiple channels 

The `select` statement lets a goroutine wait on multiple channel operations.

Go runs a case when its channel operation is ready. If multiple cases are ready, Go chooses one of them.

A `select` is useful when:
- A function is waiting for data from more than one channel.
- A task should stop when a timeout occurs.
- A program needs to respond to whichever event happens first.

### Example: Select with a timeout

Create `select_timeout.go`:

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "Data received"
	}()

	select {
	case message := <-ch:
		fmt.Println(message)

	case <-time.After(1 * time.Second):
		fmt.Println("Operation timed out")
	}
}
```

### Explanation

- A goroutine waits for two seconds before sending a message.
- `time.After(1 * time.Second)` creates a channel that becomes ready after one second.
- The `select` waits for either the message or the timeout.
- Since the timeout happens first, the program prints `Operation timed out`.

**Note:** The example demonstrates timeout selection. In production code, also consider how to stop or clean up work that is no longer needed.

---

## 5. Sync package — coordinating and protecting goroutines

The `sync` package provides tools to coordinate goroutines and safely access shared data.

Two commonly used tools are:

### 5.1 WaitGroup — wait for goroutines to finish

A `sync.WaitGroup` lets the main goroutine wait until a group of goroutines has completed.

Create `waitgroup.go`:

```go
package main

import (
	"fmt"
	"sync"
)

func printMessage(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Goroutine", id, "finished")
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go printMessage(i, &wg)
	}

	wg.Wait()
	fmt.Println("All goroutines finished")
}
```

### Explanation

- `wg.Add(1)` tells the WaitGroup that one more task needs to finish.
- `go printMessage(...)` starts a goroutine.
- `defer wg.Done()` marks that goroutine as finished when the function returns.
- `wg.Wait()` blocks until the counter becomes zero.

**Important:** Call `Add` before starting the goroutine. Each `Add(1)` must be matched by a `Done()`.

### 5.2 Mutex — protect shared data

A **mutex** (mutual exclusion lock) makes sure only one goroutine at a time can access a protected section of code.

This is useful when multiple goroutines update the same variable. Without protection, their updates can overlap and cause a **data race**.

Create `mutex.go`:

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}
```

### Explanation

- All three goroutines update the same `counter`.
- `mu.Lock()` locks the critical section.
- `counter++` runs while the lock is held.
- `mu.Unlock()` releases the lock so another goroutine can enter.
- `WaitGroup` makes sure all updates finish before the final value is printed.

The output is:

```text
Counter: 3
```

**Important:** Always release a mutex. In longer functions, `defer mu.Unlock()` immediately after `mu.Lock()` is often a safer pattern.

---

## 6. Quick comparison and recap

| Concept | What it does | Simple way to remember |
|---|---|---|
| Goroutine | Runs a function concurrently | `go functionName()` |
| Channel | Sends data between goroutines | `ch <- value`, `<-ch` |
| Unbuffered channel | Send and receive synchronize | No waiting room |
| Buffered channel | Holds a limited number of values | Small waiting room |
| Select | Waits on multiple channel operations | Whichever is ready |
| WaitGroup | Waits for goroutines to finish | Tracks completed tasks |
| Mutex | Protects shared data | One goroutine at a time |

### Key takeaways

1. Use goroutines to run independent work concurrently.
2. Use channels when goroutines need to communicate.
3. Use `select` to handle multiple channel operations and timeouts.
4. Use `WaitGroup` when one part of the program must wait for goroutines to finish.
5. Use `Mutex` when multiple goroutines access shared data and at least one of them modifies it.

A useful Go principle is: **“Don't communicate by sharing memory; share memory by communicating.”**

This means channels are often a good choice for passing data between goroutines. A mutex is still appropriate when protecting shared state is the simpler design.

---


## How to run the examples

Each code example is a separate Go program. Save each snippet in its own `.go` file and run it from the terminal:

```bash
go run goroutine.go
go run unbuffered_channel.go
go run buffered_channel.go
go run select_timeout.go
go run waitgroup.go
go run mutex.go
```

