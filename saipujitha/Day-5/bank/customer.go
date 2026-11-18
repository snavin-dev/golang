package bank

type Customer struct {
	Name  string
	Email string
}

type BankCustomer struct {
	Customer
	Account
}
