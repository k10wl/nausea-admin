package firestore

import "nausea-admin/internal/models"

func (f *Firestore) CreateMedia(m models.Media) error {
	doc := f.collectionMedia().Doc(m.ID.ID)
	_, err := doc.Set(f.ctx, m)
	return err
}
