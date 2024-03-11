package db

import "nausea-admin/internal/models"

type IDB interface {
	GetFolderByID(string) (models.Folder, error)
	CreateFolder(models.Folder) (models.Folder, models.FolderContent, error)
	MarkFolderDeletedByID(string) (models.Folder, error)
	MarkFolderRestoredByID(string) (models.Folder, error)
}

type DB struct {
	client IDB
}

func NewDB(db IDB) *DB {
	return &DB{
		client: db,
	}
}

func (db DB) GetFolderByID(id string) (models.Folder, error) {
	return db.client.GetFolderByID(id)
}

func (db DB) CreateFolder(folder models.Folder) (models.Folder, models.FolderContent, error) {
	return db.client.CreateFolder(folder)
}

func (db DB) MarkFolderDeletedByID(id string) (models.Folder, error) {
	return db.client.MarkFolderDeletedByID(id)
}

func (db DB) MarkFolderRestoredByID(id string) (models.Folder, error) {
	return db.client.MarkFolderRestoredByID(id)
}
