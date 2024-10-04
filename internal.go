package psutilgolangmohamedfadel

import (
	"fmt"
	"os"
	"strings"
)

func openAndReadFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadingFile, err)
	}

	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func getFieldValue(lines []string, field string) string {
	for _, l := range lines {
		line := strings.Split(l, ":")

		if strings.TrimSpace(line[0]) != field {
			continue
		}

		return strings.TrimSpace(line[1])
	}

	return ""
}
