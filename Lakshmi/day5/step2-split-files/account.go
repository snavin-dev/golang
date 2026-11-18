package main

import (
	"errors"
	"fmt"
)

type Account struct {
	Owner   string
	Balance float64
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if amount > a.Balance {
		return errors.New("not enough money in the account")
	}
	a.Balance -= amount
	return nil
}

func (a Account) String() string {
	return fmt.Sprintf("%s's balance: $%.2f", a.Owner, a.Balance)
}
