
package main

import (
	"fmt"
)

// 1. STRUCT
type Account struct {
	AccountHolder string
	AccountNumber int
	Balance       float64
}

// 2. VALUE RECEIVER METHOD
func (a Account) DisplayDetails() {
	fmt.Println("\n--- Account Details ---")
	fmt.Println("Holder:", a.AccountHolder)
	fmt.Println("Account Number:", a.AccountNumber)
	fmt.Println("Balance:", a.Balance)
}

// 3. POINTER RECEIVER METHOD - DEPOSIT
func (a *Account) Deposit(amount float64) {
	if amount <= 0 {
		fmt.Println("Invalid deposit amount!")
		return
	}

	a.Balance += amount
	fmt.Println("Deposited:", amount)
	fmt.Println("Updated Balance:", a.Balance)
}

// 4. POINTER RECEIVER METHOD - WITHDRAW
func (a *Account) Withdraw(amount float64) {
	if amount <= 0 {
		fmt.Println("Invalid withdrawal amount!")
		return
	}

	if amount > a.Balance {
		fmt.Println("Insufficient Balance!")
		return
	}

	a.Balance -= amount
	fmt.Println("Withdrawn:", amount)
	fmt.Println("Remaining Balance:", a.Balance)
}

// 5. EMBEDDING
type Customer struct {
	Name  string
	Email string
}

type BankCustomer struct {
	Customer // Embedded struct
	Account
}

// 6. INTERFACE
type Stringer interface {
	String() string
}

// Implementing Stringer interface
func (a Account) String() string {
	return fmt.Sprintf(
		"Account Holder: %s | Balance: %.2f",
		a.AccountHolder,
		a.Balance,
	)
}

func main() {

	// Creating struct literal
	account := Account{
		AccountHolder: "Pujitha",
		AccountNumber: 12345678,
		Balance:       5000,
	}

	// Display account details
	account.DisplayDetails()

	// Deposit money
	fmt.Println("\n Deposit ")
	account.Deposit(2000)

	// Withdraw money
	fmt.Println("\n Withdrawal ")
	account.Withdraw(1500)

	// Invalid withdrawal
	fmt.Println("\n Invalid Withdrawal ")
	account.Withdraw(10000)

	// Embedding example
	customer := BankCustomer{
		Customer: Customer{
			Name:  "Pujitha",
			Email: "pujitha@gmail.com",
		},
		Account: account,
	}

	fmt.Println("\n Bank Customer ")
	fmt.Println("Name:", customer.Name)
	fmt.Println("Email:", customer.Email)
	customer.DisplayDetails()

	// Interface example
	var s fmt.Stringer = account

	fmt.Println("\n Interface Output ")
	fmt.Println(s.String())
}
