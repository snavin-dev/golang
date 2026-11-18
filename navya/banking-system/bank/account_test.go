
package bank

import "testing"

func TestWithdraw(t *testing.T) {
    account := Account{
        Name:    "Navya",
        Balance: 1000,
    }

    err := account.Withdraw(500) 

    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }

    if account.Balance != 500 {
        t.Errorf("Expected balance 500, got %v", account.Balance)
    }
}

func TestOverdraft(t *testing.T) {
    account := Account{
        Name:    "Navya",
        Balance: 1000,
    }

    err := account.Withdraw(1500)

    if err == nil {
        t.Error("Expected overdraft error, but got nil")
    }

    if account.Balance != 1000 {
        t.Errorf("Balance should remain 1000, got %v", account.Balance)
    }
}
