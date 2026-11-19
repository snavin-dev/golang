# Employee CRUD API

A small PostgreSQL-backed HTTP API using Go's standard `net/http` server and `pgx`.
Go serves HTTP requests concurrently; each database operation uses the request
context so cancelled or timed-out requests do not continue waiting on the database.

## Requirements

- Go 1.22 or newer
- PostgreSQL

Set `DATABASE_URL` to a PostgreSQL connection string. `HTTP_ADDR` is optional and
defaults to `localhost:8080`. This binds the server to the local machine.

PowerShell example:

```powershell
$env:DATABASE_URL = "postgres://postgres:password@localhost:5432/employees?sslmode=disable"
$env:HTTP_ADDR = "localhost:8080"
go run .
```

The application creates the `employees` table at startup if it does not exist.

## Endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/healthz` | Check database connectivity |
| `GET` | `/api/employees` | List employees |
| `POST` | `/api/employees` | Create an employee |
| `GET` | `/api/employees/{id}` | Get an employee |
| `PUT` | `/api/employees/{id}` | Replace an employee's editable fields |
| `DELETE` | `/api/employees/{id}` | Delete an employee |

Create/update JSON:

```json
{
  "name": "Alex Morgan",
  "email": "alex@example.com",
  "department": "Engineering"
}
```

Example request:

```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/employees `
  -ContentType "application/json" `
  -Body '{"name":"Alex Morgan","email":"alex@example.com","department":"Engineering"}'
```

## Project layout

```text
crud-operations/
├── main.go
├── cmd/
│   └── cmd.go
└── internal/
    ├── config/
    ├── database/
    ├── handlers/
    ├── models/
    └── repository/
```
