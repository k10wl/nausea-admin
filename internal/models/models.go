package models

import "time"

type Timestamps struct {
	CreatedAt time.Time  `firestore:"createdAt" json:"createdAt"`
	UpdatedAt time.Time  `firestore:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `firestore:"deletedAt" json:"deletedAt"`
}

type ID struct {
	ID string
}

type Link struct {
	ID
	URL  string `firestore:"url" json:"url"`
	Text string `firestore:"text" json:"text"`
	Timestamps
}

type Contacts struct {
	Email string `firestore:"email" json:"email"`
	Links []Link `firestore:"links" json:"links"`
	Timestamps
}

type About struct {
	Bio string `firestore:"bio"`
	Timestamps
}

func NewTimestamps() Timestamps {
	return Timestamps{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
}

func (t *Timestamps) Update() {
	t.UpdatedAt = time.Now()
}

func (t *Timestamps) Delete() {
	now := time.Now()
	t.DeletedAt = &now
}
