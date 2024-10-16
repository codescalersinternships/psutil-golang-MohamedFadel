package psutils


import (
	"strconv"
	"strings"
)

type MemoryInfo struct {
	TotalMemory     int
	UsedMemory      int
	AvailableMemory int
	FreeMemory      int
	Buffers         int
	Cached          int
}

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
