package psutils

import (
	"strconv"
	"strings"
)

// MemoryInfo contains various statistics about system memory usage.
type MemoryInfo struct {
	TotalMemory     int // Total amount of physical RAM, in KB
	UsedMemory      int // Total used memory, in KB
	AvailableMemory int // Available memory, in KB
	FreeMemory      int // Free memory, in KB
	Buffers         int // Memory used by kernel buffers, in KB
	Cached          int // Memory used by the page cache and slabs, in KB
}

/*
GetMemoryInfo retrieves and returns information about the system's memory usage.
It reads from /proc/meminfo to gather this information.
Returns a pointer to MemoryInfo and an error if any occurred during the process.
*/
func GetMemoryInfo() (*MemoryInfo, error) {
	data, err := openAndReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	var memoryInfo MemoryInfo
	memoryInfo.TotalMemory, err = strconv.Atoi(strings.Split(getFieldValue(data, "MemTotal"), " ")[0])
	if err != nil {
		return nil, err
	}
	memoryInfo.AvailableMemory, err = strconv.Atoi(strings.Split(getFieldValue(data, "MemAvailable"), " ")[0])
	if err != nil {
		return nil, err
	}
	memoryInfo.Buffers, err = strconv.Atoi(strings.Split(getFieldValue(data, "Buffers"), " ")[0])
	if err != nil {
		return nil, err
	}
	memoryInfo.Cached, err = strconv.Atoi(strings.Split(getFieldValue(data, "Cached"), " ")[0])
	if err != nil {
		return nil, err
	}
	memoryInfo.FreeMemory, err = strconv.Atoi(strings.Split(getFieldValue(data, "MemFree"), " ")[0])
	if err != nil {
		return nil, err
	}

	memoryInfo.UsedMemory = memoryInfo.TotalMemory - memoryInfo.FreeMemory - memoryInfo.Buffers - memoryInfo.Cached

	return &memoryInfo, nil
}
