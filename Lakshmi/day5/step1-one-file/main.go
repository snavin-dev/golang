package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const balanceFile = "balance.txt"

// ---------- account ----------

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

// ---------- saving to a file ----------

func readBalance() (float64, error) {
	data, err := os.ReadFile(balanceFile)
	if err != nil {
		return 0, fmt.Errorf("could not read %s: %w", balanceFile, err)
	}
	balance, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return 0, fmt.Errorf("could not understand the saved balance: %w", err)
	}
	return balance, nil
}

func writeBalance(balance float64) error {
	text := fmt.Sprint(balance)
	return os.WriteFile(balanceFile, []byte(text), 0644)
}

// ---------- menu ----------

func showMenu() {
	fmt.Println()
	fmt.Println("What do you want to do?")
	fmt.Println("1. Check balance")
	fmt.Println("2. Deposit money")
	fmt.Println("3. Withdraw money")
	fmt.Println("4. Exit")
}

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
