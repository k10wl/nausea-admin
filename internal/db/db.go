package db

import "nausea-admin/internal/models"

type IDB interface {
	GetInfo() (*models.Info, error)
	WriteInfo(info models.Info) error
}

type DB struct {
	client IDB
}

func NewDB(db IDB) *DB {
	return &DB{
		client: db,
	}
}

func (db *DB) GetInfo() (*models.Info, error) {
	return db.client.GetInfo()
}

func (db *DB) WriteInfo(info models.Info) error {
	return db.client.WriteInfo(info)
}
