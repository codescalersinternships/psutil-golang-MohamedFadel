package psutils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CPUCoreFreq represents the frequency information for a CPU core.
type CPUCoreFreq struct {
	Core    string // The name of the CPU core
	MinFreq int    // The minimum frequency of the core in Hz
	MaxFreq int    // The maximum frequency of the core in Hz
}

// CPUInfo contains various information about the CPU.
type CPUInfo struct {
	NumOfCores int           // The number of CPU cores
	ModelName  string        // The model name of the CPU
	CacheSize  int           // The total cache size in KB
	CPUMHz     float64       // The average CPU frequency in MHz
	Frequency  []CPUCoreFreq // Frequency information for each core
}

/*
GetCPUInfo retrieves and returns information about the CPU.
It reads various system files to gather this information.
Returns a pointer to CPUInfo and an error if any occurred during the process.
*/
func GetCPUInfo() (*CPUInfo, error) {
	var cpuInfo CPUInfo

	data, err := openAndReadFile("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}

	cpuInfo.ModelName = getFieldValue(data, "model name")

	for _, line := range data {
		if strings.Contains(line, "processor") {
			cpuInfo.NumOfCores++
		}
	}

	var totalMHz float64
	for _, line := range data {
		if strings.Contains(line, "cpu MHz") {
			if mhz, err := strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]), 64); err == nil {
				totalMHz += mhz
			}
		}
	}
	cpuInfo.CPUMHz = totalMHz / float64(cpuInfo.NumOfCores)

	d, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index0/size")
	if err != nil {
		return nil, err
	}

	value := strings.Split(string(d), "\n")[0]
	value = value[:len(value)-1]
	l1iCache, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}

	d2, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index1/size")
	if err != nil {
		return nil, err
	}

	value2 := strings.Split(string(d2), "\n")[0]
	value2 = value2[:len(value2)-1]
	l1dCache, err := strconv.Atoi(value2)
	if err != nil {
		return nil, err
	}

	d3, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index2/size")
	if err != nil {
		return nil, err
	}

	value3 := strings.Split(string(d3), "\n")[0]
	value3 = value3[:len(value3)-1]
	l2Cache, err := strconv.Atoi(value3)
	if err != nil {
		return nil, err
	}

	d4, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cache/index3/size")
	if err != nil {
		return nil, err
	}

	value4 := strings.Split(string(d4), "\n")[0]
	value4 = value4[:len(value4)-1]
	l3Cache, err := strconv.Atoi(value4)
	if err != nil {
		return nil, err
	}

	cpuInfo.CacheSize = cpuInfo.NumOfCores*(l1iCache+l1dCache+l2Cache) + l3Cache

	var coreNames []string
	for i := range cpuInfo.NumOfCores {
		coreName := []string{"cpu", strconv.Itoa(i)}
		coreNames = append(coreNames, strings.Join(coreName, ""))
	}

	for i := range cpuInfo.NumOfCores {
		var cpuCoreFreq CPUCoreFreq
		cpuCoreFreq.Core = coreNames[i]

		dataMin, err := os.ReadFile(fmt.Sprintf("/sys/devices/system/cpu/%s/cpufreq/cpuinfo_min_freq", coreNames[i]))
		if err != nil {
			return nil, err
		}
		if minFreq, err := strconv.Atoi(strings.Split(string(dataMin), "\n")[0]); err == nil {
			cpuCoreFreq.MinFreq = minFreq
		}

		dataMax, err := os.ReadFile(fmt.Sprintf("/sys/devices/system/cpu/%s/cpufreq/cpuinfo_max_freq", coreNames[i]))
		if err != nil {
			return nil, err
		}
		if maxFreq, err := strconv.Atoi(strings.Split(string(dataMax), "\n")[0]); err == nil {
			cpuCoreFreq.MaxFreq = maxFreq
		}

		cpuInfo.Frequency = append(cpuInfo.Frequency, cpuCoreFreq)

	}

	return &cpuInfo, nil
}
