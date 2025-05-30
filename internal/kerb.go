package internal

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	kerbHeader = "# This file is managed by kerb. DO NOT EDIT."
)

// SyncKerbFiles removes files containing KerbHeader from dstDirectoryPath,
// then copies files containing KerbHeader from srcDirectoryPath to dstDirectoryPath,
// preserving their relative paths.
func SyncKerbFiles(srcDirectoryPath, dstDirectoryPath string) error {
	// 1. Remove KerbHeader files from dstDirectoryPath
	dstFiles, err := findFilesWithString(dstDirectoryPath, kerbHeader)
	if err != nil {
		return err
	}
	if err := removeFiles(dstFiles); err != nil {
		return err
	}

	// 2. Find KerbHeader files in srcDirectoryPath
	srcFiles, err := findFilesWithString(srcDirectoryPath, kerbHeader)
	if err != nil {
		return err
	}
	for _, srcFile := range srcFiles {
		relPath, err := filepath.Rel(srcDirectoryPath, srcFile)
		if err != nil {
			return err
		}
		dstFile := filepath.Join(dstDirectoryPath, relPath)

		// Ensure the destination directory exists
		dstDir := filepath.Dir(dstFile)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return err
		}

		// Copy the file
		srcF, err := os.Open(srcFile)
		if err != nil {
			return err
		}
		dstF, err := os.Create(dstFile)
		if err != nil {
			srcF.Close()
			return err
		}
		_, err = io.Copy(dstF, srcF)
		srcF.Close()
		dstF.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// findFilesWithString recursively searches all files under rootDir and
// returns the paths of files that contain targetStr.
func findFilesWithString(rootDir string, targetStr string) ([]string, error) {
	var result []string
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		file, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(file), targetStr) {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// removeFiles deletes all files specified in the paths slice.
// It returns an error if any file could not be deleted.
func removeFiles(paths []string) error {
	for _, path := range paths {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

// HasKerbHeader returns true if the file at path contains kerbHeader.
func HasKerbHeader(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	return strings.Contains(string(data), kerbHeader), nil
}

// InsertKerbHeader inserts kerbHeader at the beginning of the file at path.
// If the file already contains kerbHeader at the beginning, it does nothing.
func InsertKerbHeader(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(data) >= len(kerbHeader) && string(data[:len(kerbHeader)]) == kerbHeader {
		// Already has kerbHeader at the beginning
		return nil
	}
	newData := append([]byte(kerbHeader+"\n"), data...)
	return os.WriteFile(path, newData, 0644)
}

// ReplaceInKerbFile replaces all occurrences of 'old' with 'new' in the file at path,
// but only if HasKerbHeader(path) returns true. If the header is not present, it does nothing.
func ReplaceInKerbFile(path, old, new string) error {
	has, err := HasKerbHeader(path)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	replaced := strings.ReplaceAll(string(data), old, new)
	return os.WriteFile(path, []byte(replaced), 0644)
}

// ListKerbFiles recursively searches all files under dir and returns the paths of files
// for which HasKerbHeader returns true.
func ListKerbFiles(dir string) ([]string, error) {
	var result []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		has, err := HasKerbHeader(path)
		if err != nil {
			return err
		}
		if has {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
