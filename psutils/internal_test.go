package psutils

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

func TestGetFieldValue(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		field    string
		expected string
	}{
		{
			name: "Field Exists",
			lines: []string{
				"vendor_id  : GenuineIntel",
				"model name : Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz",
				"cpu MHz : 3399.975",
			},
			field:    "model name",
			expected: "Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz",
		},
		{
			name: "Field Does Not Exist",
			lines: []string{
				"vendor_id  : GenuineIntel",
				"model name : Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz",
				"cpu MHz : 3399.975",
			},
			field:    "cache size",
			expected: "",
		},
		{
			name: "Empty Field Name",
			lines: []string{
				"vendor_id  : GenuineIntel",
				"model name : Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz",
			},
			field:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFieldValue(tt.lines, tt.field)
			if result != tt.expected {
				t.Fatalf("got %v want %v", result, tt.expected)
			}
		})
	}
}
