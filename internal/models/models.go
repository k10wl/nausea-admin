package models

import "time"

type Timestamps struct {
	CreatedAt time.Time  `firestore:"createdAt"`
	UpdatedAt time.Time  `firestore:"updatedAt"`
	DeletedAt *time.Time `firestore:"deletedAt"`
}

type ID struct {
	ID string `firestore:"-"`
}

type Link struct {
	ID
	URL  string `firestore:"url"`
	Text string `firestore:"text"`
	Timestamps
}

type Contacts struct {
	Email string `firestore:"email"`
	Links []Link `firestore:"links"`
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

func (t *Timestamps) Restore() {
	t.DeletedAt = nil
}
