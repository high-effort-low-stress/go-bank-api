// Package testutil provides utility functions for testing, such as reading file contents for templates or mock data.
package testutil

import (
	"fmt"
	"io"
	"os"
)

func ReadFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	fmt.Printf("Successfully opened file: %s\n", filePath)

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	return string(content), nil
}
