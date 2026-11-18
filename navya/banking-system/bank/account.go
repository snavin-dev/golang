
package bank

import "errors"

type Account struct {
    Name    string
    Balance float64
}

func (a *Account) Deposit(amount float64) {
    if amount > 0 {
        a.Balance += amount
    }
}

func (a *Account) Withdraw(amount float64) error {
    if amount <= 0 {
        return errors.New("invalid amount")
    }

    if amount > a.Balance {
        return errors.New("insufficient balance")
    }

    a.Balance -= amount
    return nil
}
