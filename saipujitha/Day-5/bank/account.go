package bank

import (
	"fmt"
)

type Account struct {
	AccountHolder string
	AccountNumber int
	balance       float64
}

func NewAccount(holder string, accountNumber int, openingBalance float64) (*Account, error) {
	if openingBalance < 0 {
		return nil, fmt.Errorf("create account: opening balance cannot be negative: %w", ErrInvalidAmount)
	}

	return &Account{
		AccountHolder: holder,
		AccountNumber: accountNumber,
		balance:       openingBalance,
	}, nil
}

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
		overdraft := &OverdraftError{
			AccountNumber: a.AccountNumber,
			Requested:     amount,
			Available:     a.balance,
		}
		return fmt.Errorf("withdraw from account %d: %w", a.AccountNumber, overdraft)
	}

	a.balance -= amount
	return nil
}

func (a Account) String() string {
	return fmt.Sprintf(
		"Account Holder: %s | Account Number: %d | Balance: %.2f",
		a.AccountHolder,
		a.AccountNumber,
		a.balance,
	)
}
