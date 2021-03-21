package utils

import (
	"io"
	"os"
)

// IsExists returns whether specified file or directory path exits
func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// WriteFile would write file according to specified content
func WriteFile(path string, data []byte) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	n, err := file.Write(data)

	if err == nil && n < len(data) {
		return io.ErrShortWrite
	}

	return err
}

// IdentReader returns the original io reader
func IdentReader(encoding string, input io.Reader) (io.Reader, error) {
	return input, nil
}
