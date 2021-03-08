package utils

import "io"

// IdentReader returns the original io reader
func IdentReader(encoding string, input io.Reader) (io.Reader, error) {
	return input, nil
}
