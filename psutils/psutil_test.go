package psutils

import (
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"strings"
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

var (
	readFile func(name string) ([]byte, error)
	readDir  func(name string) ([]fs.DirEntry, error)
)

func init() {
	readFile = os.ReadFile
	readDir = os.ReadDir
}

func TestGetCPUInfo(t *testing.T) {
	originalOpenAndReadFile := openAndReadFile
	originalReadFile := readFile
	defer func() {
		openAndReadFile = originalOpenAndReadFile
		readFile = originalReadFile
	}()

	t.Run("Successful CPU info retrieval", func(t *testing.T) {
		openAndReadFile = func(path string) ([]string, error) {
			if path != "/proc/cpuinfo" {
				t.Fatalf("Unexpected file path: %s", path)
			}
			return []string{
				"processor	: 0",
				"model name	: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz",
				"cpu MHz		: 3400.000",
				"processor	: 1",
				"model name	: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz",
				"cpu MHz		: 3400.000",
			}, nil
		}

		readFile = func(name string) ([]byte, error) {
			switch {
			case strings.Contains(name, "cache/index"):
				return []byte("32K\n"), nil
			case strings.Contains(name, "cpufreq/cpuinfo_min_freq"),
				strings.Contains(name, "cpufreq/cpuinfo_max_freq"):
				return nil, os.ErrNotExist
			default:
				return nil, fmt.Errorf("unexpected file: %s", name)
			}
		}

		result, err := GetCPUInfo()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.NumOfCores != 2 {
			t.Errorf("Expected NumOfCores to be 2, got %d", result.NumOfCores)
		}
		if result.ModelName != "Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz" {
			t.Errorf("Unexpected ModelName: %s", result.ModelName)
		}
		if result.CPUMHz != 3400.0 {
			t.Errorf("Expected CPUMHz to be 3400.0, got %f", result.CPUMHz)
		}

		if result.CacheSize != 0 && result.CacheSize != 6784 {
			t.Errorf("Unexpected CacheSize: %d", result.CacheSize)
		}

		if len(result.Frequency) > 0 {
			if len(result.Frequency) != 2 {
				t.Errorf("Expected 2 Frequency entries, got %d", len(result.Frequency))
			}
			for i, freq := range result.Frequency {
				if freq.Core != fmt.Sprintf("cpu%d", i) {
					t.Errorf("Unexpected Core name: %s", freq.Core)
				}
				if freq.MinFreq != 400000 && freq.MinFreq != 0 {
					t.Errorf("Unexpected MinFreq: %d", freq.MinFreq)
				}
				if freq.MaxFreq != 3400000 && freq.MaxFreq != 0 {
					t.Errorf("Unexpected MaxFreq: %d", freq.MaxFreq)
				}
			}
		}
	})

	t.Run("File read error", func(t *testing.T) {
		openAndReadFile = func(path string) ([]string, error) {
			return nil, os.ErrNotExist
		}

		_, err := GetCPUInfo()
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})
}

func TestGetProcessInfo(t *testing.T) {
	if _, err := os.Stat("/proc"); os.IsNotExist(err) {
		t.Skip("Skipping test as /proc doesn't exist on this system")
	}

	t.Run("BasicFunctionality", func(t *testing.T) {
		procInfo, err := GetProcessInfo()
		if err != nil {
			t.Fatalf("GetProcessInfo() returned an error: %v", err)
		}

		if procInfo == nil {
			t.Fatal("GetProcessInfo() returned nil")
		}

		if len(*procInfo) == 0 {
			t.Fatal("GetProcessInfo() returned an empty slice")
		}

		found := make(map[string]bool)
		for _, info := range *procInfo {
			if info.Name == "systemd" {
				found["systemd"] = true
			}

		}

		if !found["systemd"] {
			t.Error("Did not find systemd process")
		}

		for _, info := range *procInfo {
			if info.PID <= 0 {
				t.Errorf("Invalid PID: %d", info.PID)
			}
			if info.Name == "" {
				t.Errorf("Empty process name for PID: %d", info.PID)
			}
			if info.State == "" {
				t.Errorf("Empty process state for PID: %d", info.PID)
			}
		}
	})

}
