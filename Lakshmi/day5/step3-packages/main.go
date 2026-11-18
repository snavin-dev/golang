package main

import (
	"errors"
	"fmt"

	"bank/account" // module name "bank" + folder "account"
	"bank/fileops" // module name "bank" + folder "fileops"
)

const balanceFile = "balance.txt"

func main() {
	balance, err := fileops.ReadFloat(balanceFile)
	if err != nil {
		fmt.Println("Note:", err)
		fmt.Println("Starting with a balance of 1000.")
		balance = 1000
	}

	acc := account.New("Asha", balance)
	// acc.balance = 1000000 // would not compile: balance is unexported

	fmt.Println("Welcome to Go Bank!")

	for {
		showMenu()

		var choice int
		fmt.Print("Your choice: ")
		if _, err := fmt.Scan(&choice); err != nil {
			fmt.Println("Please enter a number. Goodbye!")
			return
		}

		switch choice {
		case 1:
			fmt.Println(acc)
		case 2:
			var amount float64
			fmt.Print("Amount to deposit: ")
			fmt.Scan(&amount)
			if err := acc.Deposit(amount); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println("Deposited. New", acc)
		case 3:
			var amount float64
			fmt.Print("Amount to withdraw: ")
			fmt.Scan(&amount)
			if err := acc.Withdraw(amount); err != nil {
				showWithdrawError(err)
				continue
			}
			fmt.Println("Withdrawn. New", acc)
		case 4:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Pick 1 to 4.")
			continue
		}

		if err := fileops.WriteFloat(balanceFile, acc.Balance()); err != nil {
			fmt.Println("Error saving balance:", err)
		}
	}
}

// showWithdrawError uses errors.Is and errors.As to give a helpful message.
func showWithdrawError(err error) {
	var fundsErr *account.InsufficientFundsError

	switch {
	case errors.Is(err, account.ErrInvalidAmount):
		fmt.Println("Please enter an amount greater than 0.")
	case errors.As(err, &fundsErr):
		short := fundsErr.Requested - fundsErr.Balance
		fmt.Printf("You are short by $%.2f.\n", short)
	default:
		fmt.Println("Error:", err)
	}
	fmt.Println("Details:", err)
}
