package psutils

import (
	"fmt"
	"os"
	"strconv"
)

type ProcessInfo struct {
	PID   int
	Name  string
	State string
}

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
