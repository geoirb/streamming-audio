package storage

import (
	"io"
)

// Storage ...
type Storage struct{}

// List storage
func (s *Storage) List() io.ReadWriteCloser {
	return &queue{}
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{}
}
