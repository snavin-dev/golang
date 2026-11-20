package account

import (
	"errors"
	"testing"
)

func TestWithdrawOverdraftReturnsErrorWithoutChangingBalance(t *testing.T) {
	a := Account{Holder: "Alex", Balance: 100}

	err := a.Withdraw(101)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf("Withdraw() error = %v, want %v", err, ErrInsufficientFunds)
	}
	if a.Balance != 100 {
		t.Errorf("balance after rejected withdrawal = %v, want 100", a.Balance)
	}
}

func TestWithdrawEntireBalance(t *testing.T) {
	a := Account{Holder: "Alex", Balance: 100}

	if err := a.Withdraw(100); err != nil {
		t.Fatalf("Withdraw() error = %v, want nil", err)
	}
	if a.Balance != 0 {
		t.Errorf("balance after withdrawal = %v, want 0", a.Balance)
	}
}
