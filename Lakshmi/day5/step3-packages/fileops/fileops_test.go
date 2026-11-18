package fileops

import (
	"errors"
	"io/fs"
	"path/filepath"
	"testing"
)

func TestWriteThenRead(t *testing.T) {
	file := filepath.Join(t.TempDir(), "balance.txt") // a temporary folder for this test

	if err := WriteFloat(file, 1250.5); err != nil {
		t.Fatalf("WriteFloat returned an error: %v", err)
	}

	got, err := ReadFloat(file)
	if err != nil {
		t.Fatalf("ReadFloat returned an error: %v", err)
	}
	if got != 1250.5 {
		t.Errorf("ReadFloat() = %v, want 1250.5", got)
	}
}

func TestReadMissingFile(t *testing.T) {
	_, err := ReadFloat(filepath.Join(t.TempDir(), "missing.txt"))

	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("expected a file-not-found error, got %v", err)
	}
}
