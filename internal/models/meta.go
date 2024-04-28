package models

type Meta struct {
	Background Media `firestore:"background,omitempty"`
	Timestamps
}
