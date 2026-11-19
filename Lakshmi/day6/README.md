# Day 6: Concurrency (Goroutines and Channels)

**Concurrency** means a program works on several tasks during the same period of time. While one task waits (for a website, a file or a timer), another one runs.

Go has two main tools for this:

| Tool | What it does | Think of it as |
|---|---|---|
| **goroutine** | runs a function alongside the rest of the program | another cook in the kitchen |
| **channel** | passes data safely between goroutines | the counter where cooks hand over dishes |

```
day6/
├── goroutines/main.go       go keyword, main does not wait
├── channels/main.go         unbuffered, buffered, close, range
├── select/main.go           select, timeout, default
├── waitgroup/main.go        waiting for goroutines to finish
├── mutex/main.go            data race and how a Mutex fixes it
├── deadlock/main.go         crashes on purpose to show a deadlock
└── urlchecker/
    ├── main.go              Day 6 goal: concurrent URL status checker
    └── main_test.go         tests that use a local test server
```

Run every command from the `practice-go` folder.

> **The order of output changes from run to run** in most of these programs. That's normal: Go decides which goroutine runs when.

---

## Goroutines

A goroutine is started by writing `go` before a function call:

```go
sayHello()      // normal call: waits until it finishes
go sayHello()   // goroutine: starts it and carries on straight away
```

- Goroutines are very cheap. A program can run thousands of them.
- `main` is a goroutine too. **When `main` ends, the program ends**, and any goroutines still running are stopped.
- A goroutine can also be a function with no name: `go func() { ... }()`.

```bash
go run ./day6/goroutines
```

Example output (the order of the middle lines changes each run):

```
--- main does not wait ---
main carries on without waiting
--- two goroutines ---
anonymous goroutine says hi
B step 1
Hello from a goroutine
A step 1
A step 2
B step 2
B step 3
A step 3
main finished
```

This program uses `time.Sleep` to wait for the goroutines. That's only for learning. The right way is a WaitGroup (below).

---

## Channels

A channel is a pipe between goroutines. One sends a value in, another receives it.

| Code | Meaning |
|---|---|
| `ch := make(chan string)` | create an unbuffered channel of strings |
| `ch := make(chan string, 3)` | create a buffered channel with room for 3 |
| `ch <- "hi"` | send |
| `msg := <-ch` | receive |
| `close(ch)` | no more values will be sent |
| `for v := range ch` | receive until the channel is closed |
| `v, ok := <-ch` | `ok` is `false` when the channel is closed |

**Unbuffered** (`make(chan T)`): the sender waits until someone receives, like handing something over in person.

**Buffered** (`make(chan T, n)`): the sender only waits when the buffer is full, like a mailbox.

**Send-only and receive-only** in function parameters: `chan<- T` can only send, `<-chan T` can only receive.

Rules:
- Only the **sender** closes a channel.
- Sending on a closed channel panics.
- A function started with `go` can't return a value, so it **sends its result on a channel** instead.

```bash
go run ./day6/channels
```

Output (the order of the three `square` lines changes each run):

```
--- unbuffered channel ---
Received: ping
--- results from goroutines ---
square: 9
square: 1
square: 4
--- buffered channel ---
len: 2 cap: 3
tea
coffee
--- close and range ---
got 1
got 2
got 3
after close: 0 false
--- send-only and receive-only ---
consumer got: hello from producer
```

### Deadlock

If every goroutine is waiting and none can continue, Go stops the program. This program does it on purpose:

```bash
go run ./day6/deadlock
```

```
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
```

With an unbuffered channel, the sender and the receiver must be in **different** goroutines.

---

## Select

`select` is like a `switch` for channels. It waits until one of its cases is ready and runs that one.

| `select` with... | Behaviour |
|---|---|
| only channel cases | waits until one is ready |
| a `case <-time.After(d):` | waits, but only up to time `d` (a timeout) |
| a `default:` case | never waits |

```bash
go run ./day6/select
```

Output:

```
--- select ---
Got reply from fast server
Got reply from slow server
--- timeout ---
timeout: took longer than 1 second
--- default ---
no message yet, moving on
```

---

## WaitGroup

A `sync.WaitGroup` counts running goroutines so `main` can wait for all of them.

| Method | Meaning |
|---|---|
| `wg.Add(1)` | one more goroutine is starting (call it **before** `go`) |
| `wg.Done()` | this goroutine has finished (use `defer wg.Done()`) |
| `wg.Wait()` | wait until all have finished |

Pass the WaitGroup as a pointer (`&wg`), so every goroutine uses the same counter.

```bash
go run ./day6/waitgroup
```

Output (the order of the dishes changes each run):

```
rice is ready
dal is ready
curry is ready
All dishes are ready
Took about 100ms
```

Each dish takes 100ms, but they cook at the same time, so all three take about 100ms in total, not 300ms.

---

## Mutex

When several goroutines change the **same variable** at the same time, some changes get lost. This is a **data race**. `count++` is really three steps (read, add 1, write), and two goroutines can read the same value:

```
goroutine 1: read 5 ── add 1 ── write 6
goroutine 2:    read 5 ── add 1 ── write 6     ← one increase is lost
```

A `sync.Mutex` is a lock: only one goroutine at a time can be between `mu.Lock()` and `mu.Unlock()`.

```bash
go run ./day6/mutex
```

Output (the first number changes each run):

```
without Mutex: 972 (expected 1000)
with Mutex:    1000 (always 1000)
```

Go can find data races for you:

```bash
go run -race ./day6/mutex
```

It prints `WARNING: DATA RACE` for `countWithoutLock`.

| Use | When |
|---|---|
| channel | passing data or results between goroutines |
| `sync.Mutex` | several goroutines update one shared value |
| `sync.WaitGroup` | waiting until all goroutines finish |

---

## Day 6 goal: concurrent URL status checker

The program checks a list of websites **at the same time** and prints each one's HTTP status code.

How it works:

1. `check` visits one URL with `client.Get` and returns the status code, or an error if the site can't be reached. It closes the response with `defer resp.Body.Close()`.
2. `main` starts **one goroutine per URL**, so all the checks run at the same time.
3. Each goroutine **sends its result on a channel**. The channel is buffered, so a goroutine can always send and finish.
4. `main` prints results as they arrive, using **`select`** with **`time.After`** so it stops waiting after 10 seconds.
5. A **`tally`** struct counts sites that are up and down. Many goroutines update it at once, so it's protected by a **`Mutex`**.
6. A **`WaitGroup`** makes sure every goroutine has finished before the totals are printed.
7. The `http.Client` has a 5-second `Timeout`, so one slow site can't hang the program.

A site counts as **up** when it replies with a status code below 400.

```bash
go run ./day6/urlchecker
```

Example output (it needs the internet, and the order and times change each run):

```
ERROR  https://no-such-site.invalid               Get "https://no-such-site.invalid": dial tcp: lookup no-such-site.invalid: no such host
200    https://github.com                         200 OK          155ms
200    https://www.google.com                     200 OK          308ms
200    https://go.dev                             200 OK          486ms
200    https://pkg.go.dev                         200 OK          563ms
404    https://go.dev/this-page-does-not-exist    404 Not Found   574ms

Checked 6 sites in 574ms: 4 up, 2 down
```

**Why concurrency helps:** the times above add up to more than 2 seconds, but the whole check took 574ms, about as long as the **slowest** site. Checking one after another would take as long as **all of them added together**.

### Testing it

The tests start a small web server on your own computer with `httptest`, so they don't need the internet:

| Test | Checks |
|---|---|
| `TestCheckStatusCodes` | a page that exists gives 200, a missing one gives 404 |
| `TestCheckTimeout` | a slow page gives a timeout error |
| `TestCheckUnreachable` | a site that isn't there gives an error |
| `TestTallyIsSafe` | 100 goroutines updating the tally at once give the right totals |

```bash
go test -v ./day6/urlchecker
go test -race ./day6/urlchecker
```

---

## What I learned

1. `go f()` starts a goroutine. `main` doesn't wait for it.
2. A channel passes data between goroutines. Unbuffered channels make the sender and receiver wait for each other; buffered channels only wait when full or empty.
3. Only the sender closes a channel, and `range` reads until it's closed.
4. `select` handles whichever channel is ready first. `time.After` adds a timeout, and `default` means don't wait.
5. `sync.WaitGroup` waits for goroutines to finish: `Add` before `go`, `defer Done()` inside, `Wait()` at the end.
6. When goroutines share a variable, protect it with `sync.Mutex`, and check with `go run -race`.
7. Doing slow work concurrently takes about as long as the slowest task, not all of them added together.
