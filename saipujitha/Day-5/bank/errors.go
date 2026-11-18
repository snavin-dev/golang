package bank

import "errors"

var (
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type OverdraftError struct {
	AccountNumber int
	Requested     float64
	Available     float64
}

func (e *OverdraftError) Error() string {
	return "withdrawal exceeds the available balance"
}

func (e *OverdraftError) Unwrap() error {
	return ErrInsufficientFunds
}
