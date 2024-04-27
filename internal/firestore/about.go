package firestore

import (
	"nausea-admin/internal/models"

	"cloud.google.com/go/firestore"
)

func (f *Firestore) GetAbout() (models.About, error) {
	var about models.About
	doc, err := f.docAbout().Get(f.ctx)
	if err != nil {
		return about, err
	}
	err = doc.DataTo(&about)
	return about, err
}

func (f *Firestore) SetAbout(about models.About) error {
	updates := []firestore.Update{{Path: "bio", Value: about.Bio}}
	if about.Image != nil {
		updates = append(
			updates,
			firestore.Update{Path: "image", Value: about.Image},
		)
	}
	_, err := f.docAbout().Update(f.ctx, updates)
	return err
}
