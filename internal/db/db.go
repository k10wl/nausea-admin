package db

import "nausea-admin/internal/models"

type IDB interface {
	GetAbout() (models.About, error)
	SetAbout(models.About) error

	GetEmail() (string, error)
	SetEmail(string) error

	GetLinks() ([]models.Link, error)
	CreateLink(models.Link) (models.Link, error)
	SetLink(models.Link) (models.Link, error)
	DeleteLink(string) error

	GetFolders() ([]models.Folder, error)
	GetFolderContents()
	CreateFolder(models.Folder) error
	SetFolder(models.Folder) error
	DeleteFolder(id string) error
}

type DB struct {
	client IDB
}

func NewDB(db IDB) *DB {
	return &DB{
		client: db,
	}
}

func (db *DB) GetAbout() (models.About, error) {
	return db.client.GetAbout()
}

func (db *DB) SetAbout(a models.About) error {
	return db.client.SetAbout(a)
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

func (db *DB) CreateLink(l models.Link) (models.Link, error) {
	return db.client.CreateLink(l)
}

func (db *DB) SetLink(l models.Link) (models.Link, error) {
	return db.client.SetLink(l)
}

func (db *DB) DeleteLink(id string) error {
	return db.client.DeleteLink(id)
}
