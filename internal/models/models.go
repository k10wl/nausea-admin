package models

import "time"

type Meta struct {
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt time.Time `firestore:"updatedAt" json:"updatedAt,omitempty"`
	DeletedAt time.Time `firestore:"deletedAt" json:"deletedAt,omitempty"`
}

type Link struct {
	ID   string
	URL  string `firestore:"url" json:"url"`
	Text string `firestore:"text" json:"text"`
}

type Contacts struct {
	Email string `firestore:"email" json:"email,omitempty"`
	Links []Link `firestore:"links" json:"links,omitempty"`
}

type About struct {
	Bio string `firestore:"bio"`
}

type Folder struct {
	ID   string
	Name string `firestore:"name" json:"name,omitempty"`
	// ROOT or valid ID of other folder
	ParentID string `firestore:"parentId" json:"parentId,omitempty"`
	// Media is copied over on folder level because fIrEbAsE fIrEsToRe reasons.
	Media []Media `firestore:"media" json:"media,omitempty"`
	Meta
}

type Media struct {
	ID string
	// Storage URL
	URL string
	Meta
}
