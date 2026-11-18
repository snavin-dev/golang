package main

import (
	"errors"
	"fmt"

	"example.com/day5bank/bank"
)

func main() {
	account, err := bank.NewAccount("Pujitha", 12345678, 5000)
	if err != nil {
		fmt.Println("Could not create account:", err)
		return
	}

	fmt.Println("Initial account:", account)

	if err := account.Deposit(2000); err != nil {
		fmt.Println("Deposit failed:", err)
		return
	}
	fmt.Printf("After deposit: %.2f\n", account.Balance())

	if err := account.Withdraw(1500); err != nil {
		fmt.Println("Withdrawal failed:", err)
		return
	}
	fmt.Printf("After withdrawal: %.2f\n", account.Balance())

	err = account.Withdraw(10000)
	if errors.Is(err, bank.ErrInsufficientFunds) {
		var overdraft *bank.OverdraftError
		if errors.As(err, &overdraft) {
			fmt.Printf(
				"Overdraft rejected for account %d: requested %.2f, available %.2f\n",
				overdraft.AccountNumber,
				overdraft.Requested,
				overdraft.Available,
			)
		}
	}

	customer := bank.BankCustomer{
		Customer: bank.Customer{Name: "Pujitha", Email: "pujitha@example.com"},
		Account:  *account,
	}
	fmt.Printf("Customer: %s <%s>\n", customer.Name, customer.Email)
	fmt.Println("Customer account:", customer.Account)
}
