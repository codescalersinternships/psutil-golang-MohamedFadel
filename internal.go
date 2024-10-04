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

