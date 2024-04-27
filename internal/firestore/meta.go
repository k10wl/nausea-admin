package firestore

import (
	"nausea-admin/internal/models"
)

func (f *Firestore) GetMeta() (models.Meta, error) {
	var meta models.Meta
	snapshot, err := f.docMeta().Get(f.ctx)
	if err != nil {
		return models.Meta{}, err
	}
	err = snapshot.DataTo(&meta)
	return meta, err
}

func (f *Firestore) SetMeta(meta models.Meta) error {
	_, err := f.docMeta().Set(f.ctx, meta)
	return err
}
