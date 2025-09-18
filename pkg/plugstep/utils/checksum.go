package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"os"
)

func CalculateFileSHA256(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	sum := hasher.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}

func CalculateFileSHA512(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha512.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	sum := hasher.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}
