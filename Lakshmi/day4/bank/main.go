package main

import "fmt"

// Account is a struct: it groups the owner and the balance together
type Account struct {
	Owner   string
	Balance float64
}

// Deposit uses a pointer receiver because it changes the balance
func (a *Account) Deposit(amount float64) bool {
	if amount <= 0 {
		return false
	}
	a.Balance += amount
	return true
}

// Withdraw uses a pointer receiver because it changes the balance
func (a *Account) Withdraw(amount float64) bool {
	if amount <= 0 || amount > a.Balance {
		return false
	}
	a.Balance -= amount
	return true
}

// String makes Account a fmt.Stringer, so fmt.Println(acc) prints this
func (a Account) String() string {
	return fmt.Sprintf("%s's balance: $%.2f", a.Owner, a.Balance)
}

func showMenu() {
	fmt.Println()
	fmt.Println("What do you want to do?")
	fmt.Println("1. Check balance")
	fmt.Println("2. Deposit money")
	fmt.Println("3. Withdraw money")
	fmt.Println("4. Exit")
}

func main() {
	acc := &Account{Owner: "Asha", Balance: 1000}
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
			if acc.Deposit(amount) {
				fmt.Println("Deposited. New", acc)
			} else {
				fmt.Println("Invalid amount. It must be greater than 0.")
			}
		case 3:
			var amount float64
			fmt.Print("Amount to withdraw: ")
			fmt.Scan(&amount)
			if acc.Withdraw(amount) {
				fmt.Println("Withdrawn. New", acc)
			} else {
				fmt.Println("Cannot withdraw that amount.", acc)
			}
		case 4:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Pick 1 to 4.")
		}
	}
}
