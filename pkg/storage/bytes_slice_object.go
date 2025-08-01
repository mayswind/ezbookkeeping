package storage

import (
	"bytes"
)

// bytesSliceObject represents a byte slice object in storage
type bytesSliceObject struct {
	*bytes.Reader
}

// Close does nothing because it does not hold any resources that need to be released
func (b *bytesSliceObject) Close() error {
	return nil
}

// newByteSliceObject creates a new byte slice object from the specified byte slice
func newByteSliceObject(data []byte) ObjectInStorage {
	return &bytesSliceObject{
		Reader: bytes.NewReader(data),
	}
}
