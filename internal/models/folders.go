package models

import (
	"github.com/google/uuid"
)

const (
	FolderContentType string = "folder"
	ImageContentType  string = "image"
)

type Content interface {
	GetType() string
}

type Folder struct {
	ID
	Timestamps
	ParentID string    `firestore:"parentID" json:"parentID"`
	Name     string    `firestore:"name" json:"name"`
	Contents []Content `firestore:"contents" json:"contents"`
}

type ContentBase struct {
	ID
	Timestamps
	Type  string `firestore:"type" json:"type"`
	RefID string `firestore:"refID" json:"refID"`
}

type FolderContent struct {
	ContentBase
	Name string `firestore:"name" json:"name"`
}

type ImageContent struct {
	ContentBase
	URL         string `firestore:"url" json:"url"`
	Name        string `firestore:"name" json:"name"`
	Description string `firestore:"description" json:"description"`
}

func NewFolder(name string, parentID string) (*Folder, error) {
	ic, _ := NewImageContent("", "", "", "")
	f := &Folder{
		Name:     name,
		ParentID: parentID,
		Contents: []Content{ic},
	}
	err := f.generateID()
	return f, err
}

func NewFolderContent(refID string, name string) (*FolderContent, error) {
	fc := FolderContent{
		Name: name,
		ContentBase: ContentBase{
			RefID:      refID,
			Type:       FolderContentType,
			Timestamps: NewTimestamps(),
		},
	}
	err := fc.generateID()
	return &fc, err
}

func NewImageContent(refID string, URL string, name string, description string) (*ImageContent, error) {
	ic := ImageContent{
		Name:        name,
		URL:         URL,
		Description: description,
		ContentBase: ContentBase{
			RefID:      refID,
			Type:       ImageContentType,
			Timestamps: NewTimestamps(),
		},
	}
	err := ic.generateID()
	return &ic, err
}

func (id *ID) generateID() error {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	id.ID = uuid.String()
	return nil
}

func (ic ImageContent) GetType() string {
	return ImageContentType
}

func (fc FolderContent) GetType() string {
	return FolderContentType
}

func (f Folder) GetFolderContents() []FolderContent {
	fc := []FolderContent{}
	for _, content := range f.Contents {
		switch content.(type) {
		case FolderContent:
			fc = append(fc, content.(FolderContent))
		default:
			continue
		}
	}
	return fc
}

func (f Folder) GetImageContents() []ImageContent {
	ic := []ImageContent{}
	for _, content := range f.Contents {
		switch content.(type) {
		case ImageContent:
			ic = append(ic, content.(ImageContent))
		default:
			continue
		}
	}
	return ic
}
