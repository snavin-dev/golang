package main

import (
	"errors"
	"fmt"
)

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

func main() {
	account := Account{Holder: "Alex", Balance: 100}

	if err := account.Deposit(50); err != nil {
		fmt.Println("Deposit failed:", err)
		return
	}
	if err := account.Withdraw(30); err != nil {
		fmt.Println("Withdrawal failed:", err)
		return
	}

	fmt.Printf("%s's balance: $%.2f\n", account.Holder, account.Balance)
}