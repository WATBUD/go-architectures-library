package utility

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProjectRootDirAndEnvPath() (string, string, error) {
	executablePath, err := os.Executable() // Get the full path of the executable file
	if err != nil {
		return "", "", err
	}

	dir := filepath.Dir(executablePath) // Get the directory where the executable is located
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); !os.IsNotExist(err) {
			envPath := filepath.Join(dir, ".env") // Combine the project root with ".env"
			return dir, envPath, nil              // Return the project root directory and the .env path
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break // Reached the root directory of the filesystem, stop traversing
		}
		dir = parentDir
	}

	return "", "", fmt.Errorf("project root directory not found")
}
