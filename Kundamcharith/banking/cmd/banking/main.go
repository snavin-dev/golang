package main

import (
	"fmt"

	"banking/account"
)

func main() {
	account := account.Account{Holder: "Alex", Balance: 100}

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
