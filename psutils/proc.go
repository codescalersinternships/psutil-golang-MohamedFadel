package psutils

import (
	"fmt"
	"os"
	"strconv"
)

// ProcessInfo contains basic information about a process.
type ProcessInfo struct {
	PID   int    // Process ID
	Name  string // Name of the process
	State string // Current state of the process
}

/*
GetProcessInfo retrieves and returns information about all running processes.
It reads from the /proc filesystem to gather this information.
Returns a pointer to a slice of ProcessInfo and an error if any occurred during the process.
*/
func GetProcessInfo() (*[]ProcessInfo, error) {
	var procInfoList []ProcessInfo
	var procInfo ProcessInfo
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if pid == 0 {
			continue
		}
		if err != nil {
			return nil, err
		}

		lines, err := openAndReadFile(fmt.Sprintf("/proc/%d/status", pid))
		if err != nil {
			return nil, err
		}
		name := getFieldValue(lines, "Name")
		state := getFieldValue(lines, "State")

		procInfo.PID = pid
		procInfo.Name = name
		procInfo.State = state
		procInfoList = append(procInfoList, procInfo)
	}

	return &procInfoList, nil
}
