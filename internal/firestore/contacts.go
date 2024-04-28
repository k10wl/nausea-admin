package firestore

import "nausea-admin/internal/models"

func (f *Firestore) GetContacts() (models.Contacts, error) {
	var contacts models.Contacts
	snapshot, err := f.docContacts().Get(f.ctx)
	if err != nil {
		return models.Contacts{}, err
	}
	err = snapshot.DataTo(&contacts)
	return contacts, err
}

func (f *Firestore) SetContacts(contacts models.Contacts) error {
	_, err := f.docContacts().Set(f.ctx, contacts)
	return err
}
