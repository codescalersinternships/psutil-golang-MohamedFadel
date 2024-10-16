package psutilgolangmohamedfadel

import (
	"io/fs"
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
	defer func() {
		openAndReadFile = originalOpenAndReadFile
		readFile = os.ReadFile
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
			switch name {
			case "/sys/devices/system/cpu/cpu0/cache/index0/size":
				return []byte("32K\n"), nil
			case "/sys/devices/system/cpu/cpu0/cache/index1/size":
				return []byte("32K\n"), nil
			case "/sys/devices/system/cpu/cpu0/cache/index2/size":
				return []byte("256K\n"), nil
			case "/sys/devices/system/cpu/cpu0/cache/index3/size":
				return []byte("6144K\n"), nil
			case "/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_min_freq":
				return []byte("400000\n"), nil
			case "/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq":
				return []byte("3400000\n"), nil
			case "/sys/devices/system/cpu/cpu1/cpufreq/cpuinfo_min_freq":
				return []byte("400000\n"), nil
			case "/sys/devices/system/cpu/cpu1/cpufreq/cpuinfo_max_freq":
				return []byte("3400000\n"), nil
			default:
				return nil, os.ErrNotExist
			}
		}

		expected := &CPUInfo{
			NumOfCores: 2,
			ModelName:  "Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz",
			CacheSize:  6784,
			CPUMHz:     3400.0,
			Frequency: []CPUCoreFreq{
				{Core: "cpu0", MinFreq: 400000, MaxFreq: 3400000},
				{Core: "cpu1", MinFreq: 400000, MaxFreq: 3400000},
			},
		}

		result, err := GetCPUInfo()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("GetCPUInfo() = %+v, want %+v", result, expected)
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
