package storage

import "io"

type IStorage interface {
	AddObject(io.Reader, string) (string, error)
	RemoveObject(string) error
	ParseURLKey(string) string
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

func (s *Storage) RemoveObject(name string) error {
	return s.storage.RemoveObject(name)
}

func (s *Storage) ParseURLKey(name string) string {
	return s.storage.ParseURLKey(name)
}
