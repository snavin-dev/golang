// Package fileops saves a number to a file and reads it back.
package fileops

import (
	"fmt"
	"os"
	"strconv"
)

// ReadFloat is exported (capital R), so main can call fileops.ReadFloat.
func ReadFloat(fileName string) (float64, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return 0, fmt.Errorf("could not read %s: %w", fileName, err)
	}
	value, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return 0, fmt.Errorf("could not understand the value in %s: %w", fileName, err)
	}
	return value, nil
}

// WriteFloat is exported (capital W).
func WriteFloat(fileName string, value float64) error {
	return os.WriteFile(fileName, []byte(fmt.Sprint(value)), 0644)
}
