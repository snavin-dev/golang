package account

import "errors"

var (
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type Account struct {
	Holder  string
	Balance float64
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if amount > a.Balance {
		return ErrInsufficientFunds
	}

	a.Balance -= amount
	return nil
}
