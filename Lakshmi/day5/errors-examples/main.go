// This program shows every feature of Go's errors package
// (https://pkg.go.dev/errors), using bank examples.
// Each example is its own function so you can explain them one by one.
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func main() {
	example1ErrorInterface()
	example2New()
	example3NewIsUnique()
	example4Errorf()
	example5WrapWithW()
	example6WrapWithV()
	example7Unwrap()
	example8Is()
	example9CustomErrorType()
	example10As()
	example11AsType()
	example12CustomUnwrap()
	example13CustomIs()
	example14CustomAs()
	example15Join()
	example16MultipleW()
	example17ErrUnsupported()
	example18StandardLibrary()
}

func title(name string) {
	fmt.Println()
	fmt.Println("=====", name, "=====")
}

// ---------------------------------------------------------------
// 1. The error interface: any type with an Error() string method
// ---------------------------------------------------------------

func checkPIN(pin string) error {
	if pin != "1234" {
		return errors.New("wrong PIN")
	}
	return nil // nil means "no error"
}

func example1ErrorInterface() {
	title("1. error interface and nil")

	if err := checkPIN("0000"); err != nil {
		fmt.Println("Error:", err)
	}
	if err := checkPIN("1234"); err == nil {
		fmt.Println("PIN accepted, err is nil")
	}
}

// ---------------------------------------------------------------
// 2. errors.New: make an error from a text message
// ---------------------------------------------------------------

func example2New() {
	title("2. errors.New")

	err := errors.New("account is locked")
	fmt.Println(err)
}

// ---------------------------------------------------------------
// 3. Each errors.New call is a different error, even with the same text.
//    That is why we create sentinel errors once and reuse them.
// ---------------------------------------------------------------

var ErrAccountLocked = errors.New("account is locked") // sentinel error

func example3NewIsUnique() {
	title("3. errors.New makes unique values (sentinel errors)")

	a := errors.New("account is locked")
	b := errors.New("account is locked")
	fmt.Println("two New calls, same text, equal?", errors.Is(a, b))

	err := ErrAccountLocked
	fmt.Println("same sentinel variable, equal?  ", errors.Is(err, ErrAccountLocked))
}

// ---------------------------------------------------------------
// 4. fmt.Errorf: make an error with values in the message
// ---------------------------------------------------------------

func example4Errorf() {
	title("4. fmt.Errorf (formatted message)")

	owner, id := "Asha", 1001
	err := fmt.Errorf("account %q (id %d) not found", owner, id)
	fmt.Println(err)
}

// ---------------------------------------------------------------
// 5. Wrapping with %w: add context and keep the original error inside
// ---------------------------------------------------------------

var ErrInsufficientFunds = errors.New("insufficient funds")

func withdraw(balance, amount float64) error {
	if amount > balance {
		return ErrInsufficientFunds
	}
	return nil
}

func payBill(balance, bill float64) error {
	if err := withdraw(balance, bill); err != nil {
		return fmt.Errorf("pay electricity bill: %w", err)
	}
	return nil
}

func example5WrapWithW() {
	title("5. wrapping with %w")

	err := payBill(100, 250)
	fmt.Println("message:", err)
	fmt.Println("still insufficient funds inside?", errors.Is(err, ErrInsufficientFunds))
}

// ---------------------------------------------------------------
// 6. %v copies only the text. The original error is lost.
// ---------------------------------------------------------------

func example6WrapWithV() {
	title("6. %v instead of %w (original error lost)")

	err := fmt.Errorf("pay electricity bill: %v", ErrInsufficientFunds)
	fmt.Println("message:", err)
	fmt.Println("still insufficient funds inside?", errors.Is(err, ErrInsufficientFunds))
}

// ---------------------------------------------------------------
// 7. errors.Unwrap: take off one layer of wrapping
// ---------------------------------------------------------------

func example7Unwrap() {
	title("7. errors.Unwrap")

	inner := ErrInsufficientFunds
	middle := fmt.Errorf("withdraw: %w", inner)
	outer := fmt.Errorf("pay bill: %w", middle)

	fmt.Println("outer:          ", outer)
	fmt.Println("unwrap once:    ", errors.Unwrap(outer))
	fmt.Println("unwrap twice:   ", errors.Unwrap(errors.Unwrap(outer)))
	fmt.Println("unwrap 3 times: ", errors.Unwrap(errors.Unwrap(errors.Unwrap(outer))))
}

// ---------------------------------------------------------------
// 8. errors.Is: is this error (or anything wrapped inside) equal to X?
//    Use it instead of ==, because == does not look inside wrapped errors.
// ---------------------------------------------------------------

func example8Is() {
	title("8. errors.Is vs ==")

	err := payBill(100, 250) // wrapped

	fmt.Println("err == ErrInsufficientFunds:        ", err == ErrInsufficientFunds)
	fmt.Println("errors.Is(err, ErrInsufficientFunds):", errors.Is(err, ErrInsufficientFunds))
	fmt.Println("errors.Is(err, ErrAccountLocked):    ", errors.Is(err, ErrAccountLocked))
}

// ---------------------------------------------------------------
// 9. A custom error type: a struct that carries extra details
// ---------------------------------------------------------------

type InsufficientFundsError struct {
	Balance   float64
	Requested float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("balance $%.2f is less than $%.2f", e.Balance, e.Requested)
}

func withdrawDetailed(balance, amount float64) error {
	if amount > balance {
		return fmt.Errorf("withdraw: %w",
			&InsufficientFundsError{Balance: balance, Requested: amount})
	}
	return nil
}

func example9CustomErrorType() {
	title("9. custom error type")

	err := withdrawDetailed(100, 250)
	fmt.Println(err)
}

// ---------------------------------------------------------------
// 10. errors.As: find an error of a certain TYPE and get its details
// ---------------------------------------------------------------

func example10As() {
	title("10. errors.As")

	err := withdrawDetailed(100, 250)

	var fundsErr *InsufficientFundsError // 1. variable of the type we want
	if errors.As(err, &fundsErr) {       // 2. pass its address
		fmt.Println("balance:  ", fundsErr.Balance) // 3. use the details
		fmt.Println("requested:", fundsErr.Requested)
		fmt.Println("short by: ", fundsErr.Requested-fundsErr.Balance)
	}
}

// ---------------------------------------------------------------
// 11. errors.AsType (Go 1.26+): same as As, but returns the error.
//     The Go docs say: for most uses, prefer AsType.
// ---------------------------------------------------------------

func example11AsType() {
	title("11. errors.AsType (Go 1.26+)")

	err := withdrawDetailed(100, 250)

	if fundsErr, ok := errors.AsType[*InsufficientFundsError](err); ok {
		fmt.Println("short by:", fundsErr.Requested-fundsErr.Balance)
	}

	// No match: ok is false and the value is the zero value (nil here)
	_, ok := errors.AsType[*fs.PathError](err)
	fmt.Println("is it a *fs.PathError?", ok)
}

// ---------------------------------------------------------------
// 12. A custom error type that wraps another error with Unwrap() error
// ---------------------------------------------------------------

type TransactionError struct {
	TxID string
	Err  error // the cause
}

func (e *TransactionError) Error() string {
	return "transaction " + e.TxID + " failed: " + e.Err.Error()
}

// Unwrap lets errors.Is, errors.As and errors.Unwrap look inside
func (e *TransactionError) Unwrap() error {
	return e.Err
}

func example12CustomUnwrap() {
	title("12. custom type with an Unwrap method")

	err := &TransactionError{TxID: "TX-42", Err: ErrInsufficientFunds}

	fmt.Println(err)
	fmt.Println("errors.Is finds the cause:", errors.Is(err, ErrInsufficientFunds))
	fmt.Println("errors.Unwrap gives:      ", errors.Unwrap(err))
}

// ---------------------------------------------------------------
// 13. A custom Is method: decide yourself what counts as a match
// ---------------------------------------------------------------

var ErrDeclined = errors.New("payment declined")

type CardError struct {
	Code string
}

func (e *CardError) Error() string {
	return "card error " + e.Code
}

// Is makes every CardError count as ErrDeclined.
// The docs say an Is method should only compare, never call Unwrap.
func (e *CardError) Is(target error) bool {
	return target == ErrDeclined
}

func example13CustomIs() {
	title("13. custom Is method")

	err := fmt.Errorf("checkout: %w", &CardError{Code: "EXPIRED"})

	fmt.Println(err)
	fmt.Println("err == ErrDeclined:        ", err == ErrDeclined)
	fmt.Println("errors.Is(err, ErrDeclined):", errors.Is(err, ErrDeclined))
}

// ---------------------------------------------------------------
// 14. A custom As method: let your error act as a different type
// ---------------------------------------------------------------

type ATMError struct {
	Machine string
	Balance float64
	Wanted  float64
}

func (e *ATMError) Error() string {
	return "ATM " + e.Machine + " refused the withdrawal"
}

// As lets errors.As / errors.AsType turn an ATMError into an InsufficientFundsError
func (e *ATMError) As(target any) bool {
	if t, ok := target.(**InsufficientFundsError); ok {
		*t = &InsufficientFundsError{Balance: e.Balance, Requested: e.Wanted}
		return true
	}
	return false
}

func example14CustomAs() {
	title("14. custom As method")

	err := fmt.Errorf("cash out: %w", &ATMError{Machine: "MG-Road", Balance: 40, Wanted: 100})

	if fundsErr, ok := errors.AsType[*InsufficientFundsError](err); ok {
		fmt.Println("got InsufficientFundsError:", fundsErr)
	}
}

// ---------------------------------------------------------------
// 15. errors.Join: combine several errors into one
// ---------------------------------------------------------------

var (
	ErrNoName   = errors.New("owner name is empty")
	ErrBadPhone = errors.New("phone number is invalid")
)

func validateNewAccount(name, phone string) error {
	var nameErr, phoneErr error
	if name == "" {
		nameErr = ErrNoName
	}
	if len(phone) != 10 {
		phoneErr = ErrBadPhone
	}
	return errors.Join(nameErr, phoneErr) // nil errors are dropped
}

func example15Join() {
	title("15. errors.Join")

	err := validateNewAccount("", "123")
	fmt.Println("message (one line per error):")
	fmt.Println(err)
	fmt.Println("has ErrNoName?   ", errors.Is(err, ErrNoName))
	fmt.Println("has ErrBadPhone? ", errors.Is(err, ErrBadPhone))
	fmt.Println("errors.Unwrap on a joined error:", errors.Unwrap(err))

	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		fmt.Println("number of errors inside:", len(joined.Unwrap()))
	}

	fmt.Println("Join with only nil errors:", errors.Join(nil, nil))
	fmt.Println("valid account gives nil:  ", validateNewAccount("Asha", "9876543210"))
}

// ---------------------------------------------------------------
// 16. fmt.Errorf with more than one %w (Go 1.20+)
// ---------------------------------------------------------------

func example16MultipleW() {
	title("16. fmt.Errorf with several %w")

	err := fmt.Errorf("open account failed: %w and %w", ErrNoName, ErrBadPhone)
	fmt.Println(err)
	fmt.Println("has ErrNoName?  ", errors.Is(err, ErrNoName))
	fmt.Println("has ErrBadPhone?", errors.Is(err, ErrBadPhone))
}

// ---------------------------------------------------------------
// 17. errors.ErrUnsupported: "this operation is not supported".
//     Do not return it bare: wrap it with context.
// ---------------------------------------------------------------

func transferAbroad(country string) error {
	if country != "India" {
		return fmt.Errorf("transfer to %s: %w", country, errors.ErrUnsupported)
	}
	return nil
}

func example17ErrUnsupported() {
	title("17. errors.ErrUnsupported")

	err := transferAbroad("Japan")
	fmt.Println(err)
	fmt.Println("unsupported?", errors.Is(err, errors.ErrUnsupported))
}

// ---------------------------------------------------------------
// 18. The standard library uses the same tools
// ---------------------------------------------------------------

func example18StandardLibrary() {
	title("18. errors from the standard library")

	_, err := os.Open("missing-statement.txt")

	fmt.Println(err)
	fmt.Println("errors.Is(err, fs.ErrNotExist):", errors.Is(err, fs.ErrNotExist))

	if pathErr, ok := errors.AsType[*fs.PathError](err); ok {
		fmt.Println("operation:", pathErr.Op)
		fmt.Println("path:     ", pathErr.Path)
	}
}
