package main

import "fmt"

type Account struct {
	Name    string
	Balance float64
}

func (a *Account) Deposit(amount float64) {
	if amount > 0 {
		a.Balance += amount
		fmt.Println("Amount deposited:", amount)
	} else {
		fmt.Println("Invalid deposit amount")
	}
}

func (a *Account) Withdraw(amount float64) {
	if amount <= 0 {
		fmt.Println("Invalid withdrawal amount")
	} else if amount > a.Balance {
		fmt.Println("Insufficient balance")
	} else {
		a.Balance -= amount
		fmt.Println("Amount withdrawn:", amount)
	}
}

func main() {
	account := Account{
		Name:    "Navya",
		Balance: 1000,
	}

	fmt.Println("Account Holder:", account.Name)
	fmt.Println("Initial Balance:", account.Balance)

	account.Deposit(500)

	account.Withdraw(200)

	fmt.Println("Final Balance:", account.Balance)
}
