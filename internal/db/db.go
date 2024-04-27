package db

import "nausea-admin/internal/models"

type IDB interface {
	GetFolderByID(string) (models.Folder, error)
	CreateFolder(models.Folder) (models.Folder, models.FolderContent, error)
	MarkFolderDeletedByID(string) (models.Folder, error)
	MarkFolderRestoredByID(string) (models.Folder, error)
	CreateMedia(models.Media) error
	UploadMediaToFolder([]models.MediaContent, string) error
	MarkMediaAsDeletedInFolder(mediaID string, folderID string) (models.MediaContent, error)
	MarkMediaAsRestoredInFolder(mediaID string, folderID string) (models.MediaContent, error)
	PatchFolder(folderID string, patch models.Folder) (models.Folder, error)
	UpdateMediaInFolder(media models.MediaContent) (models.MediaContent, error)
	GetAbout() (models.About, error)
	SetAbout(about models.About) error
	GetMeta() (models.Meta, error)
	SetMeta(models.Meta) error
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

func (db DB) UploadMediaToFolder(media []models.MediaContent, folderID string) error {
	return db.client.UploadMediaToFolder(media, folderID)
}

func (db DB) MarkMediaAsDeletedInFolder(mediaID string, folderID string) (models.MediaContent, error) {
	return db.client.MarkMediaAsDeletedInFolder(mediaID, folderID)
}

func (db DB) MarkMediaAsRestoredInFolder(mediaID string, folderID string) (models.MediaContent, error) {
	return db.client.MarkMediaAsRestoredInFolder(mediaID, folderID)
}

func (db DB) PatchFolder(folderID string, patch models.Folder) (models.Folder, error) {
	return db.client.PatchFolder(folderID, patch)
}

func (db DB) CreateMedia(media models.Media) error {
	return db.client.CreateMedia(media)
}

func (db DB) UpdateMediaInFolder(media models.MediaContent) (models.MediaContent, error) {
	return db.client.UpdateMediaInFolder(media)
}

func (db DB) GetAbout() (models.About, error) {
	return db.client.GetAbout()
}

func (db DB) SetAbout(about models.About) error {
	return db.client.SetAbout(about)
}

func (db DB) GetMeta() (models.Meta, error) {
	return db.client.GetMeta()
}

func (db DB) SetMeta(meta models.Meta) error {
	return db.client.SetMeta(meta)
}
