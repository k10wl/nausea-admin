package models

import (
	"github.com/google/uuid"
)

const (
	RootFolderID string = "--ROOT--"
)

type Folder struct {
	ID
	Timestamps
	ParentID       string          `firestore:"parentID"`
	Name           string          `firestore:"name"`
	FolderContents []FolderContent `firestore:"folders"`
	MediaContents  []MediaContent  `firestore:"media"`
}

type ContentBase struct {
	ID
	Timestamps
	RefID string `firestore:"refID"`
}

type FolderContent struct {
	ContentBase
	Name string
}

type MediaContent struct {
	ContentBase
	URL         string
	Name        string
	Description string
}

func NewFolder(parentID string, name string) (*Folder, error) {
	f := &Folder{
		Name:           name,
		ParentID:       parentID,
		FolderContents: []FolderContent{},
		MediaContents:  []MediaContent{},
		Timestamps:     NewTimestamps(),
	}
	err := f.generateID()
	return f, err
}

func NewFolderContent(refID string, name string) (*FolderContent, error) {
	fc := FolderContent{
		Name: name,
		ContentBase: ContentBase{
			RefID:      refID,
			Timestamps: NewTimestamps(),
		},
	}
	err := fc.generateID()
	return &fc, err
}

func (id *ID) generateID() error {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	id.ID = uuid.String()
	return nil
}

func (f Folder) AsContent() (FolderContent, error) {
	folderContent := FolderContent{
		ContentBase: ContentBase{
			Timestamps: f.Timestamps,
			RefID:      f.ID.ID,
		},
		Name: f.Name,
	}
	err := folderContent.generateID()
	return folderContent, err
}

func (f *Folder) MarkDeletedFolderContents(id string) {
	for i, v := range f.FolderContents {
		if v.RefID == id {
			v.Delete()
			v.UpdatedAt = *v.DeletedAt
			f.FolderContents[i] = v
		}
	}
}

func (f *Folder) MarkRestoredFolderContents(id string) {
	for i, v := range f.FolderContents {
		if v.RefID == id {
			v.Restore()
			v.Update()
			f.FolderContents[i] = v
		}
	}
}
