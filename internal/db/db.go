package db

import "nausea-admin/internal/models"

type IDB interface {
	GetBio() (string, error)
	SetBio(string) error

	GetEmail() (string, error)
	SetEmail(string) error

	GetLinks() ([]models.Link, error)
	SetLink(models.Link) (models.Link, error)
	DeleteLink(string) error
}

type DB struct {
	client IDB
}

func NewDB(db IDB) *DB {
	return &DB{
		client: db,
	}
}

func (db *DB) GetBio() (string, error) {
	return db.client.GetBio()
}

func (db *DB) SetBio(b string) error {
	return db.client.SetBio(b)
}

func (db *DB) GetEmail() (string, error) {
	return db.client.GetEmail()
}

func (db *DB) SetEmail(b string) error {
	return db.client.SetEmail(b)
}

func (db *DB) GetLinks() ([]models.Link, error) {
	return db.client.GetLinks()
}

func (db *DB) SetLink(l models.Link) (models.Link, error) {
	return db.client.SetLink(l)
}

func (db *DB) DeleteLink(id string) error {
	return db.client.DeleteLink(id)
}
