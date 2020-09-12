package storage

// List storage
type List interface {
	Push([]byte)
	Pop() []byte
}

// StorageFactory ...
type StorageFactory struct{}

// NewList storage
func (s *StorageFactory) NewList() List {
	return &list{}
}

// NewStorageFactory ...
func NewStorageFactory() *StorageFactory {
	return &StorageFactory{}
}
