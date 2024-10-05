package psutilgolangmohamedfadel

import (
	"os"
	"reflect"
	"testing"
)

func TestGetMemoryInfo(t *testing.T) {
	originalOpenAndReadFile := openAndReadFile
	defer func() {
		openAndReadFile = originalOpenAndReadFile
	}()

	t.Run("Successful memory info retrieval", func(t *testing.T) {
		openAndReadFile = func(path string) ([]string, error) {
			if path != "/proc/meminfo" {
				t.Fatalf("Unexpected file path: %s", path)
			}
			return []string{
				"MemTotal:       2025 kB",
				"MemFree:         518 kB",
				"MemAvailable:   1298 kB",
				"Buffers:          103 kB",
				"Cached:          745 kB",
			}, nil
		}

		expected := &MemoryInfo{
			TotalMemory:     2025,
			UsedMemory:      659,
			AvailableMemory: 1298,
			FreeMemory:      518,
			Buffers:         103,
			Cached:          745,
		}

		result, err := GetMemoryInfo()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("GetMemoryInfo() = %+v, want %+v", result, expected)
		}
	})

	t.Run("File read error", func(t *testing.T) {
		openAndReadFile = func(path string) ([]string, error) {
			return nil, os.ErrNotExist
		}

		_, err := GetMemoryInfo()
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})

	t.Run("Parsing error", func(t *testing.T) {
		openAndReadFile = func(path string) ([]string, error) {
			return []string{
				"MemTotal:       InvalidValue kB",
			}, nil
		}

		_, err := GetMemoryInfo()
		if err == nil {
			t.Errorf("Expected a parsing error, got nil")
		}
	})
}
