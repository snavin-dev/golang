package account

import (
	"errors"
	"testing"
)

func TestDeposit(t *testing.T) {
	acc := New("Asha", 100)

	if err := acc.Deposit(50); err != nil {
		t.Fatalf("Deposit(50) returned an error: %v", err)
	}
	if acc.Balance() != 150 {
		t.Errorf("Balance() = %v, want 150", acc.Balance())
	}
}

func TestDepositInvalidAmount(t *testing.T) {
	acc := New("Asha", 100)

	err := acc.Deposit(-5)

	if !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("Deposit(-5) error = %v, want ErrInvalidAmount", err)
	}
	if acc.Balance() != 100 {
		t.Errorf("Balance() = %v, want 100 (unchanged)", acc.Balance())
	}
}

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name        string
		start       float64
		amount      float64
		wantBalance float64
		wantErr     bool
	}{
		{"normal withdrawal", 100, 40, 60, false},
		{"withdraw everything", 100, 100, 0, false},
		{"overdraft", 100, 150, 100, true},
		{"zero amount", 100, 0, 100, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			acc := New("Asha", tc.start)

			err := acc.Withdraw(tc.amount)

			if (err != nil) != tc.wantErr {
				t.Fatalf("Withdraw(%v) error = %v, wantErr %v", tc.amount, err, tc.wantErr)
			}
			if acc.Balance() != tc.wantBalance {
				t.Errorf("Balance() = %v, want %v", acc.Balance(), tc.wantBalance)
			}
		})
	}
}

func TestOverdraftError(t *testing.T) {
	acc := New("Asha", 100)

	err := acc.Withdraw(150)

	var fundsErr *InsufficientFundsError
	if !errors.As(err, &fundsErr) {
		t.Fatalf("expected InsufficientFundsError, got %v", err)
	}
	if fundsErr.Balance != 100 || fundsErr.Requested != 150 {
		t.Errorf("got Balance=%v Requested=%v, want 100 and 150",
			fundsErr.Balance, fundsErr.Requested)
	}
}
