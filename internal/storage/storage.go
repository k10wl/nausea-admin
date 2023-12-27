package storage

import "io"

type IStorage interface {
	AddObject(file io.Reader, name string) (string, error)
}

type Storage struct {
	storage IStorage
}

func NewStorage(storage IStorage) *Storage {
	return &Storage{
		storage: storage,
	}
}

func (s *Storage) AddObject(file io.Reader, name string) (string, error) {
	return s.storage.AddObject(file, name)
}
