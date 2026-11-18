// Package account holds the bank account and its rules.
package account

import (
	"errors"
	"fmt"
)

// ErrInvalidAmount is a sentinel error. Callers check it with errors.Is.
var ErrInvalidAmount = errors.New("amount must be greater than 0")

// InsufficientFundsError is a custom error. Callers read it with errors.As.
type InsufficientFundsError struct {
	Balance   float64
	Requested float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("not enough money: balance is $%.2f, you asked for $%.2f",
		e.Balance, e.Requested)
}

// Account starts with a capital A, so other packages can use it.
type Account struct {
	Owner   string  // exported: other packages can read and change it
	balance float64 // unexported: only this package can change it
}

// New creates an account. Other packages must use it, because they
// cannot set the unexported balance field themselves.
func New(owner string, balance float64) *Account {
	return &Account{Owner: owner, balance: balance}
}

// Balance lets other packages read the balance.
func (a *Account) Balance() float64 {
	return a.balance
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit: %w", ErrInvalidAmount)
	}
	a.balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdraw: %w", ErrInvalidAmount)
	}
	if amount > a.balance {
		return fmt.Errorf("withdraw: %w",
			&InsufficientFundsError{Balance: a.balance, Requested: amount})
	}
	a.balance -= amount
	return nil
}

func (a *Account) String() string {
	return fmt.Sprintf("%s's balance: $%.2f", a.Owner, a.balance)
}
