package firestore

import (
	"nausea-admin/internal/models"

	"cloud.google.com/go/firestore"
)

type FirestoreLink struct {
	URL  string `firestore:"url"`
	Text string `firestore:"text"`
}

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
		link.ID = models.ID{}
		links = append(links, link)
	}
	return links, err
}

func (f *Firestore) CreateLink(l models.Link) (models.Link, error) {
	_, _, err := f.colLinks().Add(f.ctx, linkToFirestoreLink(l))
	return models.Link{
		ID:   models.ID{},
		URL:  l.URL,
		Text: l.Text,
	}, err
}

func (f *Firestore) SetLink(l models.Link) (models.Link, error) {
	u, err := structToUpdate(l)
	if err != nil {
		return models.Link{}, err
	}
	_, err = f.colLinks().Doc(l.ID.ID).Update(f.ctx, *u)
	return l, err
}

func (f *Firestore) DeleteLink(id string) error {
	_, err := f.colLinks().Doc(id).Delete(f.ctx)
	return err
}
