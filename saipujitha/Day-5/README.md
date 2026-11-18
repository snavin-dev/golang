# Day 5: Error Handling and Packages

This example refactors the Day 4 banking app into a `bank` package and a
command-line program. Account operations return errors; overdrafts use a custom
`OverdraftError` that can be inspected with `errors.As` and matched to
`ErrInsufficientFunds` with `errors.Is`.

Run the app and its unit tests from this directory:

```text
go run .
go test ./...
```
