package storage

import (
	"io"
)

// ObjectInStorage represents the object instance in the storage
type ObjectInStorage interface {
	io.ReadCloser
	io.Seeker
}
