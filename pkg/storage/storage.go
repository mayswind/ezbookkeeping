package storage

import "github.com/mayswind/ezbookkeeping/pkg/core"

// ObjectStorage represents an object storage to store file object
type ObjectStorage interface {
	Exists(ctx core.Context, path string) (bool, error)
	Read(ctx core.Context, path string) (ObjectInStorage, error)
	Save(ctx core.Context, path string, object ObjectInStorage) error
	Delete(ctx core.Context, path string) error
}
