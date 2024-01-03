package firestore

import (
	"nausea-admin/internal/models"

	"cloud.google.com/go/firestore"
)

type FirestoreLink struct {
	URL  string `firestore:"url"`
	Text string `firestore:"text"`
}

func (f *Firestore) GetBio() (string, error) {
	var about models.About
	doc, err := f.docAbout().Get(f.ctx)
	if err != nil {
		return "", err
	}
	err = doc.DataTo(&about)
	return about.Bio, err
}

func (f *Firestore) SetBio(b string) error {
	_, err := f.docAbout().Update(f.ctx, []firestore.Update{{Path: "bio", Value: b}})
	return err
}

func (f *Firestore) GetEmail() (string, error) {
	var contacts models.Contacts
	doc, err := f.docContacts().Get(f.ctx)
	if err != nil {
		return "", err
	}
	err = doc.DataTo(&contacts)
	return contacts.Email, err
}

func (f *Firestore) SetEmail(e string) error {
	_, err := f.docContacts().Update(f.ctx, []firestore.Update{{Path: "email", Value: e}})
	return err
}

func (f *Firestore) GetLinks() ([]models.Link, error) {
	var links []models.Link
	docs, err := f.colLinks().Documents(f.ctx).GetAll()
	if err != nil {
		return []models.Link{}, err
	}
	for _, v := range docs {
		var link models.Link
		err = v.DataTo(&link)
		if err != nil {
			return []models.Link{}, err
		}
		link.ID = v.Ref.ID
		links = append(links, link)
	}
	return links, err
}

func (f *Firestore) CreateLink(l models.Link) (models.Link, error) {
	ref, _, err := f.colLinks().Add(f.ctx, linkToFirestoreLink(l))
	return models.Link{
		ID:   ref.ID,
		URL:  l.URL,
		Text: l.Text,
	}, err
}

func (f *Firestore) SetLink(l models.Link) (models.Link, error) {
	u, err := structToUpdate(l)
	if err != nil {
		return models.Link{}, err
	}
	_, err = f.colLinks().Doc(l.ID).Update(f.ctx, *u)
	return l, err
}

func (f *Firestore) DeleteLink(id string) error {
	_, err := f.colLinks().Doc(id).Delete(f.ctx)
	return err
}
