package main

import "fmt"

// main.go only has the program flow.
// Account, readBalance, writeBalance and showMenu live in the other files,
// but they are all in package main, so main can use them directly.

func main() {
	balance, err := readBalance()
	if err != nil {
		fmt.Println("Note:", err)
		fmt.Println("Starting with a balance of 1000.")
		balance = 1000
	}

	acc := &Account{Owner: "Asha", Balance: balance}
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
				fmt.Println("Error:", err)
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

		if err := writeBalance(acc.Balance); err != nil {
			fmt.Println("Error saving balance:", err)
		}
	}
}
