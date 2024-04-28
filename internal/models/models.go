package models

import "time"

type Timestamps struct {
	CreatedAt time.Time  `firestore:"createdAt"`
	UpdatedAt time.Time  `firestore:"updatedAt"`
	DeletedAt *time.Time `firestore:"deletedAt"`
}

type ID struct {
	ID string `firestore:"id"`
}

type Contacts struct {
	Links string `firestore:"links"`
	Timestamps
}

type About struct {
	Bio   string `firestore:"bio,omitempty"`
	Image *Media `firestore:"image,omitempty"`
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
