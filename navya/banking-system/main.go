package main

import (
	"banking-system/bank"
	"fmt"
)

func main() {
	account := bank.Account{
		Name:    "Navya",
		Balance: 1000,
	}

	fmt.Println("Account Holder:", account.Name)
	fmt.Println("Initial Balance:", account.Balance)

	account.Deposit(500)
	fmt.Println("After Deposit:", account.Balance)

	err := account.Withdraw(200)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("After Withdrawal:", account.Balance)
	}
}
