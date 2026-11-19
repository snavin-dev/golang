# Employee CRUD API: Project Walkthrough

This document explains the project from process startup to an HTTP request and
its PostgreSQL query. It is intended as a guide for presenting the design and
request flow to a team.

## 1. What the application does

The application exposes an HTTP API for creating, reading, updating, and
deleting employee records. It uses:

- Go's standard `net/http` package for the HTTP server and routing.
- PostgreSQL to persist employee data.
- `pgx` and `pgxpool` to communicate with PostgreSQL and reuse connections.
- Go contexts to bound startup and health-check work and to propagate request
  cancellation to database calls.

The high-level request path is:

```text
HTTP client
    │
    ▼
net/http server and ServeMux router
    │
    ▼
EmployeeHandler (HTTP decoding, validation, status codes, JSON)
    │
    ▼
EmployeeRepository interface
    │
    ▼
PostgresEmployeeRepository (SQL and row scanning)
    │
    ▼
PostgreSQL employees table
```

The response follows the same path in reverse: PostgreSQL returns rows or an
error, the repository converts that result into Go data or a repository error,
and the handler writes an HTTP status and JSON response.

## 2. Project layout and responsibilities

```text
crud-operations/
├── main.go
├── go.mod
├── go.sum
├── cmd/
│   └── cmd.go
├── docs/
│   └── PROJECT_WALKTHROUGH.md
└── internal/
    ├── config/
    │   └── config.go
    ├── database/
    │   └── postgres.go
    ├── handlers/
    │   └── employee.go
    ├── models/
    │   └── employee.go
    └── repository/
        └── postgres_employee.go
```

| Location | Responsibility |
|---|---|
| `main.go` | Small executable entry point; calls `cmd.Execute`. |
| `cmd/cmd.go` | Loads config, opens the database, wires dependencies, registers routes, starts the server, and handles shutdown. |
| `internal/config` | Reads required and optional settings from environment variables. |
| `internal/database` | Creates the PostgreSQL pool, checks connectivity, and ensures the employees table exists. |
| `internal/models` | Defines employee API/data structures and input normalization/validation. |
| `internal/repository` | Defines employee data operations and implements them with PostgreSQL. |
| `internal/handlers` | Converts HTTP requests into application/repository calls and writes HTTP responses. |
| `go.mod` / `go.sum` | Declare the module/dependencies and record dependency checksums. |

`internal` is a Go visibility convention enforced by the compiler: packages
inside it can be imported by this module, but cannot be imported by unrelated
modules outside the parent project tree.

## 3. Startup: from `main.go` to the listening server

### 3.1 The entry point

`main.go` is intentionally small:

```go
func main() {
    if err := cmd.Execute(); err != nil {
        log.Fatal(err)
    }
}
```

When the program is launched with `go run .` or as a compiled executable, Go
calls `main`. It delegates all application setup to `cmd.Execute`.

`Execute` returns an `error` rather than exiting the process itself. This keeps
startup and server lifecycle logic in one place. `log.Fatal` reports an
unexpected returned error and exits with a non-zero status. (Like other calls
to `os.Exit`, `log.Fatal` does not run deferred functions in `main`; normal
resource cleanup is performed inside `Execute` before it returns.)

### 3.2 Load environment configuration

The first operation in `cmd.Execute` is:

```go
cfg, err := config.Load()
if err != nil {
    return err
}
```

`config.Load` reads:

- `DATABASE_URL`: required PostgreSQL connection string.
- `HTTP_ADDR`: optional HTTP bind address. It defaults to `localhost:8080`.

The configuration loader trims surrounding whitespace. An absent database URL
is a startup error: the app does not guess which database to use.

Example PowerShell configuration:

```powershell
$env:DATABASE_URL = "postgres://postgres:your-password@localhost:5432/employees?sslmode=disable"
$env:HTTP_ADDR = "localhost:8080"
```

Use the real local connection details for your PostgreSQL setup. Do not commit
actual credentials to source control.

### 3.3 Bound startup work and open PostgreSQL

```go
startupCtx, cancelStartup := context.WithTimeout(
    context.Background(),
    10*time.Second,
)
defer cancelStartup()

db, err := database.Open(startupCtx, cfg.DatabaseURL)
if err != nil {
    return fmt.Errorf("initialize database: %w", err)
}
defer db.Close()
```

`context.Background()` is the root context for this startup operation.
`context.WithTimeout` derives a context that is canceled after ten seconds.
Passing it into `database.Open` bounds the database ping and table creation
performed at startup.

`cancelStartup` is deferred so the context timer/resources are released when
`Execute` exits, including when startup fails early.

`database.Open` returns a `*pgxpool.Pool`. A pool manages reusable database
connections, allowing simultaneous operations to use connections from one
shared pool rather than creating a new pool for every request.

`defer db.Close()` makes the pool close when `Execute` returns, after the
server has stopped. `fmt.Errorf` with `%w` adds startup context while retaining
the original error for further inspection.

### 3.4 Database initialization

`database.Open` performs these steps:

1. Creates a `pgxpool.Pool` using the connection string.
2. Calls `pool.Ping(ctx)` to check that PostgreSQL is reachable.
3. Executes `CREATE TABLE IF NOT EXISTS employees (...)`.
4. Returns the ready pool.

If the ping or schema statement fails, it closes the pool before returning the
error. The schema contains:

| Column | Meaning |
|---|---|
| `id` | Database-generated `BIGINT` primary key. |
| `name` | Required employee name. |
| `email` | Required and unique email address. |
| `department` | Required department. |
| `created_at` | Timestamp set by PostgreSQL when inserted. |
| `updated_at` | Timestamp defaulted on insert and refreshed by the update query. |

`IF NOT EXISTS` makes this initialization safe to run on each application
start when the table already exists. This is a learning-project convenience;
larger applications commonly use explicit versioned migrations.

### 3.5 Dependency wiring

The application constructs the data implementation and injects it into the
handler:

```go
employeeHandler := handlers.NewEmployeeHandler(
    repository.NewPostgresEmployeeRepository(db),
)
```

The repository receives the shared pool. The handler receives the repository
through an interface. The handler therefore does not need to know about SQL or
the PostgreSQL driver.

This is dependency injection: a component is given the service it depends on
instead of constructing that dependency internally.

### 3.6 Register routes

The `http.ServeMux` is the application's router. Each registration associates
an HTTP method and URL pattern with a handler:

| Route registration | Handler | Purpose |
|---|---|---|
| `GET /healthz` | Inline function | Check database connectivity. |
| `GET /api/employees` | `employeeHandler.List` | List employees. |
| `POST /api/employees` | `employeeHandler.Create` | Create an employee. |
| `GET /api/employees/{id}` | `employeeHandler.Get` | Fetch one employee. |
| `PUT /api/employees/{id}` | `employeeHandler.Update` | Replace editable fields. |
| `DELETE /api/employees/{id}` | `employeeHandler.Delete` | Delete an employee. |

The `"{id}"` part is a Go 1.22+ `ServeMux` path wildcard. In the handler it is
available as `r.PathValue("id")`.

The explicit HTTP method is part of each route. For example, `GET` and `POST`
to `/api/employees` reach different handlers.

### 3.7 HTTP server timeouts

The `http.Server` is configured with:

- `ReadHeaderTimeout`: maximum time to receive request headers.
- `ReadTimeout`: maximum time to read a request.
- `WriteTimeout`: maximum time to write a response.
- `IdleTimeout`: maximum time an idle keep-alive connection remains open.

These settings help prevent slow connections from holding server resources
indefinitely. They are HTTP server settings; they are not a dedicated timeout
for every SQL operation.

### 3.8 Start the server goroutine and channel

The server is launched in a goroutine:

```go
serverErr := make(chan error, 1)

go func() {
    log.Printf("HTTP server listening on %s", cfg.HTTPAddr)
    serverErr <- server.ListenAndServe()
}()
```

#### What is a goroutine?

A goroutine is a lightweight unit of execution managed by the Go runtime. The
`go` keyword starts the function asynchronously, allowing `Execute` to continue
running while `ListenAndServe` waits for HTTP connections.

Without a separate goroutine, `ListenAndServe` would occupy the current flow
and the code could not proceed to the `select` that also listens for shutdown
signals.

#### What is the channel doing?

`serverErr` is a channel whose values are of type `error`. It carries the
eventual return value from `ListenAndServe` back to the main `Execute`
goroutine.

It is **buffered with capacity 1**. The server goroutine can send its one
result even if `Execute` has not yet started receiving from the channel. This
avoids leaving the server goroutine stuck on that send during a brief scheduling
gap.

This channel is not used to send HTTP requests or database data. Its sole
purpose is to communicate the server's eventual exit/error.

### 3.9 Wait for server exit or an operating-system signal

Before starting the server, the code creates a signal-aware context:

```go
ctx, stop := signal.NotifyContext(
    context.Background(),
    os.Interrupt,
    syscall.SIGTERM,
)
defer stop()
```

The context is canceled when the process receives Ctrl+C (`os.Interrupt`) or a
termination signal (`SIGTERM`). `stop` is deferred to unregister the signal
notification when `Execute` returns.

The `select` then waits for one of two events:

```go
select {
case err := <-serverErr:
    // The HTTP server returned.
case <-ctx.Done():
    // A shutdown signal canceled the context.
}
```

`select` blocks until one of its channel operations is ready. Here, the
operations are:

1. Receive the server's return value from `serverErr`.
2. Receive from `ctx.Done()`, which becomes ready when the signal context is
   canceled.

If `ListenAndServe` returns an error other than `http.ErrServerClosed`, the app
wraps and returns it. `ErrServerClosed` is the expected result when the server
is shut down intentionally.

On a shutdown signal, the app calls `server.Shutdown` with a ten-second
deadline. Shutdown stops accepting new work and waits for active HTTP requests
to complete. If the deadline expires or shutdown otherwise fails, the error is
returned. After a successful shutdown, `Execute` receives the server goroutine's
result from `serverErr` and logs that the server stopped.

Because `db.Close()` was deferred earlier, the database pool closes after this
server lifecycle is complete.

### Startup/lifecycle sequence

```text
main
 └─ cmd.Execute
     ├─ config.Load
     ├─ database.Open (pool, ping, create table)
     ├─ construct repository and handler
     ├─ register routes and configure http.Server
     ├─ create signal context and serverErr channel
     ├─ go ListenAndServe ───────────────┐
     └─ select: serverErr or signal      │
          ├─ server error/exit           │
          └─ signal → graceful Shutdown  │
                                         ┘
     └─ deferred cleanup: stop signal listener, close DB pool
```

## 4. Concurrency and “non-blocking I/O”

It is useful to be precise when describing this project as non-blocking:

- The application starts `ListenAndServe` in its own goroutine so startup code
  can concurrently wait for a shutdown signal.
- Go's `net/http` server is designed to serve multiple connections/requests
  concurrently. It arranges handler execution so one request waiting on I/O
  does not ordinarily prevent unrelated requests from being handled.
- The `pgxpool` can manage multiple database connections, subject to its pool
  configuration and available PostgreSQL capacity.
- Each handler passes its request context to repository methods, and each
  repository method passes that context into `pgx` calls. If the request is
  canceled, database work can observe that cancellation.
- A call such as `QueryRow` is still written synchronously: that handler
  goroutine waits for its result. The Go runtime can schedule other goroutines
  while network I/O is waiting; the code does not start a goroutine per SQL
  statement or use an async/await API.

Thus, this is a concurrent HTTP server using Go's I/O/runtime model, rather
than code that makes every operation explicitly asynchronous. The request
context also does not currently have an application-defined per-query
`WithTimeout`; the HTTP server timeouts and client/request cancellation should
not be mistaken for an explicit SQL deadline.

## 5. Employee model and input validation

`internal/models/employee.go` defines two structs.

### `Employee`

`Employee` is the complete data returned by the API:

- `ID`, `CreatedAt`, and `UpdatedAt` are generated or controlled by the
  application/database.
- `Name`, `Email`, and `Department` are editable fields.
- The struct's `json` tags control the JSON property names, such as
  `CreatedAt` becoming `created_at`.

### `EmployeeInput`

`EmployeeInput` represents only fields clients may provide for create/update.
Separating it from `Employee` prevents clients from supplying an ID or
timestamps in the request body.

`Normalize` trims leading/trailing whitespace from the three strings and
returns the normalized struct. `Validate` checks that all three fields are
present and that the email contains `@`. This is a deliberately basic check,
not full email-address verification.

## 6. Repository: separating data access

`internal/repository/postgres_employee.go` contains both the repository
interface and its PostgreSQL implementation, which is a reasonable layout for
this small project.

### The interface

`EmployeeRepository` lists the operations the handler needs: `List`, `Get`,
`Create`, `Update`, and `Delete`. Each method accepts a `context.Context` so
request cancellation can flow down to database calls.

Go interfaces are satisfied implicitly. `PostgresEmployeeRepository`
implements the interface because it has methods with the required names and
signatures; there is no `implements` declaration.

The package-level sentinel errors `ErrNotFound` and `ErrConflict` give
repository operations stable meanings:

- `ErrNotFound`: no row matched the requested ID.
- `ErrConflict`: an email violated the database's unique constraint.

The handler uses these errors to choose HTTP status codes without depending on
PostgreSQL-specific error details.

### The PostgreSQL implementation

`PostgresEmployeeRepository` holds a pointer to the shared `pgxpool.Pool`.
Its constructor receives that pool from `cmd.Execute`.

The shared `employeeColumns` SQL fragment selects employee columns in the
expected scan order and formats timestamps as UTC strings. Every SQL operation
uses the `employees` table.

| Method | SQL behavior | Result |
|---|---|---|
| `List` | `SELECT ... ORDER BY id` | Reads every row into a slice. |
| `Get` | `SELECT ... WHERE id = $1` | Reads one row; maps no row to `ErrNotFound`. |
| `Create` | `INSERT ... RETURNING ...` | Inserts fields and returns the generated row. |
| `Update` | `UPDATE ... WHERE id = $4 RETURNING ...` | Replaces editable fields and refreshes `updated_at`. |
| `Delete` | `DELETE ... WHERE id = $1` | Checks affected-row count; zero means not found. |

The `$1`, `$2`, etc. SQL placeholders are parameterized values. The user-provided
values are sent separately from the SQL text, rather than concatenated into it.
This is important for correctness and helps prevent SQL injection.

`List` calls `Query` because it can return many rows. It closes the rows with
`defer`, loops through them with `Next`, scans each one, and checks `rows.Err`
after the loop for errors that occurred during iteration.

`Get`, `Create`, and `Update` use `QueryRow` because each should return at most
one row. `RETURNING` lets create/update return the newly written database
record without a second query.

`Delete` uses `Exec` because it does not need row contents. It checks
`RowsAffected` to distinguish a successful deletion from a missing ID.

The `scanEmployee` helper takes a small `rowScanner` interface. Both a single
row and an item from a multi-row result provide the `Scan` method, so the helper
can populate the same struct in both cases. Scan destination order must match
the selected column order.

`isUniqueViolation` uses `errors.As` to find a PostgreSQL error in the error
chain and checks code `23505` (unique-constraint violation). Create and update
translate that database error into `ErrConflict`.

## 7. HTTP handlers: request and response boundary

`EmployeeHandler` stores an `EmployeeRepository`. Its methods accept
`http.ResponseWriter` and `*http.Request`, the standard Go HTTP handler inputs:

- `ResponseWriter` is where the handler sets headers, status, and response
  body.
- `Request` contains the method, URL, body, path values, and request context.

The handler's job is HTTP-specific: parse IDs and JSON, validate request data,
call the repository, and map results/errors to HTTP responses. It should not
contain SQL.

### JSON decoding

`decodeInput` is shared by create and update:

1. `http.MaxBytesReader` caps the body at 1 MiB.
2. `json.NewDecoder` reads JSON from the request body.
3. `DisallowUnknownFields` rejects fields not present in `EmployeeInput`.
4. The first decode parses one employee input object.
5. A second decode must return `io.EOF`, ensuring there is no second JSON value
   or trailing extra object.
6. The input is normalized and validated.

Malformed JSON, unknown fields, oversized bodies, or invalid employee values
result in `400 Bad Request`.

### ID parsing

`employeeID` gets the `id` wildcard with `r.PathValue("id")`, then parses it as
a base-10 signed 64-bit integer. It rejects invalid, zero, and negative IDs with
`400 Bad Request`.

### Error mapping and JSON helpers

`writeRepositoryError` maps repository errors as follows:

| Error/result | HTTP status |
|---|---:|
| Invalid ID or request body | `400 Bad Request` |
| Employee does not exist | `404 Not Found` |
| Duplicate email | `409 Conflict` |
| Unexpected repository failure | `500 Internal Server Error` |

`writeError` encodes errors as JSON, for example:

```json
{"error":"employee not found"}
```

`writeJSON` sets `Content-Type` to JSON, writes the status code, and encodes
the supplied Go value as a JSON response body.

### Endpoint behavior

| Endpoint | Success response |
|---|---|
| `GET /api/employees` | `200 OK`, JSON array (empty array if no rows). |
| `POST /api/employees` | `201 Created`, created employee JSON, and a `Location` header. |
| `GET /api/employees/{id}` | `200 OK` and employee JSON. |
| `PUT /api/employees/{id}` | `200 OK` and updated employee JSON. |
| `DELETE /api/employees/{id}` | `204 No Content`, no response body. |

`PUT` replaces all editable fields, so `name`, `email`, and `department` are all
required. A partial update would normally be designed as a separate `PATCH`
operation with different validation rules.

## 8. Example request traced end to end

For this request:

```http
POST /api/employees
Content-Type: application/json
```

```json
{
  "name": "Alex Morgan",
  "email": "alex@example.com",
  "department": "Engineering"
}
```

the execution is:

1. The HTTP server receives the request and its `ServeMux` matches
   `POST /api/employees`.
2. The router invokes `EmployeeHandler.Create`.
3. `decodeInput` limits and decodes the JSON body.
4. `Normalize` trims whitespace; `Validate` checks required fields and the
   basic email condition.
5. The handler calls `EmployeeRepository.Create` with the request context.
6. The PostgreSQL implementation runs a parameterized `INSERT`.
7. PostgreSQL generates the ID and timestamps, then returns the inserted row.
8. The repository scans the row into `models.Employee`.
9. The handler sets `Location: /api/employees/{id}` and responds with
   `201 Created` and the employee JSON.

If the email already exists, PostgreSQL returns a unique-constraint error; the
repository maps it to `ErrConflict`, and the handler responds with `409
Conflict`. If any unexpected database error occurs, the handler responds with
`500 Internal Server Error`.

## 9. Health check

`GET /healthz` is registered directly in `cmd/cmd.go`. It creates a context
with a two-second timeout and calls `db.Ping`.

- If PostgreSQL responds, it returns `204 No Content`.
- If the ping fails or times out, it returns `503 Service Unavailable` with
  `database unavailable`.

This route checks the database too, so it indicates more than merely that the
HTTP process is running.

## 10. Run and try the API

Start PostgreSQL, set its connection URL, then run the service from the project
root:

```powershell
$env:DATABASE_URL = "postgres://postgres:your-password@localhost:5432/employees?sslmode=disable"
go run .
```

The HTTP address defaults to `localhost:8080`. To override it:

```powershell
$env:HTTP_ADDR = "localhost:8080"
```

Create an employee:

```powershell
Invoke-RestMethod -Method Post `
  -Uri http://localhost:8080/api/employees `
  -ContentType "application/json" `
  -Body '{"name":"Alex Morgan","email":"alex@example.com","department":"Engineering"}'
```

List employees:

```powershell
Invoke-RestMethod -Method Get -Uri http://localhost:8080/api/employees
```

Fetch one:

```powershell
Invoke-RestMethod -Method Get -Uri http://localhost:8080/api/employees/1
```

Update one (all three editable fields are required):

```powershell
Invoke-RestMethod -Method Put `
  -Uri http://localhost:8080/api/employees/1 `
  -ContentType "application/json" `
  -Body '{"name":"Alex Morgan","email":"alex@example.com","department":"Platform"}'
```

Delete one:

```powershell
Invoke-RestMethod -Method Delete -Uri http://localhost:8080/api/employees/1
```

Check application/database health:

```powershell
Invoke-WebRequest -Method Get -Uri http://localhost:8080/healthz
```

## 11. Design choices and current scope

- **Small entry point:** `main.go` only delegates to the application setup.
- **Separated responsibilities:** HTTP code, SQL code, configuration, and data
  structures are kept in focused packages.
- **Interface at the repository boundary:** handlers can use another
  implementation later, such as a fake repository in unit tests.
- **Request contexts:** contexts flow from handlers into database operations.
- **Parameterized SQL:** values are not concatenated into SQL statements.
- **Database constraints:** PostgreSQL enforces required fields and email
  uniqueness even if a future caller bypasses HTTP validation.
- **Graceful shutdown:** active requests have an opportunity to complete before
  the process closes the database pool.

This is a teaching/sample application, not a fully production-hardened service.
It currently has no authentication/authorization, pagination, structured
request logging, rate limiting, versioned migration system, or dedicated
per-query timeout. Email validation is intentionally basic. Those would be
separate design decisions for a production deployment.

## 12. Useful checks

From the project root:

```powershell
go test ./...
go vet ./...
```

`go test ./...` compiles and runs tests in all packages. If a package has no
test files, Go reports `[no test files]`; that still confirms the package
compiled. `go vet ./...` checks for common suspicious code patterns. These
commands do not replace integration tests against a running PostgreSQL
instance.
