package storage

import (
	"container/list"
	"io"
)

// Storage ...
type Storage struct{}

// List storage
func (s *Storage) List() io.ReadWriteCloser {
	return &queue{
		list: list.New(),
	}
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{}
}
