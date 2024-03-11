package firestore

import (
	"context"

	"nausea-admin/internal/models"

	"cloud.google.com/go/firestore"
)

type FirestoreLink struct {
	URL  string `firestore:"url"`
	Text string `firestore:"text"`
}

func (f *Firestore) SetAbout(models.About) error {
	return nil
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

func (f *Firestore) GetFolderByID(id string) (models.Folder, error) {
	var folder models.Folder
	doc, err := f.collectionFolders().Doc(id).Get(f.ctx)
	if err != nil {
		return folder, err
	}
	err = doc.DataTo(&folder)
	folder.ID.ID = id
	return folder, err
}

func (f *Firestore) CreateFolder(
	folder models.Folder,
) (models.Folder, models.FolderContent, error) {
	folderDoc := f.collectionFolders().Doc(folder.ID.ID)
	parentDoc := f.collectionFolders().Doc(folder.ParentID)
	asContent, err := folder.AsContent()
	if err != nil {
		return folder, asContent, err
	}
	err = f.client.RunTransaction(f.ctx, func(_ context.Context, tx *firestore.Transaction) error {
		parentDocumentSnapshot, err := tx.Get(parentDoc)
		if err != nil {
			return err
		}
		var parentFolder models.Folder
		err = parentDocumentSnapshot.DataTo(&parentFolder)
		if err != nil {
			return err
		}
		err = tx.Create(folderDoc, folder)
		if err != nil {
			return err
		}
		contents := append(parentFolder.FolderContents, asContent)
		err = tx.Update(
			parentDoc,
			[]firestore.Update{{Path: "folders", Value: contents}},
		)

		return err
	})
	return folder, asContent, err
}

func (f *Firestore) MarkFolderDeletedByID(id string) (models.Folder, error) {
	var folder models.Folder
	err := f.client.RunTransaction(f.ctx, func(_ context.Context, tx *firestore.Transaction) error {
		toDeleteDoc := f.collectionFolders().Doc(id)
		toDeleteSnapshot, err := tx.Get(toDeleteDoc)
		if err != nil {
			return err
		}
		var toDeleteFolder models.Folder
		err = toDeleteSnapshot.DataTo(&toDeleteFolder)
		if err != nil {
			return err
		}
		parentDoc := f.collectionFolders().Doc(toDeleteFolder.ParentID)
		parentSnapshot, err := tx.Get(parentDoc)
		if err != nil {
			return err
		}
		var parentFolder models.Folder
		err = parentSnapshot.DataTo(&parentFolder)
		if err != nil {
			return err
		}
		parentFolder.MarkDeletedFolderContents(id)
		err = tx.Update(
			parentDoc,
			[]firestore.Update{{Path: "folders", Value: parentFolder.FolderContents}},
		)
		if err != nil {
			return err
		}
		toDeleteFolder.Delete()
		toDeleteFolder.UpdatedAt = *toDeleteFolder.DeletedAt
		err = tx.Update(
			toDeleteDoc,
			[]firestore.Update{
				{Path: "deletedAt", Value: toDeleteFolder.DeletedAt},
				{Path: "updatedAt", Value: toDeleteFolder.UpdatedAt},
			},
		)
		folder = toDeleteFolder
		folder.ID.ID = id
		return err
	})
	return folder, err
}

func (f *Firestore) MarkFolderRestoredByID(id string) (models.Folder, error) {
	var folder models.Folder
	err := f.client.RunTransaction(f.ctx, func(_ context.Context, tx *firestore.Transaction) error {
		toDeleteDoc := f.collectionFolders().Doc(id)
		toDeleteSnapshot, err := tx.Get(toDeleteDoc)
		if err != nil {
			return err
		}
		var toRestoreFolder models.Folder
		err = toDeleteSnapshot.DataTo(&toRestoreFolder)
		if err != nil {
			return err
		}
		parentDoc := f.collectionFolders().Doc(toRestoreFolder.ParentID)
		parentSnapshot, err := tx.Get(parentDoc)
		if err != nil {
			return err
		}
		var parentFolder models.Folder
		err = parentSnapshot.DataTo(&parentFolder)
		if err != nil {
			return err
		}
		parentFolder.MarkRestoredFolderContents(id)
		err = tx.Update(
			parentDoc,
			[]firestore.Update{
				{Path: "folders", Value: parentFolder.FolderContents},
			},
		)
		if err != nil {
			return err
		}
		toRestoreFolder.Restore()
		toRestoreFolder.Update()
		err = tx.Update(
			toDeleteDoc,
			[]firestore.Update{
				{Path: "deletedAt", Value: toRestoreFolder.DeletedAt},
				{Path: "updatedAt", Value: toRestoreFolder.UpdatedAt},
			},
		)
		folder = toRestoreFolder
		folder.ID.ID = id
		return err
	})
	return folder, err
}
