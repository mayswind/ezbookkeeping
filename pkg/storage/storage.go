package storage

// ObjectStorage represents an object storage to store file object
type ObjectStorage interface {
	Exists(path string) (bool, error)
	Read(path string) (ObjectInStorage, error)
	Save(path string, object ObjectInStorage) error
	Delete(path string) error
}
