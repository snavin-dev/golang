package main

import (
	"fmt"
	"os"
	"strconv"
)

const balanceFile = "balance.txt"

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
