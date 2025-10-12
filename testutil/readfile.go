package testutil

import (
	"fmt"
	"io"
	"os"
)

func ReadFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Error opening file: %v\n", err)
	}
	defer file.Close()

	fmt.Printf("Successfully opened file: %s\n", filePath)

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error reading file: %v\n", err)
	}

	return string(content), nil
}
