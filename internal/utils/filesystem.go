package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsValidDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %s does not exist", path)
		}
		return fmt.Errorf("error checking the path %s: %v", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}
	return nil
}

func PathToAbsolute(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("error converting path %s to absolute path: %s", path, err.Error())
	}

	return absolutePath, nil
}

func FileExistAndItIsAFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("empty file path")
	}

	currentDir, _ := os.Getwd()

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist in current directory %s", filePath, currentDir)
		}

		return fmt.Errorf("error checking the file %s: %v", filePath, err)
	}

	return nil
}

func FileIsNotEmpty(filepath string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", file, err)
	}

	if len(file) == 0 {
		return fmt.Errorf("file %s is empty", filepath)
	}

	return nil
}

func IsGitRepository(path string) error {
	if path == "" {
		return fmt.Errorf("empty path")
	}

	if err := IsValidDir(path); err != nil {
		return err
	}

	_, err := os.Stat(filepath.Join(path, ".git"))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path %s is not a git repository", path)
		}

		return fmt.Errorf("error checking the path %s: %v", path, err)
	}

	return nil
}

func FindGitRepoDir(pathname string, levels int) (string, error) {
	absPath, err := filepath.Abs(pathname)
	if err != nil {
		return "", fmt.Errorf("error converting path %s to absolute path: %s", pathname, err.Error())
	}
	for i := 0; i < levels; i++ {
		gitPath := filepath.Join(absPath, ".git")
		if stat, err := os.Stat(gitPath); err == nil && stat.IsDir() {
			return absPath, nil
		}
		parentPath := filepath.Dir(absPath)

		// To avoid going beyond the root ("/" or "C:\"), check if we're already at the root
		if parentPath == absPath {
			return "", fmt.Errorf("path %s is not a git repository", pathname)
		}

		absPath = parentPath
	}
	return "", fmt.Errorf("path %s is not a git repository", pathname)
}

// GetFileContent returns the content of a file as a string.
func GetFileContent(filePath string) (string, error) {
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	return string(contentBytes), nil
}
