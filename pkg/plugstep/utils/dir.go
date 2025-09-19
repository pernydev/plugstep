package utils

import "os"

func EnsureDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
