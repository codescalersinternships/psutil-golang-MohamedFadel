package psutils

import (
	"io/fs"
	"os"
	"path/filepath"
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

var testDir string

func init() {
	readFile = os.ReadFile
	readDir = os.ReadDir
}

func setupTestFiles() error {
	testDir = "./mock_sys/devices/system/cpu/"

	err := os.MkdirAll(testDir+"cpu0/cpufreq", os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(testDir+"cpu1/cpufreq", os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(testDir+"cpu0/cache/index0", os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(testDir+"cpu0/cache/index1", os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(testDir+"cpu0/cache/index2", os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(testDir+"cpu0/cache/index3", os.ModePerm)
	if err != nil {
		return err
	}

	files := map[string]string{
		filepath.Join(testDir, "cpu0/cpufreq/cpuinfo_min_freq"): "400000\n",
		filepath.Join(testDir, "cpu0/cpufreq/cpuinfo_max_freq"): "3400000\n",
		filepath.Join(testDir, "cpu1/cpufreq/cpuinfo_min_freq"): "400000\n",
		filepath.Join(testDir, "cpu1/cpufreq/cpuinfo_max_freq"): "3400000\n",
		filepath.Join(testDir, "cpu0/cache/index0/size"):        "32K\n",
		filepath.Join(testDir, "cpu0/cache/index1/size"):        "32K\n",
		filepath.Join(testDir, "cpu0/cache/index2/size"):        "256K\n",
		filepath.Join(testDir, "cpu0/cache/index3/size"):        "6144K\n",
	}

	for path, content := range files {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func cleanupTestFiles() {
	os.RemoveAll("./mock_sys/")
}

func TestGetCPUInfo(t *testing.T) {
	err := setupTestFiles()
	if err != nil {
		t.Fatalf("Failed to set up test files: %v", err)
	}
	defer cleanupTestFiles()

	originalOpenAndReadFile := openAndReadFile
	originalReadFile := readFile
	originalReadDir := readDir
	defer func() {
		openAndReadFile = originalOpenAndReadFile
		readFile = originalReadFile
		readDir = originalReadDir
	}()

	readFile = func(name string) ([]byte, error) {
		mockedPath := strings.Replace(name, "/sys/devices/system/cpu/", testDir, 1)
		return os.ReadFile(mockedPath)
	}

	readDir = func(name string) ([]fs.DirEntry, error) {
		return []fs.DirEntry{
			&testDirEntry{name: "cpu0", isDir: true},
			&testDirEntry{name: "cpu1", isDir: true},
		}, nil
	}

	t.Run("Successful CPU info retrieval", func(t *testing.T) {
		openAndReadFile = func(path string) ([]string, error) {
			if path != "/proc/cpuinfo" {
				t.Fatalf("Unexpected file path: %s", path)
			}
			return []string{
				"processor : 0",
				"model name : Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz",
				"cpu MHz : 3400.000",
				"processor : 1",
				"model name : Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz",
				"cpu MHz : 3400.000",
			}, nil
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
		readFile = func(name string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		_, err := GetCPUInfo()
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})
}

type testDirEntry struct {
	name  string
	isDir bool
}

func (e *testDirEntry) Name() string               { return e.name }
func (e *testDirEntry) IsDir() bool                { return e.isDir }
func (e *testDirEntry) Type() fs.FileMode          { return fs.ModeDir }
func (e *testDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

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
