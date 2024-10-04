package psutilgolangmohamedfadel

import (
	"os"
	"reflect"
	"testing"
)

func TestOpenAndReadFile(t *testing.T) {
	t.Run("File exists", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatalf("Failed to create temp file : %v", err)
		}
		defer os.Remove(tmpFile.Name())

		testContent := "line1\nline2\nline3"
		if _, err := tmpFile.WriteString(testContent); err != nil {
			t.Fatalf("Failed to write test content to temporary file: %v", err)
		}

		if err := tmpFile.Close(); err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}

		lines, err := openAndReadFile(tmpFile.Name())
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		expected := []string{"line1", "line2", "line3"}

		if !reflect.DeepEqual(lines, expected) {
			t.Fatalf("Expected %v got %v", expected, lines)
		}
	})

	t.Run("File does not exist", func(t *testing.T) {
		_, err := openAndReadFile("nonexistentfile.txt")
		if err == nil {
			t.Errorf("Expected an error for a non-existent file, got nil")
		}
	})
}
