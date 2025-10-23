package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func AbsPath(p string) (string, error) {
	absolutePath, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}
	return absolutePath, nil
}

func CheckDestDirExists(dir string) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("destination directory does not exist: %s", dir)
	}
	if err != nil {
		return fmt.Errorf("failed checking destination directory: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("destination path is not a directory: %s", dir)
	}
	return nil
}

func CleanupTestFiles(s Settings) []error {
	var errs []error

	for i := range s.Streams {
		dstFile := fmt.Sprintf("%s/iospeedtest_%d.txt", s.DestDir, i)
		if err := os.Remove(dstFile); err != nil {
			errs = append(errs, fmt.Errorf("error removing destination file %s: %v", dstFile, err))
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
